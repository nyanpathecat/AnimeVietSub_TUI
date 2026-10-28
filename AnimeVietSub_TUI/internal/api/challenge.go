package api

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
)

// Lỗi đặc biệt giúp tầng TUI phân biệt "bị Cloudflare chặn" với lỗi mạng thông thường.
// Khi gặp các lỗi này, thông báo cho người dùng cần hướng dẫn cụ thể (bật FlareSolverr)
// thay vì chỉ hiện "HTTP 403".
var (
	// ErrChallenge: máy chủ trả về trang xác minh của Cloudflare thay vì dữ liệu thật.
	ErrChallenge = errors.New("máy chủ yêu cầu xác minh Cloudflare (Xác Minh An Toàn)")
	// ErrSolverUnavailable: không gọi được FlareSolverr để giải challenge.
	ErrSolverUnavailable = errors.New("không kết nối được FlareSolverr")
)

// challengeMarkers là các dấu hiệu đặc trưng trong HTML của trang challenge.
// Cloudflare chèn _cf_chl_opt vào mọi trang managed challenge; các chuỗi còn lại
// là tiêu đề trang (bản tiếng Việt và tiếng Anh).
var challengeMarkers = []string{
	"_cf_chl_opt",
	"cf_chl_",
	"cf-challenge",
	"cf_challenge",
	"xác minh an toàn",
	"just a moment",
	"attention required",
	"enable javascript and cookies",
}

// maxChallengeScan giới hạn số byte quét đầu trang: trang challenge rất lớn,
// còn dữ liệu thật thì luôn có marker ngay trong phần đầu tài liệu.
const maxChallengeScan = 16 * 1024

// looksLikeChallengeBody kiểm tra nội dung response có phải trang challenge không.
func looksLikeChallengeBody(body []byte) bool {
	if len(body) == 0 {
		return false
	}

	head := body
	if len(head) > maxChallengeScan {
		head = head[:maxChallengeScan]
	}

	lower := strings.ToLower(string(head))
	for _, marker := range challengeMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// isChallengeResponse xác định một response là trang challenge của Cloudflare.
//
// Ba trường hợp được coi là challenge:
//   - Nội dung chứa marker của Cloudflare (chắc chắn nhất).
//   - HTTP 403: toàn bộ host animevietsub.news nằm sau managed challenge nên 403
//     ở đây luôn là challenge, không phải "hết quyền truy cập".
//   - 200 nhưng thân rỗng với endpoint bắt buộc phải có dữ liệu (requiresBody):
//     /ajax/item và /ajax/player trả 0 byte khi bị chặn, đây chính là biểu hiện
//     khiến ứng dụng hiện "không có phim nào" thay vì báo lỗi.
func isChallengeResponse(status int, body []byte, requiresBody bool) bool {
	if looksLikeChallengeBody(body) {
		return true
	}
	if status == http.StatusForbidden {
		return true
	}
	if requiresBody && status == http.StatusOK && len(bytes.TrimSpace(body)) == 0 {
		return true
	}
	return false
}
