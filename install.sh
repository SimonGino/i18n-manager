#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检测操作系统和架构
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    "Darwin")
        OS="darwin"
        case "$ARCH" in
            "x86_64") ARCH="x86_64" ;;
            "arm64") ARCH="arm64" ;;
            *) echo -e "${RED}不支持的架构: $ARCH${NC}" && exit 1 ;;
        esac
        ;;
    "Linux")
        OS="linux"
        case "$ARCH" in
            "x86_64") ARCH="x86_64" ;;
            "aarch64") ARCH="arm64" ;;
            "armv7l") ARCH="arm" ;;
            "i386"|"i686") ARCH="i386" ;;
            *) echo -e "${RED}不支持的架构: $ARCH${NC}" && exit 1 ;;
        esac
        ;;
    "MINGW"*|"MSYS"*|"CYGWIN"*)
        OS="windows"
        case "$ARCH" in
            "x86_64") ARCH="x86_64" ;;
            "i386"|"i686") ARCH="i386" ;;
            *) echo -e "${RED}不支持的架构: $ARCH${NC}" && exit 1 ;;
        esac
        ;;
    *)
        echo -e "${RED}不支持的操作系统: $OS${NC}"
        exit 1
        ;;
esac

# 设置安装目录
if [ "$OS" = "windows" ]; then
    INSTALL_DIR="$USERPROFILE/bin"
else
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ]; then
        INSTALL_DIR="$HOME/bin"
    fi
fi
mkdir -p "$INSTALL_DIR"

# 获取最新版本
echo -e "${GREEN}正在获取最新版本...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/SimonGino/i18n-manager/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}无法获取最新版本信息${NC}"
    exit 1
fi

echo -e "${GREEN}最新版本: $LATEST_VERSION${NC}"

# 构建下载URL
ARCHIVE_EXT=".tar.gz"
if [ "$OS" = "windows" ]; then
    ARCHIVE_EXT=".zip"
fi

DOWNLOAD_URL="https://github.com/SimonGino/i18n-manager/releases/download/${LATEST_VERSION}/i18n-manager_${OS}_${ARCH}${ARCHIVE_EXT}"

# 下载并安装
echo -e "${GREEN}正在下载...${NC}"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

if [ "$OS" = "windows" ]; then
    curl -L -o "i18n-manager.zip" "$DOWNLOAD_URL"
    unzip "i18n-manager.zip"
    mv "i18n-manager.exe" "$INSTALL_DIR/"
else
    curl -L "$DOWNLOAD_URL" | tar xz
    mv "i18n-manager" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/i18n-manager"
fi

# 清理
cd ..
rm -rf "$TMP_DIR"

# 验证安装
if [ "$OS" = "windows" ]; then
    if [ -f "$INSTALL_DIR/i18n-manager.exe" ]; then
        echo -e "${GREEN}i18n-manager 安装成功！${NC}"
    else
        echo -e "${RED}安装失败${NC}"
        exit 1
    fi
else
    if command -v i18n-manager &> /dev/null; then
        echo -e "${GREEN}i18n-manager 安装成功！${NC}"
    else
        echo -e "${RED}安装失败${NC}"
        exit 1
    fi
fi

echo -e "${YELLOW}请运行以下命令配置API密钥：${NC}"
echo -e "i18n-manager config --set-ai-provider <deepseek|qwen>"
echo -e "i18n-manager config --set-api-key YOUR_API_KEY" 