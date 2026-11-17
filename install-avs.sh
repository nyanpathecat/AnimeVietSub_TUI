#!/usr/bin/env bash
# AnimeVietSub TUI (avs) - quick installer
set -e

REPO="nyanpathecat/AnimeVietSub_TUI"
INSTALL_DIR="${HOME}/.local/bin"
TMP_DIR="$(mktemp -d)"
BINARY_NAME="avs"

info() { echo -e "\033[1;34m[INFO]\033[0m $1"; }
success() { echo -e "\033[1;32m[OK]\033[0m $1"; }
error() { echo -e "\033[1;31m[ERR]\033[0m $1"; exit 1; }

mkdir -p "$INSTALL_DIR"

# Detect OS/Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) error "Kiến trúc không được hỗ trợ: $ARCH" ;;
esac

RELEASE_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}_${OS}_${ARCH}"
FALLBACK_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}_${OS}"
TARGET_BIN="$TMP_DIR/$BINARY_NAME"

# Prefer building the latest source when Go is available: the published
# release binary may be older and still point to a retired domain.
if command -v go >/dev/null 2>&1 && command -v git >/dev/null 2>&1; then
  if git clone --depth 1 "https://github.com/${REPO}.git" "$TMP_DIR/src" >/dev/null 2>&1 \
    && (cd "$TMP_DIR/src" && go build -ldflags="-s -w" -o "$TARGET_BIN" ./cmd/avs); then
    info "Đã biên dịch từ mã nguồn mới nhất"
    cp -f "$TARGET_BIN" "$INSTALL_DIR/$BINARY_NAME"
    chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    success "Đã cài đặt tại $INSTALL_DIR/$BINARY_NAME"
    rm -rf "$TMP_DIR"
    echo ""
    echo "Chạy bằng lệnh: avs"
    echo "Nếu 'avs' không tìm thấy, thêm vào shell: export PATH=\"\$HOME/.local/bin:\$PATH\""
    exit 0
  fi
  info "Không build được từ mã nguồn, chuyển sang tải bản phát hành..."
fi

# Fall back to downloading a pre-built release
if command -v curl >/dev/null 2>&1; then
  if curl -fsSL -L "$RELEASE_URL" -o "$TARGET_BIN" 2>/dev/null || curl -fsSL -L "$FALLBACK_URL" -o "$TARGET_BIN" 2>/dev/null; then
    info "Đã tải binary từ GitHub Release"
    cp -f "$TARGET_BIN" "$INSTALL_DIR/$BINARY_NAME"
    chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    success "Đã cài đặt tại $INSTALL_DIR/$BINARY_NAME"
    rm -rf "$TMP_DIR"
    echo ""
    echo "Chạy bằng lệnh: avs"
    echo "Nếu 'avs' không tìm thấy, thêm vào shell: export PATH=\"\$HOME/.local/bin:\$PATH\""
    exit 0
  fi
elif command -v wget >/dev/null 2>&1; then
  if wget -q -O "$TARGET_BIN" "$RELEASE_URL" 2>/dev/null || wget -q -O "$TARGET_BIN" "$FALLBACK_URL" 2>/dev/null; then
    info "Đã tải binary từ GitHub Release"
    cp -f "$TARGET_BIN" "$INSTALL_DIR/$BINARY_NAME"
    chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    success "Đã cài đặt tại $INSTALL_DIR/$BINARY_NAME"
    rm -rf "$TMP_DIR"
    echo ""
    echo "Chạy bằng lệnh: avs"
    echo "Nếu 'avs' không tìm thấy, thêm vào shell: export PATH=\"\$HOME/.local/bin:\$PATH\""
    exit 0
  fi
fi

error "Chưa có release cho $OS/$ARCH. Hãy dùng script install.sh trong repo hoặc build từ source."
