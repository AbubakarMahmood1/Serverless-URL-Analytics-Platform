@echo off
REM Kill process using port 8080

echo ========================================
echo Finding process on port 8080...
echo ========================================
echo.

REM Find the process using port 8080
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8080 ^| findstr LISTENING') do (
    set PID=%%a
    goto :found
)

echo No process found listening on port 8080
echo.
pause
exit /b 0

:found
echo Found process with PID: %PID%
echo.

REM Get process name
for /f "tokens=1" %%b in ('tasklist /FI "PID eq %PID%" /FO LIST ^| findstr "Image"') do (
    echo Process: %%b
)

echo.
set /p CONFIRM="Kill this process? (Y/N): "
if /i "%CONFIRM%" NEQ "Y" (
    echo Cancelled.
    pause
    exit /b 0
)

echo.
echo Killing process %PID%...
taskkill /F /PID %PID%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo SUCCESS: Process killed successfully!
    echo You can now run: .\run.bat
) else (
    echo.
    echo ERROR: Failed to kill process
    echo Try running this script as Administrator
)

echo.
pause
