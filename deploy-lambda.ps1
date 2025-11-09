# PowerShell script for AWS Lambda deployment

param(
    [string]$Environment = "production",
    [string]$RedisEndpoint = "",
    [string]$RedisPassword = ""
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   AWS Lambda Deployment" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

# Check SAM CLI
if (-not (Get-Command sam -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: AWS SAM CLI not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install SAM CLI:" -ForegroundColor Yellow
    Write-Host "  https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html" -ForegroundColor Yellow
    exit 1
}

# Check Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "For Lambda deployment, you need Go installed." -ForegroundColor Yellow
    Write-Host "Download from: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Check AWS CLI
if (-not (Get-Command aws -ErrorAction SilentlyContinue)) {
    Write-Host "WARNING: AWS CLI not found!" -ForegroundColor Yellow
    Write-Host "You may need it for post-deployment commands." -ForegroundColor Yellow
}

Write-Host "✓ All prerequisites found" -ForegroundColor Green
Write-Host ""

# Build
Write-Host "Step 1: Building Lambda function..." -ForegroundColor Cyan
Write-Host ""

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Lambda function built successfully" -ForegroundColor Green
Write-Host ""

# Package
Write-Host "Step 2: Creating deployment package..." -ForegroundColor Cyan
Write-Host ""

Compress-Archive -Force -Path bootstrap -DestinationPath lambda.zip

Write-Host "✓ Deployment package created" -ForegroundColor Green
Write-Host ""

# Deploy
Write-Host "Step 3: Deploying to AWS..." -ForegroundColor Cyan
Write-Host "Environment: $Environment" -ForegroundColor Gray
Write-Host ""

$params = @(
    "--template-file", "template.yaml",
    "--stack-name", "url-shortener",
    "--capabilities", "CAPABILITY_IAM",
    "--resolve-s3",
    "--parameter-overrides",
    "Environment=$Environment"
)

if ($RedisEndpoint) {
    $params += "RedisEndpoint=$RedisEndpoint"
}

if ($RedisPassword) {
    $params += "RedisPassword=$RedisPassword"
}

& sam deploy @params

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "   Deployment Successful!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""

    # Get outputs
    Write-Host "Getting API URL..." -ForegroundColor Cyan
    $outputs = aws cloudformation describe-stacks `
        --stack-name url-shortener `
        --query "Stacks[0].Outputs" `
        --output json | ConvertFrom-Json

    foreach ($output in $outputs) {
        if ($output.OutputKey -eq "ApiUrl") {
            Write-Host ""
            Write-Host "Your API is now live at:" -ForegroundColor Green
            Write-Host "  $($output.OutputValue)" -ForegroundColor Cyan
            Write-Host ""
            Write-Host "Test it:" -ForegroundColor Yellow
            Write-Host "  curl $($output.OutputValue)/health" -ForegroundColor Gray
            Write-Host ""
        }
    }
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "   Deployment Failed!" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
}

# Cleanup
Remove-Item bootstrap -ErrorAction SilentlyContinue

Write-Host "Deployment process complete!" -ForegroundColor Cyan
