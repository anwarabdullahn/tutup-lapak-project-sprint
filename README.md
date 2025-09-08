# TutupLapak Project Sprint

Microservices architecture for TutupLapak application with Go backend services and NestJS profile service.

## 📋 Prerequisites

- **Docker & Docker Compose** - for infrastructure
- **Node.js** (v18+) - for profile-service
- **Go** (v1.23+) - for auth-service and backend-infra
- **PostgreSQL** - via Docker container

## 🏗️ Architecture

```
├── backend-infra/     # API Gateway (Go + Fiber)
├── auth-service/      # Authentication Service (Go + Fiber + GORM)
├── profile-service/   # Profile Service (NestJS + Prisma)
└── docs/             # Documentation
```

## 🚀 Quick Start

### 1. Start Infrastructure (Database & Redis)

```bash
cd backend-infra
docker-compose up -d
```

**Verify infrastructure:**
```bash
# Check containers
docker ps

# Test PostgreSQL
docker exec infra_postgres psql -U postgres -c "SELECT version();"

# Test Redis
docker exec infra_redis redis-cli ping
```

### 2. Setup Profile Service

```bash
cd profile-service

# Install dependencies
npm install

# Generate Prisma client
npm run prisma:generate

# Sync database schema (recommended for development)
npm run prisma:db:push

# OR apply migrations (alternative)
npm run prisma:migrate:reset

# Start development server
npm run start:dev
```

**Profile Service will run on: http://localhost:3002**

### 3. Setup Auth Service

```bash
cd auth-service

# Install Go dependencies
go mod tidy

# Build and run
go run cmd/server/main.go

# OR build binary first
go build -o bin/auth-service cmd/server/main.go
./bin/auth-service
```

**Auth Service will run on: http://localhost:3001**

### 4. Setup Backend Infra (API Gateway)

```bash
cd backend-infra

# Install Go dependencies
go mod tidy

# Build and run
go run cmd/server/main.go

# OR build binary first
go build -o bin/gateway cmd/server/main.go
./bin/gateway
```

**API Gateway will run on: http://localhost:3000**

## 🔍 Health Checks

After all services are running, test health endpoints:

```bash
# Infrastructure
curl http://localhost:3000/healthz  # API Gateway
curl http://localhost:3001/healthz  # Auth Service  
curl http://localhost:3002/healthz  # Profile Service

# Service Info
curl http://localhost:3001/         # Auth Service info
curl http://localhost:3002/         # Profile Service info
```

## 📊 Database Management

### Profile Service (Prisma)

```bash
cd profile-service

# Generate Prisma client
npm run prisma:generate

# Apply schema changes (development)
npm run prisma:db:push

# Create new migration
npm run prisma:migrate

# Reset database (CAUTION: deletes all data)
npm run prisma:migrate:reset

# Open Prisma Studio (GUI)
npm run prisma:studio
```

### Auth Service (GORM)

```bash
cd auth-service

# Run migrations (if available)
go run cmd/migrate/main.go

# OR connect to database directly
docker exec -it infra_postgres psql -U postgres -d postgres
```

## 🌐 API Endpoints

### Profile Service (port 3002)
- `GET /healthz` - Health check
- `GET /` - Service info
- `POST /profile` - Create profile
- `GET /profile` - Get all profiles
- `GET /profile/:id` - Get profile by ID
- `GET /profile/user/:userId` - Get profile by user ID
- `PATCH /profile/:id` - Update profile
- `DELETE /profile/:id` - Delete profile

### Auth Service (port 3001)
- `GET /healthz` - Health check
- `GET /` - Service info
- Authentication endpoints (to be implemented)

### Backend Infra - API Gateway (port 3000)
- `GET /healthz` - Health check
- `/auth/*` - Proxy to Auth Service
- `/profile/*` - Proxy to Profile Service

## 🛠️ Development Scripts

### Profile Service
```bash
npm run start:dev        # Development server
npm run build           # Build TypeScript
npm run start           # Production server
npm run prisma:studio   # Database GUI
```

### Go Services (Auth & Backend-infra)
```bash
go run cmd/server/main.go  # Development
go build cmd/server/main.go # Build binary
go mod tidy               # Install dependencies
```

## 🔧 Environment Variables

### Profile Service (.env)
```env
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/postgres?schema=public"
PORT=3002
```

### Auth Service
```env
PORT=3001
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/postgres"
```

### Backend Infra
```env
PORT=3000
SERVER_PORT=3000
```

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Test connection
docker exec infra_postgres psql -U postgres -c "SELECT version();"

# Restart infrastructure
cd backend-infra && docker-compose restart
```

### Prisma Issues
```bash
# Reset and regenerate
npm run prisma:generate
npm run prisma:db:push

# If migration fails, use direct schema sync
npm run prisma:db:push
```

### Go Module Issues
```bash
# Clean and reinstall
go clean -modcache
go mod tidy
go mod download
```

## 📚 Technology Stack

- **API Gateway**: Go + Fiber + Viper
- **Auth Service**: Go + Fiber + GORM + Viper  
- **Profile Service**: Node.js + NestJS + Prisma + TypeScript
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Containerization**: Docker + Docker Compose

## 🚀 Production Deployment

1. Build all services:
```bash
# Profile Service
cd profile-service && npm run build

# Auth Service  
cd auth-service && go build -o bin/auth-service cmd/server/main.go

# Backend Infra
cd backend-infra && go build -o bin/gateway cmd/server/main.go
```

2. Run migrations:
```bash
cd profile-service && npm run prisma:migrate:deploy
```

3. Start services in order:
   1. Infrastructure (PostgreSQL, Redis)
   2. Auth Service
   3. Profile Service  
   4. API Gateway

## 📝 Notes

- For development, use `prisma:db:push` instead of migrations
- API Gateway will proxy requests to appropriate services
- All services use Viper for configuration management
- Database is shared across all services via PostgreSQL container
