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
BINARY_URL="https://github.com/zdev0x/wxpush/releases/latest/download/wxpush_Linux_x86_64.tar.gz"

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用root权限运行此脚本${NC}"
    exit 1
fi

echo -e "${YELLOW}开始安装 ${APP_NAME}...${NC}"

# 创建必要的目录
echo -e "${GREEN}创建目录...${NC}"
mkdir -p ${INSTALL_DIR}
mkdir -p ${CONFIG_DIR}
mkdir -p ${LOG_DIR}

# 下载最新版本
echo -e "${GREEN}下载最新版本...${NC}"
TMP_DIR=$(mktemp -d)
curl -L ${BINARY_URL} -o ${TMP_DIR}/app.tar.gz
tar xzf ${TMP_DIR}/app.tar.gz -C ${TMP_DIR}
mv ${TMP_DIR}/wxpush ${INSTALL_DIR}/
rm -rf ${TMP_DIR}

# 复制配置文件
echo -e "${GREEN}配置文件...${NC}"
if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
    if [ -f "${INSTALL_DIR}/config.example.yaml" ]; then
        cp ${INSTALL_DIR}/config.example.yaml ${CONFIG_DIR}/config.yaml
        echo -e "${YELLOW}配置文件已创建: ${CONFIG_DIR}/config.yaml${NC}"
        echo -e "${YELLOW}请修改配置文件后再启动服务${NC}"
    else
        echo -e "${RED}警告: 未找到配置文件示例${NC}"
    fi
else
    echo -e "${YELLOW}配置文件已存在，跳过创建${NC}"
fi

# 创建systemd服务文件
echo -e "${GREEN}创建systemd服务...${NC}"
cat > ${SYSTEMD_DIR}/${SERVICE_NAME}.service << EOF
[Unit]
Description=WeChat Push Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=${INSTALL_DIR}/wxpush -c ${CONFIG_DIR}/config.yaml
Restart=always
RestartSec=10
Environment=TZ=Asia/Shanghai

[Install]
WantedBy=multi-user.target
EOF

# 设置权限
echo -e "${GREEN}设置权限...${NC}"
chmod +x ${INSTALL_DIR}/wxpush
chmod 644 ${SYSTEMD_DIR}/${SERVICE_NAME}.service
chmod 755 ${CONFIG_DIR}
chmod 755 ${LOG_DIR}

# 重新加载systemd配置
echo -e "${GREEN}重新加载systemd配置...${NC}"
systemctl daemon-reload

echo -e "${GREEN}安装完成!${NC}"
echo -e "\n使用说明:"
echo -e "1. 编辑配置文件:"
echo -e "   ${YELLOW}vim ${CONFIG_DIR}/config.yaml${NC}"
echo -e "\n2. 管理服务:"
echo -e "   启动: ${YELLOW}systemctl start ${SERVICE_NAME}${NC}"
echo -e "   停止: ${YELLOW}systemctl stop ${SERVICE_NAME}${NC}"
echo -e "   重启: ${YELLOW}systemctl restart ${SERVICE_NAME}${NC}"
echo -e "   状态: ${YELLOW}systemctl status ${SERVICE_NAME}${NC}"
echo -e "   开机启动: ${YELLOW}systemctl enable ${SERVICE_NAME}${NC}"
echo -e "\n3. 查看日志:"
echo -e "   服务日志: ${YELLOW}journalctl -u ${SERVICE_NAME} -f${NC}"
echo -e "   应用日志: ${YELLOW}tail -f ${LOG_DIR}/push.log${NC}"