# 🛡️ TruckGuard Auth Service

### 1. What is it?

The **Auth Service** is the central security "gatekeeper" for the TruckGuard ecosystem. It is a high-performance Go-based microservice designed to handle authentication for both human users and IoT devices (cameras/adapters).

### 2. Purpose & How it Works

The service ensures that only authorized entities can access the system's internal resources. It operates on two levels:

- **User Authentication:** Validates login/password and issues **JWT tokens** for the SvelteKit frontend.
- **Machine Authentication:** Validates **X-API-Keys** for cameras and ingestion adapters.
- **Nginx Integration:** Works with the Nginx `auth_request` module. Before a request reaches the backend, Nginx makes a sub-request to this service to verify the token or key.

The codebase is organized into modular packages under `src/`:
- `src/api`: Handlers and middleware.
- `src/models`: Database GORM models.
- `src/repository`: Database and Redis logic.

### 3. Tech Stack

- **Language**: [Go 1.23+](https://go.dev/)
- **ORMs**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Cache/Session**: [Redis/Valkey](https://valkey.io/)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)

### 4. Getting Started

#### **Prerequisites**

- Go (v1.23 or higher)
- PostgreSQL
- Redis/Valkey
- Environment Variables (see below)

#### **Run Commands**

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```
2. **Start the service:**
   ```bash
   go run .
   ```
3. **Build:**
   ```bash
   go build -o auth-service
   ./auth-service
   ```

### 5. Configuration (Environment Variables)

Create a `.env` file or set the following variables:

```env
PORT=8081
DATABASE_URL=postgres://user:pass@localhost:5432/truckguard
REDIS_ADDR=localhost:6379
JWT_SECRET=your_secret_key
ADMIN_DEFAULT_PASSWORD=admin123
```
