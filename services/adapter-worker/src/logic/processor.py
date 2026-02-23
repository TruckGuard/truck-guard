import json
from src.utils.logging_utils import logger
from src.logic.camera_processor import CameraProcessor
from src.logic.weight_processor import WeightProcessor

class EventProcessor:
    def __init__(self, core_client, parser, minio_client, anpr_client):
        self.camera_processor = CameraProcessor(core_client, parser, minio_client, anpr_client)
        self.weight_processor = WeightProcessor(core_client, parser)

    def process(self, raw_data_str: str):
        data = json.loads(raw_data_str)
        event_type = data.get("type")

        if event_type == "camera":
            self.camera_processor.process(data)
        elif event_type == "weight":
            self.weight_processor.process(data)
        else:
            logger.warning(f"Unknown event type received: {event_type}")
