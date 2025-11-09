@echo off
REM Build script for Windows
REM This builds the Docker image for the URL Shortener

echo ========================================
echo Building URL Shortener Docker Image
echo ========================================

docker build -t url-shortener:latest .

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Build completed successfully!
    echo ========================================
    echo.
    echo To run the application:
    echo   run.bat
    echo.
    echo Or with docker-compose:
    echo   docker-compose up
) else (
    echo.
    echo ========================================
    echo Build failed! Check the error above.
    echo ========================================
)

pause
