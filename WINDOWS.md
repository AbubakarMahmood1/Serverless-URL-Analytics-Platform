# Windows Setup Guide

## Prerequisites

- **Docker Desktop** for Windows
- **Git** (optional, for cloning)
- **AWS Account** with DynamoDB access

> **Note**: You do NOT need Go or Make installed! Everything runs in Docker.

## Quick Start

### 1. Get the Code

Either clone the repository or download and extract the ZIP file:

```powershell
git clone https://github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform.git
cd Serverless-URL-Analytics-Platform
```

### 2. Configure Environment

Copy `.env.example` to `.env`:

```powershell
copy .env.example .env
```

Edit `.env` with your AWS credentials using Notepad or any text editor:

```powershell
notepad .env
```

Required settings:
```env
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key_here
AWS_SECRET_ACCESS_KEY=your_secret_key_here
DYNAMODB_URLS_TABLE=urls
DYNAMODB_ANALYTICS_TABLE=analytics
```

### 3. Choose Your Method

#### Option A: PowerShell Script (Recommended)

**Interactive Menu:**
```powershell
.\dev.ps1
```

This provides a menu with options to:
1. Build Docker image
2. Start services
3. Stop services
4. View logs
5. Test API
6. And more...

**Direct Commands:**
```powershell
# Build the Docker image
.\dev.ps1 build

# Start services
.\dev.ps1 start

# View logs
.\dev.ps1 logs

# Stop services
.\dev.ps1 stop

# Test API
.\dev.ps1 test
```

#### Option B: Batch Files (Simple)

Double-click these files or run from Command Prompt:

1. **Build the image:**
   ```cmd
   build.bat
   ```

2. **Start the services:**
   ```cmd
   run.bat
   ```

3. **Stop the services:**
   ```cmd
   stop.bat
   ```

#### Option C: Docker Compose (Manual)

```powershell
# Start services in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Setting Up AWS DynamoDB

### Using AWS Console (Easiest)

1. Go to [AWS DynamoDB Console](https://console.aws.amazon.com/dynamodb)
2. Click "Create table"

**Table 1: urls**
- Table name: `urls`
- Partition key: `short_code` (String)
- Click "Add global secondary index"
  - Index name: `created_by-index`
  - Partition key: `created_by` (String)
  - Sort key: `created_at` (Number)
- Click "Create table"

**Table 2: analytics**
- Table name: `analytics`
- Partition key: `short_code` (String)
- Sort key: `sort_key` (String)
- Click "Create table"

### Using AWS CLI (If Installed)

```powershell
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
  --region us-east-1

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
  --region us-east-1
```

## Testing the API

### Using PowerShell

```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8080/health"

# Shorten a URL
$body = @{
    url = "https://example.com/very/long/url"
    custom_alias = "mylink"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/shorten" `
    -Method Post `
    -Body $body `
    -ContentType "application/json"

# Get analytics
Invoke-RestMethod -Uri "http://localhost:8080/api/analytics/mylink"
```

### Using curl (if installed)

```powershell
# Health check
curl http://localhost:8080/health

# Shorten a URL
curl -X POST http://localhost:8080/api/shorten `
  -H "Content-Type: application/json" `
  -d '{\"url\":\"https://example.com\",\"custom_alias\":\"test\"}'
```

### Using Browser

- Health: http://localhost:8080/health
- Redirect test: http://localhost:8080/test (after creating a short URL)
- QR Code: http://localhost:8080/api/qr/test

## Common Commands

### View Running Containers
```powershell
docker ps
```

### View Container Logs
```powershell
# API logs
docker logs url-shortener-api -f

# Redis logs
docker logs url-shortener-redis -f

# All logs with docker-compose
docker-compose logs -f
```

### Restart a Service
```powershell
docker-compose restart api
# or
docker-compose restart redis
```

### Check Service Health
```powershell
docker-compose ps
```

### Connect to Redis CLI
```powershell
docker exec -it url-shortener-redis redis-cli
```

Then try:
```redis
PING
GET url:test123
KEYS *
```

## Troubleshooting

### Port Already in Use

If port 8080 or 6379 is already in use:

1. Edit `docker-compose.yml`
2. Change the port mapping:
   ```yaml
   ports:
     - "9090:8080"  # Use port 9090 instead
   ```

### Docker Not Starting

1. Make sure Docker Desktop is running
2. Check Docker settings → Resources → WSL Integration (if using WSL)
3. Restart Docker Desktop

### Cannot Connect to DynamoDB

1. Verify AWS credentials in `.env`
2. Check IAM permissions
3. Verify tables exist in AWS Console
4. Check the region matches in `.env`

### Redis Connection Failed

1. Make sure Redis container is running:
   ```powershell
   docker ps | findstr redis
   ```

2. Restart Redis:
   ```powershell
   docker-compose restart redis
   ```

### API Returns 500 Error

Check the logs:
```powershell
docker logs url-shortener-api
```

Common issues:
- DynamoDB tables not created
- AWS credentials invalid
- Redis not connected

## Stopping Everything

```powershell
# Stop services but keep data
docker-compose down

# Stop and remove all data
docker-compose down -v
```

## Updating the Code

After pulling new changes:

```powershell
# Rebuild the image
.\dev.ps1 build
# or
docker-compose build

# Restart services
.\dev.ps1 restart
# or
docker-compose up -d
```

## Development Workflow

1. **Start services:**
   ```powershell
   .\dev.ps1 start
   ```

2. **Make code changes** in your editor

3. **Rebuild and restart:**
   ```powershell
   .\dev.ps1 build
   .\dev.ps1 restart
   ```

4. **View logs:**
   ```powershell
   .\dev.ps1 logs
   ```

5. **Test:**
   ```powershell
   .\dev.ps1 test
   ```

## PowerShell Execution Policy

If you get an error running `dev.ps1`, you may need to allow script execution:

```powershell
# Check current policy
Get-ExecutionPolicy

# Allow scripts (run as Administrator)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Next Steps

- Test the API endpoints
- Deploy to Railway or Fly.io (see SETUP.md)
- Build the Next.js frontend
- Add custom domain

## Support

- Check [SETUP.md](SETUP.md) for deployment instructions
- Check [CLAUDE.md](CLAUDE.md) for architecture details
- Open an issue on GitHub for bugs
