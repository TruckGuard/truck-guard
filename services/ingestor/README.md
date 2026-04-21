# 🚀 TruckGuard Ingestor Service

### 1. What is it?

The **Ingestor Service** is the entry point for all IoT data (cameras, scales, sensors) in the TruckGuard ecosystem. Written in **Go**, it handles high-volume ingestion with minimal latency.

### 2. Purpose & How it Works

- **Unified Ingestion**:
  - `/ingest/camera`: Images + metadata ingestion.
  - `/ingest/weight`: Weight sensor data ingestion.
- **Async Streaming**: All events are pushed into a single Redis Stream: `events:adapter`.
- **Blob Storage**: JPG frames from cameras are stored in **Garage** (S3-compatible).
- **Self-Describing Events**: Every event includes a `type` field so downstream workers know how to process it.

### 3. Tech Stack

- **Language**: [Go 1.23+](https://go.dev/)
- **Infrastructure**: [Valkey/Redis](https://valkey.io/), [Garage Storage](https://garagehq.deuxfleurs.fr/)
- **Observability**: [OpenTelemetry](https://opentelemetry.io/)

### 4. Getting Started

#### **Prerequisites**

- Go (v1.23 or higher)
- Redis/Valkey
- Access to Garage Storage

#### **Run Commands**

1.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
2.  **Start the service:**
    ```bash
    go run .
    ```

### 5. Configuration (Environment Variables)

```env
PORT=8080
VALKEY_ADDR=localhost:6379
STORAGE_ENDPOINT=localhost:3900
STORAGE_ACCESS_KEY=your_access_key
STORAGE_SECRET_KEY=your_secret_key
BUCKET_NAME=truckguard-images
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```
