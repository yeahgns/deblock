set -euo pipefail

REPO="yeahgns/deblock"
BINARY_NAME="mc-deblock"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported operating system: $OS"; exit 1 ;;
esac

ASSET="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

echo "Downloading ${ASSET}..."
TMP_DIR="$(mktemp -d)"
curl -sSL "$URL" -o "${TMP_DIR}/${ASSET}"

echo "Extracting..."
tar -xzf "${TMP_DIR}/${ASSET}" -C "${TMP_DIR}"

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="${HOME}/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
rm -rf "$TMP_DIR"

echo ""
echo "Installed in: ${INSTALL_DIR}/${BINARY_NAME}"
if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
  echo "Warning: ${INSTALL_DIR} is not in your PATH."
  echo "Add this line to your ~/.bashrc or ~/.zshrc:"
  echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
echo "Run with: ${BINARY_NAME}"