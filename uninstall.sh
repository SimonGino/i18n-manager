#!/bin/bash

# 确保脚本在交互模式下运行
exec < /dev/tty

# 检测操作系统
if [[ "$OSTYPE" == "darwin"* ]]; then
    CONFIG_DIR="$HOME/.config/i18n-manager"
    BINARY_PATH=$(which i18n-manager)
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    CONFIG_DIR="$HOME/.config/i18n-manager"
    BINARY_PATH=$(which i18n-manager)
else
    echo "Unsupported operating system"
    exit 1
fi

# 删除二进制文件
if [ -f "$BINARY_PATH" ]; then
    sudo rm "$BINARY_PATH"
    if [ $? -ne 0 ]; then
        echo "Failed to remove i18n-manager binary. Please try again with sudo."
        exit 1
    fi
fi

# 询问是否删除配置文件
echo -n "Do you want to remove configuration files as well? (y/N) "
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    if [ -d "$CONFIG_DIR" ]; then
        rm -rf "$CONFIG_DIR"
        echo "Configuration files removed."
    else
        echo "No configuration files found."
    fi
else
    echo "Configuration files kept at $CONFIG_DIR"
fi

echo "i18n-manager has been uninstalled."

# 恢复标准输入
exec <&-