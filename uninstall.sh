#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 删除安装的文件
rm -rf "$HOME/.local/lib/i18n-manager"
rm -f "$HOME/.local/bin/i18n-manager"

echo -e "${GREEN}i18n-manager has been uninstalled.${NC}" 