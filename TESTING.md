# Testing Guide

## Overview

This project includes comprehensive testing at multiple levels:
- **Unit Tests**: Test individual functions and components
- **API Tests**: Test HTTP endpoints end-to-end
- **Integration Tests**: Test the complete system

## Windows Users (No Go Required!)

### Quick Test

Run all API tests with one command:

```powershell
# Option 1: Batch file
.\test.bat

# Option 2: PowerShell directly
.\test-api.ps1

# Option 3: Via dev menu
.\dev.ps1
# Then select option 9 (Test API)
```

### What Gets Tested

The automated test suite checks:

1. ✓ **Health Check** - Is the service running?
2. ✓ **URL Shortening** - Can create short URLs (random codes)
3. ✓ **Custom Aliases** - Can create custom short codes
4. ✓ **URL Validation** - Rejects invalid URLs
5. ✓ **Duplicate Prevention** - Prevents duplicate aliases
6. ✓ **URL Redirect** - Short URLs redirect correctly
7. ✓ **404 Handling** - Returns 404 for non-existent codes
8. ✓ **Link Information** - Can retrieve link details
9. ✓ **Analytics** - Records and retrieves click data
10. ✓ **QR Code Generation** - Generates QR codes
11. ✓ **List Links** - Can list all user links
12. ✓ **Rate Limiting** - Prevents abuse
13. ✓ **CORS Headers** - Proper cross-origin support

### Understanding Test Output

**Green ✓** = Test Passed
**Red ✗** = Test Failed

Example:
```
✓ PASS: Health Check
  Service is healthy
  Response: {"status":"healthy","service":"url-shortener"}

✗ FAIL: URL Redirect
  Request failed: Connection refused
```

### Test Results File

Results are saved to `test-results.csv` for detailed analysis:
```csv
Test,Status,Message,Data
Health Check,PASS,Service is healthy,"{""status"":""healthy""}"
Shorten URL (Random),PASS,Created short URL: http://localhost:8080/abc123,"{""short_url"":""...""}"
```

## Running Unit Tests (Docker)

Unit tests check individual functions without needing a running server:

```powershell
# Windows
.\run-tests.bat

# Or manually
docker build -f Dockerfile.test -t url-shortener-tests .
docker run --rm url-shortener-tests
```

### Unit Test Coverage

- **Shortener**
  - Code generation
  - Uniqueness
  - Custom alias validation

- **Validator**
  - URL validation
  - IP hashing
  - URL normalization

## Manual Testing

### Using PowerShell

```powershell
# 1. Health Check
Invoke-RestMethod -Uri "http://localhost:8080/health"

# 2. Create Short URL
$body = @{
    url = "https://example.com"
    custom_alias = "test"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/shorten" `
    -Method Post `
    -Body $body `
    -ContentType "application/json"

# 3. Test Redirect
Start-Process "http://localhost:8080/test"

# 4. Get Analytics
Invoke-RestMethod -Uri "http://localhost:8080/api/analytics/test"

# 5. Generate QR Code
Invoke-WebRequest -Uri "http://localhost:8080/api/qr/test" `
    -OutFile "qr-code.png"
```

### Using curl (if installed)

```bash
# Health Check
curl http://localhost:8080/health

# Create Short URL
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://example.com\",\"custom_alias\":\"test\"}"

# Get Analytics
curl http://localhost:8080/api/analytics/test

# Download QR Code
curl http://localhost:8080/api/qr/test?size=512 -o qr.png
```

### Using Browser

Visit these URLs in your browser:
- Health: http://localhost:8080/health
- QR Code: http://localhost:8080/api/qr/test
- Redirect: http://localhost:8080/test (after creating)

## Troubleshooting Failed Tests

### "Services not running"

**Solution:**
```powershell
# Start services
.\run.bat
# Or
docker-compose up -d

# Check status
docker-compose ps
```

### "Connection refused"

**Solution:**
```powershell
# Check if API is listening
docker logs url-shortener-api

# Restart services
docker-compose restart
```

### "DynamoDB errors"

**Solution:**
1. Verify AWS credentials in `.env`
2. Check DynamoDB tables exist:
   ```powershell
   aws dynamodb list-tables --region us-east-1
   ```
3. Verify IAM permissions

### "Redis connection failed"

**Solution:**
```powershell
# Check Redis is running
docker ps | findstr redis

# Test Redis
docker exec -it url-shortener-redis redis-cli ping
# Should return: PONG

# Restart Redis
docker-compose restart redis
```

### "Rate limit tests failing"

This is usually **OK** - your rate limit might be configured higher than the test expects. Check `.env`:
```env
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

## Continuous Testing

### Watch Mode (Auto-run tests)

```powershell
# Run tests every 30 seconds
while ($true) {
    Clear-Host
    .\test-api.ps1
    Start-Sleep -Seconds 30
}
```

### Before Deployment

Always run the full test suite before deploying:

```powershell
# 1. Run unit tests
.\run-tests.bat

# 2. Start services
.\run.bat

# 3. Run API tests
.\test.bat

# 4. Check logs for errors
docker-compose logs
```

## Adding Your Own Tests

### PowerShell API Test

Add to `test-api.ps1`:

```powershell
# Test 14: Your Custom Test
Write-TestHeader "Your Custom Feature"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/your/endpoint" -Method Get

    if ($response.expected_field) {
        Write-TestResult "Your Test" $true "Success message" $response
    } else {
        Write-TestResult "Your Test" $false "Failure message" $response
    }
} catch {
    Write-TestResult "Your Test" $false "Request failed: $($_.Exception.Message)"
}
```

### Go Unit Test

Create `*_test.go` file:

```go
package yourpackage

import "testing"

func TestYourFunction(t *testing.T) {
    result := YourFunction("input")
    expected := "expected output"

    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

## Test Environments

### Local Development

```powershell
# .env should point to local/dev resources
ENV=development
BASE_URL=http://localhost:8080
```

### Staging

```powershell
# Point to staging resources
ENV=staging
BASE_URL=https://staging.yourapp.com
```

### Production

**Don't run destructive tests against production!**

Only run read-only tests:
- Health checks
- Analytics retrieval
- Link information

## Performance Testing

### Load Test with PowerShell

```powershell
$iterations = 100
$results = @()

for ($i = 0; $i -lt $iterations; $i++) {
    $start = Get-Date
    Invoke-RestMethod -Uri "http://localhost:8080/health" | Out-Null
    $duration = (Get-Date) - $start
    $results += $duration.TotalMilliseconds
}

$avg = ($results | Measure-Object -Average).Average
Write-Host "Average response time: $avg ms"
```

### Stress Test

Create many short URLs:

```powershell
for ($i = 0; $i -lt 100; $i++) {
    $body = @{
        url = "https://example.com/page-$i"
    } | ConvertTo-Json

    Invoke-RestMethod -Uri "http://localhost:8080/api/shorten" `
        -Method Post `
        -Body $body `
        -ContentType "application/json"
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Run Tests
  run: |
    docker-compose up -d
    sleep 10
    docker build -f Dockerfile.test -t tests .
    docker run --network=host tests
```

## Getting Help

If tests are failing and you're not sure why:

1. **Check service logs:**
   ```powershell
   docker-compose logs api
   ```

2. **Run tests with verbose output:**
   ```powershell
   $VerbosePreference = "Continue"
   .\test-api.ps1
   ```

3. **Share test results:**
   - Send `test-results.csv`
   - Copy console output
   - Include `docker-compose ps` output

## Next Steps

- ✓ Run tests after every code change
- ✓ Add tests for new features
- ✓ Monitor test results in CI/CD
- ✓ Set up automated testing on commits
