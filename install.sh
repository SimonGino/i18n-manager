#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 检测系统类型和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 确定下载URL
GITHUB_REPO="SimonGino/i18n-manager"
VERSION="latest"

case "$OS" in
    linux)
        case "$ARCH" in
            x86_64)
                BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/i18n-manager-linux-amd64"
                ;;
            aarch64)
                BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/i18n-manager-linux-arm64"
                ;;
            *)
                echo -e "${RED}Unsupported architecture: $ARCH${NC}"
                exit 1
                ;;
        esac
        ;;
    darwin)
        case "$ARCH" in
            x86_64)
                BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/i18n-manager-darwin-amd64"
                ;;
            arm64)
                BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/i18n-manager-darwin-arm64"
                ;;
            *)
                echo -e "${RED}Unsupported architecture: $ARCH${NC}"
                exit 1
                ;;
        esac
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# 创建安装目录
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

# 下载二进制文件
echo "Downloading i18n-manager..."
curl -L "$BINARY_URL" -o "$INSTALL_DIR/i18n-manager"
chmod +x "$INSTALL_DIR/i18n-manager"

# 检查 PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc" 2>/dev/null || true
fi

echo -e "${GREEN}Installation completed!${NC}"
echo -e "Please restart your terminal or run: source ~/.bashrc"
echo -e "Then you can use 'i18n-manager' command anywhere." 