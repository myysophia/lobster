#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
DIST_ROOT="$PROJECT_ROOT/dist"

if [[ -z "${GOOS:-}" || -z "${GOARCH:-}" ]]; then
  echo "GOOS/GOARCH 必须提供"
  exit 1
fi

VERSION=${VERSION:-"dev-$(date -u +%Y%m%d%H%M%S)"}
PACKAGE_NAME="lobster_${VERSION}_${GOOS}_${GOARCH}"
WORK_DIR="$DIST_ROOT/work/${GOOS}_${GOARCH}"
BUILD_DIR="$WORK_DIR/build"
PACKAGE_DIR="$WORK_DIR/$PACKAGE_NAME"
CHECKSUMS_FILE="$DIST_ROOT/SHA256SUMS"
PACKAGE_BASENAME=""

mkdir -p "$DIST_ROOT"
rm -rf "$WORK_DIR"
mkdir -p "$BUILD_DIR"

export GOOS GOARCH

build_binary() {
  local target=$1
  local extension=""
  if [[ "$GOOS" == "windows" ]]; then
    extension=".exe"
  fi
  local output="$BUILD_DIR/${target}${extension}"
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$output" "./cmd/$target"
}

build_binary lobster
build_binary wb

mkdir -p "$PACKAGE_DIR"

cp "$BUILD_DIR/"* "$PACKAGE_DIR/"
cp "$PROJECT_ROOT/README.md" "$PACKAGE_DIR/"

if [[ "$GOOS" == "windows" ]]; then
  PACKAGE_FILE="$DIST_ROOT/${PACKAGE_NAME}.zip"
  (
    cd "$WORK_DIR"
    zip -qr "$PACKAGE_FILE" "$PACKAGE_NAME"
  )
else
  PACKAGE_FILE="$DIST_ROOT/${PACKAGE_NAME}.tar.gz"
  tar -czf "$PACKAGE_FILE" -C "$WORK_DIR" "$PACKAGE_NAME"
fi

PACKAGE_BASENAME=$(basename "$PACKAGE_FILE")

HASH_CMD=""
if command -v sha256sum >/dev/null 2>&1; then
  HASH_CMD="sha256sum"
else
  HASH_CMD="shasum -a 256"
fi

CHECKSUM_FILE="$DIST_ROOT/${PACKAGE_NAME}.sha256"
(
  cd "$DIST_ROOT"
  $HASH_CMD "$PACKAGE_BASENAME" > "$(basename "$CHECKSUM_FILE")"
)

if [[ -f "$CHECKSUMS_FILE" ]]; then
  grep -v "  ${PACKAGE_BASENAME}\$" "$CHECKSUMS_FILE" > "${CHECKSUMS_FILE}.tmp" || true
  mv "${CHECKSUMS_FILE}.tmp" "$CHECKSUMS_FILE"
fi
cat "$CHECKSUM_FILE" >> "$CHECKSUMS_FILE"

echo "打包完成：$PACKAGE_FILE"
echo "校验文件：$CHECKSUM_FILE"
echo "汇总校验：$CHECKSUMS_FILE"
