#!/bin/bash

# Docker 28.0.1 安装+开机自启配置脚本 for Ubuntu 20.04
# 使用方法：chmod +x install_docker.sh && sudo ./install_docker.sh

set -e

function ensure_docker_in_path() {
  # 检查 docker 命令是否存在于 PATH 中
  if command -v docker &> /dev/null; then
    echo -e "Found Docker at $(which docker)"
    
    # 获取 Docker 版本号
    local docker_version=$(docker --version | awk '{print $3}' | tr -d ',')
 
    # 解析版本号
    local major=$(echo $docker_version | cut -d. -f1)
    local minor=$(echo $docker_version | cut -d. -f2)
    local patch=$(echo $docker_version | cut -d. -f3)
    
    # 检查版本是否大于或等于 28.x.x
    if [ "$major" -gt 28 ]; then
      echo -e "Docker version $docker_version is compatible (>= 28.x.x)"
      return 0
    elif [ "$major" -eq 28 ]; then
      echo -e "Docker version $docker_version is compatible (28.x.x)"
      return 0
    else
      echo -e "Docker version $docker_version is too old (< 28.x.x)"
      echo -e "Please upgrade Docker to version 28 or higher"
      echo -e "See https://docs.docker.com/engine/install/ for upgrade instructions"
      return 1
    fi
  else
    echo -e "Can't find 'docker' in PATH"
    echo -e "Please install Docker or add it to your PATH"
    echo -e "See https://docs.docker.com/engine/install/ for installation instructions"
    return 1
  fi
}

# 调用函数并在失败时退出
if ensure_docker_in_path; then
  echo ""
  exit 0
fi

echo -e "\nstart install Docker 28.0.1 ..."

# 1. 卸载旧版本
echo "Step 1/8: uninstall old version ..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc || true

# 2. 更新软件包索引并安装依赖
echo "Step 2/8: install deps ..."
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

# 3. 添加 Docker 官方 GPG 密钥
echo "Step 3/8: add GPG key ..."
# 尝试阿里云镜像源
if ! sudo curl -fSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -; then
    echo "方案1失败，尝试备用镜像源..."
    
    # 尝试清华源
    if ! sudo curl -fSL https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/gpg | sudo apt-key add -; then
        echo "方案2失败，尝试手动下载..."
        
        # 手动下载密钥
        if ! curl -o /tmp/docker.gpg https://download.docker.com/linux/ubuntu/gpg --retry 3; then
            echo "错误：所有密钥获取方式均失败"
            echo "请检查网络连接或手动下载密钥文件"
            exit 1
        fi
        
        sudo apt-key add /tmp/docker.gpg
        rm /tmp/docker.gpg
    fi
fi

# 4. 添加 Docker 软件仓库
echo "Step 4/8: add repo ..."
sudo add-apt-repository \
   "deb [arch=$(ARCH)] https://mirrors.aliyun.com/docker-ce/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

# 5. 安装特定版本 Docker 28.0.1
echo "Step 5/8: install Docker 28.0.1..."
sudo apt-get update
VERSION="5:28.0.1-1~ubuntu.20.04~focal"
sudo apt-get install -y \
    docker-ce=$VERSION 

# 6. 配置开机自启
echo "Step 6/8: config to running at startup ..."
sudo systemctl enable docker.service
sudo systemctl enable containerd.service
echo "已设置 Docker 服务开机自启"

# 7. 配置用户组（可选）
echo "Step 7/8: config docker user group ..."
sudo usermod -aG docker $USER
echo "!note!: NEED RESTART for without 'sudo' !!!"

# 8. 验证安装
echo "Step 8/8: verify ..."
ensure_docker_in_path
sudo docker version

echo ""
echo "Docker 28.0.1 install completed!"

