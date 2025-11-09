# Serverless URL Analytics Platform

A modern, scalable URL shortening service with comprehensive analytics tracking, built on serverless infrastructure.

## 🚀 Features

- **URL Shortening**: Generate short URLs with random codes or custom aliases
- **Click Analytics**: Track geographic location, device type, browser, OS, and referrer data
- **QR Code Generation**: Create QR codes for shortened URLs in multiple formats
- **Rate Limiting**: Redis-based rate limiting to prevent abuse
- **Real-time Dashboard**: Analytics dashboard with visualizations (Coming Soon)
- **RESTful API**: Clean, documented API endpoints
- **Serverless**: Deploy on Railway, Fly.io, or any container platform

## 🛠️ Tech Stack

- **Backend**: Go (Fiber framework)
- **Database**: AWS DynamoDB (NoSQL)
- **Cache**: Redis
- **Deployment**: Railway / Fly.io
- **Frontend**: Next.js (Coming Soon)

## 📋 Quick Start

### Prerequisites
- **Docker Desktop** (Windows/Mac/Linux)
- **AWS Account** (Optional - can test locally without AWS!)
- Go 1.21+ (optional, if not using Docker)

> **💡 No AWS? No Problem!** See [LOCAL-TESTING.md](LOCAL-TESTING.md) to test without an AWS account using DynamoDB Local.

### Windows Users (No Go/Make Required!)

See **[WINDOWS.md](WINDOWS.md)** for detailed Windows setup with Docker.

**Quick start:**
```powershell
# 1. Configure environment
copy .env.example .env
# Edit .env with AWS credentials

# 2. Start with PowerShell (recommended)
.\dev.ps1 start

# Or use batch files
.\build.bat
.\run.bat
```

### Linux/Mac Users

#### With Docker (Recommended)
```bash
# 1. Configure environment
cp .env.example .env
# Edit .env with AWS credentials

# 2. Start services
docker-compose up -d

# View logs
docker-compose logs -f
```

#### Without Docker (Native Go)
```bash
# 1. Install dependencies
go mod download
go mod tidy

# 2. Configure environment
cp .env.example .env
# Edit .env

# 3. Run Redis locally
docker run -d -p 6379:6379 redis:alpine
# or: brew install redis && redis-server

# 4. Run the server
make run
# or: go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Shorten URL
```bash
POST /api/shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "custom_alias": "mylink"  // optional
}
```

### Redirect
```bash
GET /:shortCode
```

### Get Analytics
```bash
GET /api/analytics/:shortCode
```

### Generate QR Code
```bash
GET /api/qr/:shortCode?size=256
```

### Link Management
```bash
GET /api/links              # List all links
GET /api/links/:shortCode   # Get link details
DELETE /api/links/:shortCode # Delete link
```

## 📖 Documentation

- **[LOCAL-TESTING.md](LOCAL-TESTING.md)** - **Test without AWS!** (DynamoDB Local setup)
- **[AWS-LAMBDA-GUIDE.md](AWS-LAMBDA-GUIDE.md)** - **Deploy to AWS Lambda** (serverless, ~FREE!)
- **[DEPLOYMENT-OPTIONS.md](DEPLOYMENT-OPTIONS.md)** - Compare all deployment options
- **[DATABASE-COMPARISON.md](DATABASE-COMPARISON.md)** - Understanding databases
- **[WINDOWS.md](WINDOWS.md)** - Windows setup guide (Docker-based, no Go required)
- **[TESTING.md](TESTING.md)** - Comprehensive testing guide with automated tests
- **[CLAUDE.md](CLAUDE.md)** - Complete implementation guide and architecture
- **[SETUP.md](SETUP.md)** - Detailed setup and deployment instructions

## 🏗️ Project Structure

```
.
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── models/          # Data models
│   ├── repository/      # Data access layer (DynamoDB, Redis)
│   ├── service/         # Business logic
│   └── utils/           # Utility functions
├── Dockerfile
├── Makefile
├── CLAUDE.md            # Implementation guide
└── SETUP.md             # Setup instructions
```

## 🚢 Deployment

### AWS Lambda (Serverless - Recommended!) 🆕
```powershell
# Deploy to AWS Lambda (~FREE for most traffic!)
.\deploy-lambda.ps1

# Cost: $0/month for up to 1M requests!
```
See **[AWS-LAMBDA-GUIDE.md](AWS-LAMBDA-GUIDE.md)** for complete guide.

### Railway
```bash
railway login
railway init
railway up
```

### Fly.io
```bash
fly launch
fly deploy
```

### Docker (Any VPS)
```bash
docker build -t url-shortener .
docker run -p 8080:8080 --env-file .env url-shortener
```

**Compare all options:** [DEPLOYMENT-OPTIONS.md](DEPLOYMENT-OPTIONS.md)

## 🧪 Testing

### Automated Tests (Windows - No Go Required!)

```powershell
# Run comprehensive API tests
.\test.bat

# Or use PowerShell directly
.\test-api.ps1

# Run unit tests in Docker
.\run-tests.bat

# Via dev menu
.\dev.ps1
# Select option 9: Test API
```

**Tests Include:**
- ✓ Health checks
- ✓ URL shortening (random & custom)
- ✓ URL validation
- ✓ Redirects
- ✓ Analytics tracking
- ✓ QR code generation
- ✓ Rate limiting
- ✓ CORS headers

**Test Output:**
- Green ✓ = Passed
- Red ✗ = Failed
- Results saved to `test-results.csv`

See **[TESTING.md](TESTING.md)** for complete testing guide.

### Development

```bash
# With Docker
docker-compose up -d
docker-compose logs -f

# Native Go (Linux/Mac)
make test
make build
go fmt ./...
```

## 🎯 Roadmap

- [x] Core URL shortening functionality
- [x] Analytics tracking
- [x] QR code generation
- [x] Rate limiting
- [x] Docker-based development (Windows compatible)
- [x] Comprehensive automated tests
- [x] Unit tests & integration tests
- [x] **AWS Lambda deployment** (serverless!)
- [x] Local testing without AWS (DynamoDB Local)
- [x] Multiple deployment options
- [ ] User authentication (JWT)
- [ ] API key management
- [ ] Next.js analytics dashboard
- [ ] GeoIP integration
- [ ] Custom domains
- [ ] Link expiration

## 📄 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or support, please open an issue on GitHub.
