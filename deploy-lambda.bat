@echo off
REM Build and deploy to AWS Lambda

echo ========================================
echo AWS Lambda Deployment Script
echo ========================================
echo.

REM Check if AWS SAM CLI is installed
where sam >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: AWS SAM CLI not found!
    echo.
    echo Please install SAM CLI:
    echo   https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html
    echo.
    pause
    exit /b 1
)

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go not found!
    echo.
    echo For Lambda deployment, you need Go installed.
    echo Download from: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo Step 1: Building Lambda function...
echo.

REM Set environment for Linux build
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM Build the Lambda function
go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    pause
    exit /b 1
)

echo OK Lambda function built successfully
echo.

echo Step 2: Creating deployment package...
echo.

REM Create zip file
powershell -command "Compress-Archive -Force -Path bootstrap -DestinationPath lambda.zip"

echo OK Deployment package created
echo.

echo Step 3: Deploying to AWS...
echo.

REM Deploy with SAM
sam deploy ^
    --template-file template.yaml ^
    --stack-name url-shortener ^
    --capabilities CAPABILITY_IAM ^
    --resolve-s3 ^
    --parameter-overrides ^
        Environment=production

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Deployment successful!
    echo ========================================
    echo.
    echo To get your API URL, run:
    echo   aws cloudformation describe-stacks --stack-name url-shortener --query "Stacks[0].Outputs"
    echo.
) else (
    echo.
    echo ========================================
    echo Deployment failed!
    echo ========================================
    echo.
)

REM Cleanup
del bootstrap

pause
