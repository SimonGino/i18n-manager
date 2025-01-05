#!/usr/bin/env bash

# 定义颜色（检查终端是否支持颜色）
if [ -t 1 ]; then
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    NC='\033[0m'
else
    GREEN=''
    YELLOW=''
    NC=''
fi

# 检测操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)     CONFIG_DIR="$HOME/.config/i18n-manager";;
    Darwin*)    CONFIG_DIR="$HOME/.config/i18n-manager";;
    *)          CONFIG_DIR="$HOME/.config/i18n-manager";;
esac

# 检测安装位置
INSTALL_LOCATIONS=(
    "$HOME/.local/bin/i18n-manager"
    "/usr/local/bin/i18n-manager"
    "/opt/homebrew/bin/i18n-manager"
)

# 删除安装的文件
for location in "${INSTALL_LOCATIONS[@]}"; do
    if [ -f "$location" ]; then
        rm -f "$location"
        echo -e "${GREEN}Removed $location${NC}"
    fi
done

# 询问是否删除配置文件
if [ -d "$CONFIG_DIR" ]; then
    printf "Do you want to remove configuration files as well? (y/N) "
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}Configuration files removed.${NC}"
    else
        echo -e "${YELLOW}Configuration files kept at $CONFIG_DIR${NC}"
    fi
fi

echo -e "${GREEN}i18n-manager has been uninstalled.${NC}"