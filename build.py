import PyInstaller.__main__
import os
import sys

def build():
    # 获取当前操作系统
    if sys.platform.startswith('win'):
        separator = ';'
        ext = '.exe'
    else:
        separator = ':'
        ext = ''

    # 基本参数
    args = [
        'src/i18_manager/i18n_manager.py',  # 主脚本
        '--name=i18n-manager',              # 输出文件名
        '--onefile',                        # 单文件模式
        '--clean',                          # 清理临时文件
        '--noconfirm',                      # 不确认覆盖
    ]

    # 添加依赖
    hidden_imports = [
        'openai',                           # 新增 openai 依赖
        'typing_extensions',                # openai 的依赖
        'tqdm',                            # openai 的依赖
        'anyio',                           # openai 的依赖
        'httpx',                           # openai 的依赖
        'distro',                          # openai 的依赖
        'sniffio',                         # openai 的依赖
    ]
    
    for imp in hidden_imports:
        args.extend(['--hidden-import', imp])

    # 添加数据文件
    datas = [
        os.path.join('src', 'i18_manager', 'config.py'),
    ]
    
    for data in datas:
        if os.path.exists(data):
            args.extend(['--add-data', f'{data}{separator}.'])

    # 执行构建
    PyInstaller.__main__.run(args)

    # 构建完成后设置权限
    if sys.platform != 'win32':
        output_path = os.path.join('dist', 'i18n-manager')
        if os.path.exists(output_path):
            os.chmod(output_path, 0o755)  # 设置可执行权限

if __name__ == '__main__':
    build() 