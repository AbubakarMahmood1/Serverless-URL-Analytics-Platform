# PowerShell script to create DynamoDB Local tables
# Run this after starting docker-compose

$endpoint = "http://localhost:8000"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Setting up DynamoDB Local Tables" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if AWS CLI is installed
if (-not (Get-Command aws -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: AWS CLI not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install AWS CLI:" -ForegroundColor Yellow
    Write-Host "  https://aws.amazon.com/cli/" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Or use the Docker method (see below)" -ForegroundColor Yellow
    exit 1
}

Write-Host "Creating 'urls' table..." -ForegroundColor Green

# Create URLs table
aws dynamodb create-table `
  --table-name urls `
  --attribute-definitions `
    AttributeName=short_code,AttributeType=S `
    AttributeName=created_by,AttributeType=S `
    AttributeName=created_at,AttributeType=N `
  --key-schema `
    AttributeName=short_code,KeyType=HASH `
  --global-secondary-indexes `
    "IndexName=created_by-index,KeySchema=[{AttributeName=created_by,KeyType=HASH},{AttributeName=created_at,KeyType=RANGE}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" `
  --provisioned-throughput `
    ReadCapacityUnits=5,WriteCapacityUnits=5 `
  --endpoint-url $endpoint `
  2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ 'urls' table created successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create 'urls' table (it may already exist)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Creating 'analytics' table..." -ForegroundColor Green

# Create Analytics table
aws dynamodb create-table `
  --table-name analytics `
  --attribute-definitions `
    AttributeName=short_code,AttributeType=S `
    AttributeName=sort_key,AttributeType=S `
  --key-schema `
    AttributeName=short_code,KeyType=HASH `
    AttributeName=sort_key,KeyType=RANGE `
  --provisioned-throughput `
    ReadCapacityUnits=5,WriteCapacityUnits=5 `
  --endpoint-url $endpoint `
  2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ 'analytics' table created successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create 'analytics' table (it may already exist)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Verifying tables..." -ForegroundColor Cyan

# List tables
$tables = aws dynamodb list-tables --endpoint-url $endpoint --output json | ConvertFrom-Json

if ($tables.TableNames -contains "urls" -and $tables.TableNames -contains "analytics") {
    Write-Host "✓ All tables verified!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Tables created:" -ForegroundColor White
    foreach ($table in $tables.TableNames) {
        Write-Host "  - $table" -ForegroundColor Gray
    }
} else {
    Write-Host "✗ Some tables are missing" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Setup complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "You can now start the API server." -ForegroundColor White
Write-Host ""
