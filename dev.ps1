# PowerShell script for Windows development
# Provides a menu-driven interface for Docker operations

param(
    [string]$Command = ""
)

function Show-Menu {
    Clear-Host
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "   URL Shortener - Development Menu    " -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "1. Build Docker Image" -ForegroundColor Yellow
    Write-Host "2. Start Services (API + Redis)" -ForegroundColor Yellow
    Write-Host "3. Stop Services" -ForegroundColor Yellow
    Write-Host "4. View Logs" -ForegroundColor Yellow
    Write-Host "5. Restart Services" -ForegroundColor Yellow
    Write-Host "6. Clean (Remove containers and volumes)" -ForegroundColor Yellow
    Write-Host "7. Check Service Status" -ForegroundColor Yellow
    Write-Host "8. Setup AWS DynamoDB Tables" -ForegroundColor Yellow
    Write-Host "9. Test API" -ForegroundColor Yellow
    Write-Host "0. Exit" -ForegroundColor Red
    Write-Host ""
}

function Build-Image {
    Write-Host "Building Docker image..." -ForegroundColor Green
    docker build -t url-shortener:latest .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
    }
    Pause
}

function Start-Services {
    Write-Host "Starting services..." -ForegroundColor Green

    # Check if .env exists
    if (-not (Test-Path .env)) {
        Write-Host "WARNING: .env file not found!" -ForegroundColor Red
        Write-Host "Creating from .env.example..." -ForegroundColor Yellow
        Copy-Item .env.example .env
        Write-Host ""
        Write-Host "Please edit .env file with your AWS credentials" -ForegroundColor Yellow
        Write-Host "Then run this script again." -ForegroundColor Yellow
        Pause
        return
    }

    docker-compose up -d
    Write-Host ""
    Write-Host "Services started!" -ForegroundColor Green
    Write-Host "API available at: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "Redis available at: localhost:6379" -ForegroundColor Cyan
    Pause
}

function Stop-Services {
    Write-Host "Stopping services..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "Services stopped!" -ForegroundColor Green
    Pause
}

function View-Logs {
    Write-Host "Viewing logs (Press Ctrl+C to exit)..." -ForegroundColor Green
    docker-compose logs -f
}

function Restart-Services {
    Write-Host "Restarting services..." -ForegroundColor Yellow
    docker-compose restart
    Write-Host "Services restarted!" -ForegroundColor Green
    Pause
}

function Clean-All {
    Write-Host "WARNING: This will remove all containers, volumes, and data!" -ForegroundColor Red
    $confirmation = Read-Host "Are you sure? (yes/no)"
    if ($confirmation -eq "yes") {
        docker-compose down -v
        docker rmi url-shortener:latest -f
        Write-Host "Cleanup completed!" -ForegroundColor Green
    } else {
        Write-Host "Cleanup cancelled." -ForegroundColor Yellow
    }
    Pause
}

function Check-Status {
    Write-Host "Service Status:" -ForegroundColor Cyan
    Write-Host ""
    docker-compose ps
    Write-Host ""
    Pause
}

function Setup-DynamoDB {
    Write-Host "Setting up DynamoDB tables..." -ForegroundColor Green
    Write-Host ""
    Write-Host "Make sure you have AWS CLI installed and configured!" -ForegroundColor Yellow
    Write-Host ""

    $region = Read-Host "Enter AWS Region (default: us-east-1)"
    if ([string]::IsNullOrWhiteSpace($region)) { $region = "us-east-1" }

    Write-Host ""
    Write-Host "Creating 'urls' table..." -ForegroundColor Cyan

    $urlsTable = @"
aws dynamodb create-table \
  --table-name urls \
  --attribute-definitions \
    AttributeName=short_code,AttributeType=S \
    AttributeName=created_by,AttributeType=S \
    AttributeName=created_at,AttributeType=N \
  --key-schema \
    AttributeName=short_code,KeyType=HASH \
  --global-secondary-indexes \
    IndexName=created_by-index,KeySchema=[{AttributeName=created_by,KeyType=HASH},{AttributeName=created_at,KeyType=RANGE}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5} \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --region $region
"@

    Write-Host $urlsTable
    Write-Host ""

    Write-Host "Creating 'analytics' table..." -ForegroundColor Cyan

    $analyticsTable = @"
aws dynamodb create-table \
  --table-name analytics \
  --attribute-definitions \
    AttributeName=short_code,AttributeType=S \
    AttributeName=sort_key,AttributeType=S \
  --key-schema \
    AttributeName=short_code,KeyType=HASH \
    AttributeName=sort_key,KeyType=RANGE \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --region $region
"@

    Write-Host $analyticsTable
    Write-Host ""
    Write-Host "Copy the above commands and run them in your terminal." -ForegroundColor Yellow
    Pause
}

function Test-API {
    Write-Host "Testing API endpoints..." -ForegroundColor Green
    Write-Host ""

    # Test health endpoint
    Write-Host "1. Testing health endpoint..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
        Write-Host "   Status: OK" -ForegroundColor Green
        Write-Host "   Response: $($response | ConvertTo-Json)" -ForegroundColor Gray
    } catch {
        Write-Host "   Status: FAILED" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }

    Write-Host ""
    Write-Host "2. Testing shorten endpoint..." -ForegroundColor Cyan
    try {
        $body = @{
            url = "https://example.com/test"
            custom_alias = "test123"
        } | ConvertTo-Json

        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/shorten" -Method Post -Body $body -ContentType "application/json"
        Write-Host "   Status: OK" -ForegroundColor Green
        Write-Host "   Short URL: $($response.short_url)" -ForegroundColor Gray
    } catch {
        Write-Host "   Status: FAILED" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }

    Write-Host ""
    Pause
}

# Main script logic
if ($Command -ne "") {
    switch ($Command.ToLower()) {
        "build" { Build-Image }
        "start" { Start-Services }
        "stop" { Stop-Services }
        "logs" { View-Logs }
        "restart" { Restart-Services }
        "clean" { Clean-All }
        "status" { Check-Status }
        "test" { Test-API }
        default {
            Write-Host "Unknown command: $Command" -ForegroundColor Red
            Write-Host "Available commands: build, start, stop, logs, restart, clean, status, test" -ForegroundColor Yellow
        }
    }
} else {
    # Interactive menu
    do {
        Show-Menu
        $choice = Read-Host "Select an option"

        switch ($choice) {
            "1" { Build-Image }
            "2" { Start-Services }
            "3" { Stop-Services }
            "4" { View-Logs }
            "5" { Restart-Services }
            "6" { Clean-All }
            "7" { Check-Status }
            "8" { Setup-DynamoDB }
            "9" { Test-API }
            "0" { Write-Host "Goodbye!" -ForegroundColor Cyan; exit }
            default {
                Write-Host "Invalid option!" -ForegroundColor Red
                Pause
            }
        }
    } while ($true)
}
