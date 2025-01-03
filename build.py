import PyInstaller.__main__
import sys
import os

# 确定当前操作系统
if sys.platform.startswith('win'):
    separator = ';'
else:
    separator = ':'

# 运行构建
PyInstaller.__main__.run([
    'src/i18_manager/__main__.py',
    '--onefile',
    '--name=i18n-manager',
    f'--add-data=src/i18_manager{separator}i18_manager',
    '--clean',
    '--hidden-import=requests',
    '--hidden-import=json',
    '--hidden-import=codecs',
    '--hidden-import=argparse',
    '--hidden-import=pathlib',
])

# 构建完成后设置权限
if sys.platform != 'win32':
    output_path = os.path.join('dist', 'i18n-manager')
    if os.path.exists(output_path):
        os.chmod(output_path, 0o755) 