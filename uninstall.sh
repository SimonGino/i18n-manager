#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 检测操作系统
OS="$(uname -s)"
case "$OS" in
    "MINGW"*|"MSYS"*|"CYGWIN"*)
        IS_WINDOWS=1
        ;;
    *)
        IS_WINDOWS=0
        ;;
esac

# 查找可能的安装位置
if [ $IS_WINDOWS -eq 1 ]; then
    POSSIBLE_LOCATIONS=(
        "$USERPROFILE/bin/i18n-manager.exe"
    )
    CONFIG_DIR="$APPDATA/i18n-manager"
else
    POSSIBLE_LOCATIONS=(
        "/usr/local/bin/i18n-manager"
        "$HOME/bin/i18n-manager"
        "$HOME/.local/bin/i18n-manager"
    )
    if [ -d "$HOME/.config" ]; then
        CONFIG_DIR="$HOME/.config/i18n-manager"
    else
        CONFIG_DIR="$HOME/.i18n-manager"
    fi
fi

# 删除二进制文件
FOUND=0
for LOCATION in "${POSSIBLE_LOCATIONS[@]}"; do
    if [ -f "$LOCATION" ]; then
        echo -e "${GREEN}找到安装文件: $LOCATION${NC}"
        rm -f "$LOCATION"
        echo -e "${GREEN}已删除二进制文件${NC}"
        FOUND=1
    fi
done

if [ $FOUND -eq 0 ]; then
    echo -e "${RED}未找到i18n-manager的安装文件${NC}"
fi

# 删除配置目录
if [ -d "$CONFIG_DIR" ]; then
    echo -e "${GREEN}找到配置目录: $CONFIG_DIR${NC}"
    rm -rf "$CONFIG_DIR"
    echo -e "${GREEN}已删除配置目录${NC}"
else
    echo -e "${RED}未找到配置目录${NC}"
fi

echo -e "${GREEN}卸载完成${NC}"