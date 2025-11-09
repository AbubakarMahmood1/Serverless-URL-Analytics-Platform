@echo off
REM Automated API Test Script for Windows
REM Tests all endpoints and provides detailed results

echo ========================================
echo URL Shortener - Automated API Tests
echo ========================================
echo.

REM Check if services are running
docker ps | findstr url-shortener-api >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Services not running!
    echo Please start services first with: run.bat
    echo Or: docker-compose up -d
    pause
    exit /b 1
)

echo Services are running. Starting tests...
echo.

REM Use PowerShell to run the detailed tests
powershell -ExecutionPolicy Bypass -File test-api.ps1

pause
