@echo off
REM Run Go unit tests in Docker (no Go installation required)

echo ========================================
echo Running Unit Tests in Docker
echo ========================================
echo.

docker build -f Dockerfile.test -t url-shortener-tests .

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Running tests...
    echo.
    docker run --rm url-shortener-tests
) else (
    echo.
    echo Build failed!
    echo.
)

pause
