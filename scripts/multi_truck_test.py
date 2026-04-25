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
RELIABILITY = float(os.getenv("RELIABILITY", "0.95")) # 95% успішних спрацювань за замовчуванням

# Valkey Client for caching
r_client = redis.Redis.from_url(f"redis://{VALKEY_ADDR}", decode_responses=True)
CACHE_PREFIX = "truckguard:test:api_key:"

def get_headers(token):
    return {"Authorization": f"Bearer {token}"}

def get_or_create_data(token, endpoint, search_key, search_val, payload):
    h = get_headers(token)
    try:
        url = f"{CORE_API_URL}/data/{endpoint}"
        resp = requests.get(url, headers=h, params={search_key: search_val}).json()
        items = resp.get('data', [])
        for item in items:
            # Handle possible integer ID vs string search
            if str(item.get(search_key)) == str(search_val):
                return item['ID']
    except Exception as e:
        print(f"⚠️ Помилка при пошуку {endpoint}: {e}")
    
    print(f"➕ Створення {endpoint}: {search_val}")
    resp = requests.post(f"{CORE_API_URL}/data/{endpoint}", headers=h, json=payload).json()
    return resp.get('ID') or resp.get('id')

def get_or_create_config(token, endpoint, name, payload):
    cache_key = f"{CACHE_PREFIX}{endpoint}:{name}"
    cached_api_key = r_client.get(cache_key)
    if cached_api_key:
        return cached_api_key

    h = get_headers(token)
    try:
        resp = requests.get(f"{CORE_API_URL}/configs/{endpoint}", headers=h).json()
        items = resp.get('data', [])
        for item in items:
            if item.get('name') == name:
                key = item.get('api_key')
                if key:
                    r_client.set(cache_key, key)
                    return key
    except Exception as e:
        print(f"⚠️ Помилка при пошуку конфігурації {endpoint}/{name}: {e}")
    
    print(f"🛠️ Створення девайса {name} ({endpoint})")
    resp = requests.post(f"{CORE_API_URL}/configs/{endpoint}", headers=h, json=payload).json()
    key = resp.get('api_key')
    if not key:
        print(f"❌ Помилка: сервер не повернув api_key для {name}. Response: {resp}")
        return None
    r_client.set(cache_key, key)
    return key

def setup_env(token):
    print("🏗️ Перевірка та налаштування інфраструктури...")
    
    post_id = get_or_create_data(token, "posts", "name", "Main Terminal", 
                                 {"name": "Main Terminal", "description": "Auto-created test post"})
    
    get_or_create_data(token, "modes", "code", "IM", {"name": "Import", "code": "IM"})
    get_or_create_data(token, "modes", "code", "EK", {"name": "Export", "code": "EK"})
    get_or_create_data(token, "vehicle-types", "code", "TRUCK", {"name": "Truck", "code": "TRUCK", "entry_price": 200})
    get_or_create_data(token, "payment-types", "code", "CASH", {"name": "Cash", "code": "CASH", "is_active": True})
    get_or_create_data(token, "companies", "edrpou", "12345678", {"name": "Test Co", "edrpou": "12345678"})

    cam_keys = {
        "IN_F": get_or_create_config(token, "cameras", "ENTRY_Front", 
                                 {"name": "ENTRY_Front", "type": "front", "customs_post_id": post_id, 
                                   "match_permit": True, "format": "json", "field_mapping": '{"plate":"plate"}'}),
        "IN_B": get_or_create_config(token, "cameras", "ENTRY_Back", 
                                 {"name": "ENTRY_Back", "type": "back", "customs_post_id": post_id, 
                                   "match_permit": True, "format": "json", "field_mapping": '{"plate":"plate"}'}),
        "OUT_F": get_or_create_config(token, "cameras", "EXIT_Front", 
                                 {"name": "EXIT_Front", "type": "front", "customs_post_id": post_id, 
                                   "match_permit": False, "format": "json", "field_mapping": '{"plate":"plate"}'}),
    }

    scale_key = get_or_create_config(token, "scales", "Main_Scale", 
                                     {"name": "Main_Scale", "customs_post_id": post_id, 
                                      "match_permit": True, "format": "json", "field_mapping": '{"weight":"weight"}'})

    return {"cam_keys": cam_keys, "scale_key": scale_key}

def send_cam(key, plate, cam_label=""):
    if random.random() > RELIABILITY:
        print(f"   ⚠️ [Skip] {cam_label}: {plate} (Simulated Failure)")
        return
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

def introduce_drift(plate, probability=0.4):
    """Штучно вносить OCR помилки у номер з заданою ймовірністю (дрифт)"""
    if random.random() > probability:
        return plate
        
    plate_list = list(plate)
    # Типові помилки OCR
    ocr_mistakes = {
        '0': 'O', 'O': '0',
        '1': 'I', 'I': '1',
        '8': 'B', 'B': '8',
        'A': '4', '4': 'A',
        'C': 'G', 'G': 'C'
    }
    
    # Вибираємо випадкову позицію для заміни
    pos = random.randint(0, len(plate) - 1)
    char = plate_list[pos]
    
    # Або типова помилка, або випадкова зміна
    if char in ocr_mistakes and random.random() < 0.7:
        plate_list[pos] = ocr_mistakes[char]
    else:
        if char.isdigit():
            plate_list[pos] = random.choice(string.digits)
        else:
            plate_list[pos] = random.choice(string.ascii_uppercase)
            
    return "".join(plate_list)

def main():
    print(f"🚀 Simulation: {TRUCK_COUNT} trucks | Reliability: {RELIABILITY*100:.0f}% | Fuzzy Matching Enabled")
    
    login_resp = requests.post(f"{AUTH_URL}/login", json={"username":ADMIN_USER, "password":ADMIN_PASS})
    token = login_resp.json().get("session_id")
    if not token: return print("❌ Auth failed")

    env = setup_env(token)
    
    parking_pool = []
    finished_trucks = []
    
    # Створюємо фури, які будуть заїжджати
    incoming_queue = []
    for i in range(TRUCK_COUNT):
        incoming_queue.append({
            "id": i+1,
            "plate_f": generate_plate(),
            "plate_b": generate_plate(),
            "weight": random.randint(15000, 40000),
            "start_time": time.time() + (i * random.uniform(5, 10))
        })

    print(f"\n🚙 Початок заїзду {TRUCK_COUNT} фур...")
    
    while incoming_queue or parking_pool:
        now = time.time()
        
        # 1. Обробка заїзду (Entry) - Групована
        if incoming_queue and incoming_queue[0]["start_time"] <= now:
            t = incoming_queue.pop(0)
            print(f"\n➡️ [Truck {t['id']}] Заїжджає: {t['plate_f']}")
            
            # Entry Front - створює перепустку (без дрифту)
            send_cam(env['cam_keys']['IN_F'], t['plate_f'], f"Entry Front")
            time.sleep(0.5)
            
            # Entry Back - може містити помилки OCR
            drifted_plate_b = introduce_drift(t['plate_b'])
            if drifted_plate_b != t['plate_b']:
                print(f"   ⚠️ [Drift] Entry Back розпізнано з помилкою: {drifted_plate_b} (Оригінал: {t['plate_b']})")
            send_cam(env['cam_keys']['IN_B'], drifted_plate_b, f"Entry Back")
            time.sleep(0.5)
            
            send_weight(env['scale_key'], t['weight'], t['plate_f'])
            
            # Фура стає на парковку
            t['exit_time'] = now + random.uniform(10, 20)
            parking_pool.append(t)
            print(f"🅿️ [Truck {t['id']}] На парковці")

        # 2. Обробка виїзду (Exit) - Рандомізована
        ready_to_exit = [t for t in parking_pool if t['exit_time'] <= now]
        if ready_to_exit:
            t = random.choice(ready_to_exit)
            parking_pool.remove(t)
            
            print(f"\n⬅️ [Truck {t['id']}] Починає виїзд: {t['plate_f']}")
            
            # Exit Front - може містити помилки OCR
            drifted_plate_f = introduce_drift(t['plate_f'])
            if drifted_plate_f != t['plate_f']:
                print(f"   ⚠️ [Drift] Exit Front розпізнано з помилкою: {drifted_plate_f} (Оригінал: {t['plate_f']}) -> Очікуйте повідомлення оператору!")
                
            send_cam(env['cam_keys']['OUT_F'], drifted_plate_f, f"Exit Front")
            
            finished_trucks.append(t)
            print(f"🎉 [Truck {t['id']}] Виїхав!")

        time.sleep(0.1)

    print("\n📊 ПЕРЕВІРКА РЕЗУЛЬТАТІВ:")
    time.sleep(2)
    h = get_headers(token)
    success = 0
    for t in finished_trucks:
        try:
            r = requests.get(f"{CORE_API_URL}/permits?plate={t['plate_f']}", headers=h).json()
            if r.get('data'):
                p = r['data'][0]
                print(f"✅ Permit {p['ID']} found for {t['plate_f']} (Weight: {p['total_weight']:.0f})")
                success += 1
            else:
                print(f"❌ Permit NOT found for {t['plate_f']}")
        except: pass
    
    print(f"\n🏆 Підсумок: {success}/{TRUCK_COUNT} перепусток знайдено.")

if __name__ == "__main__":
    main()
