@echo off
REM Run script for Windows
REM This runs the URL Shortener with docker-compose (includes Redis)

echo ========================================
echo Starting URL Shortener with Redis
echo ========================================
echo.

REM Check if .env file exists
if not exist .env (
    echo WARNING: .env file not found!
    echo Creating from .env.example...
    copy .env.example .env
    echo.
    echo Please edit .env file with your AWS credentials
    echo Then run this script again.
    echo.
    pause
    exit /b 1
)

echo Starting services...
docker-compose up

pause
