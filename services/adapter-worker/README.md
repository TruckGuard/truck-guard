# ⚙️ TruckGuard Adapter Worker

### 1. What is it?

The **Adapter Worker** is a high-performance asynchronous processing engine written in **Python**. It acts as a unified "bridge" between raw hardware data (cameras, scales) and the backend, orchestrating recognition, parsing, and data transformation.

### 2. Purpose & How it Works

It operates as a background consumer for all hardware events:

1.  **Consume**: Listens to the `events:adapter` Redis Stream for all incoming ingestion events.
2.  **Route**: Dispatches processing based on the event `type` (`camera` or `weight`).
3.  **Process (Camera)**:
    - Forwards images to the **ANPR Service** for hardware-agnostic license plate recognition.
    - Decodes manufacturer-specific payloads.
4.  **Process (Weight)**:
    - Extracts weight values from sensor payloads based on configuration.
5.  **Finalize**: Sends enriched data to the **Core Service**.
6.  **Reliability**: Implements a Dead Letter Queue (`events:dlq`) for handling processing failures.

### 3. How to Run (Standalone)

#### **Prerequisites**

- **Python 3.12+**
- **Redis**
- **Access to Core & ANPR APIs**
- **Access to MinIO Storage**

#### **Configuration**

```env
VALKEY_ADDR=localhost:6379
CORE_URL=http://localhost:8081
ANPR_URL=http://localhost:8000
STORAGE_ENDPOINT=localhost:9000
...
```

#### **Run Commands**

1.  **Install dependencies:** `pip install -r requirements.txt`
2.  **Start the worker:** `python main.py`
