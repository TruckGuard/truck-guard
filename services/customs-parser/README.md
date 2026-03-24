# 📦 TruckGuard Customs Parser

### 1. What is it?

The **Customs Parser** is a specialized Go-based microservice that acts as a bridge between the TruckGuard ecosystem and external customs data sources.

### 2. Purpose & How it Works

In its current implementation, this service serves as a high-fidelity **Mock** of the Ukrainian Customs **"Unified Window" (Єдине вікно)** API. It allows the system to simulate fetching real-world declaration data during the vehicle registration process.

#### **Key Endpoints**

- **`GET /customs/declaration/:number`**: Retrieves declaration details (declarant, goods, sender, receiver, VMD number) for a given declaration identifier.

### 3. Tech Stack

- **Language**: [Go 1.23+](https://go.dev/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Observability**: [OpenTelemetry](https://opentelemetry.io/)
- **Logging**: [slog](https://pkg.go.dev/log/slog) with custom telemetry wrappers.

### 4. Getting Started

#### **Prerequisites**

- Go (v1.23 or higher)
- Environment Variables (see below)

#### **Run Commands**

1.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
2.  **Start the service:**
    ```bash
    go run main.go
    ```
3.  **Build:**
    ```bash
    go build -o customs-parser
    ./customs-parser
    ```

### 5. Configuration (Environment Variables)

```env
PORT=8085
```
