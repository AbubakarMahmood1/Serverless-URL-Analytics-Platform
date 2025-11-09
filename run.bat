@echo off
REM Run script for Windows
REM This runs the URL Shortener with docker-compose (includes Redis)

echo ========================================
echo Starting URL Shortener
echo ========================================
echo.

REM Check if .env file exists
if not exist .env (
    echo WARNING: .env file not found!
    echo Creating from .env.example...
    copy .env.example .env
    echo.
    echo .env file created with local DynamoDB settings
    echo You can start testing immediately!
    echo.
)

REM Check if go.sum exists, generate if needed
if not exist go.sum (
    echo INFO: Generating go.sum file...
    echo This only needs to happen once.
    echo.
    docker run --rm -v "%cd%":/app -w /app golang:1.21-alpine go mod tidy
    if %ERRORLEVEL% EQU 0 (
        echo OK go.sum created successfully
        echo.
    ) else (
        echo WARNING: Could not create go.sum, but Docker will handle it
        echo.
    )
)

echo Starting services...
docker-compose up -d

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Services Started Successfully!
    echo ========================================
    echo.
    echo API: http://localhost:8080
    echo DynamoDB Local: http://localhost:8000
    echo Redis: localhost:6379
    echo.
    echo Next steps:
    echo   1. Create tables: .\setup-local-tables-docker.bat
    echo   2. Run tests: .\test.bat
    echo.
    echo View logs: docker-compose logs -f
    echo Stop: docker-compose down
    echo.
)

pause
