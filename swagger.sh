#!/bin/bash

# Exit immediately if a command exits with a non-zero status (optional but recommended)
# set -e

# Go to the directory where the script is located
cd "$(dirname "$0")" || exit 1

echo "===== 检查 swag 是否已安装 ====="
if ! command -v swag >/dev/null 2>&1; then
    echo "swag 未安装，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

echo "===== 生成 Swagger 文档 ====="
swag init -g cmd/main.go
if [ $? -ne 0 ]; then
    echo "swag init 执行失败"
    read -p "按 Enter 键退出..."
    exit 1
fi

echo "===== 验证文档生成 ====="
if [ -f "docs/swagger.json" ] && [ -f "docs/swagger.yaml" ] && [ -f "docs/docs.go" ]; then
    echo "✅ Swagger 文档生成成功"
    echo "   - docs/swagger.json"
    echo "   - docs/swagger.yaml"
    echo "   - docs/docs.go"
else
    echo "❌ 文档生成失败或文件缺失"
    read -p "按 Enter 键退出..."
    exit 1
fi

read -p "按 Enter 键结束..."
