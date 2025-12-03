# Koda Shortlink - Backend

RESTful API for URL shortener service built with Go, Gin, PostgreSQL, and Redis.

## 🎯 App Description

Koda Shortlink Backend is a high-performance RESTful API service that provides URL shortening functionality with advanced features including user authentication, analytics tracking, and real-time caching.

### Key Features

- **JWT Authentication** - Secure user authentication with access and refresh tokens
- **URL Shortening** - Generate short, unique codes for long URLs
- **Analytics Dashboard** - Track clicks, views, and user statistics
- **Redis Caching** - Fast link resolution with Redis cache
- **Click Tracking** - Detailed analytics including IP, device, browser, and location
- **Auto Migration** - Database migrations run automatically on startup
- **Swagger Documentation** - Interactive API documentation
- **Rate Limiting** - Protect API from abuse
- **CORS Support** - Cross-origin resource sharing enabled

## 🏗 Backend Architecture

### Database Schema (Mermaid ERD)

```mermaid
erDiagram
    users ||--o{ sessions : has
    users ||--o{ short_links : creates
    users ||--o{ clicks : tracks
    short_links ||--o{ clicks : receives

    users {
        serial id PK
        varchar fullName
        varchar email UK
        text password
        timestamp created_at
        timestamp updated_at
        int created_by FK
        int updated_by FK
    }

    sessions {
        serial id PK
        int user_id FK
        text refresh_token
        timestamp login_time
        timestamp logout_time
        timestamp expired_at
        varchar ip_address
        text user_agent
        bool is_active
        timestamp created_at
        timestamp updated_at
        int created_by FK
        int updated_by FK
    }

    short_links {
        serial id PK
        int user_id FK
        varchar short_code UK
        text original_url
        varchar title
        bool is_active
        timestamp expired_at
        int click_count
        timestamp last_clicked_at
        timestamp created_at
        timestamp updated_at
        int created_by FK
        int updated_by FK
    }

    clicks {
        serial id PK
        int short_link_id FK
        timestamp clicked_at
        varchar ip_address
        text referer
        text user_agent
        varchar country
        varchar city
        varchar device_type
        varchar browser
        varchar os
        timestamp created_at
        timestamp updated_at
        int created_by FK
        int updated_by FK
    }
```

## 🚀 How to Run Backend

### Prerequisites

- Go >= 1.23
- PostgreSQL >= 16
- Redis >= 7

### Installation Steps

1. **Clone the repository**

```bash
   git clone https://github.com/ItsnaMaulanaHasan/koda-b4-final-backend.git
   cd koda-shortlink-backend
```

2. **Install dependencies**

```bash
   go mod download
```

3. **Configure environment variables**

   Create a `.env` file in the root directory:

```env
   # Application
   APP_URL=http://localhost:8080/
   ORIGIN_URL=http://localhost:3000

   # JWT Secrets
   APP_SECRET=your-access-token-secret-key
   REFRESH_SECRET=your-refresh-token-secret-key

   # Database
   DATABASE_URL=postgres://postgres:password@localhost:5432/koda_shortlink?sslmode=disable

   # Redis
   REDIS_URL=redis://localhost:6379/0
```

4. **Run the application**

```bash
   go run main.go
```

5. **Access the application**
   - API Base URL: http://localhost:8080
   - Swagger Documentation: http://localhost:8080/swagger/index.html

### Available Scripts

```bash
# Run development server
go run main.go

# Build for production
go build -o main main.go

# Run with hot reload (requires air)
air

# Generate Swagger docs
swag init
```

## 📦 How to Run Migrations

Migrations run **automatically** when the backend starts. The system uses `golang-migrate` to manage database schema changes.

### Automatic Migration (On Startup)

When you run the backend, it will:

1. Check for pending migrations
2. Apply new migrations automatically
3. Log the migration status

```bash
# Start backend (migrations run automatically)
go run main.go

# You'll see logs like:
✅ Migrations applied successfully
# OR
✅ No new migrations to apply
```

## 🧪 How to Test Endpoints

### Using Swagger UI (Recommended)

1. Start the backend server
2. Open browser: http://localhost:8080/swagger/index.html
3. Click "Authorize" and enter your JWT token
4. Test endpoints interactively

## 🔄 Redis Flushing Mechanism

The backend uses Redis for caching short links to improve performance. Cache is automatically managed:

### Cache Strategy

**What is Cached:**

- Short link destinations (15 minutes TTL)
- Click counts

**Cache Keys:**

```
link:{shortCode}:destination  → Full link object
link:{shortCode}:clicks       → Click counter
```

### Automatic Cache Invalidation

Cache is automatically cleared when:

1. **Creating a new link** → Clears `link:{shortCode}:destination`
2. **Updating a link** → Clears `link:{shortCode}:destination`
3. **Deleting a link** → Clears `link:{shortCode}:destination`
4. **Click increment** → Increments `link:{shortCode}:clicks`

## 📚 Tech Stack

### Core

- **Go 1.23** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL 16** - Primary database
- **Redis 7** - Caching layer

### Libraries

- **pgx/v5** - PostgreSQL driver
- **go-redis/v9** - Redis client
- **golang-jwt/v5** - JWT authentication
- **golang-migrate/v4** - Database migrations
- **swaggo/swag** - Swagger documentation
- **argon2** - Password hashing
- **user_agent** - User agent parsing

### Tools

- **Docker** - Containerization
- **Swagger** - API documentation
- **godotenv** - Environment management

## 📡 API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh access token

### User

- `GET /api/v1/users` - Get current user profile

### Short Links

- `POST /api/v1/links` - Create short link
- `GET /api/v1/links` - Get all user links (with filters)
- `GET /api/v1/links/:shortCode` - Get link by code
- `PUT /api/v1/links/:shortCode` - Update link
- `DELETE /api/v1/links/:shortCode` - Delete link
- `GET /:shortCode` - Redirect to original URL

### Dashboard

- `GET /api/v1/dashboard/stats` - Get dashboard statistics

## 🔗 Related Repositories

- [Frontend Repository](https://github.com/ItsnaMaulanaHasan/koda-b4-final-frontend.git) - React + Vite web application

## 📝 Environment Variables

| Variable         | Description                  | Example                             |
| ---------------- | ---------------------------- | ----------------------------------- |
| `APP_URL`        | Backend base URL             | `http://localhost:8080/`            |
| `ORIGIN_URL`     | Frontend URL for CORS        | `http://localhost:3000`             |
| `APP_SECRET`     | JWT access token secret      | `your-secret-key`                   |
| `REFRESH_SECRET` | JWT refresh token secret     | `your-refresh-secret`               |
| `DATABASE_URL`   | PostgreSQL connection string | `postgres://user:pass@host:5432/db` |
| `REDIS_URL`      | Redis connection string      | `redis://localhost:6379/0`          |
