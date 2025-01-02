#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 检测系统类型和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 转换架构名称
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# 确定下载URL
GITHUB_REPO="SimonGino/i18n-manager"
VERSION=$(curl -s https://api.github.com/repos/${GITHUB_REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Could not determine latest version${NC}"
    exit 1
fi

# 构建下载URL
BINARY_URL="https://github.com/$GITHUB_REPO/releases/download/${VERSION}/i18n-manager-${OS}-${ARCH}"

# 创建临时目录
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/i18n-manager"

# 创建安装目录
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

# 下载二进制文件
echo -e "Downloading i18n-manager ${VERSION} for ${OS}-${ARCH}..."
echo -e "From: ${BINARY_URL}"

if ! curl -L -o "$TMP_FILE" "$BINARY_URL"; then
    echo -e "${RED}Download failed${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 检查文件类型
FILE_TYPE=$(file "$TMP_FILE")
if [[ "$OS" == "darwin" && ! "$FILE_TYPE" =~ "Mach-O 64-bit executable arm64" ]]; then
    echo -e "${RED}Error: Invalid binary format${NC}"
    echo -e "${RED}Expected: Mach-O 64-bit executable arm64${NC}"
    echo -e "${RED}Got: $FILE_TYPE${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 设置执行权限
chmod +x "$TMP_FILE"

# 简单测试二进制文件
echo -e "Verifying binary..."
if ! "$TMP_FILE" 2>&1 | grep -q "usage\|help\|i18n-manager"; then
    echo -e "${RED}Warning: Binary might not be working correctly${NC}"
    echo -e "${RED}File type: $(file "$TMP_FILE")${NC}"
    # 继续安装，但显示警告
    echo -e "${RED}Continuing installation despite verification warning...${NC}"
fi

# 移动到安装目录
mv "$TMP_FILE" "$INSTALL_DIR/i18n-manager"

# 清理临时目录
rm -rf "$TMP_DIR"

# 检查 PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc" 2>/dev/null || true
fi

echo -e "${GREEN}Installation completed!${NC}"
echo -e "Location: $INSTALL_DIR/i18n-manager"
echo -e "\nPlease restart your terminal or run: source ~/.bashrc"
echo -e "Then you can use 'i18n-manager' command anywhere."

# 尝试运行帮助命令
echo -e "\nTrying to run help command..."
"$INSTALL_DIR/i18n-manager" --help || true 