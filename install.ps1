<#
.SYNOPSIS
    Cài đặt AnimeVietSub TUI (avs) trên Windows.

.DESCRIPTION
    - Ưu tiên tải binary từ GitHub Releases và KIỂM TRA file có phải executable
      Windows hợp lệ (header MZ, kích thước hợp lý) trước khi cài — tránh lỗi
      "The specified executable is not a valid application for this OS platform"
      do tải nhầm nội dung lỗi 404 hoặc file hỏng.
    - Nếu Releases chưa có binary (hoặc tải thất bại): tự build từ mã nguồn
      (yêu cầu Go 1.27+; Git nếu cần clone).
    - Cài vào %LOCALAPPDATA%\AnimeVietSub\bin và thêm mục nhập vào PATH của user.

.PARAMETER InstallDir
    Thư mục cài đặt. Mặc định: %LOCALAPPDATA%\AnimeVietSub\bin

.PARAMETER BuildFromSource
    Bỏ qua bước tải binary, luôn build từ mã nguồn.

.PARAMETER NoPathUpdate
    Không chỉnh sửa biến môi trường PATH của user.

.EXAMPLE
    .\install.ps1

.EXAMPLE
    .\install.ps1 -BuildFromSource

.EXAMPLE
    irm https://raw.githubusercontent.com/nyanpathecat/AnimeVietSub_TUI/main/install.ps1 | iex
#>
[CmdletBinding()]
param(
    [string]$InstallDir = (Join-Path $env:LOCALAPPDATA 'AnimeVietSub\bin'),
    [switch]$BuildFromSource,
    [switch]$NoPathUpdate
)

$ErrorActionPreference = 'Stop'

$RepoUrl = 'https://github.com/nyanpathecat/AnimeVietSub_TUI'
$ReleaseBase = "$RepoUrl/releases/latest/download"

function Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Ok($msg)   { Write-Host "==> $msg" -ForegroundColor Green }
function Warn($msg) { Write-Warning $msg }

# Kiểm tra file có phải executable Windows hợp lệ không (header "MZ", >= 1MB).
function Test-PeFile([string]$Path) {
    if ([string]::IsNullOrWhiteSpace($Path)) { return $false }
    if (-not (Test-Path -LiteralPath $Path)) { return $false }
    $item = Get-Item -LiteralPath $Path
    if ($item.Length -lt 1MB) { return $false }
    $fs = [System.IO.File]::OpenRead($item.FullName)
    try {
        $b1 = $fs.ReadByte()
        $b2 = $fs.ReadByte()
        return ($b1 -eq 0x4D -and $b2 -eq 0x5A)
    } finally {
        $fs.Dispose()
    }
}

function Get-TargetArch {
    if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { return 'arm64' }
    return 'amd64'
}

# Thử tải lần lượt các asset trên Releases; chỉ chấp nhận file PE hợp lệ.
function Get-ReleaseAsset {
    param(
        [string[]]$Names,
        [string]$OutFile
    )
    foreach ($name in $Names) {
        $url = "$ReleaseBase/$name"
        Step "Tải $name từ GitHub Releases..."
        try {
            Invoke-WebRequest -Uri $url -OutFile $OutFile -UseBasicParsing -TimeoutSec 120 | Out-Null
        } catch {
            Warn "Tải thất bại: $($_.Exception.Message)"
            continue
        }
        if (Test-PeFile $OutFile) { return $true }
        Warn "File '$name' không phải file thực thi Windows hợp lệ (release có thể chưa đính kèm asset)."
        Remove-Item -LiteralPath $OutFile -Force -ErrorAction SilentlyContinue
    }
    return $false
}

function Invoke-Build {
    param([string]$OutFile)

    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw "Không tìm thấy 'go'. Hãy cài Go 1.27+ (https://go.dev/dl/) rồi chạy lại, hoặc tải avs.exe thủ công từ $RepoUrl/releases"
    }

    $arch = Get-TargetArch
    $srcDir = $null

    $rootDir = Get-Variable -Name 'PSScriptRoot' -ValueOnly -ErrorAction SilentlyContinue
    if ($rootDir -and (Test-Path -LiteralPath (Join-Path $rootDir 'cmd\avs\main.go'))) {
        $srcDir = $rootDir
    } elseif (Test-Path -LiteralPath (Join-Path (Get-Location).Path 'cmd\avs\main.go')) {
        $srcDir = (Get-Location).Path
    } else {
        if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
            throw "Không có mã nguồn tại thư mục hiện tại và cũng không tìm thấy 'git'. Hãy clone repo về rồi chạy lại script (hoặc cài Git)."
        }
        Step 'Clone mã nguồn mới nhất từ GitHub...'
        $srcDir = Join-Path ([System.IO.Path]::GetTempPath()) ('avs-src-' + [guid]::NewGuid().ToString('N'))
        & git clone --depth 1 "$RepoUrl.git" $srcDir
        if ($LASTEXITCODE -ne 0) { throw "git clone thất bại (mã lỗi $LASTEXITCODE)" }
    }

    Step "Build avs.exe cho kiến trúc $arch từ mã nguồn..."
    $old = @{
        CGO_ENABLED = $env:CGO_ENABLED
        GOOS        = $env:GOOS
        GOARCH      = $env:GOARCH
    }
    $env:CGO_ENABLED = '0'
    $env:GOOS = 'windows'
    $env:GOARCH = $arch
    Push-Location -LiteralPath $srcDir
    try {
        & go build -ldflags='-s -w' -o $OutFile ./cmd/avs
        $buildCode = $LASTEXITCODE
    } finally {
        Pop-Location
        $env:CGO_ENABLED = $old.CGO_ENABLED
        $env:GOOS = $old.GOOS
        $env:GOARCH = $old.GOARCH
    }
    if ($buildCode -ne 0) { throw "go build thất bại (mã lỗi $buildCode)" }
    if (-not (Test-PeFile $OutFile)) { throw "Build xong nhưng file '$OutFile' không hợp lệ." }
}

Write-Host ''
Write-Host 'Cài đặt AnimeVietSub TUI (avs) trên Windows' -ForegroundColor Magenta
Write-Host ''

$arch = Get-TargetArch
$workDir = Join-Path ([System.IO.Path]::GetTempPath()) ('avs-install-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $workDir -Force | Out-Null
$exePath = Join-Path $workDir 'avs.exe'

try {
    $haveBinary = $false

    if (-not $BuildFromSource) {
        if ($arch -eq 'arm64') {
            $candidates = @('avs_windows_arm64.exe', 'avs.exe')
        } else {
            $candidates = @('avs.exe')
        }
        $haveBinary = Get-ReleaseAsset -Names $candidates -OutFile $exePath
        if (-not $haveBinary) {
            Warn 'Chưa có binary dùng được trên Releases — chuyển sang build từ mã nguồn.'
        }
    }

    if (-not $haveBinary) {
        Invoke-Build -OutFile $exePath
    }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    $dest = Join-Path $InstallDir 'avs.exe'
    Copy-Item -LiteralPath $exePath -Destination $dest -Force
    Ok "Đã cài: $dest"

    if (-not $NoPathUpdate) {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        $parts = @()
        if ($userPath) { $parts = @($userPath -split ';' | Where-Object { $_ }) }
        $haveEntry = $parts | Where-Object { $_.TrimEnd('\') -ieq $InstallDir.TrimEnd('\') }
        if (-not $haveEntry) {
            $parts += $InstallDir
            [Environment]::SetEnvironmentVariable('Path', ($parts -join ';'), 'User')
            Ok 'Đã thêm thư mục cài đặt vào PATH của user (mở terminal MỚI để nhận biến môi trường).'
        }
    }

    if (-not (Get-Command mpv -ErrorAction SilentlyContinue)) {
        Warn "Chưa tìm thấy mpv (dùng để phát video). Cài bằng: winget install mpv.mpv"
    }

    Write-Host ''
    Ok "Xong! Chạy ngay: $dest"
    Write-Host "   Hoặc mở terminal mới và gõ: avs" -ForegroundColor Cyan
    Write-Host ''
} finally {
    Remove-Item -LiteralPath $workDir -Recurse -Force -ErrorAction SilentlyContinue
}
