package player

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// FinishedMsg gửi về Bubble Tea update loop khi mpv thoát.
type FinishedMsg struct {
	Err error
}

// ValidateStreamURL kiểm tra URL có scheme hợp lệ (http/https) và có host không rỗng.
func ValidateStreamURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("URL không hợp lệ: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("chỉ chấp nhận giao thức http hoặc https, nhận được: %s", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("URL không hợp lệ (thiếu host)")
	}
	return nil
}

// IsMPVInstalled kiểm tra xem binary mpv đã cài trên hệ thống chưa.
func IsMPVInstalled() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

// GetInstallInstructions trả về hướng dẫn cài đặt mpv.
func GetInstallInstructions() string {
	switch runtime.GOOS {
	case "linux":
		return "mpv chưa được cài đặt. Bạn có thể cài đặt bằng lệnh:\n  - Ubuntu/Debian: sudo apt install mpv\n  - Fedora: sudo dnf install mpv\n  - Arch Linux: sudo pacman -S mpv"
	case "darwin":
		return "mpv chưa được cài đặt. Bạn có thể cài đặt bằng lệnh:\n  brew install mpv"
	case "windows":
		return "mpv chưa được cài đặt. Vui lòng tải mpv từ https://mpv.io hoặc cài qua winget/scoop:\n  winget install shinchiro.mpv"
	default:
		return "mpv chưa được cài đặt. Vui lòng cài đặt mpv từ https://mpv.io"
	}
}

// BuildMPVCommand xây dựng lệnh mpv với các tham số tối ưu và header cần thiết.
// Đã loại bỏ --http-header-fields để tránh lỗi ngắt header khi User-Agent chứa dấu phẩy.
// Đặt flag "--" trước URL để ngăn injection tham số.
func BuildMPVCommand(videoURL, title string) *exec.Cmd {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	referrer := "https://animevietsub.news/"

	args := []string{
		"--really-quiet",
		fmt.Sprintf("--user-agent=%s", ua),
		fmt.Sprintf("--referrer=%s", referrer),
	}

	if title != "" {
		args = append(args, fmt.Sprintf("--title=AnimeVietSub - %s", title))
	}

	args = append(args, "--", videoURL)
	return exec.Command("mpv", args...)
}

// PlayStream trả về tea.Cmd thực thi mpv qua tea.ExecProcess sau khi kiểm tra URL.
func PlayStream(videoURL, title string) tea.Cmd {
	if err := ValidateStreamURL(videoURL); err != nil {
		return func() tea.Msg {
			return FinishedMsg{Err: err}
		}
	}
	c := BuildMPVCommand(videoURL, title)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return FinishedMsg{Err: err}
	})
}

// OpenBrowser mở URL bằng trình duyệt mặc định của hệ điều hành.
func OpenBrowser(targetURL string) error {
	if err := ValidateStreamURL(targetURL); err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", targetURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	case "darwin":
		cmd = exec.Command("open", targetURL)
	default:
		return fmt.Errorf("hệ điều hành không được hỗ trợ để mở trình duyệt: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		_ = cmd.Wait()
	}()
	return nil
}

// FindChromiumBinary tìm binary của trình duyệt nhân Chromium hỗ trợ chế độ --app.
func FindChromiumBinary() string {
	candidates := []string{
		"chromium",
		"google-chrome",
		"google-chrome-stable",
		"brave-browser",
		"brave",
		"microsoft-edge",
	}

	// Windows không có các tên binary kiểu Unix ở trên; Edge/Chrome/Brave nằm trong
	// Program Files hoặc LocalAppData. Thêm các đường dẫn tuyệt đối để chế độ
	// Standalone Web Player (--app) hoạt động trên Windows.
	if runtime.GOOS == "windows" {
		candidates = append(candidates, "chrome.exe", "msedge.exe", "brave.exe")
		for _, base := range []string{
			os.Getenv("ProgramFiles"),
			os.Getenv("ProgramFiles(x86)"),
			os.Getenv("LocalAppData"),
		} {
			if base == "" {
				continue
			}
			candidates = append(candidates,
				filepath.Join(base, `Google\Chrome\Application\chrome.exe`),
				filepath.Join(base, `Microsoft\Edge\Application\msedge.exe`),
				filepath.Join(base, `BraveSoftware\Brave-Browser\Application\brave.exe`),
			)
		}
	}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

// iframePageTemplate là trang bọc player trong iframe.
// Nhiều CDN player (vd. stream.googleapiscdn.com) kiểm tra window.self === window.top
// và từ chối phát khi bị mở trực tiếp ở top-level, nên phải nhúng qua iframe.
const iframePageTemplate = `<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>AnimeVietSub Player</title>
<style>
  html, body { margin: 0; padding: 0; height: 100%; background: #000; overflow: hidden; }
  iframe { border: 0; width: 100%; height: 100%; }
</style>
</head>
<body>
<iframe src="{{URL}}" allow="autoplay; fullscreen; encrypted-media; picture-in-picture" allowfullscreen></iframe>
</body>
</html>
`

var (
	playerServerOnce sync.Once
	playerServerBase string
	playerServerErr  error
)

// ensurePlayerServer khởi động HTTP server nội bộ (127.0.0.1, cổng ngẫu nhiên) một lần duy nhất
// để phục vụ trang iframe. Dùng http:// thay vì file:// để player/CDN hoạt động ổn định hơn.
func ensurePlayerServer() (string, error) {
	playerServerOnce.Do(func() {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			playerServerErr = err
			return
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := r.URL.Query().Get("url")
			u, err := url.Parse(target)
			if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
				http.Error(w, "URL không hợp lệ", http.StatusBadRequest)
				return
			}

			page := strings.Replace(iframePageTemplate, "{{URL}}", html.EscapeString(target), 1)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(page))
		})

		server := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
		playerServerBase = "http://" + listener.Addr().String() + "/"
		go func() {
			_ = server.Serve(listener)
		}()
	})

	return playerServerBase, playerServerErr
}

// BuildIframePageURL tạo URL trang player cục bộ có nhúng stream trong iframe.
func BuildIframePageURL(targetURL string) (string, error) {
	if err := ValidateStreamURL(targetURL); err != nil {
		return "", err
	}

	base, err := ensurePlayerServer()
	if err != nil {
		return "", fmt.Errorf("không thể khởi động trình phát iframe cục bộ: %w", err)
	}

	return base + "?url=" + url.QueryEscape(targetURL), nil
}

// OpenBrowserEmbedded mở stream qua trang iframe cục bộ bằng trình duyệt mặc định.
// Dùng cho các player CDN từ chối chạy ở top-level.
func OpenBrowserEmbedded(targetURL string) error {
	launchURL, err := BuildIframePageURL(targetURL)
	if err != nil {
		launchURL = targetURL
	}
	return OpenBrowser(launchURL)
}

// PlayStandaloneWebPlayer mở URL dưới dạng cửa sổ ứng dụng Web Player độc lập (không address bar/tabs).
// Stream được bọc trong iframe cục bộ vì nhiều player CDN chặn truy cập top-level.
// Nếu không tìm thấy Chromium-based browser, tự động fallback sang OpenBrowser.
func PlayStandaloneWebPlayer(targetURL string) (string, error) {
	if err := ValidateStreamURL(targetURL); err != nil {
		return "", err
	}

	launchURL, err := BuildIframePageURL(targetURL)
	if err != nil {
		launchURL = targetURL
	}

	chromePath := FindChromiumBinary()
	if chromePath != "" {
		cmd := exec.Command(chromePath, fmt.Sprintf("--app=%s", launchURL), "--window-size=1280,720")
		if err := cmd.Start(); err == nil {
			go func() {
				_ = cmd.Wait()
			}()
			return "Cửa sổ Standalone Web Player", nil
		}
	}

	// Fallback sang trình duyệt mặc định (vẫn dùng trang iframe nếu tạo được)
	if err := OpenBrowser(launchURL); err != nil {
		return "", err
	}
	return "Trình duyệt mặc định", nil
}
