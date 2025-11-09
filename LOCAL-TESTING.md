# Local Testing Without AWS

## Overview

You can test the **entire application** without an AWS account using **DynamoDB Local** - a local version of DynamoDB that runs in Docker!

## ✅ What You Need

- **Docker Desktop** (that's it!)
- **NO AWS account required**
- **NO AWS CLI required** (optional but helpful)

## 🚀 Quick Start (5 Minutes)

### Step 1: Configure for Local DynamoDB

The default configuration already uses local DynamoDB. Just copy the example:

```powershell
copy .env.example .env
```

The `.env` file is already configured for local testing:
```env
USE_LOCAL_DYNAMODB=true
DYNAMODB_LOCAL_ENDPOINT=http://localhost:8000
AWS_ACCESS_KEY_ID=dummy
AWS_SECRET_ACCESS_KEY=dummy
```

### Step 2: Start All Services

```powershell
# Start DynamoDB Local, Redis, and API
docker-compose up -d

# Check that all services are running
docker-compose ps
```

You should see:
- ✅ `url-shortener-dynamodb` (DynamoDB Local)
- ✅ `url-shortener-redis` (Redis)
- ✅ `url-shortener-api` (API Server)

### Step 3: Create Tables

Choose **ONE** method:

#### Method A: PowerShell Script (If you have AWS CLI)

```powershell
.\setup-local-tables.ps1
```

#### Method B: Batch File (If you have AWS CLI)

```cmd
.\setup-local-tables.bat
```

#### Method C: Docker Method (No AWS CLI needed!)

```cmd
.\setup-local-tables-docker.bat
```

#### Method D: Manual Docker Commands

```powershell
# Create urls table
docker run --rm --network serverless-url-analytics-platform_default `
  amazon/aws-cli dynamodb create-table `
  --table-name urls `
  --attribute-definitions `
    AttributeName=short_code,AttributeType=S `
    AttributeName=created_by,AttributeType=S `
    AttributeName=created_at,AttributeType=N `
  --key-schema AttributeName=short_code,KeyType=HASH `
  --global-secondary-indexes `
    "IndexName=created_by-index,KeySchema=[{AttributeName=created_by,KeyType=HASH},{AttributeName=created_at,KeyType=RANGE}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" `
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 `
  --endpoint-url http://dynamodb-local:8000 `
  --region us-east-1

# Create analytics table
docker run --rm --network serverless-url-analytics-platform_default `
  amazon/aws-cli dynamodb create-table `
  --table-name analytics `
  --attribute-definitions `
    AttributeName=short_code,AttributeType=S `
    AttributeName=sort_key,AttributeType=S `
  --key-schema `
    AttributeName=short_code,KeyType=HASH `
    AttributeName=sort_key,KeyType=RANGE `
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 `
  --endpoint-url http://dynamodb-local:8000 `
  --region us-east-1
```

### Step 4: Test the API

```powershell
# Run automated tests
.\test.bat

# Or test manually
Invoke-RestMethod -Uri "http://localhost:8080/health"
```

## 🎉 That's It!

You now have a fully functional URL shortener with:
- ✅ DynamoDB Local (no AWS account)
- ✅ Redis cache
- ✅ Complete API
- ✅ Analytics tracking
- ✅ QR code generation

## 🔍 Verifying Setup

### Check Services

```powershell
# View running containers
docker-compose ps

# View logs
docker-compose logs api
docker-compose logs dynamodb-local
docker-compose logs redis
```

### Check Tables

If you have AWS CLI:
```powershell
# List tables
aws dynamodb list-tables `
  --endpoint-url http://localhost:8000

# Describe a table
aws dynamodb describe-table `
  --table-name urls `
  --endpoint-url http://localhost:8000
```

### Access DynamoDB Local Admin

DynamoDB Local runs on port 8000:
- http://localhost:8000

## 📊 Testing Everything

### 1. Automated Test Suite

```powershell
.\test.bat
```

This runs 13 comprehensive tests covering all features.

### 2. Manual Testing

```powershell
# Create a short URL
$body = @{
    url = "https://example.com"
    custom_alias = "test123"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/shorten" `
    -Method Post `
    -Body $body `
    -ContentType "application/json"

Write-Host "Short URL: $($response.short_url)"

# Test redirect
Start-Process "http://localhost:8080/test123"

# Get analytics
Invoke-RestMethod -Uri "http://localhost:8080/api/analytics/test123"
```

## 🗄️ Data Persistence

DynamoDB Local stores data in a Docker volume, so your data persists even if you stop the containers:

```powershell
# Stop services (data preserved)
docker-compose down

# Start again (data still there)
docker-compose up -d
```

To completely reset:
```powershell
# Remove all data
docker-compose down -v

# Start fresh
docker-compose up -d

# Recreate tables
.\setup-local-tables.bat
```

## 🐛 Troubleshooting

### Issue: "Connection refused" when creating tables

**Solution**: Make sure DynamoDB Local is running
```powershell
docker ps | findstr dynamodb
```

If not running:
```powershell
docker-compose up -d dynamodb-local
```

### Issue: Tables already exist

This is normal! If you run the setup script twice, it will warn that tables exist. You can ignore this or reset:

```powershell
# Delete and recreate everything
docker-compose down -v
docker-compose up -d
.\setup-local-tables.bat
```

### Issue: API can't connect to DynamoDB

Check logs:
```powershell
docker-compose logs api
```

Look for:
```
Using Local DynamoDB at http://dynamodb-local:8000
```

If you see AWS errors, check `.env`:
```env
USE_LOCAL_DYNAMODB=true
```

### Issue: Port 8000 already in use

Something else is using port 8000. Change it in `docker-compose.yml`:

```yaml
dynamodb-local:
  ports:
    - "8001:8000"  # Use 8001 instead
```

Then update `.env`:
```env
DYNAMODB_LOCAL_ENDPOINT=http://localhost:8001
```

## 🔄 Switching to AWS DynamoDB

When ready to use real AWS DynamoDB:

1. **Update `.env`**:
   ```env
   USE_LOCAL_DYNAMODB=false
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your_real_key
   AWS_SECRET_ACCESS_KEY=your_real_secret
   ```

2. **Create AWS tables** (see SETUP.md)

3. **Restart API**:
   ```powershell
   docker-compose restart api
   ```

## 📖 Understanding DynamoDB Local

### What is it?

DynamoDB Local is an official AWS tool that simulates DynamoDB on your computer. It's perfect for:
- ✅ Development without AWS costs
- ✅ Testing before deploying to AWS
- ✅ Learning DynamoDB
- ✅ Offline development

### Differences from AWS DynamoDB

**Same:**
- ✅ Same API
- ✅ Same data model
- ✅ Same queries
- ✅ Code works exactly the same

**Different:**
- ❌ No automatic scaling
- ❌ No backups
- ❌ Data is local only
- ❌ No IAM permissions
- ❌ No encryption at rest

For local testing, these differences don't matter!

## 🎯 Development Workflow

1. **Develop locally** with DynamoDB Local
2. **Test** with automated tests
3. **Deploy** to AWS when ready
4. **No code changes needed** - just update `.env`!

## 💡 Tips

### View Data

```powershell
# Scan urls table
aws dynamodb scan `
  --table-name urls `
  --endpoint-url http://localhost:8000

# Query specific short code
aws dynamodb get-item `
  --table-name urls `
  --key '{"short_code":{"S":"test123"}}' `
  --endpoint-url http://localhost:8000
```

### Backup Local Data

```powershell
# Export volume
docker run --rm `
  -v serverless-url-analytics-platform_dynamodb-data:/data `
  -v ${PWD}:/backup `
  alpine tar czf /backup/dynamodb-backup.tar.gz /data
```

### Delete Specific Items

```powershell
# Delete a short URL
aws dynamodb delete-item `
  --table-name urls `
  --key '{"short_code":{"S":"test123"}}' `
  --endpoint-url http://localhost:8000
```

## 🚀 Next Steps

- ✅ Test all features locally
- ✅ Develop new features
- ✅ Run the test suite
- ✅ When ready, deploy to AWS (see SETUP.md)

## 📚 Additional Resources

- [DynamoDB Local Documentation](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)
- [AWS CLI DynamoDB Commands](https://docs.aws.amazon.com/cli/latest/reference/dynamodb/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

---

**You're all set! No AWS account needed for testing!** 🎉
