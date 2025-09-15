# Purchase Service

A microservice for handling purchase operations in the TutupLapak application.

## Features

- **POST /api/v1/purchase** - Create a new purchase order
- Product information copying to prevent race conditions
- Integration with User and Product services
- Payment details aggregation by seller

## API Endpoints

### Create Purchase

**POST** `/api/v1/purchase`

Creates a new purchase order with items from the cart.

#### Request Body

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

#### Response

```json
{
  "purchaseId": "string",
  "purchasedItems": [
    {
      "productId": "string",
      "name": "string",
      "category": "string",
      "qty": 1,
      "price": 100,
      "sku": "string",
      "fileId": "string",
      "fileUri": "string",
      "fileThumbnailUri": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ],
  "totalPrice": 100,
  "paymentDetails": [
    {
      "bankAccountName": "string",
      "bankAccountHolder": "string",
      "bankAccountNumber": "string",
      "totalPrice": 100
    }
  ]
}
```

## Environment Variables

- `USER_SERVICE_URL` - URL of the user service (default: http://localhost:3002)
- `PRODUCT_SERVICE_URL` - URL of the product service (default: http://localhost:3003)
- `SERVER_PORT` - Port to run the service on (default: 3001)
- `JWT_SECRET` - Secret for JWT token validation
- `INTERNAL_SECRET` - Secret for internal service communication

## Database

The service uses PostgreSQL with GORM for database operations. It creates the following tables:

- `purchases` - Main purchase records
- `purchase_items` - Items in each purchase
- `purchase_senders` - Sender information for purchases

## Running the Service

1. Set up environment variables in `.env` file
2. Run database migrations: `make migrate:up`
3. Start the service: `go run cmd/server/main.go`

## Architecture

The service follows a clean architecture pattern:

- **Handlers** - HTTP request/response handling
- **Services** - Business logic
- **Repositories** - Data access layer
- **Entities** - Domain models
- **DTOs** - Data transfer objects

## External Service Integration

The service integrates with:

1. **User Service** - Fetches seller bank account information
2. **Product Service** - Fetches product details and seller information

Both integrations are done via HTTP calls with proper error handling and timeout configuration.
