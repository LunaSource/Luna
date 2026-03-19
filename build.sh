set -euo pipefail

SRC="main.go"
OUT_DIR="./bin"

mkdir -p "$OUT_DIR"

echo "Building linux/amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$OUT_DIR/Luna_linux_amd64" "$SRC"

echo "Building linux/arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o "$OUT_DIR/Luna_linux_arm64" "$SRC"

echo "Building windows/amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "$OUT_DIR/Luna_windows_amd64.exe" "$SRC"

echo "Building darwin/amd64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o "$OUT_DIR/Luna_darwin_amd64" "$SRC"

echo "Building darwin/arm64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o "$OUT_DIR/Luna_darwin_arm64" "$SRC"

echo "Build success. Artifacts:"
