import requests
import json
import time
import io
import os
import random
import string
import redis
from PIL import Image

# Конфігурація
BASE_URL = os.getenv("BASE_URL", "http://localhost")
AUTH_URL = f"{BASE_URL}/auth"
CORE_API_URL = f"{BASE_URL}/api"
INGEST_CAMERA_URL = f"{BASE_URL}/ingest/camera"
INGEST_WEIGHT_URL = f"{BASE_URL}/ingest/weight"
VALKEY_ADDR = os.getenv("VALKEY_ADDR", "localhost:6379")

ADMIN_USER = os.getenv("ADMIN_USER", "admin")
ADMIN_PASS = os.getenv("ADMIN_DEFAULT_PASSWORD", "secret123")

TRUCK_COUNT = int(os.getenv("TRUCK_COUNT", "3"))

# Valkey Client for caching
r_client = redis.Redis.from_url(f"redis://{VALKEY_ADDR}", decode_responses=True)
CACHE_PREFIX = "truckguard:test:api_key:"

def get_headers(token):
    return {"Authorization": f"Bearer {token}"}

def get_or_create_data(token, endpoint, search_key, search_val, payload):
    h = get_headers(token)
    # Check if exists
    try:
        url = f"{CORE_API_URL}/data/{endpoint}"
        resp = requests.get(url, headers=h, params={search_key: search_val}).json()
        items = resp.get('data', [])
        for item in items:
            if str(item.get(search_key)) == str(search_val):
                return item['ID']
    except Exception as e:
        print(f"⚠️ Помилка при пошуку {endpoint}: {e}")
    
    # Create if not found
    print(f"➕ Створення {endpoint}: {search_val}")
    resp = requests.post(f"{CORE_API_URL}/data/{endpoint}", headers=h, json=payload).json()
    if 'ID' in resp: return resp['ID']
    if 'id' in resp: return resp['id']
    return None

def get_or_create_config(token, endpoint, name, payload):
    # Try cache first
    cache_key = f"{CACHE_PREFIX}{endpoint}:{name}"
    cached_api_key = r_client.get(cache_key)
    if cached_api_key:
        return cached_api_key

    h = get_headers(token)
    # Check Core API (using simple GET all since no search params for configs usually)
    try:
        resp = requests.get(f"{CORE_API_URL}/configs/{endpoint}", headers=h).json()
        items = resp.get('data', [])
        for item in items:
            if item.get('name') == name:
                # Some endpoints return 'camera' or 'scale' object, some return top level
                key = item.get('api_key')
                if key:
                    r_client.set(cache_key, key)
                    return key
    except Exception as e:
        print(f"⚠️ Помилка при пошуку конфігурації {endpoint}/{name}: {e}")
    
    # Create new
    print(f"🛠️ Створення девайса {name} ({endpoint})")
    resp = requests.post(f"{CORE_API_URL}/configs/{endpoint}", headers=h, json=payload).json()
    key = resp.get('api_key')
    if not key:
        print(f"❌ Помилка: сервер не повернув api_key для {name}. Response: {resp}")
        return None
    r_client.set(cache_key, key)
    return key

def setup_env(token):
    print("🏗️ Перевірка та налаштування інфраструктури (ідемпотентно)...")
    
    # 1. Довідники
    post_id = get_or_create_data(token, "posts", "name", "Main Terminal", 
                                 {"name": "Main Terminal", "description": "Auto-created test post"})
    
    get_or_create_data(token, "modes", "code", "IM", {"name": "Import", "code": "IM", "description": "Importing goods"})
    get_or_create_data(token, "modes", "code", "EK", {"name": "Export", "code": "EK", "description": "Exporting goods"})
    
    get_or_create_data(token, "vehicle-types", "code", "TRUCK", 
                       {"name": "Truck", "code": "TRUCK", "entry_price": 200, "daily_price": 50})
    
    get_or_create_data(token, "payment-types", "code", "CASH", 
                       {"name": "Cash", "code": "CASH", "is_active": True})
    
    get_or_create_data(token, "companies", "edrpou", "12345678", 
                       {"name": "TransLogistic", "edrpou": "12345678", "details": {"address":"Kyiv"}})

    # 2. Камери та Ваги
    # В цій версії ми прив'язуємось безпосередньо до Посту (бо Gates не імплементовані в Core API)
    cam_keys = {
        "IN": [
            get_or_create_config(token, "cameras", "ENTRY_Front", 
                                 {"name": "ENTRY_Front", "type": "front", "customs_post_id": post_id, 
                                   "trigger_permit_creation": True, "format": "json", "field_mapping": '{"plate":"plate"}'}),
            get_or_create_config(token, "cameras", "ENTRY_Back", 
                                 {"name": "ENTRY_Back", "type": "back", "customs_post_id": post_id, 
                                   "format": "json", "field_mapping": '{"plate":"plate"}'})
        ],
        "OUT": [
            get_or_create_config(token, "cameras", "EXIT_Front", 
                                 {"name": "EXIT_Front", "type": "front", "customs_post_id": post_id, 
                                   "format": "json", "field_mapping": '{"plate":"plate"}'}),
        ]
    }

    scale_key = get_or_create_config(token, "scales", "Main_Scale", 
                                     {"name": "Main_Scale", "customs_post_id": post_id, 
                                      "format": "json", "field_mapping": '{"weight":"weight"}'})

    return {
        "cam_keys": cam_keys,
        "scale_key": scale_key
    }

def send_cam(key, plate, cam_label=""):
    f = io.BytesIO()
    color = (random.randint(0,255), random.randint(0,255), random.randint(0,255))
    Image.new('RGB', (100, 100), color=color).save(f, 'jpeg')
    f.seek(0)
    requests.post(INGEST_CAMERA_URL, headers={'X-API-Key':key}, files={'image':('p.jpg',f)}, 
                  data={'device_id':'SIM','payload':json.dumps({"plate":plate})})
    print(f"   📸 {cam_label}: {plate}")

def send_weight(key, val, truck_plate):
    requests.post(INGEST_WEIGHT_URL, headers={'X-API-Key':key}, 
                  data={'device_id':'SCALE','payload':json.dumps({"weight":val})})
    print(f"   ⚖️  Вага для {truck_plate}: {val} kg")

def generate_plate():
    return f"{''.join(random.choices(string.ascii_uppercase, k=2))}{''.join(random.choices(string.digits, k=4))}{''.join(random.choices(string.ascii_uppercase, k=2))}"

def main():
    print(f"🚀 Запуск симуляції на {TRUCK_COUNT} вантажівок")
    
    # Auth
    try:
        login_resp = requests.post(f"{AUTH_URL}/login", json={"username":ADMIN_USER, "password":ADMIN_PASS})
        token = login_resp.json().get("token")
    except Exception as e:
        print(f"❌ Помилка авторизації: {e}")
        return

    if not token: 
        print("❌ Не отримано токен")
        return

    # Setup environment idempotently
    env = setup_env(token)
    keys = env['cam_keys']
    scale_key = env['scale_key']

    # Generate Trucks
    trucks = []
    base_start_time = time.time()
    
    print("\n⏱️  Генеруємо розклад...")
    for i in range(TRUCK_COUNT):
        p_front = generate_plate()
        p_back = generate_plate()
        weight = random.randint(15000, 40000)
        
        start_delay = i * random.uniform(5.0, 10.0)
        
        truck = {
            "id": i + 1,
            "plate_f": p_front,
            "plate_b": p_back,
            "weight": weight,
            "next_action_time": base_start_time + start_delay,
            "tasks": [
                ("CAM", keys['IN'][0], p_front, "ENTRY Front"),
                ("CAM", keys['IN'][1], p_back,  "ENTRY Back"),
                ("WAIT", random.uniform(3, 5)),
                ("CAM", keys['IN'][0], p_front, "WEIGHT Front (Scale)"),
                ("WEIGHT", scale_key, weight, ""),
                ("WAIT", random.uniform(8, 15)),
                ("CAM", keys['OUT'][0], p_front, "EXIT Front"),
            ]
        }
        trucks.append(truck)
        print(f"🚛 Truck {truck['id']}: {truck['plate_f']} (Старт через: {start_delay:.1f}s)")

    print("\n🏁 Рух почався...")
    
    unfinished_trucks = [t for t in trucks if len(t['tasks']) > 0]
    
    while unfinished_trucks:
        now = time.time()
        ready = [t for t in unfinished_trucks if t['next_action_time'] <= now]
        if not ready:
            time.sleep(0.1)
            continue

        t = random.choice(ready)
        task = t['tasks'].pop(0)
        action_type = task[0]
        
        if action_type == "WAIT":
            t['next_action_time'] = now + task[1]
        elif action_type == "CAM":
            send_cam(task[1], task[2], f"[{t['id']}] {task[3]}")
            t['next_action_time'] = now + 0.5
        elif action_type == "WEIGHT":
            send_weight(task[1], task[2], t['plate_f'])
            t['next_action_time'] = now + 1.0
            
        if len(t['tasks']) == 0:
            unfinished_trucks.remove(t)
            print(f"🎉 Truck {t['id']} finished!")

    # Verify Results
    print("\n📊 ПЕРЕВІРКА РЕЗУЛЬТАТІВ:")
    time.sleep(2) 
    h = get_headers(token)
    
    success_count = 0
    for t in trucks:
        print(f"\n🔎 Перевірка Truck {t['id']} ({t['plate_f']})...")
        try:
            r = requests.get(f"{CORE_API_URL}/permits?plate={t['plate_f']}", headers=h).json()
            if r.get('data'):
                p = r['data'][0]
                print(f"    ✅  Перепустка знайдена. ID: {p['ID']}")
                print(f"    ⚖️  Вага: {p['total_weight']} кг")
                success_count += 1
            else:
                print("   ❌ Перепустку НЕ знайдено!")
        except Exception as e:
            print(f"   ⚠️ Помилка перевірки: {e}")

    print(f"\n📈 Результат: {success_count}/{TRUCK_COUNT} успішних проїздів.")

if __name__ == "__main__":
    main()
