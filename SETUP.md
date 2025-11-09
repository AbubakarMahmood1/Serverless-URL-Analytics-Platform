# Setup Guide - Serverless URL Analytics Platform

## Prerequisites

- **Go** 1.21 or higher
- **AWS Account** with DynamoDB access
- **Redis** instance (local or cloud)
- **Docker** (optional, for containerized deployment)

## Installation Steps

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Serverless-URL-Analytics-Platform
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

If you encounter network issues, try:
```bash
GOPROXY=https://goproxy.io,direct go mod tidy
```

### 3. Configure Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```bash
# Server Configuration
PORT=8080
ENV=development

# AWS DynamoDB Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
DYNAMODB_URLS_TABLE=urls
DYNAMODB_ANALYTICS_TABLE=analytics

# Redis Configuration
REDIS_URL=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Application Settings
BASE_URL=http://localhost:8080
SHORT_CODE_LENGTH=6
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

### 4. Set Up AWS DynamoDB Tables

#### Create URLs Table

```bash
aws dynamodb create-table \
  --table-name urls \
  --attribute-definitions \
    AttributeName=short_code,AttributeType=S \
    AttributeName=created_by,AttributeType=S \
    AttributeName=created_at,AttributeType=N \
  --key-schema \
    AttributeName=short_code,KeyType=HASH \
  --global-secondary-indexes \
    IndexName=created_by-index,KeySchema=["{AttributeName=created_by,KeyType=HASH}","{AttributeName=created_at,KeyType=RANGE}"],Projection="{ProjectionType=ALL}",ProvisionedThroughput="{ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --region us-east-1
```

#### Create Analytics Table

```bash
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
  --region us-east-1
```

**Note**: For production, consider using on-demand billing instead of provisioned throughput.

### 5. Set Up Redis

#### Local Redis (Docker)
```bash
docker run -d -p 6379:6379 redis:alpine
```

#### Or install locally
```bash
# macOS
brew install redis
redis-server

# Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis
```

### 6. Run the Application

#### Development Mode
```bash
make run
# or
go run cmd/server/main.go
```

#### Build and Run
```bash
make build
./bin/server
```

#### Docker
```bash
make docker-build
make docker-run
```

## Testing the API

### Health Check
```bash
curl http://localhost:8080/health
```

### Shorten a URL
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/very/long/url",
    "custom_alias": "mylink"
  }'
```

Response:
```json
{
  "short_url": "http://localhost:8080/mylink",
  "short_code": "mylink",
  "qr_code_url": "http://localhost:8080/api/qr/mylink",
  "original_url": "https://example.com/very/long/url",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Access Shortened URL
```bash
curl -L http://localhost:8080/mylink
```

### Get Analytics
```bash
curl http://localhost:8080/api/analytics/mylink
```

### Generate QR Code
```bash
curl http://localhost:8080/api/qr/mylink?size=512 -o qr.png
```

## Deployment

### Railway

1. Install Railway CLI:
```bash
npm i -g @railway/cli
```

2. Login and initialize:
```bash
railway login
railway init
```

3. Add environment variables:
```bash
railway variables set AWS_REGION=us-east-1
railway variables set AWS_ACCESS_KEY_ID=your_key
railway variables set AWS_SECRET_ACCESS_KEY=your_secret
railway variables set REDIS_URL=your_redis_url
```

4. Deploy:
```bash
railway up
```

### Fly.io

1. Install Fly CLI:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Create app:
```bash
fly launch
```

3. Set secrets:
```bash
fly secrets set AWS_REGION=us-east-1
fly secrets set AWS_ACCESS_KEY_ID=your_key
fly secrets set AWS_SECRET_ACCESS_KEY=your_secret
```

4. Deploy:
```bash
fly deploy
```

## Troubleshooting

### Issue: "Failed to connect to DynamoDB"
- Verify AWS credentials are correct
- Check IAM permissions for DynamoDB access
- Ensure tables exist in the specified region

### Issue: "Failed to connect to Redis"
- Verify Redis is running: `redis-cli ping` should return `PONG`
- Check Redis URL and port
- Verify Redis password if authentication is enabled

### Issue: "Rate limit not working"
- Ensure Redis connection is established
- Check rate limit configuration in `.env`

### Issue: "QR code generation failed"
- Verify short code exists
- Check application logs for detailed error messages

## Development

### Project Structure
```
.
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # Custom middleware
│   │   └── routes.go    # Route definitions
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   │   ├── dynamodb/    # DynamoDB implementation
│   │   └── redis/       # Redis cache implementation
│   ├── service/         # Business logic
│   └── utils/           # Utility functions
├── Dockerfile
├── Makefile
└── README.md
```

### Running Tests
```bash
make test
```

### Code Formatting
```bash
go fmt ./...
```

### Build for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server
```

## Next Steps

- [ ] Implement user authentication (JWT)
- [ ] Add API key management
- [ ] Build Next.js frontend dashboard
- [ ] Add comprehensive test coverage
- [ ] Set up CI/CD pipeline
- [ ] Add monitoring and logging (Prometheus, Grafana)
- [ ] Implement GeoIP lookup for location analytics

## Support

For issues and questions, please check:
- CLAUDE.md for detailed implementation notes
- README.md for project overview
- GitHub Issues (if repository is public)
