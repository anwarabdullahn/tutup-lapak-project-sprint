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
├── purchase-service/  # Purchase Service (Go + Fiber + GORM)
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

### 4. Setup Purchase Service

```bash
cd purchase-service

# Install Go dependencies
go mod tidy

# Run database migrations
make migrate:up

# Build and run
go run cmd/server/main.go

# OR build binary first
go build -o bin/purchase-service cmd/server/main.go
./bin/purchase-service
```

**Purchase Service will run on: http://localhost:3004**

### 5. Setup Backend Infra (API Gateway)

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
curl http://localhost:3004/healthz  # Purchase Service

# Service Info
curl http://localhost:3001/         # Auth Service info
curl http://localhost:3002/         # Profile Service info
curl http://localhost:3004/         # Purchase Service info
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

### Purchase Service (GORM + Migrate)

```bash
cd purchase-service

# Apply all migrations
make migrate:up

# Apply one migration step
make migrate:up-1

# Rollback all migrations
make migrate:down

# Rollback one migration step
make migrate:down-1

# Check current migration version
make migrate:version

# Create new migration
make migrate:create name=create_table_name

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

### Purchase Service (port 3004)
- `GET /healthz` - Health check
- `GET /` - Service info
- `POST /api/v1/purchase` - Create a new purchase order
- `GET /api/v1/purchase` - List user's purchases (paginated)
- `GET /api/v1/purchase/:purchaseId` - Get purchase by ID
- `POST /api/v1/purchase/:purchaseId` - Upload payment proof

#### Purchase Service API Details

**Create Purchase** - `POST /api/v1/purchase`
```json
{
  "purchasedItems": [
    {
      "productId": "string",
      "qty": 1
    }
  ],
  "senderName": "string",
  "senderContactType": "email|phone",
  "senderContactDetail": "string"
}
```

**Upload Payment Proof** - `POST /api/v1/purchase/:purchaseId`
```json
{
  "fileIds": ["string"]
}
```

**List Purchases** - `GET /api/v1/purchase?page=1&limit=10`
- Query parameters: `page` (default: 1), `limit` (default: 10, max: 100)

### Backend Infra - API Gateway (port 3000)
- `GET /healthz` - Health check
- `/v1/login/*` - Auth endpoints (email/phone login/register)
- `/v1/profile/*` - Profile endpoints (JWT protected)
- `/v1/purchase/*` - Purchase endpoints (JWT protected)
  - `POST /v1/purchase` - Create purchase
  - `GET /v1/purchase` - List purchases
  - `GET /v1/purchase/:id` - Get purchase by ID
  - `POST /v1/purchase/:id` - Upload payment proof

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
AUTH_SERVICE_URL="http://localhost:3001"
PURCHASE_SERVICE_URL="http://localhost:3004"
JWT_SECRET="your-jwt-secret-key"
```

### Purchase Service
```env
PORT=3004
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/postgres"
USER_SERVICE_URL="http://localhost:3002"
PRODUCT_SERVICE_URL="http://localhost:3003"
JWT_SECRET="your-jwt-secret-key"
INTERNAL_SECRET="your-internal-service-secret"
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
- **Purchase Service**: Go + Fiber + GORM + Viper + Migrate
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

# Purchase Service
cd purchase-service && go build -o bin/purchase-service cmd/server/main.go

# Backend Infra
cd backend-infra && go build -o bin/gateway cmd/server/main.go
```

2. Run migrations:
```bash
# Profile Service
cd profile-service && npm run prisma:migrate:deploy

# Purchase Service
cd purchase-service && make migrate:up
```

3. Start services in order:
   1. Infrastructure (PostgreSQL, Redis)
   2. Auth Service
   3. Profile Service
   4. Purchase Service
   5. API Gateway

## 🛒 Purchase Service Features

The Purchase Service handles the complete purchase workflow for the TutupLapak application:

### Core Features
- **Purchase Order Creation**: Create purchase orders with multiple items from cart
- **Product Information Snapshot**: Copies product details to prevent race conditions
- **Payment Proof Upload**: Handle payment proof file uploads
- **Multi-Seller Support**: Aggregates payment details by seller
- **Purchase History**: Paginated list of user purchases
- **External Service Integration**: Fetches data from User and Product services

### Database Schema
- **purchases**: Main purchase records with UUID v7 primary keys
- **purchase_items**: Individual items in each purchase (with product snapshots)
- **purchase_senders**: Sender contact information for each purchase

### External Dependencies
- **User Service**: Fetches seller bank account information
- **Product Service**: Fetches product details and seller information
- **File Service**: Handles payment proof file storage

### Security Features
- JWT token validation for user authentication
- Internal service communication with secret validation
- User ownership validation for purchase access

### Gateway Integration
- All purchase endpoints are accessible through the API Gateway at `/v1/purchase/*`
- JWT authentication is handled at the gateway level
- User context is automatically forwarded to the purchase service
- Request/response proxying with proper error handling

## 📝 Notes

- For development, use `prisma:db:push` instead of migrations
- API Gateway will proxy requests to appropriate services
- All services use Viper for configuration management
- Database is shared across all services via PostgreSQL container
- Purchase Service uses UUID v7 for better performance and ordering
- Product information is copied during purchase to prevent data inconsistency
