@echo off
chcp 65001 >nul

REM 切换到当前脚本所在目录（避免路径问题）
cd /d %~dp0

echo ===== go mod tidy =====
go mod tidy
IF %ERRORLEVEL% NEQ 0 (
    echo go mod tidy 执行失败
    pause
    exit /b 1
)

echo ===== 启动服务 =====
go run cmd\main.go -addr=8.210.176.105:2379
IF %ERRORLEVEL% NEQ 0 (
    echo 程序运行错误！
    pause
    exit /b 1
)

pause