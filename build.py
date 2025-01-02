import PyInstaller.__main__
import sys
import os
import platform

def get_platform_options():
    if sys.platform == 'darwin':
        if platform.machine() == 'arm64':
            return [
                '--target-arch=arm64',
                '--codesign-identity=-',  # 跳过代码签名
            ]
    return []

# 确定当前操作系统
if sys.platform.startswith('win'):
    separator = ';'
else:
    separator = ':'

# 基础选项
options = [
    'src/i18_manager/i18n_manager.py',
    '--onefile',
    '--name=i18n-manager',
    f'--add-data=src/i18_manager/lang{separator}lang',
    '--clean',
    '--noupx',
    '--strip',
    '--log-level=DEBUG',
    '--hidden-import=requests',
    '--hidden-import=json',
    '--hidden-import=codecs',
    '--hidden-import=argparse',
    '--console',
]

# 添加平台特定选项
options.extend(get_platform_options())

# 运行构建
PyInstaller.__main__.run(options)

# 构建完成后设置权限
if sys.platform != 'win32':
    output_path = os.path.join('dist', 'i18n-manager')
    if os.path.exists(output_path):
        os.chmod(output_path, 0o755)
        # 在 macOS 上运行文件验证
        if sys.platform == 'darwin':
            os.system(f'codesign --force --deep --sign - {output_path}') 