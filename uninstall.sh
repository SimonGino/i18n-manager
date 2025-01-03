#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 删除安装的文件
rm -f "$HOME/.local/bin/i18n-manager"

# 询问是否删除配置文件
read -p "Do you want to remove configuration files as well? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    rm -rf "$HOME/.config/i18n-manager"
    echo -e "${GREEN}Configuration files removed.${NC}"
else
    echo -e "${YELLOW}Configuration files kept at $HOME/.config/i18n-manager${NC}"
fi

echo -e "${GREEN}i18n-manager has been uninstalled.${NC}" 