# CLAUDE.md - Implementation Guide

## Project Overview
**Serverless URL Analytics Platform**

A scalable URL shortening service with comprehensive analytics tracking, built on modern serverless infrastructure.

### Tech Stack
- **Backend**: Go (Fiber/Gin framework)
- **Deployment**: Railway or Fly.io
- **Database**: AWS DynamoDB (NoSQL)
- **Cache**: Redis
- **Frontend**: Next.js (React)

---

## Core Features

### 1. URL Shortening Service
- Generate short URLs with random codes (e.g., `abc123`)
- Support custom aliases (e.g., `/my-link`)
- Validate URLs before shortening
- Collision handling for generated codes

### 2. Click Analytics
- **Geographic Data**: Country, city, region
- **Device Information**: Desktop, mobile, tablet, OS, browser
- **Referrer Tracking**: Source of traffic
- **Timestamp**: Click time for temporal analysis
- **IP Address**: (anonymized or hashed for privacy)

### 3. QR Code Generation
- Generate QR codes for shortened URLs
- Multiple formats (PNG, SVG)
- Customizable size and error correction

### 4. API with Rate Limiting
- RESTful API endpoints
- Token-based authentication
- Rate limiting per API key/IP
- Redis-backed rate limiter

### 5. Real-time Analytics Dashboard
- Overview metrics (total clicks, unique visitors)
- Geographic distribution map
- Device/browser breakdown
- Traffic over time graphs
- Referrer sources

---

## Architecture

### Backend Structure (Go)
```
/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── url.go          # URL shortening endpoints
│   │   │   ├── analytics.go   # Analytics endpoints
│   │   │   └── qr.go           # QR code generation
│   │   ├── middleware/          # Middleware (auth, rate limit, CORS)
│   │   └── routes.go            # Route definitions
│   ├── models/                  # Data models
│   │   ├── url.go
│   │   └── analytics.go
│   ├── repository/              # Database layer
│   │   ├── dynamodb/
│   │   │   ├── url_repo.go
│   │   │   └── analytics_repo.go
│   │   └── redis/
│   │       └── cache.go
│   ├── service/                 # Business logic
│   │   ├── url_service.go
│   │   ├── analytics_service.go
│   │   └── qr_service.go
│   └── utils/                   # Utilities
│       ├── shortener.go         # ID generation
│       ├── validator.go         # URL validation
│       └── geoip.go             # GeoIP lookup
├── config/
│   └── config.go                # Configuration management
├── go.mod
├── go.sum
└── Dockerfile
```

### Frontend Structure (Next.js)
```
frontend/
├── app/
│   ├── page.tsx                 # Home/shortener page
│   ├── dashboard/
│   │   └── [id]/
│   │       └── page.tsx         # Analytics dashboard
│   └── api/                     # API routes (if needed)
├── components/
│   ├── URLShortener.tsx
│   ├── AnalyticsDashboard.tsx
│   ├── GeographicMap.tsx
│   └── DeviceChart.tsx
├── lib/
│   └── api.ts                   # API client
├── package.json
└── next.config.js
```

---

## Database Schema

### DynamoDB Tables

#### Table: `urls`
- **PK**: `short_code` (String) - Partition key
- **Attributes**:
  - `original_url` (String)
  - `short_code` (String)
  - `custom_alias` (Boolean)
  - `created_at` (Number - Unix timestamp)
  - `created_by` (String - User ID/API key)
  - `expires_at` (Number - Optional TTL)
  - `is_active` (Boolean)

#### Table: `analytics`
- **PK**: `short_code` (String) - Partition key
- **SK**: `timestamp#click_id` (String) - Sort key
- **Attributes**:
  - `short_code` (String)
  - `timestamp` (Number - Unix timestamp)
  - `ip_hash` (String - Hashed IP)
  - `country` (String)
  - `city` (String)
  - `device_type` (String - mobile/desktop/tablet)
  - `os` (String)
  - `browser` (String)
  - `referrer` (String)
  - `user_agent` (String)

#### Global Secondary Index (GSI): `created_by-index`
- **PK**: `created_by`
- **SK**: `created_at`
- Purpose: Query all URLs created by a specific user/API key

### Redis Cache Structure
```
Key Pattern                    | Value Type | Purpose
-------------------------------|------------|----------------------------------
url:{short_code}              | String     | Cached original URL (TTL: 1h)
analytics:{short_code}:count  | Number     | Total click count
ratelimit:{ip}:{window}       | Number     | Rate limit counter
ratelimit:{api_key}:{window}  | Number     | API key rate limit
```

---

## API Endpoints

### URL Shortening
```
POST   /api/shorten
Body: {
  "url": "https://example.com/very/long/url",
  "custom_alias": "my-link" (optional)
}
Response: {
  "short_url": "https://short.ly/abc123",
  "qr_code_url": "https://short.ly/qr/abc123"
}
```

### URL Redirection
```
GET    /{short_code}
Response: 302 Redirect to original URL
```

### Analytics
```
GET    /api/analytics/{short_code}
Response: {
  "total_clicks": 1234,
  "unique_visitors": 567,
  "geographic_data": [...],
  "device_breakdown": {...},
  "traffic_over_time": [...]
}
```

### QR Code
```
GET    /api/qr/{short_code}?size=256&format=png
Response: QR code image
```

### Link Management
```
GET    /api/links                  # List all user's links
GET    /api/links/{short_code}     # Get link details
DELETE /api/links/{short_code}     # Delete/deactivate link
```

---

## Implementation Phases

### Phase 1: Backend Foundation ✓
- [x] Set up Go project structure
- [ ] Configure DynamoDB connection
- [ ] Configure Redis connection
- [ ] Implement URL model and repository
- [ ] Create short code generator
- [ ] Build URL shortening endpoint
- [ ] Build redirect endpoint

### Phase 2: Analytics & Tracking
- [ ] Implement analytics model and repository
- [ ] Add click tracking middleware
- [ ] Integrate GeoIP lookup (MaxMind or similar)
- [ ] Parse User-Agent for device info
- [ ] Build analytics aggregation service
- [ ] Create analytics API endpoints

### Phase 3: Rate Limiting & Security
- [ ] Implement Redis-based rate limiter
- [ ] Add API key authentication
- [ ] Add CORS middleware
- [ ] Input validation and sanitization
- [ ] Add URL blacklist checking

### Phase 4: QR Code Generation
- [ ] Integrate QR code library (go-qrcode)
- [ ] Build QR code generation endpoint
- [ ] Add caching for generated QR codes
- [ ] Support multiple formats/sizes

### Phase 5: Frontend (Next.js)
- [ ] Set up Next.js project
- [ ] Build URL shortener UI
- [ ] Create analytics dashboard layout
- [ ] Implement geographic map visualization
- [ ] Add device/browser charts
- [ ] Build traffic timeline graph
- [ ] Add link management interface

### Phase 6: Deployment & DevOps
- [ ] Create Dockerfile for Go backend
- [ ] Set up Railway/Fly.io deployment
- [ ] Configure environment variables
- [ ] Set up DynamoDB tables in AWS
- [ ] Configure Redis instance
- [ ] Deploy frontend to Vercel/Railway
- [ ] Set up custom domain
- [ ] Add monitoring and logging

### Phase 7: Testing & Optimization
- [ ] Write unit tests for services
- [ ] Write integration tests for API
- [ ] Load testing with k6 or Artillery
- [ ] Optimize DynamoDB queries
- [ ] Implement connection pooling
- [ ] Add caching strategies
- [ ] Performance monitoring

---

## Environment Variables

```bash
# Server
PORT=8080
ENV=production

# DynamoDB
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
DYNAMODB_URLS_TABLE=urls
DYNAMODB_ANALYTICS_TABLE=analytics

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Application
BASE_URL=https://short.ly
SHORT_CODE_LENGTH=6
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60 # seconds

# GeoIP
GEOIP_DB_PATH=/path/to/GeoLite2-City.mmdb

# Frontend
NEXT_PUBLIC_API_URL=https://api.short.ly
```

---

## Key Libraries

### Go Backend
```go
// go.mod
require (
    github.com/gofiber/fiber/v2        // or github.com/gin-gonic/gin
    github.com/aws/aws-sdk-go-v2
    github.com/redis/go-redis/v9
    github.com/skip2/go-qrcode
    github.com/oschwald/geoip2-golang
    github.com/mileusna/useragent
    github.com/google/uuid
    github.com/joho/godotenv
)
```

### Next.js Frontend
```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "recharts": "^2.10.0",
    "react-map-gl": "^7.1.0",
    "axios": "^1.6.0",
    "tailwindcss": "^3.4.0"
  }
}
```

---

## Development Commands

```bash
# Backend
cd backend
go mod init github.com/yourusername/url-shortener
go mod tidy
go run cmd/server/main.go

# Frontend
cd frontend
npm install
npm run dev

# Docker
docker build -t url-shortener .
docker run -p 8080:8080 url-shortener
```

---

## Next Steps

1. **Initialize Go project** - Set up module and folder structure
2. **Configure AWS DynamoDB** - Create tables with proper indexes
3. **Set up Redis** - Local or cloud instance
4. **Implement core shortening logic** - Generate codes, store URLs
5. **Add analytics tracking** - Capture click data
6. **Build REST API** - All endpoints with proper error handling
7. **Create Next.js dashboard** - User interface for analytics
8. **Deploy to Railway/Fly** - Production deployment

---

## Resources

- [AWS DynamoDB Go SDK](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Fiber Framework](https://docs.gofiber.io/)
- [Redis Go Client](https://redis.uptrace.dev/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Railway Documentation](https://docs.railway.app/)

---

**Last Updated**: 2025-11-09
**Status**: Planning Complete - Ready for Implementation
