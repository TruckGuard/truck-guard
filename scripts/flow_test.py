import requests
import json
import time
import io
import os
import random
from PIL import Image
from datetime import datetime

# Конфігурація
BASE_URL = os.getenv("BASE_URL", "http://localhost")
AUTH_URL = f"{BASE_URL}/auth"
CORE_API_URL = f"{BASE_URL}/api"
INGEST_CAMERA_URL = f"{BASE_URL}/ingest/camera"
INGEST_WEIGHT_URL = f"{BASE_URL}/ingest/weight"

ADMIN_USER = os.getenv("ADMIN_USER", "admin")
# Пріоритет на пароль з .env, інакше дефолтний
ADMIN_PASS = os.getenv("ADMIN_DEFAULT_PASSWORD", "secret123")

# Дані фури
TRUCK = {"f": "BC7777EX", "b": "BC7777EE", "w": 32500}

def get_headers(token):
    return {"Authorization": f"Bearer {token}"}

def cleanup(token):
    """Повна очистка перед тестом"""
    print("🧹 Повна очистка конфігурацій та довідників...")
    h = get_headers(token)
    
    # 1. Configs
    for ep in ["flows", "gates", "scales", "cameras"]:
        try:
            resp = requests.get(f"{CORE_API_URL}/configs/{ep}", headers=h)
            items = resp.json().get('data', []) if isinstance(resp.json(), dict) else resp.json()
            if items:
                for i in items:
                    requests.delete(f"{CORE_API_URL}/configs/{ep}/{i['ID']}", headers=h)
        except: pass

    # 2. Data entities
    for ep in ["posts", "modes", "vehicle-types", "payment-types", "companies"]:
        try:
            resp = requests.get(f"{CORE_API_URL}/data/{ep}", headers=h)
            items = resp.json().get('data', []) if isinstance(resp.json(), dict) else resp.json()
            if items:
                for i in items:
                    requests.delete(f"{CORE_API_URL}/data/{ep}/{i['ID']}", headers=h)
        except: pass
    
    print("✨ Система чиста.")

def setup_data(token):
    """Створення довідників"""
    h = get_headers(token)
    print("📦 Наповнення довідників...")
    
    # 1. Митний пост
    post_id = requests.post(f"{CORE_API_URL}/data/posts", headers=h, 
                            json={"name": "Головний Термінал", "description": "Автоматично створений пост"}).json()['ID']
    
    # 2. Режими
    requests.post(f"{CORE_API_URL}/data/modes", headers=h, json={"name": "Імпорт", "code": "IM", "description": "Ввезення"})
    requests.post(f"{CORE_API_URL}/data/modes", headers=h, json={"name": "Експорт", "code": "EK", "description": "Вивезення"})
    
    # 3. Типи ТЗ
    requests.post(f"{CORE_API_URL}/data/vehicle-types", headers=h, 
                  json={"name": "Вантажівка", "code": "TRUCK", "entry_price": 200, "daily_price": 50})
    
    # 4. Компанії
    requests.post(f"{CORE_API_URL}/data/companies", headers=h, 
                  json={"name": "ТрансЛогістик", "edrpou": "12345678", "details": '{"address":"Kyiv"}'})
    
    # 5. Типи оплати
    requests.post(f"{CORE_API_URL}/data/payment-types", headers=h, 
                  json={"name": "Готівка", "code": "CASH", "is_active": True})

    return post_id

def setup_infra(token, post_id):
    """Налаштування обладнання"""
    h = get_headers(token)
    print(f"🏗️ Створення обладнання для посту {post_id}...")
    
    # Створюємо камери
    # Камера на в'їзд (Front) - вона тригерить створення перепустки
    cam_in_f = requests.post(f"{CORE_API_URL}/configs/cameras", headers=h, 
                             json={"name": "ENTRY_Front", "type": "front", "customs_post_id": post_id, 
                                   "match_permit": True, "format": "json", "field_mapping": '{"plate":"plate"}'}).json()
    
    # Камера на в'їзд (Back)
                                   "match_permit": True, "format": "json", "field_mapping": '{"plate":"plate"}'}).json()
    
    # Ваги
    scale = requests.post(f"{CORE_API_URL}/configs/scales", headers=h, 
                          json={"name": "Main_Scale", "customs_post_id": post_id, 
                                "match_permit": True, "format": "json", "field_mapping": '{"weight":"weight"}'}).json()
    
    return {
        "cam_in_f": cam_in_f['api_key'],
        "cam_in_b": cam_in_b['api_key'],
        "scale_key": scale['api_key']
    }

def send_cam(key, plate, cam_label=""):
    f = io.BytesIO()
    # Створюємо унікальну картинку кожного разу (шляхом випадкового кольору)
    color = (random.randint(0,255), random.randint(0,255), random.randint(0,255))
    Image.new('RGB', (100, 100), color=color).save(f, 'jpeg')
    f.seek(0)
    
    print(f"   📸 Відправка фото ({cam_label}): {plate}")
    requests.post(INGEST_CAMERA_URL, headers={'X-API-Key':key}, files={'image':('p.jpg',f)}, 
                  data={'device_id':'SIM','payload':json.dumps({"plate":plate})})

def send_weight(key, val):
    print(f"   ⚖️  Зчитування ваги: {val} kg")
    requests.post(INGEST_WEIGHT_URL, headers={'X-API-Key':key}, 
                  data={'device_id':'SCALE','payload':json.dumps({"weight":val})})

def main():
    print("🚀 Початок Flow-тесту...")
    
    # 0. Авторизація
    login_resp = requests.post(f"{AUTH_URL}/login", json={"username":ADMIN_USER, "password":ADMIN_PASS})
    if login_resp.status_code != 200:
        print(f"❌ Помилка входу: {login_resp.status_code} {login_resp.text}")
        return
    token = login_resp.json().get("session_id")

    # 1. Підготовка
    cleanup(token)
    post_id = setup_data(token)
    keys = setup_infra(token, post_id)
    
    # 2. Симуляція
    print("\n--- 🚛 СЦЕНАРІЙ: ЗАЇЗД ТА ЗВАЖУВАННЯ ---")
    
    # ПЕРЕДНЯ КАМЕРА (Створює перепустку)
    send_cam(keys['cam_in_f'], TRUCK['f'], "ENTRY Front")
    time.sleep(1)
    
    # ЗАДНЯ КАМЕРА (Доповнює перепустку)
    send_cam(keys['cam_in_b'], TRUCK['b'], "ENTRY Back")
    time.sleep(1)
    
    # ВАГА
    send_weight(keys['scale_key'], TRUCK['w'])

    # 3. Перевірка
    print("\n📊 ПЕРЕВІРКА РЕЗУЛЬТАТІВ...")
    time.sleep(3) # Час на асинхронну обробку
    
    h = get_headers(token)
    r = requests.get(f"{CORE_API_URL}/permits?plate={TRUCK['f']}", headers=h).json()
    
    if r.get('data') and len(r['data']) > 0:
        p = r['data'][0]
        print(f"✅ Перепустку знайдено! [ID: {p['ID']}]")
        print(f"   🚚 Номери: {p['plate_front']} / {p['plate_back']}")
        print(f"   ⚖️  Вага: {p['total_weight']} кг")
        # Перевірка подій
        events_resp = requests.get(f"{CORE_API_URL}/events/plate?permit_id={p['ID']}", headers=h).json()
        events_count = len(events_resp.get('data', []))
        print(f"   📸 Кількість Plate-подій: {events_count}")
        
        if p['plate_front'] == TRUCK['f'] and p['total_weight'] == TRUCK['w']:
            print("\n🏆 ТЕСТ УСПІШНО ПРОЙДЕНО!")
        else:
            print("\n❌ Дані перепустки не співпадають.")
    else:
        print(f"❌ Перепустку для {TRUCK['f']} не знайдено.")

if __name__ == "__main__":
    main()