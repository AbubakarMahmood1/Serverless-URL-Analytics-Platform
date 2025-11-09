@echo off
REM Batch script to create DynamoDB Local tables

echo ========================================
echo Setting up DynamoDB Local Tables
echo ========================================
echo.

REM Check if AWS CLI is installed
where aws >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: AWS CLI not found!
    echo.
    echo Please install AWS CLI from:
    echo   https://aws.amazon.com/cli/
    echo.
    echo Or use Docker method ^(see setup-local-tables-docker.bat^)
    pause
    exit /b 1
)

set ENDPOINT=http://localhost:8000

echo Creating 'urls' table...
aws dynamodb create-table ^
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
  --endpoint-url %ENDPOINT% >nul 2>&1

if %ERRORLEVEL% EQU 0 (
    echo OK 'urls' table created successfully
) else (
    echo WARN Failed to create 'urls' table ^(may already exist^)
)

echo.
echo Creating 'analytics' table...
aws dynamodb create-table ^
  --table-name analytics ^
  --attribute-definitions ^
    AttributeName=short_code,AttributeType=S ^
    AttributeName=sort_key,AttributeType=S ^
  --key-schema ^
    AttributeName=short_code,KeyType=HASH ^
    AttributeName=sort_key,KeyType=RANGE ^
  --provisioned-throughput ^
    ReadCapacityUnits=5,WriteCapacityUnits=5 ^
  --endpoint-url %ENDPOINT% >nul 2>&1

if %ERRORLEVEL% EQU 0 (
    echo OK 'analytics' table created successfully
) else (
    echo WARN Failed to create 'analytics' table ^(may already exist^)
)

echo.
echo ========================================
echo Setup complete!
echo ========================================
echo.
echo Tables have been created in DynamoDB Local
echo You can now start testing the API
echo.

pause
