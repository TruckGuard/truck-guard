#!/bin/bash

# Завантажуємо змінні з .env
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
else
    echo "Помилка: Файл .env не знайдено!"
    exit 1
fi

COMPOSE_CMD="docker compose -f docker-compose.dev.yaml"
GARAGE_EXEC="$COMPOSE_CMD exec -T garage /garage"

echo "--- 1. Отримання Node ID ---"
NODE_ID=$($GARAGE_EXEC status | awk '/ID/ {getline; print $1}')
if [ -z "$NODE_ID" ]; then
    echo "Помилка: Не вдалося отримати Node ID. Перевір, чи запущений контейнер garage."
    exit 1
fi
echo "Знайдено Node ID: $NODE_ID"

echo "--- 2. Налаштування Layout ---"
$GARAGE_EXEC layout assign $NODE_ID -z garage -c 1G -t garage || true

# Отримання поточної версії для layout apply
CURRENT_VERSION=$($GARAGE_EXEC layout show | grep "Staged edit" | awk '{print $4}' | sed 's/v//')
if [ -z "$CURRENT_VERSION" ]; then
    # Якщо немає staged edit, спробуємо отримати останню застосовану версію
    CURRENT_VERSION=$($GARAGE_EXEC layout apply --version 0 2>&1 | grep "expected" | awk '{print $NF}')
    if [ -z "$CURRENT_VERSION" ]; then CURRENT_VERSION=1; fi
fi

echo "Застосування Layout (версія $CURRENT_VERSION)..."
$GARAGE_EXEC layout apply --version "$CURRENT_VERSION" || true

echo "--- 3. Створення/отримання ключа ---"
# Очікуємо готовності лейауту
echo "Очікування готовності кластера..."
until $GARAGE_EXEC key list; do
    sleep 2
done

EXISTING_KEY=$($GARAGE_EXEC key list | grep 'local-key' | awk '{print $1}')

if [ -z "$EXISTING_KEY" ]; then
    echo "Створюємо новий ключ..."
    KEY_INFO=$($GARAGE_EXEC key create local-key)
    ACCESS_KEY=$(echo "$KEY_INFO" | grep "Key ID" | awk '{print $3}')
    SECRET_KEY=$(echo "$KEY_INFO" | grep "Secret key" | awk '{print $3}')
else
    echo "Ключ уже існує: $EXISTING_KEY"
    ACCESS_KEY=$EXISTING_KEY
    SECRET_KEY="ДИВИСЬ_НИЖЧЕ (вже створено)"
fi

echo "--- 4. Створення бакета ($STORAGE_BUCKET_NAME) ---"
$GARAGE_EXEC bucket create $STORAGE_BUCKET_NAME || true

echo "--- 5. Надання прав доступу для $ACCESS_KEY ---"
$GARAGE_EXEC bucket allow $STORAGE_BUCKET_NAME --key $ACCESS_KEY --read --write
$GARAGE_EXEC bucket website --allow $STORAGE_BUCKET_NAME

echo ""
echo "===================================================="
echo "ГОРЯЧО! ГАРАЖ ГОТОВИЙ."
echo "Тобі потрібно оновити .env файл цими значеннями:"
echo "===================================================="
echo "STORAGE_ACCESS_KEY=$ACCESS_KEY"
if [ "$SECRET_KEY" != "ДИВИСЬ_НИЖЧЕ (вже створено)" ]; then
  echo "STORAGE_SECRET_KEY=$SECRET_KEY"
fi
echo "===================================================="
echo "Після оновлення .env зроби: docker compose up -d"

