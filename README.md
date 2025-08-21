# SMS Gateway

A robust, scalable SMS gateway service built with Go, featuring message queuing, transaction management, and multi-provider support.

## 🚀 Features

- **RESTful API** for SMS sending and management
- **Message Queuing** with RabbitMQ for reliable message delivery
- **Multi-Provider Support** with mock provider implementation
- **Credit-based System** for user balance management
- **Transaction Tracking** for all SMS operations
- **Database Migrations** for schema management
- **Docker Support** for easy deployment
- **Graceful Shutdown** handling
- **Structured Logging** with configurable levels

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

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

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

The application uses RabbitMQ with multiple queues for message processing:

- **Primary Queue**: High-priority messages (90% weight)
- **Secondary Queue**: Standard messages (10% weight)
- **Dead Letter Queue**: Failed messages for retry logic

### Queue Strategy
- Messages are distributed based on configured weights
- Consumers process messages with prefetch limits
- Failed messages are handled with retry mechanisms

## 🧪 Testing

The project includes a mock SMS provider for testing purposes:

```bash
# Start mock provider
docker-compose up sms-provider-mock

# The mock provider runs on port 8081
curl http://localhost:8081/health
```

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

## 🔧 Troubleshooting

### Common Issues

**Database Connection Failed**
- Verify MySQL is running and accessible
- Check database credentials in `.env`
- Ensure database `sms_gateway` exists

**RabbitMQ Connection Failed**
- Verify RabbitMQ service is running
- Check RabbitMQ credentials and host configuration
- Ensure RabbitMQ management plugin is enabled

**Port Already in Use**
- Check if services are already running
- Modify port configurations in docker-compose.yml
- Use `docker ps` to see running containers

### Logs

View application logs:
```bash
make logs
```

View specific service logs:
```bash
docker-compose logs -f [service-name]
```
