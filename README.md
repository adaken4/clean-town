# CleanTown (WIP)
The purpose of CleanTown is to connect volunteers, organizations, and sponsors in coordinated garbage collection campaigns across Kenyan towns and beyond.

# CleanTown Backend API

A production-ready Go REST API for volunteer cleanup coordination, donation management, and environmental impact tracking. Built following the patterns and best practices from [Let's Go Further](https://lets-go-further.alexedwards.net/) by Alex Edwards.

## 🌍 About CleanTown

CleanTown is a platform that connects volunteers, organizers, and sponsors to coordinate community cleanup events, track environmental impact, and facilitate donations. Our backend API powers volunteer registration, event management, real-time check-ins, M-Pesa/Stripe payment processing, and impact analytics.

## 🚀 Features

### MVP (Current)
- ✅ Volunteer registration and authentication (JWT)
- ✅ Event creation, listing, and RSVP management
- ✅ Geolocation-based event discovery
- ✅ Mobile check-in with QR codes and geofencing
- ✅ Dual payment processing (M-Pesa + Stripe)
- ✅ Impact tracking with photo uploads
- ✅ Campaign management for awareness
- ✅ Admin dashboard endpoints
- ✅ Prometheus metrics and health checks

### Planned (v1+)
- 📧 Email workflows and notifications
- 🏆 Gamification and leaderboards  
- 🔍 Full-text search and advanced filtering
- 🌐 Multi-language support (EN/SW)
- 🔗 Partner webhooks and integrations

## 🏗️ Architecture

This API is built with:

- **Language**: Go 1.21+ with minimal dependencies
- **HTTP**: `net/http` standard library with custom middleware
- **Database**: PostgreSQL with connection pooling
- **Migrations**: `golang-migrate` for schema management
- **Auth**: JWT tokens with role-based authorization
- **Payments**: Stripe (global) + M-Pesa Daraja API (Kenya)
- **Storage**: Signed URLs for S3/Supabase file uploads
- **Observability**: Prometheus metrics, structured logging
- **Deployment**: Docker containers with systemd

## 📋 Prerequisites

- **Go**: 1.21 or later ([download](https://golang.org/dl/))
- **PostgreSQL**: 14+ ([installation guide](https://www.postgresql.org/download/))
- **golang-migrate**: For database migrations ([installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- **Make**: For build automation (usually pre-installed on Unix systems)

### Optional for Development
- **Docker & Docker Compose**: For containerized development
- **curl or httpie**: For API testing
- **Postman**: For API exploration ([collection link](#api-documentation))

## 🛠️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/adaken4/clean-town.git
cd clean-town
```

### 2. Environment Configuration

Copy the example environment file and customize:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
PORT=8080
ENV=development
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Database Configuration  
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cleantown_dev
DB_USER=cleantown_user
DB_PASSWORD=your_secure_password
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_IDLE_TIME=15m

# Authentication
JWT_SECRET=your_jwt_secret_key_here_minimum_32_characters
JWT_EXPIRY=24h
BCRYPT_COST=12

# Payment Configuration
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
MPESA_CONSUMER_KEY=your_mpesa_consumer_key
MPESA_CONSUMER_SECRET=your_mpesa_consumer_secret
MPESA_PASSKEY=your_mpesa_passkey
MPESA_SHORTCODE=174379

# File Storage
STORAGE_PROVIDER=supabase  # or 's3'
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_supabase_anon_key
SUPABASE_SERVICE_KEY=your_supabase_service_key

# Email (Optional - for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

# Observability
LOG_LEVEL=info
METRICS_ENABLED=true
```

### 3. Database Setup

Create a PostgreSQL database and user:

```sql
CREATE USER cleantown_user WITH PASSWORD 'your_secure_password';
CREATE DATABASE cleantown_dev OWNER cleantown_user;
GRANT ALL PRIVILEGES ON DATABASE cleantown_dev TO cleantown_user;
```

### 4. Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 5. Run Database Migrations

```bash
# Run all up migrations
make migrate-up

# Check migration status
make migrate-status
```

### 6. Build and Run

```bash
# Build the application
make build

# Run in development mode
make run

# Or run directly with go
go run ./cmd/cleantown
```

The server will start on `http://localhost:8080` (or your configured PORT).

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run only unit tests
make test-unit

# Run integration tests (requires Docker)
make test-integration
```

### Linting and Formatting

```bash
# Run linter
make lint

# Format code
make fmt

# Run security checks
make security
```

## 📊 API Documentation

### Health Checks

- `GET /health` - Basic health check
- `GET /ready` - Readiness check (includes DB connectivity)  
- `GET /metrics` - Prometheus metrics

### Authentication

- `POST /v1/auth/register` - Register new volunteer
- `POST /v1/auth/login` - Login and get JWT token
- `POST /v1/auth/refresh` - Refresh JWT token
- `POST /v1/auth/logout` - Logout and invalidate token

### Users & Profiles

- `GET /v1/users/{id}` - Get public user profile
- `PUT /v1/users/{id}` - Update user profile (authenticated)
- `GET /v1/users/me/events` - Get user's events

### Events

- `POST /v1/events` - Create event (organizer only)
- `GET /v1/events` - List events (with filtering)
- `GET /v1/events/{id}` - Get event details  
- `POST /v1/events/{id}/join` - RSVP to event
- `POST /v1/events/{id}/checkin` - Check-in to event
- `POST /v1/events/{id}/impact` - Record cleanup impact

### Donations

- `POST /v1/donations/stripe/checkout` - Create Stripe checkout session
- `POST /v1/donations/mpesa/initiate` - Start M-Pesa STK push
- `POST /v1/webhooks/stripe` - Stripe webhook handler
- `POST /v1/webhooks/mpesa` - M-Pesa callback handler

### Campaigns

- `GET /v1/campaigns` - List public campaigns
- `POST /v1/campaigns` - Create campaign (organizer/admin only)

### Admin (Admin role required)

- `GET /v1/admin/stats` - Platform statistics
- `POST /v1/admin/organizers/approve` - Approve organizer applications

### Example Requests

#### Register a New User
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com", 
    "phone": "+254700123456",
    "town": "Nairobi",
    "password": "securePassword123"
  }'
```

#### Create an Event
```bash
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -d '{
    "title": "Karura Forest Cleanup",
    "description": "Monthly forest cleanup drive",
    "location": {
      "latitude": -1.2421,
      "longitude": 36.8617
    },
    "address": "Karura Forest, Nairobi",
    "start_time": "2025-10-15T09:00:00Z",
    "end_time": "2025-10-15T13:00:00Z",
    "capacity": 50
  }'
```

#### List Events with Filtering
```bash
curl "http://localhost:8080/v1/events?town=Nairobi&limit=10&offset=0&sort=date_asc"
```

## 🐳 Docker Development

### Quick Start with Docker Compose

```bash
# Start all services (API + PostgreSQL)
docker-compose up

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f cleantown-api

# Stop services
docker-compose down
```

### Build Docker Image

```bash
# Build production image
make docker-build

# Run container
make docker-run
```

## 🚀 Deployment

### Option 1: Traditional Server (Ubuntu/CentOS)

1. **Build for production:**
   ```bash
   make build-prod
   ```

2. **Copy binary and files to server:**
   ```bash
   scp -r bin/ migrations/ deploy/ user@your-server:/opt/cleantown/
   ```

3. **Set up systemd service:**
   ```bash
   sudo cp deploy/cleantown.service /etc/systemd/system/
   sudo systemctl enable cleantown
   sudo systemctl start cleantown
   ```

### Option 2: Container Deployment

```bash
# Build and push to registry
docker build -t your-registry/cleantown-api:latest .
docker push your-registry/cleantown-api:latest

# Deploy with your orchestrator (k8s, Docker Swarm, etc.)
```

### Environment Variables for Production

Set these additional variables for production:

```env
ENV=production
DB_SSL_MODE=require
CORS_ORIGINS=https://cleantown.org,https://app.cleantown.org
JWT_SECRET=very_long_random_secret_for_production
LOG_LEVEL=warn
```

## 📊 Monitoring & Observability

### Metrics

Access Prometheus metrics at `/metrics`:

- HTTP request duration and status codes
- Database connection pool stats
- Custom business metrics (events created, donations processed)
- Go runtime metrics (memory, goroutines, GC)

### Logging

All logs are structured JSON format:

```json
{
  "time": "2025-09-26T10:30:00Z",
  "level": "info", 
  "msg": "HTTP request completed",
  "method": "POST",
  "path": "/v1/events",
  "status": 201,
  "duration": "45ms",
  "user_id": "123",
  "trace_id": "abc-def-ghi"
}
```

### Health Monitoring

- `/health` - Always returns 200 if service is running
- `/ready` - Returns 200 only if database is accessible
- Use these endpoints for load balancer health checks

## 🧩 Development

### Project Structure

```
.
├── cmd/cleantown/           # Application entry point
├── internal/                # Private application code
│   ├── app/                # Application setup and lifecycle  
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and helpers
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models and validation
│   ├── services/           # Business logic
│   └── workers/            # Background job processing
├── pkg/                    # Public, reusable packages
├── migrations/             # Database migration files
├── scripts/                # Deployment and utility scripts
├── deploy/                 # Deployment configurations
├── docs/                   # Documentation
└── tests/                  # Integration tests
```

### Code Style

This project follows:

- Standard Go formatting (`gofmt`)
- [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- [Let's Go Further](https://lets-go-further.alexedwards.net/) patterns
- 100-character line length
- Comprehensive error handling
- Structured logging

### Adding New Features

1. **Create database migration** (if needed)
2. **Add/update models** in `internal/models/`  
3. **Implement business logic** in `internal/services/`
4. **Create HTTP handlers** in `internal/handlers/`
5. **Add routes** to main router
6. **Write tests** for all layers
7. **Update API documentation**

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make build-prod        # Build for production (optimized)
make run               # Run in development mode
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make lint              # Run linter
make fmt               # Format code
make security          # Run security checks
make migrate-up        # Run database migrations
make migrate-down      # Rollback last migration
make migrate-status    # Show migration status
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make clean             # Clean build artifacts
```

## 🤝 Contributing

1. **Fork the repository**
2. **Create feature branch**: `git checkout -b feature/amazing-feature`
3. **Follow coding standards** and add tests
4. **Commit changes**: `git commit -m 'Add amazing feature'`
5. **Push to branch**: `git push origin feature/amazing-feature`
6. **Open Pull Request**

### Pull Request Requirements

- [ ] All tests pass (`make test`)
- [ ] Code is properly formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] New features include tests
- [ ] Documentation updated (if applicable)

## 🔐 Security

- **Authentication**: JWT tokens with configurable expiry
- **Authorization**: Role-based access control
- **Input validation**: All inputs sanitized and validated
- **Rate limiting**: Configurable per-endpoint limits
- **CORS**: Strict origin controls
- **Secrets**: Never commit credentials; use environment variables
- **Dependencies**: Regular security audits with `go mod audit`

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support & Contact

- **Issues**: [GitHub Issues](https://github.com/yourusername/cleantown-backend/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/cleantown-backend/discussions)
- **Email**: support@cleantown.org
- **Documentation**: [Full API docs](https://docs.cleantown.org)

## 🙏 Acknowledgments

- [Alex Edwards](https://www.alexedwards.net/) for the excellent [Let's Go Further](https://lets-go-further.alexedwards.net/) book
- The Go community for amazing tools and libraries
- All contributors and volunteers making CleanTown possible

---

**Built with ❤️ for a cleaner planet 🌍**