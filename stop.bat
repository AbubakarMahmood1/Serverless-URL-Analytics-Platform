@echo off
REM Stop script for Windows
REM This stops all running containers

echo ========================================
echo Stopping URL Shortener Services
echo ========================================

docker-compose down

echo.
echo Services stopped successfully!
echo.

pause
