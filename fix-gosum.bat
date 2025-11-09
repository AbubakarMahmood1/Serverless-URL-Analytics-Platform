@echo off
REM Generate go.sum file using Docker (no Go installation needed!)

echo ========================================
echo Generating go.sum file
echo ========================================
echo.

echo This will download dependencies and create go.sum...
echo.

REM Use Go Docker image to run go mod tidy
docker run --rm -v "%cd%":/app -w /app golang:1.21-alpine go mod tidy

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Success! go.sum file created
    echo ========================================
    echo.
    echo You can now run:
    echo   docker-compose build
    echo   docker-compose up -d
    echo.
) else (
    echo.
    echo ========================================
    echo Failed to generate go.sum
    echo ========================================
    echo.
)

pause
