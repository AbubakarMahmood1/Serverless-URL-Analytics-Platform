# AWS Lambda Deployment Guide

Deploy your URL Shortener as a **serverless application** on AWS Lambda with near-zero cost for low traffic!

## 🎯 Why AWS Lambda?

### Cost Comparison

| Traffic/Month | Traditional Server | AWS Lambda |
|---------------|-------------------|------------|
| 1,000 requests | $5 (always on) | **$0** (free tier) |
| 10,000 requests | $5 | **$0** (free tier) |
| 100,000 requests | $5-10 | **$0** (free tier) |
| 1,000,000 requests | $10-20 | **$0** (free tier!) |
| 10,000,000 requests | $50+ | **~$2** |

**AWS Lambda Free Tier (Permanent):**
- ✅ 1 million requests per month
- ✅ 400,000 GB-seconds of compute time
- ✅ Never expires!

### Benefits

- 💰 **Pay-per-use**: Only charged when someone uses your app
- 📈 **Auto-scaling**: Handles 1 to millions of requests automatically
- 🔒 **Secure**: AWS manages security patches
- 🌍 **Global**: Deploy in multiple regions easily
- ⚡ **Fast**: Single-digit millisecond response (after warmup)
- 🛠️ **Zero maintenance**: No servers to manage

### Trade-offs

- ❄️ **Cold starts**: First request after idle can be slower (~1-2 seconds)
- 🕐 **Timeout**: Max 30 seconds per request (fine for our app)
- 💾 **Stateless**: Can't store files on disk (we use DynamoDB + Redis)
- 📚 **Learning curve**: AWS concepts to understand

---

## 📋 Prerequisites

### Required Tools

1. **Go** (for building)
   - Download: https://golang.org/dl/
   - Version: 1.21+

2. **AWS CLI**
   - Download: https://aws.amazon.com/cli/
   - Configure: `aws configure`

3. **AWS SAM CLI**
   - Download: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html
   - Test: `sam --version`

### AWS Account Setup

1. **Create AWS Account** (if you don't have one)
   - Go to: https://aws.amazon.com/free/
   - Free tier includes Lambda!

2. **Create IAM User** (for deployment)
   ```bash
   # In AWS Console:
   # 1. Go to IAM
   # 2. Create User
   # 3. Attach policies:
   #    - AWSLambdaFullAccess
   #    - AmazonDynamoDBFullAccess
   #    - AWSCloudFormationFullAccess
   #    - IAMFullAccess
   #    - AmazonAPIGatewayAdministrator
   ```

3. **Configure AWS CLI**
   ```powershell
   aws configure
   # AWS Access Key ID: (your key)
   # AWS Secret Access Key: (your secret)
   # Default region: us-east-1
   # Default output format: json
   ```

### Optional: Redis

Lambda works best with **serverless Redis**:
- **Upstash** (Recommended): https://upstash.com/ - FREE tier available
- **ElastiCache Serverless**: AWS managed, pay-per-use

---

## 🚀 Deployment Steps

### Method 1: PowerShell Script (Easiest)

```powershell
# Deploy to AWS
.\deploy-lambda.ps1

# With custom Redis
.\deploy-lambda.ps1 `
    -RedisEndpoint "your-upstash-endpoint:6379" `
    -RedisPassword "your-redis-password"
```

### Method 2: Batch Script

```cmd
.\deploy-lambda.bat
```

### Method 3: Manual SAM Commands

```powershell
# 1. Build the Lambda function
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -tags lambda.norpc -o bootstrap cmd/lambda/main.go

# 2. Deploy with SAM
sam deploy `
    --template-file template.yaml `
    --stack-name url-shortener `
    --capabilities CAPABILITY_IAM `
    --resolve-s3 `
    --parameter-overrides Environment=production

# 3. Get API URL
aws cloudformation describe-stacks `
    --stack-name url-shortener `
    --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" `
    --output text
```

---

## 🔧 Configuration

### Environment Variables

Edit `template.yaml` to customize:

```yaml
Environment:
  Variables:
    SHORT_CODE_LENGTH: "6"
    RATE_LIMIT_REQUESTS: "100"
    RATE_LIMIT_WINDOW: "60"
```

### Redis Setup

#### Option A: Upstash (Recommended - FREE Tier)

1. Go to https://upstash.com/
2. Create account
3. Create Redis database
4. Copy endpoint and password
5. Deploy:
   ```powershell
   .\deploy-lambda.ps1 `
       -RedisEndpoint "your-endpoint.upstash.io:6379" `
       -RedisPassword "your-password"
   ```

#### Option B: Without Redis (Budget Option)

Works but slower (no caching):
1. Comment out Redis parts in code
2. Deploy normally

#### Option C: AWS ElastiCache Serverless

1. Create ElastiCache Serverless cluster
2. Note endpoint
3. Deploy with endpoint

### Custom Domain

After deployment, add custom domain:

```powershell
# 1. In AWS Console: API Gateway
# 2. Custom domain names → Create
# 3. Enter your domain (e.g., short.yourdomain.com)
# 4. Create API mapping
# 5. Add DNS record (CNAME or A record)
```

---

## 📊 Monitoring & Logs

### View Logs

```powershell
# Real-time logs
sam logs --stack-name url-shortener --tail

# Specific time range
sam logs --stack-name url-shortener `
    --start-time '10min ago' `
    --end-time '5min ago'
```

### CloudWatch Dashboard

1. Go to AWS CloudWatch
2. View Lambda metrics:
   - Invocations
   - Duration
   - Errors
   - Throttles

### X-Ray Tracing (Optional)

Enable in `template.yaml`:
```yaml
Globals:
  Function:
    Tracing: Active
```

---

## 🧪 Testing

### Test Deployed API

```powershell
# Get API URL
$API_URL = aws cloudformation describe-stacks `
    --stack-name url-shortener `
    --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" `
    --output text

# Health check
Invoke-RestMethod -Uri "$API_URL/health"

# Create short URL
$body = @{
    url = "https://example.com"
    custom_alias = "test"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "$API_URL/api/shorten" `
    -Method Post `
    -Body $body `
    -ContentType "application/json"

# Test redirect
Start-Process "$API_URL/test"
```

### Local Testing with SAM

```powershell
# Start local API (simulates Lambda)
sam local start-api

# Test locally
Invoke-RestMethod -Uri "http://localhost:3000/health"
```

---

## 💰 Cost Optimization

### DynamoDB

**On-Demand Pricing** (default in template):
- Pay per request
- No minimum cost
- Best for variable traffic

**Provisioned Capacity** (if steady traffic):
```yaml
BillingMode: PROVISIONED
ProvisionedThroughput:
  ReadCapacityUnits: 5
  WriteCapacityUnits: 5
```

### Lambda

**Memory Optimization**:
- Default: 512MB (good for most cases)
- Lower traffic: Try 256MB (cheaper)
- High traffic: Try 1024MB (faster, may be cheaper overall)

Edit in `template.yaml`:
```yaml
Globals:
  Function:
    MemorySize: 256  # or 512, 1024
```

### API Gateway

**HTTP API** (what we use):
- ~70% cheaper than REST API
- $1 per million requests

**Caching** (optional):
- Reduces Lambda invocations
- Faster responses
- ~$0.02/hour per GB

---

## 🔄 Updates & Rollbacks

### Update Lambda Function

```powershell
# Make code changes, then:
.\deploy-lambda.ps1
```

### Rollback

```powershell
# List previous versions
aws lambda list-versions-by-function `
    --function-name url-shortener

# Rollback to version
aws lambda update-alias `
    --function-name url-shortener `
    --name production `
    --function-version 2  # previous version
```

### Blue/Green Deployment

Built into SAM:
```yaml
AutoPublishAlias: live
DeploymentPreference:
  Type: Canary10Percent10Minutes
```

---

## 🗑️ Cleanup

### Delete Everything

```powershell
# Delete stack (removes all resources)
sam delete --stack-name url-shortener

# Confirm deletion
# This removes:
# - Lambda function
# - API Gateway
# - DynamoDB tables (⚠️ deletes all data!)
# - CloudFormation stack
```

### Delete But Keep Data

```yaml
# In template.yaml, add to tables:
DeletionPolicy: Retain
```

---

## 🐛 Troubleshooting

### Issue: "AccessDeniedException"

**Solution:** Check IAM permissions
```powershell
# Ensure your IAM user has necessary permissions
aws iam get-user
```

### Issue: "Cold start too slow"

**Solutions:**
1. Enable **Provisioned Concurrency** (keeps function warm):
   ```yaml
   ProvisionedConcurrencyConfig:
     ProvisionedConcurrentExecutions: 1
   ```
   Cost: ~$15/month, eliminates cold starts

2. Use **Lambda SnapStart** (Java only, not for Go yet)

3. Optimize code:
   - Reduce dependencies
   - Initialize outside handler
   - Use connection pooling

### Issue: "Function timeout"

**Solution:** Increase timeout in template.yaml:
```yaml
Globals:
  Function:
    Timeout: 30  # Max for API Gateway
```

### Issue: "Redis connection failed"

**Solutions:**
1. Check Redis endpoint in parameters
2. Verify Redis is accessible from Lambda (VPC settings if using ElastiCache)
3. Check Redis password

### Issue: "DynamoDB throttling"

**Solution:** Switch to provisioned capacity or increase on-demand limits

---

## 📈 Scaling

### Automatic Scaling

Lambda scales automatically:
- 0 to 1,000 concurrent executions by default
- Can request increase to 100,000+
- Scales in seconds

### Regional Failover

Deploy to multiple regions:

```powershell
# Deploy to us-east-1
sam deploy --region us-east-1

# Deploy to eu-west-1
sam deploy --region eu-west-1 --stack-name url-shortener-eu

# Use Route 53 for global routing
```

---

## 📊 Cost Calculator

**Example Monthly Costs:**

```
Scenario: Personal URL Shortener
- 10,000 requests/month
- 500ms average duration
- 512MB memory

Lambda: FREE (within free tier)
DynamoDB: $0 - $1 (25GB free tier)
API Gateway: $0.01
Data Transfer: $0 - $1
Total: ~$0 - $2/month
```

```
Scenario: Small Business
- 500,000 requests/month
- 400ms average duration
- 512MB memory

Lambda: FREE (within free tier)
DynamoDB: $2 - $5
API Gateway: $0.50
Data Transfer: $1 - $2
Redis (Upstash): FREE (free tier)
Total: ~$3 - $8/month
```

```
Scenario: High Traffic
- 10,000,000 requests/month
- 300ms average duration
- 512MB memory

Lambda: $2 - $3
DynamoDB: $10 - $20
API Gateway: $10
Data Transfer: $5 - $10
Redis (Upstash): $10
Total: ~$37 - $53/month
```

Still way cheaper than running servers 24/7!

---

## 🎓 Best Practices

1. **Use Environment Variables** for configuration
2. **Enable CloudWatch Logs** for debugging
3. **Set up Alarms** for errors and throttling
4. **Use API keys** for rate limiting (in API Gateway)
5. **Enable X-Ray** for tracing in production
6. **Tag Resources** for cost tracking
7. **Use Secrets Manager** for sensitive data (Redis passwords, etc.)
8. **Set up CI/CD** with GitHub Actions or AWS CodePipeline

---

## 🔗 Useful Links

- [AWS Lambda Pricing](https://aws.amazon.com/lambda/pricing/)
- [DynamoDB Pricing](https://aws.amazon.com/dynamodb/pricing/)
- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [Upstash Redis](https://upstash.com/)
- [AWS Free Tier](https://aws.amazon.com/free/)

---

## 🆘 Need Help?

1. Check CloudWatch Logs: `sam logs --tail`
2. View stack status: `aws cloudformation describe-stacks --stack-name url-shortener`
3. Test locally: `sam local start-api`
4. AWS Support: https://console.aws.amazon.com/support/

---

**Enjoy your serverless URL shortener! 🚀**

**Estimated cost for most users: $0 - $5/month**
