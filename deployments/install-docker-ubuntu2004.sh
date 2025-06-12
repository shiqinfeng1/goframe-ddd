#!/bin/bash

# Docker 28.0.1 安装+开机自启配置脚本 for Ubuntu 20.04
# 使用方法：chmod +x install_docker.sh && sudo ./install_docker.sh

set -e

echo "正在安装 Docker 28.0.1 并配置开机启动..."

# 1. 卸载旧版本
echo "步骤 1/8: 卸载旧版本..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc || true

# 2. 更新软件包索引并安装依赖
echo "步骤 2/8: 安装依赖..."
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

# 3. 添加 Docker 官方 GPG 密钥
echo "步骤 3/8: 添加 GPG 密钥..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# 4. 添加 Docker 软件仓库
echo "步骤 4/8: 添加软件仓库..."
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

# 5. 安装特定版本 Docker 28.0.1
echo "步骤 5/8: 安装 Docker 28.0.1..."
sudo apt-get update
VERSION="5:28.0.1-1~ubuntu.20.04~focal"
sudo apt-get install -y \
    docker-ce=$VERSION \
    docker-ce-cli=$VERSION \
    containerd.io

# 6. 配置开机自启
echo "步骤 6/8: 配置开机启动..."
sudo systemctl enable docker.service
sudo systemctl enable containerd.service
echo "已设置 Docker 服务开机自启"

# 7. 验证安装
echo "步骤 7/8: 验证安装..."
sudo docker run --rm hello-world

# 8. 配置用户组（可选）
echo "步骤 8/8: 配置用户组..."
sudo usermod -aG docker $SUDO_USER
echo "请注意：需要注销并重新登录才能使docker用户组生效"

echo ""
echo "Docker 28.0.1 安装完成！已配置开机自动启动。"
echo "验证版本：docker --version"
echo "当前版本应为：Docker version 28.0.1"
echo "检查开机启动状态：systemctl is-enabled docker"