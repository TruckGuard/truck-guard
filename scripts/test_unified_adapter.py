import requests
import json
import time

# Base URL for ingestor
BASE_URL = "http://localhost:8082" 
# Use a valid API Key from .env if running against a real system
API_KEY = "worker_system_secret_2025" 

def test_camera_ingest():
    print("Testing Camera Ingest...")
    url = f"{BASE_URL}/ingest/camera"
    headers = {
        "X-API-Key": API_KEY,
        "X-Source-ID": "cam-test-01",
        "X-Source-Name": "Test Camera"
    }
    files = {
        'image': ('test.jpg', b'fake-image-data', 'image/jpeg')
    }
    data = {
        'device_id': 'dev-cam-01',
        'payload': json.dumps({"plate": "AA1234BB"})
    }
    resp = requests.post(url, headers=headers, files=files, data=data)
    print(f"Response: {resp.status_code}, {resp.text}")

def test_weight_ingest():
    print("Testing Weight Ingest...")
    url = f"{BASE_URL}/ingest/weight"
    headers = {
        "X-API-Key": API_KEY,
        "X-Source-ID": "scale-test-01",
        "X-Source-Name": "Test Scale"
    }
    data = {
        'device_id': 'dev-scale-01',
        'payload': json.dumps({"weight": 12500.5})
    }
    resp = requests.post(url, headers=headers, data=data)
    print(f"Response: {resp.status_code}, {resp.text}")

if __name__ == "__main__":
    # Note: These tests assume the services are running and accessible
    # In a real environment, you'd check Redis streams and Core logs
    # test_camera_ingest()
    # test_weight_ingest()
    pass
