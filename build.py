import PyInstaller.__main__
import sys
import os

# 确定当前操作系统
if sys.platform.startswith('win'):
    separator = ';'
else:
    separator = ':'

PyInstaller.__main__.run([
    'src/i18_manager/i18n_manager.py',
    '--onefile',
    '--name=i18n-manager',
    f'--add-data=src/i18_manager/lang{separator}lang',
    '--clean'
]) 