package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// DefaultSolverURL là địa chỉ mặc định của FlareSolverr.
// Docker map cổng 8191 ra localhost trên cả Windows và Linux nên dùng chung một giá trị.
const DefaultSolverURL = "http://127.0.0.1:8191"

const (
	// solverMaxTimeout là thời gian tối đa FlareSolverr dành để giải challenge (mili-giây).
	solverMaxTimeout = 45000
	// sessionTTL là thời gian coi cf_clearance còn hiệu lực.
	// Cloudflare thường cấp cookie sống vài giờ, nhưng giữ ngắn để tránh dùng cookie
	// đã bị vô hiệu sau khi đổi IP/mạng và chỉ phải giải lại khi thật sự cần.
	sessionTTL = 25 * time.Minute
	// maxSolveRounds: số lần giải challenge tối đa cho một request.
	// Lần 2 dành cho trường hợp đổi router/IP khiến cookie cũ mất hiệu lực.
	maxSolveRounds = 2
	// maxSolverResponse giới hạn kích thước response từ FlareSolverr.
	maxSolverResponse = 8 << 20
)

// solverRequest là body gửi tới FlareSolverr POST /v1.
type solverRequest struct {
	Cmd        string `json:"cmd"`
	URL        string `json:"url"`
	MaxTimeout int    `json:"maxTimeout"`
}

type solverCookie struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Domain string  `json:"domain"`
	Path   string  `json:"path"`
	Expiry float64 `json:"expiry"`
}

// solverResponse là response của FlareSolverr POST /v1.
type solverResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Solution struct {
		URL       string         `json:"url"`
		Status    int            `json:"status"`
		Response  string         `json:"response"`
		UserAgent string         `json:"userAgent"`
		Cookies   []solverCookie `json:"cookies"`
	} `json:"solution"`
}

// solverSession là phiên đã vượt challenge: cookie cf_clearance kèm User-Agent đã dùng.
// Cloudflare gắn cf_clearance với cặp (IP, User-Agent), nên các request sau bắt buộc
// phải gửi đúng UA này, nếu không cookie sẽ bị coi là không hợp lệ.
type solverSession struct {
	cookies   []*http.Cookie
	userAgent string
	expires   time.Time
}

func (s solverSession) valid() bool {
	return len(s.cookies) > 0 && !s.expires.IsZero() && time.Now().Before(s.expires)
}

// solverURLFromEnv đọc biến môi trường AVS_SOLVER_URL.
// Không đặt biến này thì dùng DefaultSolverURL; đặt thành "off" để tắt solver.
func solverURLFromEnv() string {
	if v, ok := os.LookupEnv("AVS_SOLVER_URL"); ok {
		return v
	}
	return DefaultSolverURL
}

// SolverURL trả về địa chỉ FlareSolverr đang dùng (rỗng nghĩa là đã tắt).
func (c *Client) SolverURL() string {
	return c.solverURL
}

// SetSolverURL đặt địa chỉ FlareSolverr.
// Chuỗi rỗng, "off" hoặc "none" để tắt hoàn toàn (chỉ dùng request trực tiếp).
func (c *Client) SetSolverURL(raw string) {
	raw = strings.TrimSpace(raw)
	switch strings.ToLower(raw) {
	case "", "off", "none", "false", "0":
		c.solverURL = ""
	default:
		c.solverURL = strings.TrimRight(raw, "/")
	}
	c.invalidateSession()
}

// solverEnabled cho biết có cấu hình FlareSolverr hay không.
func (c *Client) solverEnabled() bool {
	return c.solverURL != ""
}

// userAgent trả về User-Agent cần dùng cho request hiện tại.
// Khi có phiên đã giải challenge, phải dùng đúng UA của phiên đó.
func (c *Client) userAgent() string {
	c.solverMu.RLock()
	defer c.solverMu.RUnlock()

	if c.session.valid() && c.session.userAgent != "" {
		return c.session.userAgent
	}
	return DefaultUA
}

// hasFreshSession kiểm tra đã có phiên còn hiệu lực chưa.
func (c *Client) hasFreshSession() bool {
	c.solverMu.RLock()
	defer c.solverMu.RUnlock()
	return c.session.valid()
}

// invalidateSession xóa phiên hiện tại, buộc lần sau phải giải lại challenge.
// Được gọi khi một request đã có cookie vẫn bị chặn (cookie hết hạn hoặc đổi IP/mạng).
func (c *Client) invalidateSession() {
	c.solverMu.Lock()
	c.session = solverSession{}
	c.solverMu.Unlock()
}

// sessionCookies trả về bản sao cookie của phiên hiện tại (dùng khi http.Client không có cookie jar).
func (c *Client) sessionCookies() []*http.Cookie {
	c.solverMu.RLock()
	defer c.solverMu.RUnlock()

	if !c.session.valid() {
		return nil
	}
	out := make([]*http.Cookie, len(c.session.cookies))
	copy(out, c.session.cookies)
	return out
}

// solve chạy FlareSolverr để lấy cf_clearance rồi nạp cookie + User-Agent vào client.
// Chỉ một request được giải tại một thời điểm; các request song song sẽ dùng lại kết quả.
func (c *Client) solve(ctx context.Context) error {
	if !c.solverEnabled() {
		return ErrSolverUnavailable
	}
	if c.hasFreshSession() {
		return nil
	}

	c.solveMu.Lock()
	defer c.solveMu.Unlock()

	// Request khác có thể vừa giải xong trong lúc ta chờ khóa.
	if c.hasFreshSession() {
		return nil
	}

	payload, err := json.Marshal(solverRequest{
		Cmd:        "request.get",
		URL:        c.baseURL + "/",
		MaxTimeout: solverMaxTimeout,
	})
	if err != nil {
		return fmt.Errorf("không thể tạo yêu cầu FlareSolverr: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.solverURL+"/v1", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("không thể tạo yêu cầu FlareSolverr: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.solverHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("%w tại %s: %v. Khởi động bằng: docker run -d --name flaresolverr -p 8191:8191 flaresolverr/flaresolverr:latest", ErrSolverUnavailable, c.solverURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSolverResponse))
	if err != nil {
		return fmt.Errorf("không đọc được phản hồi FlareSolverr: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FlareSolverr trả về HTTP %d: %s", resp.StatusCode, truncateForLog(string(body), 200))
	}

	var parsed solverResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("phản hồi FlareSolverr không phải JSON hợp lệ: %w", err)
	}
	if !strings.EqualFold(parsed.Status, "ok") {
		return fmt.Errorf("FlareSolverr không giải được challenge (%s): %s", parsed.Status, truncateForLog(parsed.Message, 200))
	}

	cookies := make([]*http.Cookie, 0, len(parsed.Solution.Cookies))
	hasClearance := false
	for _, sc := range parsed.Solution.Cookies {
		if sc.Name == "" {
			continue
		}
		ck := &http.Cookie{Name: sc.Name, Value: sc.Value, Path: sc.Path}
		ck.Path = strings.TrimSpace(ck.Path)
		if ck.Path == "" {
			ck.Path = "/"
		}
		if sc.Domain != "" {
			ck.Domain = sc.Domain
		}
		if sc.Expiry > 0 {
			ck.Expires = time.Unix(int64(sc.Expiry), 0)
		}
		cookies = append(cookies, ck)

		if sc.Name == "cf_clearance" && sc.Value != "" {
			hasClearance = true
		}
	}
	if !hasClearance {
		return fmt.Errorf("FlareSolverr không trả về cookie cf_clearance (trang trả về: %s)", truncateForLog(parsed.Solution.Response, 160))
	}

	if jar := c.currentHTTPClient().Jar; jar != nil {
		if base, err := url.Parse(c.baseURL); err == nil {
			jar.SetCookies(base, cookies)
		}
	}

	ua := strings.TrimSpace(parsed.Solution.UserAgent)
	if ua == "" {
		ua = DefaultUA
	}

	c.solverMu.Lock()
	c.session = solverSession{
		cookies:   cookies,
		userAgent: ua,
		expires:   time.Now().Add(sessionTTL),
	}
	c.solverMu.Unlock()

	return nil
}

// truncateForLog rút gọn chuỗi dài để nhét vừa thông báo lỗi.
func truncateForLog(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
