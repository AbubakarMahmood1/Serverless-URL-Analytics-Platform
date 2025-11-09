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
- Go 1.21+
- AWS Account (DynamoDB)
- Redis instance
- Docker (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform.git
   cd Serverless-URL-Analytics-Platform
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your AWS and Redis credentials
   ```

4. **Set up DynamoDB tables**
   See [SETUP.md](SETUP.md) for detailed AWS DynamoDB setup instructions

5. **Run the server**
   ```bash
   make run
   # or
   go run cmd/server/main.go
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

### Docker
```bash
docker build -t url-shortener .
docker run -p 8080:8080 --env-file .env url-shortener
```

## 🧪 Development

```bash
# Run tests
make test

# Build binary
make build

# Format code
go fmt ./...

# Run with hot reload (requires air)
make dev
```

## 🎯 Roadmap

- [x] Core URL shortening functionality
- [x] Analytics tracking
- [x] QR code generation
- [x] Rate limiting
- [ ] User authentication (JWT)
- [ ] API key management
- [ ] Next.js analytics dashboard
- [ ] GeoIP integration
- [ ] Custom domains
- [ ] Link expiration
- [ ] Comprehensive test coverage

## 📄 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or support, please open an issue on GitHub.
