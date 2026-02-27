set -euo pipefail

SRC="main.go"
OUT_DIR="./bin"
WINDOWS_OUT="$OUT_DIR/Luna.exe"
LINUX_OUT="$OUT_DIR/Luna_linux_amd64" 

mkdir -p "$OUT_DIR"

echo "Building linux/amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$LINUX_OUT" "$SRC"

echo "Building windows/amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "$WINDOWS_OUT" "$SRC"

echo "Build success. Artifacts:"
ls -lh "$OUT_DIR"
