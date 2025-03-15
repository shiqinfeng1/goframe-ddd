
#!/bin/bash

# 检测当前的 Linux 发行版
detect_distribution() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo $ID
    elif [ -f /etc/redhat-release ]; then
        echo "rhel"
    else
        echo "unknown"
    fi
}

# 根据发行版安装 curl
install_curl() {
    local distro=$(detect_distribution)
    case $distro in
        ubuntu|debian)
            sudo apt-get update
            sudo apt-get install -y curl
            ;;
        centos|rhel)
            sudo yum install -y curl
            ;;
        fedora)
            sudo dnf install -y curl
            ;;
        arch)
            sudo pacman -Syu --noconfirm curl
            ;;
        alpine)
            sudo apk add --no-cache curl
            ;;
        *)
            echo "不支持的 Linux 发行版: $distro"
            exit 1
            ;;
    esac
}
install_unzip() {
    local distro=$(detect_distribution)
    case $distro in
        ubuntu|debian)
            sudo apt-get update
            sudo apt-get install -y unzip zip
            ;;
        centos|rhel)
            sudo yum install -y unzip zip
            ;;
        fedora)
            sudo dnf install -y unzip zip
            ;;
        arch)
            sudo pacman -Syu --noconfirm unzip zip
            ;;
        alpine)
            sudo apk add --no-cache unzip zip
            ;;
        *)
            echo "不支持的 Linux 发行版: $distro"
            exit 1
            ;;
    esac
}
# 定义 protoc 工具的下载地址
PROTOBUF_RELEASES_URL="https://github.com/protocolbuffers/protobuf/releases"

# 如果下载太慢，
# 方法1 ：使用代理, 先确认代理可用
# 方法2 ：替换github域名解析
# export http_proxy="10.17.11.17:1080"
# export https_proxy="10.17.11.17:1080"

# 获取当前操作系统名称
OS=$(uname -s)

# 根据操作系统选择合适的下载和安装方法
case $OS in
  Linux)
    if which protoc >/dev/null 2>&1; then
      exit 0
    fi

    if [ -z "$(which curl 2>/dev/null)" ]; then
      echo "'curl' tool not installed, try to install it."
      install_curl
    fi
    if [ -z "$(which unzip 2>/dev/null)" ]; then
      echo "'unzip' tool not installed, try to install it."
      install_unzip
    fi
    # 获取最新版本的下载链接
    DOWNLOAD_URL="$PROTOBUF_RELEASES_URL/download/v29.3/protoc-29.3-linux-x86_64.zip"
    # 下载并解压
    curl -LO -c - $DOWNLOAD_URL
    unzip $(basename $DOWNLOAD_URL) -d protoc
    sudo cp protoc/bin/protoc /usr/local/bin/
    sudo chmod +x /usr/local/bin/protoc
    rm -rf protoc $(basename $DOWNLOAD_URL)
    echo "protoc 已成功安装在 /usr/local/bin 目录下。"
    ;;
  Darwin)
    # 获取最新版本的下载链接
    DOWNLOAD_URL="$PROTOBUF_RELEASES_URL/download/v29.3/protoc-29.3-osx-x86_64.zip"

    # 下载并解压
    curl -LO $DOWNLOAD_URL
    unzip $(basename $DOWNLOAD_URL) -d protoc
    sudo cp protoc/bin/protoc /usr/local/bin/
    sudo chmod +x /usr/local/bin/protoc
    rm -rf protoc $(basename $DOWNLOAD_URL)
    echo "protoc 已成功安装在 /usr/local/bin 目录下。"
    ;;
  MINGW*|CYGWIN*)
    # 获取最新版本的下载链接
    DOWNLOAD_URL="$PROTOBUF_RELEASES_URL/download/latest/v29.3/protoc-29.3-win64.zip"

    # 下载并解压
    curl -LO $DOWNLOAD_URL
    unzip $(basename $DOWNLOAD_URL) -d protoc
    cp protoc/bin/protoc.exe /c/Windows/System32/
    rm -rf protoc $(basename $DOWNLOAD_URL)
    echo "protoc 已成功安装在 C:/Windows/System32 目录下。"
    ;;
  *)
    echo "不支持的操作系统: $OS"
    exit 1
    ;;
esac
