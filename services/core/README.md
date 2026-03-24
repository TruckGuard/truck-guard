# 🧠 TruckGuard Core Service

### 1. What is it?

The **Core Service** is the central management unit of the TruckGuard system. It acts as the primary API for managing system configurations, tracking camera states, and logging high-level events.

### 2. Purpose & How it Works

It manages the persistence layer for the entire system's configuration and historical data:

- **Configuration Manager:** Stores and serves camera, weight scale, and gate configurations, presets, and system-wide settings.
- **Event Orchestrator & Correlation:** Receives processed data from cameras and scales (via Adapters/Ingestor), correlates them into unified **Permits** (passes), and handles multi-plate vehicle detection.
- **Ignore List Management:** Maintains a list of excluded license plates.
- **Integration:** Communicates with the **Auth Service** to automatically provision API keys for new cameras and scales.

The project follows a modular Go structure:
- `src/api`: REST handlers and validation middleware.
- `src/models`: Domain entities (Cameras, Presets, Events).
- `src/repository`: GORM-based data access layer.

### 3. Tech Stack

- **Language**: [Go 1.23+](https://go.dev/)
- **ORMs**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)

### 4. Getting Started

#### **Prerequisites**

- Go (v1.23 or higher)
- PostgreSQL
- Environment Variables (see below)

#### **Run Commands**

1.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
2.  **Start the service:**
    ```bash
    go run .
    ```
3.  **Build:**
    ```bash
    go build -o core-service
    ./core-service
    ```

### 5. Configuration (Environment Variables)

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/truckguard
```
