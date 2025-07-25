#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
APP_NAME="wxpush"
INSTALL_DIR="/opt/${APP_NAME}"
CONFIG_DIR="/etc/${APP_NAME}"
LOG_DIR="/var/log/${APP_NAME}"
SERVICE_NAME="${APP_NAME}"
SYSTEMD_DIR="/etc/systemd/system"

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用root权限运行此脚本${NC}"
    exit 1
fi

echo -e "${YELLOW}开始卸载 ${APP_NAME}...${NC}"

# 停止并禁用服务
echo -e "${GREEN}停止服务...${NC}"
systemctl stop ${SERVICE_NAME}
systemctl disable ${SERVICE_NAME}

# 删除服务文件
echo -e "${GREEN}删除服务文件...${NC}"
rm -f ${SYSTEMD_DIR}/${SERVICE_NAME}.service
systemctl daemon-reload

# 删除应用文件
echo -e "${GREEN}删除应用文件...${NC}"
rm -rf ${INSTALL_DIR}

# 询问是否删除配置和日志
read -p "是否删除配置文件和日志? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}删除配置文件和日志...${NC}"
    rm -rf ${CONFIG_DIR}
    rm -rf ${LOG_DIR}
    echo -e "${YELLOW}配置文件和日志已删除${NC}"
else
    echo -e "${YELLOW}保留配置文件和日志:${NC}"
    echo -e "  配置文件: ${CONFIG_DIR}"
    echo -e "  日志文件: ${LOG_DIR}"
fi

echo -e "${GREEN}卸载完成!${NC}" 