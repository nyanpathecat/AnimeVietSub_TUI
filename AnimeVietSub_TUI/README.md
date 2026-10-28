<div align="center">

# 🎬 AnimeVietSub TUI (avs)

**Trình xem Anime Vietsub mượt mà, tiện lợi trực tiếp trên Terminal.**  
*Tận hưởng trọn vẹn thế giới Anime không quảng cáo, tối ưu bàn phím và siêu nhẹ.*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D_1.22-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows%20(WSL)-lightgrey)](https://github.com/nyanpathecat/AnimeVietSub_TUI)
[![TUI Framework](https://img.shields.io/badge/Built%20with-Bubble%20Tea-f38ba8)](https://github.com/charmbracelet/bubbletea)

</div>

---

## 📑 Mục Lục

- [Giới Thiệu](#-giới-thiệu)
- [Tính Năng Nổi Bật](#-tính-năng-nổi-bật)
- [Yêu Cầu Hệ Thống](#-yêu-cầu-hệ-thống)
- [Hướng Dẫn Cài Đặt](#-hướng-dẫn-cài-đặt)
  - [Cách 1: Cài đặt tự động bằng 1 lệnh (Build sẵn / Script Installer)](#cách-1-cài-đặt-tự-động-bằng-1-lệnh-khuyến-nghị)
  - [Cách 2: Tự build từ mã nguồn (Build from Source)](#cách-2-tự-build-từ-mã-nguồn-build-from-source)
- [Hướng Dẫn Sử Dụng](#-hướng-dẫn-sử-dụng)
  - [Khởi chạy](#khởi-chạy)
  - [Bảng phím tắt điều hướng](#bảng-phím-tắt-điều-hướng)
- [Vượt Cloudflare (FlareSolverr)](#vượt-cloudflare-flaresolverr)
- [Cấu Hình & Tùy Biến](#-cấu-hình--tùy-biến)
- [Cấu Trúc Thư Mục](#-cấu-trúc-thư-mục)
- [Hướng Dẫn Gỡ Cài Đặt](#-hướng-dẫn-gỡ-cài-đặt)
  - [Cách 1: Gỡ cài đặt tự động bằng 1 lệnh](#cách-1-gỡ-cài-đặt-tự-động-bằng-1-lệnh)
  - [Cách 2: Gỡ cài đặt thủ công / Từ mã nguồn](#cách-2-gỡ-cài-đặt-thủ-công--từ-mã-nguồn)
- [Giấy Phép (License)](#-giấy-phép-license)
- [Tuyên Bố Miễn Trừ Trách Nhiệm](#-tuyên-bố-miễn-trừ-trách-nhiệm)

---

## 🌟 Giới Thiệu

**AnimeVietSub TUI (`avs`)** là ứng dụng dòng lệnh (Terminal User Interface - TUI) được xây dựng bằng ngôn ngữ **Golang**, sử dụng hệ sinh thái thư viện giao diện terminal hiện đại **Bubble Tea** và **Lipgloss** từ Charmbracelet.

Ứng dụng giúp bạn duyệt kho phim phong phú từ AnimeVietSub, chọn tập, chuyển đổi server linh hoạt và phát video trực tiếp thông qua **MPV** hoặc **Cửa sổ Web Player độc lập**, mang lại trải nghiệm xem phim:
- 🚀 **Siêu nhanh và nhẹ:** Khởi động tức thì, tiêu tốn rất ít tài nguyên máy tính.
- 🚫 **Không quảng cáo gây phiền:** Không pop-up, không banner phiền toái.
- ⌨️ **Tối ưu bàn phím:** Điều hướng mượt mà theo phong cách Vim (`h`, `j`, `k`, `l`) hoặc các phím mũi tên.

---

## ✨ Tính Năng Nổi Bật

- 🔍 **Tìm kiếm nhanh chóng:** Tìm kiếm phim với thanh Search trực quan ngay trên đỉnh giao diện (`/`).
- 📺 **Trang chủ đầy đủ chuyên mục:** Phim đề cử, Anime mới cập nhật, Phim bộ, Phim lẻ, Bảng xếp hạng.
- 🏷️ **Duyệt theo thể loại (`g`):** Khám phá anime theo hàng chục thể loại (Hành động, Isekai, Shounen, Tình cảm, v.v.).
- 🎬 **Hỗ trợ phát đa chế độ:**
  - **MPV Player (Khuyến nghị):** Phát trực tiếp stream chất lượng cao (1080p, 720p) trong cửa sổ MPV native với đầy đủ tính năng tua, âm lượng.
  - **Standalone Web Player:** Mở trình phát web qua Chromium/Brave/Edge ở chế độ `--app` không thanh địa chỉ/tab.
  - **Browser Fallback:** Mở nhanh bằng trình duyệt mặc định khi cần (`o`).
- ⚡ **Quản lý Server & Tập phim:** Tự động nạp danh sách server (VIP, Dự phòng, Iframe) và danh sách tập phim nhanh gọn.
- 🕒 **Lịch sử xem phim (`H`):** Tự động lưu vết các tập phim đã xem dở để bạn có thể tiếp tục bất cứ lúc nào.
- ⚙️ **Menu Cài đặt tương tác (`c`):** Đổi ngôn ngữ, đổi chế độ phát, chỉnh chất lượng video ưu tiên, xóa lịch sử ngay trong giao diện.
- 🌐 **Đa ngôn ngữ (i18n):** Hỗ trợ chuyển đổi mượt mà giữa **Tiếng Việt** và **English**.

---

## 📋 Yêu Cầu Hệ Thống

1. **Hệ điều hành:** Linux, macOS, hoặc Windows (Windows Terminal / PowerShell / Git Bash — không bắt buộc WSL).
2. **Trình phát video (Khuyến nghị):** Cài đặt **[mpv](https://mpv.io/)** để có trải nghiệm xem phim trực tiếp tốt nhất.
   - **Ubuntu / Debian:** `sudo apt update && sudo apt install -y mpv`
   - **Arch Linux:** `sudo pacman -S mpv`
   - **Fedora:** `sudo dnf install -y mpv`
   - **macOS:** `brew install mpv`
   - **Windows:** `winget install shinchiro.mpv`
3. **Golang (Chỉ cần nếu tự build từ mã nguồn):** Phiên bản Go `>= 1.22`.
4. **FlareSolverr (khuyến nghị):** Cần khi website bật lớp xác minh Cloudflare — xem [Vượt Cloudflare](#vượt-cloudflare-flaresolverr).

---

## 🚀 Hướng Dẫn Cài Đặt

Bạn có thể lựa chọn 1 trong 2 cách cài đặt dưới đây:

### Cách 1: Cài đặt tự động bằng 1 lệnh (Khuyến nghị)

Cách này phù hợp nhất nếu bạn muốn cài đặt nhanh chóng vào hệ thống mà không cần cài đặt Go. Script sẽ tự động nhận diện hệ điều hành (Linux/macOS), kiến trúc CPU (x86_64, arm64), kiểm tra `mpv` và cài đặt binary vào `/usr/local/bin` (hoặc `~/.local/bin`).

#### Chạy lệnh cài đặt One-Liner qua `curl`:
```bash
curl -fsSL https://raw.githubusercontent.com/nyanpathecat/AnimeVietSub_TUI/main/install.sh | bash
```

#### Hoặc sử dụng `wget`:
```bash
wget -qO- https://raw.githubusercontent.com/nyanpathecat/AnimeVietSub_TUI/main/install.sh | bash
```

#### Windows (Windows Terminal / PowerShell)

`install.sh` là script bash, trên Windows hãy dùng `install.ps1` — script cài đặt chính thức cho PowerShell:

```powershell
irm https://raw.githubusercontent.com/nyanpathecat/AnimeVietSub_TUI/main/install.ps1 | iex
```

Script sẽ:

- Tải `avs.exe` từ GitHub Releases và **kiểm tra file có phải executable Windows hợp lệ** trước khi cài (không bao giờ lưu nhầm trang lỗi 404 vào `avs.exe` — nguyên nhân lỗi *"The specified executable is not a valid application for this OS platform"*).
- Nếu Releases chưa có binary: tự build từ mã nguồn (cần [Go 1.27+](https://go.dev/dl/), và Git nếu chưa clone repo).
- Cài vào `%LOCALAPPDATA%\AnimeVietSub\bin` và thêm vào PATH của user.

Nếu đã clone repo về máy, chạy trực tiếp:

```powershell
git clone https://github.com/nyanpathecat/AnimeVietSub_TUI.git
cd AnimeVietSub_TUI
.\install.ps1
```

Sau khi cài, **mở terminal mới** rồi chạy `avs`.

##### Tự build thủ công (không cần script)

```powershell
git clone https://github.com/nyanpathecat/AnimeVietSub_TUI.git
cd AnimeVietSub_TUI

go build -ldflags="-s -w" -o avs.exe cmd/avs/main.go

# Chạy thử
.\avs.exe
```

> **Lưu ý khi tải binary trực tiếp bằng `curl`:** luôn kèm `-f` (vd. `curl.exe -fL -o avs.exe <url>`) để curl dừng khi server trả 404 thay vì ghi nội dung lỗi vào `avs.exe`. Binary Windows được đính kèm tự động vào trang [Releases](https://github.com/nyanpathecat/AnimeVietSub_TUI/releases) qua GitHub Actions khi push tag.

Hoặc dùng one-liner ở trên nếu bạn đang chạy trong Git Bash / WSL.

#### Hoặc nếu bạn đã tải/clone repo về máy:
```bash
chmod +x install.sh
./install.sh
```

> **Ghi chú:** Script sẽ ưu tiên cài vào `/usr/local/bin` (sử dụng sudo nếu cần) để bạn có thể gõ lệnh `avs` ở bất kỳ thư mục nào. Nếu không có quyền root, script sẽ tự động cài vào `~/.local/bin`.

---

### Cách 2: Tự build từ mã nguồn (Build from Source)

Cách này dành cho các bạn lập trình viên muốn tự tay tùy biến mã nguồn hoặc kiểm soát quy trình biên dịch.

#### Bước 1: Clone kho mã nguồn
```bash
git clone https://github.com/nyanpathecat/AnimeVietSub_TUI.git
cd AnimeVietSub_TUI
```

#### Bước 2: Biên dịch và cài đặt

##### Lựa chọn 2.1: Sử dụng `Makefile` (Tiện lợi nhất)
```bash
# Biên dịch binary ra thư mục hiện tại
make build

# Cài đặt trực tiếp vào hệ thống (/usr/local/bin)
sudo make install

# Trên Windows Terminal / PowerShell: build ra avs.exe
make build-windows

# Biên dịch cho cả Windows và Linux
make build-all
```

> Mã nguồn không dùng cgo, nên bạn có thể cross-compile từ Linux/macOS ra `avs.exe` cho Windows và ngược lại mà không cần MinGW.

##### Lựa chọn 2.2: Sử dụng lệnh `go build`
```bash
# Tải dependencies và biên dịch tối ưu (strip debug info)
go build -ldflags="-s -w" -o avs cmd/avs/main.go

# Di chuyển binary vào PATH hệ thống
sudo mv avs /usr/local/bin/
```

##### Lựa chọn 2.3: Sử dụng `go install`
```bash
# Biên dịch và cài trực tiếp vào $GOPATH/bin
go install ./cmd/avs
```
*(Hãy chắc chắn `$GOPATH/bin` hoặc `$HOME/go/bin` đã nằm trong biến môi trường `$PATH` của bạn).*

---

## 🎮 Hướng Dẫn Sử Dụng

### Khởi chạy

Sau khi cài đặt xong, bạn chỉ cần mở terminal và chạy:

```bash
avs
```

---

### Bảng phím tắt điều hướng

Ứng dụng hỗ trợ cả phím điều hướng tiêu chuẩn và phong cách Vim:

| Phím | Chức năng |
| :---: | :--- |
| `↑` / `k` | Di chuyển lên |
| `↓` / `j` | Di chuyển xuống |
| `←` / `h` | Di chuyển sang trái / Quay về tab trước |
| `→` / `l` | Di chuyển sang phải / Sang tab kế tiếp |
| `Enter` | Xác nhận / Xem phim / Chọn tập / Chọn server |
| `Esc` | Quay lại màn hình trước đó |
| `/` | Mở thanh tìm kiếm phim |
| `g` | Mở danh mục Thể loại Anime |
| `H` *(Shift + h)* | Xem Lịch sử xem phim |
| `c` | Mở Menu Cấu hình & Cài đặt |
| `o` | Mở phim đang chọn bằng trình duyệt web |
| `r` | Tải lại dữ liệu (Reload) |
| `?` | Bật / tắt thanh gợi ý phím tắt trợ giúp |
| `q` / `Ctrl + C` | Thoát ứng dụng |

---

## Vượt Cloudflare (FlareSolverr)

Website thường bật lớp xác minh **Cloudflare Managed Challenge** ("Xác Minh An Toàn"). Khi đó mọi request HTTP thông thường đều nhận `HTTP 403`, hoặc tệ hơn là `HTTP 200` với nội dung rỗng khiến ứng dụng hiện "không có phim nào". Đây không phải lỗi mạng hay lỗi retry: chỉ trình duyệt thực thi JavaScript mới vượt được lớp này.

`avs` tự xử lý qua **FlareSolverr** — một proxy chạy Chromium ẩn, chạy được giống nhau trên Windows Terminal và Linux:

1. Gửi request trực tiếp trước (nhanh, không cần solver).
2. Nếu nhận trang xác minh → gọi FlareSolverr để lấy cookie `cf_clearance` kèm đúng User-Agent đã dùng.
3. Gửi lại request với phiên vừa có → nhận dữ liệu thật.
4. Phiên được cache trong ~25 phút. Khi đổi router/mạng khiến cookie mất hiệu lực, ứng dụng tự xóa phiên và giải lại — nhờ vậy chạy được trên mọi mạng, không chỉ mạng đã dùng để đăng nhập.

### Bật FlareSolverr

**Windows (Windows Terminal / PowerShell) và Linux** dùng chung một lệnh:

```bash
docker run -d --name flaresolverr -p 8191:8191 flaresolverr/flaresolverr:latest
```

- **Windows:** cài Docker Desktop trước, sau đó chạy lệnh trên trong PowerShell.
- **Ubuntu / Debian:** `sudo apt update && sudo apt install -y docker.io`
- **Kiểm tra solver đã sẵn sàng:** `curl http://127.0.0.1:8191/health`

Nếu Docker khởi động kèm máy, thêm `--restart unless-stopped` vào lệnh `docker run` để solver luôn sẵn sàng.

### Tắt hoặc đổi địa chỉ solver

- Biến môi trường (ưu tiên cao nhất):
  ```bash
  # Windows PowerShell: $env:AVS_SOLVER_URL = "http://127.0.0.1:8191"
  export AVS_SOLVER_URL="http://127.0.0.1:8191"
  # Tắt hoàn toàn: export AVS_SOLVER_URL=off
  ```
- Hoặc khóa `solver_url` trong `.avs_config.json` (xem mục dưới).

Khi solver tắt hoặc không chạy, ứng dụng báo rõ *"máy chủ yêu cầu xác minh Cloudflare"* kèm lệnh cần chạy, thay vì chỉ hiện `HTTP 403`.

---

## ⚙️ Cấu Hình & Tùy Biến

Ứng dụng lưu cấu hình tại file `.avs_config.json` và lịch sử xem tại `.avs_history.json`.

Bạn có thể thay đổi thiết lập trực tiếp trong màn hình **Cấu hình** (nhấn phím `c`) hoặc chỉnh sửa file `.avs_config.json`:

```json
{
  "language": "vi",
  "auto_open_browser_iframe": false,
  "player_mode": "auto",
  "preferred_quality": "1080p",
  "theme_mode": "terminal",
  "solver_url": "http://127.0.0.1:8191"
}
```

- **`language`**: Ngôn ngữ hiển thị (`vi` cho Tiếng Việt, `en` cho English).
- **`player_mode`**: 
  - `auto`: Tự động nhận diện (dùng MPV cho link stream direct, Standalone Window cho iframe).
  - `mpv`: Luôn ưu tiên dùng MPV Player.
  - `app_window`: Mở cửa sổ Web Player độc lập (Chromium-based app window).
  - `browser`: Luôn mở qua trình duyệt mặc định.
- **`preferred_quality`**: Chất lượng phát mặc định (`1080p`, `720p`, hoặc `auto`).
- **`theme_mode`**: Giao diện màu (`terminal`, `dark`, hoặc `light`).
- **`solver_url`**: Địa chỉ FlareSolverr dùng để vượt Cloudflare (mặc định `http://127.0.0.1:8191`). Đặt `"off"` để tắt solver.

---

## 📂 Cấu Trúc Thư Mục

```text
.
├── cmd/
│   └── avs/
│       └── main.go             # Điểm khởi động ứng dụng (Entrypoint)
├── internal/
│   ├── api/                    # Xử lý HTTP Client, parser dữ liệu HTML/JSON & mô hình
│   ├── config/                 # Quản lý cấu hình ứng dụng (.avs_config.json)
│   ├── history/                # Quản lý lịch sử xem anime (.avs_history.json)
│   ├── i18n/                   # Hệ thống dịch song ngữ (Tiếng Việt & English)
│   ├── player/                 # Tích hợp MPV, Chromium App Player và Browser
│   └── tui/                    # Toàn bộ giao diện Bubble Tea (Home, Detail, Search, ...)
├── docs/
│   └── API.md                  # Tài liệu chi tiết kỹ thuật hệ thống API
├── install.sh                  # Trình cài đặt tự động 1 lệnh vào hệ thống
├── uninstall.sh                # Trình gỡ cài đặt tự động khỏi hệ thống
├── Makefile                    # Tập lệnh build và cài đặt từ source
├── LICENSE                     # Giấy phép mã nguồn mở MIT
└── README.md                   # Tài liệu hướng dẫn sử dụng
```

---

## 🗑️ Hướng Dẫn Gỡ Cài Đặt

Tương tự như cài đặt, bạn có 2 cách để gỡ bỏ hoàn toàn `avs` khỏi hệ thống:

### Cách 1: Gỡ cài đặt tự động bằng 1 lệnh

Cách nhanh nhất và sạch sẽ nhất. Script sẽ tự động quét mọi vị trí cài đặt của binary `avs` (`/usr/local/bin`, `~/.local/bin`, `/usr/bin`) và hỗ trợ xóa sạch dữ liệu.

#### Chạy lệnh gỡ cài đặt trực tiếp qua `curl`:
```bash
curl -fsSL https://raw.githubusercontent.com/nyanpathecat/AnimeVietSub_TUI/main/uninstall.sh | bash
```

#### Hoặc chạy file script có sẵn trong repo:
```bash
# Gỡ cài đặt thông thường (sẽ hỏi nếu muốn xóa dữ liệu cấu hình & lịch sử)
./uninstall.sh

# Gỡ cài đặt hoàn toàn (xóa cả binary, cấu hình và lịch sử xem phim mà không cần hỏi)
./uninstall.sh --purge
```

---

### Cách 2: Gỡ cài đặt thủ công / Từ mã nguồn

Nếu bạn đã cài đặt từ mã nguồn hoặc muốn tự tay xóa tệp:

##### Lựa chọn 2.1: Sử dụng `Makefile`
```bash
# Chỉ gỡ bỏ binary khỏi hệ thống
sudo make uninstall

# Gỡ bỏ binary VÀ dọn sạch toàn bộ tệp cấu hình, lịch sử xem
sudo make purge
```

##### Lựa chọn 2.2: Xóa thủ công bằng lệnh shell
```bash
# Xóa binary khỏi PATH
sudo rm -f /usr/local/bin/avs
rm -f ~/.local/bin/avs

# (Tùy chọn) Xóa tệp cấu hình và lịch sử xem
rm -f ~/.avs_config.json ~/.avs_history.json .avs_config.json .avs_history.json
```

---

## 📄 Giấy Phép (License)

Dự án này được phân phối dưới giấy phép **MIT License**. Bạn hoàn toàn có quyền tự do sử dụng, chỉnh sửa và phân phối lại.  
Xem thông tin chi tiết tại tệp [LICENSE](LICENSE).

---

## ⚠️ Tuyên Bố Miễn Trừ Trách Nhiệm

- Ứng dụng **AnimeVietSub TUI** được phát triển nhằm mục đích học tập, giải trí cá nhân và nghiên cứu công nghệ phát triển giao diện TUI với Golang.
- Ứng dụng không lưu trữ hay sở hữu bất kỳ tệp video hoặc dữ liệu nội dung phim nào trên máy chủ riêng. Toàn bộ hình ảnh, video và dữ liệu liên quan thuộc quyền sở hữu của dịch vụ gốc và các nhà phát hành bản quyền tương ứng.

---

<div align="center">
  <b>Chúc bạn có những giờ phút xem Anime thật thư giãn! 🍿✨</b><br/>
  <i>Đừng quên tặng một ⭐️ nếu dự án hữu ích với bạn!</i>
</div>
