import time
from redis import Redis
from src.config import cfg
from src.utils.logging_utils import logger
from src.clients.core_client import CoreClient
from src.logic.payload_parser import PayloadParser
from src.logic.processor import EventProcessor
from src.clients.minio_client import MinioStorage
from src.clients.anpr_client import ANPRClient
from src.telemetry import init_telemetry

def main():
    # Initialize OpenTelemetry
    init_telemetry("truckguard-adapter-worker")
    
    logger.info("Starting Adapter Worker...")
    
    redis = Redis.from_url(f"redis://{cfg.REDIS_ADDR}", decode_responses=True) 
    core = CoreClient()
    parser = PayloadParser()
    storage = MinioStorage()  
    anpr = ANPRClient()

    processor = EventProcessor(core, parser, storage, anpr)

    last_id = "0"
    while True:
        try:
            streams = redis.xread({cfg.STREAM_RAW: last_id}, count=1, block=5000)
            if not streams:
                continue

            for _, messages in streams:
                for msg_id, data in messages:
                    raw = data.get("data", "")
                    try:
                        processor.process(raw)
                    except Exception as e:
                        logger.error(f"Failed to process message {msg_id}: {e}")
                        try:
                            redis.xadd(cfg.STREAM_DLQ, {
                                "msg_id": msg_id,
                                "stream": cfg.STREAM_RAW,
                                "data": raw,
                                "error": str(e),
                            })
                        except Exception as dlq_err:
                            logger.critical(f"Failed to write message {msg_id} to DLQ: {dlq_err}")
                            # Do not advance last_id — the message will be retried on next startup.
                            continue

                    last_id = msg_id

        except Exception as e:
            logger.critical(f"Redis connection error: {e}")
            time.sleep(5)

if __name__ == "__main__":
    main()