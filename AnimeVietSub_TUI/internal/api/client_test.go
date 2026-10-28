package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const testSolverUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"

// challengePage là nội dung tối giản mô phỏng trang xác minh của Cloudflare.
const challengePage = `<html><head><title>Xác Minh An Toàn</title></head><body><script>window._cf_chl_opt={}</script></body></html>`

// newSolverStub tạo FlareSolverr giả trả về cf_clearance + User-Agent.
func newSolverStub(t *testing.T, calls *int32) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(calls, 1)

		if r.URL.Path != "/v1" {
			http.NotFound(w, r)
			return
		}

		var req solverRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("body gửi FlareSolverr không phải JSON: %v", err)
		}
		if req.Cmd == "" || req.URL == "" {
			t.Errorf("thiếu cmd/url khi gọi FlareSolverr: %+v", req)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"status":"ok","message":"","solution":{"url":%q,"status":200,"userAgent":%q,"cookies":[{"name":"cf_clearance","value":"good-cookie","domain":"","path":"/"}]}}`, req.URL, testSolverUA)
	}))

	return srv
}

// TestFetchEscalatesToSolverAndReusesSession kiểm tra luồng đầy đủ:
// request trực tiếp bị chặn → giải challenge → gửi lại bằng phiên mới → có dữ liệu,
// và request sau đó tái sử dụng phiên (không giải lại).
func TestFetchEscalatesToSolverAndReusesSession(t *testing.T) {
	var solverCalls int32

	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie("cf_clearance")
		if err != nil || ck.Value != "good-cookie" || r.Header.Get("User-Agent") != testSolverUA {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(challengePage))
			return
		}
		_, _ = w.Write([]byte("DATA"))
	}))
	defer site.Close()

	solver := newSolverStub(t, &solverCalls)
	defer solver.Close()

	client, err := NewClient(site.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetSolverURL(solver.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := client.doRequest(ctx, &ClientRequest{
		Method:       http.MethodGet,
		URL:          site.URL + "/ajax/item",
		RequiresBody: true,
	})
	if err != nil {
		t.Fatalf("request đầu tiên phải thành công sau khi giải challenge, nhận lỗi: %v", err)
	}
	if string(data) != "DATA" {
		t.Fatalf("dữ liệu trả về sai: %q", string(data))
	}
	if got := atomic.LoadInt32(&solverCalls); got != 1 {
		t.Fatalf("FlareSolverr phải được gọi đúng 1 lần, thực tế %d", got)
	}

	// Request thứ hai: phiên còn hiệu lực nên không được gọi solver nữa.
	if _, err := client.doRequest(ctx, &ClientRequest{
		Method:       http.MethodGet,
		URL:          site.URL + "/ajax/item",
		RequiresBody: true,
	}); err != nil {
		t.Fatalf("request thứ hai phải thành công: %v", err)
	}
	if got := atomic.LoadInt32(&solverCalls); got != 1 {
		t.Fatalf("phiên phải được tái sử dụng, nhưng FlareSolverr bị gọi %d lần", got)
	}
}

// TestEmptyAjaxBodyIsChallenge xử lý đúng trường hợp thực tế đang gặp:
// Cloudflare trả HTTP 200 với body rỗng cho /ajax/item, khiến app hiện "không có phim".
func TestEmptyAjaxBodyIsChallenge(t *testing.T) {
	var solverCalls int32

	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ck, err := r.Cookie("cf_clearance"); err == nil && ck.Value == "good-cookie" {
			_, _ = w.Write([]byte("<div>phim</div>"))
			return
		}
		// Giống Cloudflare: 200 nhưng không có dữ liệu.
		w.WriteHeader(http.StatusOK)
	}))
	defer site.Close()

	solver := newSolverStub(t, &solverCalls)
	defer solver.Close()

	client, err := NewClient(site.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetSolverURL(solver.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := client.doRequest(ctx, &ClientRequest{
		Method:       http.MethodPost,
		URL:          site.URL + "/ajax/item",
		Body:         "widget=list-film&type=anime-new",
		IsAjax:       true,
		RequiresBody: true,
	})
	if err != nil {
		t.Fatalf("body rỗng phải được coi là challenge và tự giải, nhận lỗi: %v", err)
	}
	if !strings.Contains(string(data), "phim") {
		t.Fatalf("dữ liệu trả về sai: %q", string(data))
	}
	if got := atomic.LoadInt32(&solverCalls); got != 1 {
		t.Fatalf("FlareSolverr phải được gọi 1 lần, thực tế %d", got)
	}
}

// TestSolverDisabledReportsChallenge: khi tắt solver, lỗi phải nêu rõ nguyên nhân
// và hướng khắc phục thay vì chỉ "HTTP 403".
func TestSolverDisabledReportsChallenge(t *testing.T) {
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(challengePage))
	}))
	defer site.Close()

	client, err := NewClient(site.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetSolverURL("off")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.doRequest(ctx, &ClientRequest{Method: http.MethodGet, URL: site.URL + "/ajax/item", RequiresBody: true})
	if !errors.Is(err, ErrChallenge) {
		t.Fatalf("mong đợi ErrChallenge, nhận: %v", err)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "flaresolverr") {
		t.Fatalf("thông báo lỗi phải hướng dẫn bật FlareSolverr, nhận: %v", err)
	}
}

// TestSolverUnreachableReturnsSolverError: solver bật nhưng không chạy.
func TestSolverUnreachableReturnsSolverError(t *testing.T) {
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(challengePage))
	}))
	defer site.Close()

	client, err := NewClient(site.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	// Cổng 1 gần như chắc chắn không có gì lắng nghe.
	client.SetSolverURL("http://127.0.0.1:1")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.doRequest(ctx, &ClientRequest{Method: http.MethodGet, URL: site.URL + "/ajax/item", RequiresBody: true})
	if !errors.Is(err, ErrSolverUnavailable) {
		t.Fatalf("mong đợi ErrSolverUnavailable, nhận: %v", err)
	}
}

func TestLooksLikeChallengeBody(t *testing.T) {
	cases := []struct {
		name string
		body string
		want bool
	}{
		{"trang xác minh tiếng Việt", challengePage, true},
		{"managed challenge", `<script>window._cf_chl_opt={cvId:"3"}</script>`, true},
		{"just a moment", `<title>Just a moment...</title>`, true},
		{"html bình thường", `<div class="item"><a href="/phim/a">A</a></div>`, false},
		{"json bình thường", `{"success":1,"link":"https://cdn/x.m3u8"}`, false},
		{"rỗng", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := looksLikeChallengeBody([]byte(tc.body)); got != tc.want {
				t.Fatalf("looksLikeChallengeBody = %v, muốn %v", got, tc.want)
			}
		})
	}
}

func TestIsChallengeResponse(t *testing.T) {
	cases := []struct {
		name         string
		status       int
		body         string
		requiresBody bool
		want         bool
	}{
		{"403 luôn là challenge", http.StatusForbidden, "Forbidden", false, true},
		{"200 body rỗng ở endpoint bắt buộc", http.StatusOK, "", true, true},
		{"200 body rỗng ở endpoint tuỳ chọn", http.StatusOK, "", false, false},
		{"200 có dữ liệu", http.StatusOK, "<div>ok</div>", true, false},
		{"200 chứa marker", http.StatusOK, challengePage, false, true},
		{"500 không phải challenge", http.StatusInternalServerError, "lỗi máy chủ", true, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isChallengeResponse(tc.status, []byte(tc.body), tc.requiresBody); got != tc.want {
				t.Fatalf("isChallengeResponse = %v, muốn %v", got, tc.want)
			}
		})
	}
}
