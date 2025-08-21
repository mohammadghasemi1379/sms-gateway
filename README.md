# 📱 SMS Gateway

A robust, scalable SMS gateway service built with Go, featuring intelligent message queuing, transaction management, and multi-provider support.

## 🏗️ Architecture

The application follows a clean architecture pattern with:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Entities**: Domain models
- **Providers**: External SMS service integrations

### Components

- **SMS Service**: Core SMS sending functionality with queue management
- **User Service**: User management and credit operations
- **Transaction Service**: Financial transaction tracking
- **Queue Manager**: Message distribution and processing
- **Mock Provider**: Simulated SMS provider for testing

## 📋 Prerequisites

- Go 1.24.6 or higher
- Docker and Docker Compose
- MySQL 8.0
- RabbitMQ 3.x
- Redis 7.x (for caching)

## 🛠️ Installation

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/mohammadghasemi1379/sms-gateway.git
cd sms-gateway
```

2. Create environment file:
```bash
cp .env.example .env
# Edit .env file with your configuration
```

3. Build and run with Docker Compose:
```bash
make build-up
```

### Manual Installation

1. Install dependencies:
```bash
go mod download
```

2. Set up the database and message queue
3. Run database migrations:
```bash
go run cmd/sms-gateway/main.go
```

## 📚 Documentation

### 📖 Comprehensive Guides
- **[🧠 Queue Management System](./Docs/README.md)** - Detailed explanation of our intelligent multi-queue architecture with Mermaid diagrams
- **[📡 API Collection](./Docs/sms-gateway.postman_collection.json)** - Complete Postman collection for testing all endpoints

### 🔗 Quick Links
- **Base URL**: `http://localhost:8080/api`
- **RabbitMQ Management**: `http://localhost:15672` (guest/guest)
- **Mock Provider**: `http://localhost:8081`

### API Reference

### Endpoints

#### User Management

**Create User**
```http
POST /api/user/create
Content-Type: application/json

{
  "name": "John Doe",
  "phone_number": "+1234567890"
}
```

**Update User Credit**
```http
POST /api/user/update-credit
Content-Type: application/json

{
  "user_id": 1,
  "amount": 1000
}
```

#### SMS Operations

**Send SMS**
```http
POST /api/sms/send
Content-Type: application/json

{
  "user_id": 1,
  "receive_number": "+1234567890",
  "message": "Hello, World!"
}
```

**Get SMS History**
```http
GET /api/sms/history?user_id=1&limit=10&offset=0
```

## 🗄️ Database Schema

### Users Table
- `id`: Primary key
- `name`: User's full name
- `phone_number`: Unique phone number
- `credit`: Available credit balance
- `created_at`, `updated_at`: Timestamps

### SMS Table
- `id`: Primary key
- `user_id`: Foreign key to users
- `receive_number`: Recipient phone number
- `message`: SMS content
- `status`: PENDING/SENT/FAILED
- `cost`: Message cost
- `created_at`, `updated_at`: Timestamps

### Transactions Table
- `id`: Primary key
- `user_id`: Foreign key to users
- `sms_id`: Foreign key to SMS (optional)
- `amount`: Transaction amount
- `type`: DEBIT/CREDIT
- `description`: Transaction description
- `created_at`, `updated_at`: Timestamps

## 🔄 Message Queue System

The application uses RabbitMQ with an intelligent multi-queue architecture for optimal message processing:

- **Main Queue** (`sms-gateway`): Primary processing under normal load
- **Primary Overflow** (`sms-gateway-primary`): 90% of overflow traffic  
- **Secondary Overflow** (`sms-gateway-secondary`): 10% of overflow traffic

### Queue Strategy
- **Intelligent Load Balancing**: Automatic overflow detection and routing
- **Weighted Distribution**: Configurable traffic distribution across queues
- **Concurrent Consumers**: Dedicated consumers for each queue
- **Adaptive Processing**: Dynamic queue selection based on real-time load

> 📖 **For detailed architecture and flow diagrams**, see our [Queue Management Documentation](./Docs/README.md)

### 📡 API Testing with Postman

Import our complete API collection to test all endpoints:

1. **Download**: [Postman Collection](./Docs/sms-gateway.postman_collection.json)
2. **Import** into Postman
3. **Configure** environment variables:
   - `baseURL`: `http://localhost:8080`

The collection includes:
- ✅ User management endpoints
- ✅ SMS sending and history
- ✅ Credit management
- ✅ Health checks
- ✅ Pre-configured test data

## 📦 Make Commands

```bash
# Build and start all services
make build-up

# Start existing services
make up

# View application logs
make logs

# Stop and remove all services with volumes
make destroy
```

## 🐳 Docker Services

The docker-compose setup includes:

- **MySQL**: Database server (port 3306)
- **Redis**: Cache server (port 6379)
- **RabbitMQ**: Message broker (port 5672, management UI: 15672)
- **SMS Gateway**: Main application (port 8080)
- **SMS Provider Mock**: Mock provider (port 8081)

## 📈 Monitoring

- RabbitMQ management interface at http://localhost:15672

### Logs

View application logs:
```bash
make logs
```

View specific service logs:
```bash
docker-compose logs -f [service-name]
```
