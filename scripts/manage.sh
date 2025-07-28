#!/bin/bash

# WxPush 服务管理脚本
# 支持安装、卸载、状态查询等操作

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
APP_NAME="wxpush"
INSTALL_DIR="/opt/${APP_NAME}"
CONFIG_DIR="/etc/${APP_NAME}"
LOG_DIR="/var/log/${APP_NAME}"
SERVICE_NAME="${APP_NAME}"
SYSTEMD_DIR="/etc/systemd/system"
BINARY_URL="https://github.com/zdev0x/wxpush/releases/latest/download/wxpush_Linux_x86_64.tar.gz"

# 显示帮助信息
show_help() {
    echo -e "${BLUE}WxPush 服务管理脚本${NC}"
    echo
    echo -e "${YELLOW}用法:${NC}"
    echo "  $0 <command> [options]"
    echo
    echo -e "${YELLOW}命令:${NC}"
    echo "  install    安装 WxPush 服务"
    echo "  uninstall  卸载 WxPush 服务"
    echo "  status     查看服务状态"
    echo "  start      启动服务"
    echo "  stop       停止服务" 
    echo "  restart    重启服务"
    echo "  enable     设置开机启动"
    echo "  disable    取消开机启动"
    echo "  logs       查看服务日志"
    echo "  update     更新到最新版本"
    echo "  config     编辑配置文件"
    echo "  help       显示帮助信息"
    echo
    echo -e "${YELLOW}选项:${NC}"
    echo "  --force    强制执行（跳过确认）"
    echo "  --keep     卸载时保留配置文件和日志"
    echo
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0 install               # 安装服务"
    echo "  $0 uninstall --keep      # 卸载但保留配置"
    echo "  $0 status                # 查看状态"
    echo "  $0 logs                  # 查看日志"
}

# 检查root权限
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${RED}错误: 请使用root权限运行此脚本${NC}"
        exit 1
    fi
}

# 检查系统
check_system() {
    if ! command -v systemctl &> /dev/null; then
        echo -e "${RED}错误: 此脚本需要systemd支持${NC}"
        exit 1
    fi
}

# 安装服务
install_service() {
    echo -e "${YELLOW}开始安装 ${APP_NAME}...${NC}"
    
    # 创建必要的目录
    echo -e "${GREEN}创建目录...${NC}"
    mkdir -p ${INSTALL_DIR}
    mkdir -p ${CONFIG_DIR}
    mkdir -p ${LOG_DIR}
    
    # 下载最新版本
    echo -e "${GREEN}下载最新版本...${NC}"
    TMP_DIR=$(mktemp -d)
    if ! curl -L ${BINARY_URL} -o ${TMP_DIR}/app.tar.gz; then
        echo -e "${RED}错误: 下载失败${NC}"
        rm -rf ${TMP_DIR}
        exit 1
    fi
    
    if ! tar xzf ${TMP_DIR}/app.tar.gz -C ${TMP_DIR}; then
        echo -e "${RED}错误: 解压失败${NC}"
        rm -rf ${TMP_DIR}
        exit 1
    fi
    
    mv ${TMP_DIR}/wxpush ${INSTALL_DIR}/
    rm -rf ${TMP_DIR}
    
    # 复制配置文件
    echo -e "${GREEN}配置文件...${NC}"
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        # 创建默认配置文件
        cat > ${CONFIG_DIR}/config.yaml << 'EOF'
# 微信公众号配置
wechat:
  app_id: "REPLACE_ME_APPID"
  app_secret: "REPLACE_ME_SECRET"
  token: "REPLACE_ME_TOKEN"

# 模板配置
templates:
  - name: "notification"
    id: "REPLACE_ME_TEMPLATE_ID"
    title: "通知消息"
    content: "内容：{{CONTENT.DATA}} 时间：{{TIME.DATA}}"

# 用户配置
users:
  - name: "user1"
    openid: "REPLACE_ME_OPENID"

# 通知组配置
notify_groups:
  - name: "default"
    users: ["user1"]

# 服务配置
server:
  api_key: "REPLACE_ME_API_KEY"
  listen_addr: ":8801"
  log_file: "/var/log/wxpush/push.log"
  mode: "release"
EOF
        echo -e "${YELLOW}配置文件已创建: ${CONFIG_DIR}/config.yaml${NC}"
        echo -e "${YELLOW}请修改配置文件后再启动服务${NC}"
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
    echo -e "\n${YELLOW}后续操作:${NC}"
    echo -e "1. 编辑配置文件: ${BLUE}$0 config${NC}"
    echo -e "2. 启动服务: ${BLUE}$0 start${NC}"
    echo -e "3. 查看状态: ${BLUE}$0 status${NC}"
    echo -e "4. 查看日志: ${BLUE}$0 logs${NC}"
}

# 卸载服务
uninstall_service() {
    local keep_data=false
    
    # 检查参数
    for arg in "$@"; do
        case $arg in
            --keep)
                keep_data=true
                ;;
        esac
    done
    
    echo -e "${YELLOW}开始卸载 ${APP_NAME}...${NC}"
    
    # 停止并禁用服务
    echo -e "${GREEN}停止服务...${NC}"
    systemctl stop ${SERVICE_NAME} 2>/dev/null || true
    systemctl disable ${SERVICE_NAME} 2>/dev/null || true
    
    # 删除服务文件
    echo -e "${GREEN}删除服务文件...${NC}"
    rm -f ${SYSTEMD_DIR}/${SERVICE_NAME}.service
    systemctl daemon-reload
    
    # 删除应用文件
    echo -e "${GREEN}删除应用文件...${NC}"
    rm -rf ${INSTALL_DIR}
    
    # 处理配置文件和日志
    if [ "$keep_data" = true ]; then
        echo -e "${YELLOW}保留配置文件和日志:${NC}"
        echo -e "  配置文件: ${CONFIG_DIR}"
        echo -e "  日志文件: ${LOG_DIR}"
    else
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
    fi
    
    echo -e "${GREEN}卸载完成!${NC}"
}

# 查看服务状态
show_status() {
    echo -e "${BLUE}=== WxPush 服务状态 ===${NC}"
    
    if systemctl is-active --quiet ${SERVICE_NAME}; then
        echo -e "服务状态: ${GREEN}运行中${NC}"
    else
        echo -e "服务状态: ${RED}已停止${NC}"
    fi
    
    if systemctl is-enabled --quiet ${SERVICE_NAME}; then
        echo -e "开机启动: ${GREEN}已启用${NC}"
    else
        echo -e "开机启动: ${RED}已禁用${NC}"
    fi
    
    if [ -f "${INSTALL_DIR}/wxpush" ]; then
        echo -e "安装目录: ${GREEN}${INSTALL_DIR}${NC}"
    else
        echo -e "安装状态: ${RED}未安装${NC}"
        return
    fi
    
    if [ -f "${CONFIG_DIR}/config.yaml" ]; then
        echo -e "配置文件: ${GREEN}${CONFIG_DIR}/config.yaml${NC}"
    else
        echo -e "配置文件: ${RED}不存在${NC}"
    fi
    
    echo -e "日志目录: ${BLUE}${LOG_DIR}${NC}"
    echo
    
    # 显示详细状态
    systemctl status ${SERVICE_NAME} --no-pager -l
}

# 启动服务
start_service() {
    echo -e "${GREEN}启动 ${APP_NAME} 服务...${NC}"
    systemctl start ${SERVICE_NAME}
    
    if systemctl is-active --quiet ${SERVICE_NAME}; then
        echo -e "${GREEN}服务启动成功${NC}"
    else
        echo -e "${RED}服务启动失败${NC}"
        echo -e "${YELLOW}查看详细日志: $0 logs${NC}"
    fi
}

# 停止服务
stop_service() {
    echo -e "${YELLOW}停止 ${APP_NAME} 服务...${NC}"
    systemctl stop ${SERVICE_NAME}
    echo -e "${GREEN}服务已停止${NC}"
}

# 重启服务
restart_service() {
    echo -e "${YELLOW}重启 ${APP_NAME} 服务...${NC}"
    systemctl restart ${SERVICE_NAME}
    
    if systemctl is-active --quiet ${SERVICE_NAME}; then
        echo -e "${GREEN}服务重启成功${NC}"
    else
        echo -e "${RED}服务重启失败${NC}"
        echo -e "${YELLOW}查看详细日志: $0 logs${NC}"
    fi
}

# 启用服务
enable_service() {
    echo -e "${GREEN}设置 ${APP_NAME} 开机启动...${NC}"
    systemctl enable ${SERVICE_NAME}
    echo -e "${GREEN}开机启动已启用${NC}"
}

# 禁用服务
disable_service() {
    echo -e "${YELLOW}取消 ${APP_NAME} 开机启动...${NC}"
    systemctl disable ${SERVICE_NAME}
    echo -e "${GREEN}开机启动已禁用${NC}"
}

# 查看日志
show_logs() {
    echo -e "${BLUE}=== ${APP_NAME} 服务日志 ===${NC}"
    echo -e "${YELLOW}按 Ctrl+C 退出日志查看${NC}"
    echo
    journalctl -u ${SERVICE_NAME} -f --no-pager
}

# 更新服务
update_service() {
    echo -e "${YELLOW}更新 ${APP_NAME} 到最新版本...${NC}"
    
    if ! systemctl is-active --quiet ${SERVICE_NAME}; then
        echo -e "${RED}服务未运行，直接更新二进制文件${NC}"
    else
        echo -e "${GREEN}停止服务...${NC}"
        systemctl stop ${SERVICE_NAME}
    fi
    
    # 备份当前版本
    if [ -f "${INSTALL_DIR}/wxpush" ]; then
        cp ${INSTALL_DIR}/wxpush ${INSTALL_DIR}/wxpush.backup
        echo -e "${GREEN}已备份当前版本${NC}"
    fi
    
    # 下载最新版本
    echo -e "${GREEN}下载最新版本...${NC}"
    TMP_DIR=$(mktemp -d)
    if ! curl -L ${BINARY_URL} -o ${TMP_DIR}/app.tar.gz; then
        echo -e "${RED}错误: 下载失败${NC}"
        rm -rf ${TMP_DIR}
        exit 1
    fi
    
    if ! tar xzf ${TMP_DIR}/app.tar.gz -C ${TMP_DIR}; then
        echo -e "${RED}错误: 解压失败${NC}"
        rm -rf ${TMP_DIR}
        exit 1
    fi
    
    mv ${TMP_DIR}/wxpush ${INSTALL_DIR}/
    chmod +x ${INSTALL_DIR}/wxpush
    rm -rf ${TMP_DIR}
    
    echo -e "${GREEN}启动服务...${NC}"
    systemctl start ${SERVICE_NAME}
    
    if systemctl is-active --quiet ${SERVICE_NAME}; then
        echo -e "${GREEN}更新完成，服务运行正常${NC}"
        rm -f ${INSTALL_DIR}/wxpush.backup
    else
        echo -e "${RED}更新后服务启动失败，恢复备份版本${NC}"
        mv ${INSTALL_DIR}/wxpush.backup ${INSTALL_DIR}/wxpush
        systemctl start ${SERVICE_NAME}
    fi
}

# 编辑配置文件
edit_config() {
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        echo -e "${RED}配置文件不存在: ${CONFIG_DIR}/config.yaml${NC}"
        echo -e "${YELLOW}请先安装服务: $0 install${NC}"
        exit 1
    fi
    
    # 检查编辑器
    if command -v vim &> /dev/null; then
        vim ${CONFIG_DIR}/config.yaml
    elif command -v nano &> /dev/null; then
        nano ${CONFIG_DIR}/config.yaml
    else
        echo -e "${RED}未找到可用的编辑器 (vim/nano)${NC}"
        echo -e "${YELLOW}配置文件路径: ${CONFIG_DIR}/config.yaml${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}配置文件已修改${NC}"
    echo -e "${YELLOW}重启服务使配置生效: $0 restart${NC}"
}

# 主函数
main() {
    case "${1:-}" in
        install)
            check_root
            check_system
            install_service
            ;;
        uninstall)
            check_root
            check_system
            shift
            uninstall_service "$@"
            ;;
        status)
            show_status
            ;;
        start)
            check_root
            start_service
            ;;
        stop)
            check_root
            stop_service
            ;;
        restart)
            check_root
            restart_service
            ;;
        enable)
            check_root
            enable_service
            ;;
        disable)
            check_root
            disable_service
            ;;
        logs)
            show_logs
            ;;
        update)
            check_root
            check_system
            update_service
            ;;
        config)
            check_root
            edit_config
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}错误: 未知命令 '${1:-}'${NC}"
            echo
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"