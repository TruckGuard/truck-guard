# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TruckGuard is a microservices-based customs weighing system for managing vehicle entry/exit at customs zones. It integrates ANPR cameras and electronic scales with a backend ecosystem to manage permits, vehicle tracking, weight measurements, and reporting.

## Development Commands

### Full Stack (Docker)
```bash
make dev-up-build     # Start dev environment with rebuild
make dev-up           # Start dev environment (no rebuild)
make dev-down         # Stop and remove volumes
make dev-down-soft    # Stop without removing volumes
make dev-rebuild      # Full teardown + rebuild + start
make dev-restart      # Soft restart
make dev-restart-core # Restart core services only (auth, core, nginx)
make dev-init         # Initialize Garage S3 storage + restart
make logs             # Stream logs from all services
```

### Frontend (SvelteKit)
```bash
cd services/frontend
yarn install
yarn dev              # Dev server at http://localhost:5173
yarn build            # Production build
yarn check            # Type check + Svelte validation
yarn check:watch      # Watch mode
```

### Go Services (auth, core, ingestor, customs-parser)
```bash
cd services/<service>
go mod tidy
go run .
```

### Python Services (adapter-worker, anpr)
```bash
cd services/adapter-worker   # or services/anpr
pip install -r requirements.txt
python main.py
```

## Architecture

### Service Map

| Service | Lang | Port | Role |
|---------|------|------|------|
| `auth` | Go | 8080 | JWT + API key auth, user management |
| `core` | Go | 8081 | Permit management, event correlation, system config |
| `ingestor` | Go | 8082 | IoT entry point, publishes to Redis Streams |
| `adapter-worker` | Python | - | Consumes Redis Streams, routes events |
| `anpr` | Python | 8000 | License plate recognition (Nomeroff-Net + PyTorch) |
| `customs-parser` | Go | 8085 | Mock Unified Window customs API |
| `frontend` | SvelteKit | 5173/80 | Web dashboard |
| `nginx` | - | 80 | API gateway with `auth_request` |
| `db` | PostgreSQL | 5432 | Two schemas: `truckguard_auth`, `truckguard_core` |
| `valkey` | Redis | 6379 | Redis Streams (`events:adapter`) + caching |
| `garage` | - | 3900-3903 | S3-compatible image storage |

### Data Flow

```
ANPR Cameras / Scales
       ↓ POST /ingest/camera  /ingest/weight
   Ingestor (Go)
       ↓ publish to Redis Stream (events:adapter)
   Adapter Worker (Python)
       ├─ Camera events → ANPR Service → plate recognition
       └─ Weight events → parse sensor payload
       ↓ POST to Core Service REST API
   Core Service (Go)
       ├─ Correlates events → Permits
       ├─ Metadata → PostgreSQL (truckguard_core)
       └─ Images → Garage S3 (truckguard-images bucket)
       ↓
   Frontend (SvelteKit) — authenticated via JWT from Auth Service
```

### Security Model

- **Nginx** validates every request via `auth_request` → Auth Service before proxying
- **Users**: username/password → JWT tokens
- **Machines** (cameras, adapter worker): API keys validated by Auth Service
- **RBAC**: Operators, Controllers, Accountants, Admins with field-level permissions

### Go Service Structure

All Go services follow the same layout:
```
services/<name>/
├── main.go
├── seed.go          # DB seeding (auth service)
├── go.mod
└── src/
    ├── api/         # Gin handlers + middleware
    ├── models/      # GORM models
    ├── repository/  # Data access layer
    └── pkg/         # Shared utilities (telemetry, etc.)
```

### Frontend Structure

```
services/frontend/src/
├── routes/          # SvelteKit file-based routing
└── lib/
    ├── components/  # Reusable UI components (Shadcn-Svelte + Bits UI)
    └── assets/
```

The frontend uses **Svelte 5 Runes** for reactivity and **Tailwind CSS v4**.

### Event Streaming

- Redis Stream key: `events:adapter`
- Dead letter queue: `events:dlq`
- Events are self-describing with a `type` field (`camera` | `weight`)

## Infrastructure Configuration

- `infra/nginx/nginx.conf` — API gateway routing (`/auth`, `/api`, `/docs`)
- `infra/postgres/init-db.sh` — Creates both databases and users
- `infra/garage/garage.toml` — S3 storage config
- `infra/otel/otel-config.yaml` — OpenTelemetry collector config
- `.env.example` — All required environment variables with descriptions

## CI/CD

**GitHub Actions** on `ghcr.io/TruckGuard/truck-guard/<service>`:
- `build-on-develop.yml`: Push to `develop` → builds only changed services → pushes `:dev` + SHA tags
- `release.yml`: Git tag `v*` → builds all services → pushes semver tags

## Observability

All services are instrumented with **OpenTelemetry** (gRPC OTLP export to collector at port 4317). Go services use `slog` for structured logging; Python services use the standard `logging` module.
