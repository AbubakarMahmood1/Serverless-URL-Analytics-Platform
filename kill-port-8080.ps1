#!/usr/bin/env pwsh

# Kill process using port 8080

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Finding process on port 8080..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Find the process using port 8080
$connections = Get-NetTCPConnection -LocalPort 8080 -State Listen -ErrorAction SilentlyContinue

if (-not $connections) {
    Write-Host "✓ No process found listening on port 8080" -ForegroundColor Green
    Write-Host ""
    Write-Host "Port 8080 is available. You can run:" -ForegroundColor Yellow
    Write-Host "  .\run.bat" -ForegroundColor White
    Write-Host ""
    Read-Host "Press Enter to exit"
    exit 0
}

$pid = $connections[0].OwningProcess
$process = Get-Process -Id $pid

Write-Host "Found process:" -ForegroundColor Yellow
Write-Host "  PID:  $pid" -ForegroundColor White
Write-Host "  Name: $($process.Name)" -ForegroundColor White
Write-Host "  Path: $($process.Path)" -ForegroundColor Gray
Write-Host ""

$confirm = Read-Host "Kill this process? (Y/N)"

if ($confirm -ne "Y" -and $confirm -ne "y") {
    Write-Host ""
    Write-Host "Cancelled." -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Press Enter to exit"
    exit 0
}

Write-Host ""
Write-Host "Killing process $pid..." -ForegroundColor Yellow

try {
    Stop-Process -Id $pid -Force
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "SUCCESS: Process killed successfully!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "You can now run:" -ForegroundColor Yellow
    Write-Host "  .\run.bat" -ForegroundColor White
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "ERROR: Failed to kill process" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Try running PowerShell as Administrator:" -ForegroundColor Yellow
    Write-Host "  Right-click PowerShell -> Run as Administrator" -ForegroundColor White
    Write-Host "  Then run: .\kill-port-8080.ps1" -ForegroundColor White
    Write-Host ""
    Write-Host "Error details: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
}

Read-Host "Press Enter to exit"
