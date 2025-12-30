#!/bin/bash
set -e

# 切换到脚本所在目录（避免路径问题）
cd "$(dirname "$0")"

echo "===== go mod tidy ====="
if ! go mod tidy; then
    echo "go mod tidy 执行失败"
    read -p "按回车键退出..."
    exit 1
fi

echo "===== 启动服务 ====="
if ! go run cmd/main.go -addr=8.210.176.105:2379; then
    echo "程序运行错误！"
    read -p "按回车键退出..."
    exit 1
fi

# 模拟 Windows pause
read -p "按回车键退出..."
