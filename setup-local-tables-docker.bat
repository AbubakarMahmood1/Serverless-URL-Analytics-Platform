@echo off
REM Create DynamoDB Local tables using Docker (no AWS CLI needed!)

echo ========================================
echo Setting up DynamoDB Local Tables
echo Using Docker ^(No AWS CLI needed^)
echo ========================================
echo.

REM Check if DynamoDB Local container is running
docker ps | findstr url-shortener-dynamodb >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: DynamoDB Local container is not running!
    echo Please start it first with: docker-compose up -d
    pause
    exit /b 1
)

echo Creating 'urls' table...
docker run --rm --network url-shortener-analytics-platform_default ^
  amazon/aws-cli dynamodb create-table ^
  --table-name urls ^
  --attribute-definitions ^
    AttributeName=short_code,AttributeType=S ^
    AttributeName=created_by,AttributeType=S ^
    AttributeName=created_at,AttributeType=N ^
  --key-schema ^
    AttributeName=short_code,KeyType=HASH ^
  --global-secondary-indexes ^
    "IndexName=created_by-index,KeySchema=[{AttributeName=created_by,KeyType=HASH},{AttributeName=created_at,KeyType=RANGE}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" ^
  --provisioned-throughput ^
    ReadCapacityUnits=5,WriteCapacityUnits=5 ^
  --endpoint-url http://dynamodb-local:8000 ^
  --region us-east-1 >nul 2>&1

if %ERRORLEVEL% EQU 0 (
    echo OK 'urls' table created
) else (
    echo WARN Table may already exist
)

echo.
echo Creating 'analytics' table...
docker run --rm --network url-shortener-analytics-platform_default ^
  amazon/aws-cli dynamodb create-table ^
  --table-name analytics ^
  --attribute-definitions ^
    AttributeName=short_code,AttributeType=S ^
    AttributeName=sort_key,AttributeType=S ^
  --key-schema ^
    AttributeName=short_code,KeyType=HASH ^
    AttributeName=sort_key,KeyType=RANGE ^
  --provisioned-throughput ^
    ReadCapacityUnits=5,WriteCapacityUnits=5 ^
  --endpoint-url http://dynamodb-local:8000 ^
  --region us-east-1 >nul 2>&1

if %ERRORLEVEL% EQU 0 (
    echo OK 'analytics' table created
) else (
    echo WARN Table may already exist
)

echo.
echo ========================================
echo Setup complete!
echo ========================================
echo.

pause
