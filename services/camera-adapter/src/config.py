import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Config:
    CORE_URL: str = os.getenv("CORE_URL", "http://gateway/api/v1")
    WORKER_API_KEY: str = os.getenv(
        "WORKER_API_KEY", "worker_internal_secret_2025_token"
    )
    ANPR_URL: str = os.getenv("ANPR_URL", "http://anpr:8000/recognize")
    REDIS_ADDR: str = os.getenv("VALKEY_ADDR", "valkey:6379")
    STORAGE_ENDPOINT: str = os.getenv("STORAGE_ENDPOINT", "garage:3900")
    STORAGE_ACCESS_KEY: str = os.getenv("STORAGE_ACCESS_KEY", "local-admin")
    STORAGE_SECRET_KEY: str = os.getenv("STORAGE_SECRET_KEY", "local-admin")
    STORAGE_BUCKET: str = os.getenv("STORAGE_BUCKET", "truck-images")

    STREAM_RAW: str = "camera:raw"
    STREAM_DLQ: str = "camera:dlq"
    CACHE_TTL: int = 60


cfg = Config()
