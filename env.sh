#!/bin/sh

# 设置环境变量
export PROTOCOL="http"
# http://{ip}:{port}
# https://{域名} , 大括号只是标记
export SUB_URL_PREFIX="看注释"
export SUB_PORT="2096"
export PORT="54321"
export MODE="prod"

# 显示环境变量值
echo "PROTOCOL: $PROTOCOL"
echo "SUB_URL_PREFIX: $SUB_URL_PREFIX"
echo "SUB_PORT: $SUB_PORT"
echo "PORT: $PORT"
echo "Mode: $MODE"

# 你可以在这里添加你需要运行的命令
# 例如：./your_application
