@echo off
chcp 65001 >nul

cd /d %~dp0

echo ===== 检查 swag 是否已安装 =====
swag --version >nul 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo swag 未安装，正在安装...
    go install github.com/swaggo/swag/cmd/swag@latest
)

echo ===== 生成 Swagger 文档 =====
swag init -g cmd\main.go
IF %ERRORLEVEL% NEQ 0 (
    echo swag init 执行失败
    pause
    exit /b 1
)

echo ===== 验证文档生成 =====
IF EXIST "docs\swagger.json" IF EXIST "docs\swagger.yaml" IF EXIST "docs\docs.go" (
    echo ✅ Swagger 文档生成成功
    echo    - docs/swagger.json
    echo    - docs/swagger.yaml
    echo    - docs/docs.go
) ELSE (
    echo ❌ 文档生成失败或文件缺失
    pause
    exit /b 1
)

pause
