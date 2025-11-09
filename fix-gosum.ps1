# PowerShell script to generate go.sum using Docker

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Generating go.sum file" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "This will download dependencies and create go.sum..." -ForegroundColor Yellow
Write-Host ""

# Use Go Docker image to run go mod tidy
docker run --rm -v "${PWD}:/app" -w /app golang:1.23-alpine go mod tidy

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Success! go.sum file created" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "You can now run:" -ForegroundColor White
    Write-Host "  docker-compose build" -ForegroundColor Cyan
    Write-Host "  docker-compose up -d" -ForegroundColor Cyan
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "Failed to generate go.sum" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
}
