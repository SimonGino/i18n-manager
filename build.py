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
    '--clean',
    '--noupx',
    '--strip',  # 减小二进制大小
    '--log-level=DEBUG',  # 添加调试信息
    # 添加运行时钩子
    '--hidden-import=requests',
    '--hidden-import=json',
    '--hidden-import=codecs',
    '--hidden-import=argparse',
    # 设置入口点
    '--console',
])

# 构建完成后设置权限
if sys.platform != 'win32':
    output_path = os.path.join('dist', 'i18n-manager')
    if os.path.exists(output_path):
        os.chmod(output_path, 0o755) 