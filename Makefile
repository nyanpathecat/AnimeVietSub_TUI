# Makefile for AnimeVietSub TUI (avs)

BINARY_NAME := avs
BINARY_WIN  := avs.exe
MAIN_SRC    := cmd/avs/main.go
PREFIX      ?= /usr/local
INSTALL_BIN := $(DESTDIR)$(PREFIX)/bin

GO          := go
LDFLAGS     := -ldflags="-s -w"

.PHONY: all build build-windows build-linux build-all run install uninstall purge clean tidy help

all: build

help:
	@echo "AnimeVietSub TUI (avs) - Build Commands:"
	@echo "  make build         - Biên dịch binary cho máy hiện tại (./avs hoặc avs.exe)"
	@echo "  make build-windows - Cross-compile ra avs.exe cho Windows Terminal / PowerShell"
	@echo "  make build-linux   - Cross-compile ra avs cho Linux (amd64)"
	@echo "  make build-all     - Biên dịch cả hai nền tảng"
	@echo "  make run         - Chạy trực tiếp từ mã nguồn"
	@echo "  make install     - Cài đặt binary vào $(INSTALL_BIN) (yêu cầu quyền sudo nếu cần)"
	@echo "  make uninstall   - Gỡ cài đặt binary khỏi $(INSTALL_BIN)"
	@echo "  make purge       - Gỡ cài đặt và xóa toàn bộ dữ liệu cấu hình, lịch sử"
	@echo "  make clean       - Xóa file binary đã biên dịch"
	@echo "  make tidy        - Dọn dẹp và cập nhật Go modules"

build:
	@echo "==> Đang biên dịch $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_SRC)
	@echo "==> Biên dịch thành công: ./$(BINARY_NAME)"

# Cross-compile không cần cgo nên chạy được trên mọi máy có Go,
# kể cả build từ Linux ra avs.exe cho Windows Terminal.
build-windows:
	@echo "==> Đang biên dịch $(BINARY_WIN) cho Windows (amd64)..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY_WIN) $(MAIN_SRC)
	@echo "==> Biên dịch thành công: ./$(BINARY_WIN)"

build-linux:
	@echo "==> Đang biên dịch $(BINARY_NAME) cho Linux (amd64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_SRC)
	@echo "==> Biên dịch thành công: ./$(BINARY_NAME)"

build-all: build-windows build-linux

run:
	$(GO) run $(MAIN_SRC)

install: build
	@echo "==> Đang cài đặt $(BINARY_NAME) vào $(INSTALL_BIN)..."
	@mkdir -p $(INSTALL_BIN)
	@install -m 755 $(BINARY_NAME) $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "==> Cài đặt thành công! Bạn có thể khởi chạy bằng lệnh: $(BINARY_NAME)"

uninstall:
	@echo "==> Đang gỡ bỏ $(BINARY_NAME) khỏi $(INSTALL_BIN)..."
	@rm -f $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "==> Đã gỡ bỏ cài đặt thành công."

purge: uninstall
	@echo "==> Đang xóa tệp cấu hình và lịch sử xem..."
	@rm -f .avs_config.json .avs_history.json $(HOME)/.avs_config.json $(HOME)/.avs_history.json
	@echo "==> Đã dọn dẹp sạch sẽ toàn bộ dữ liệu."

clean:
	@echo "==> Đang dọn dẹp..."
	@rm -f $(BINARY_NAME) $(BINARY_WIN)
	@echo "==> Hoàn tất."

tidy:
	$(GO) mod tidy
