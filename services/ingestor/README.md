# 🚀 TruckGuard Ingestor Service

### 1. What is it?

The **Ingestor Service** is the entry point for all IoT data (cameras, scales, sensors) in the TruckGuard ecosystem. Written in **Go**, it handles high-volume ingestion with minimal latency.

### 2. Purpose & How it Works

- **Unified Ingestion**:
  - `/ingest/camera`: Images + metadata ingestion.
  - `/ingest/weight`: Weight sensor data ingestion.
- **Async Streaming**: All events are pushed into a single Redis Stream: `events:adapter`.
- **Blob Storage**: JPG frames from cameras are stored in **MinIO**.
- **Self-Describing Events**: Every event includes a `type` field so downstream workers know how to process it.

### 3. How to Run (Standalone)

#### **Run Commands**

1.  **Install dependencies:** `go mod tidy`
2.  **Start the service:** `go run .`
