.PHONY: env-up
env-up:
	cp .env.example .env
	${EDITOR} .env

.PHONY: up-build
up-build:
	docker compose up -d --build

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down -v

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: rebuild
rebuild: down up-build

.PHONY: dev-up
dev-up:
	docker compose -f docker-compose.dev.yaml up -d

dev-up-build:
	docker compose -f docker-compose.dev.yaml up -d --build

.PHONY: dev-down
dev-down:
	docker compose -f docker-compose.dev.yaml down -v

.PHONY: dev-down-soft
dev-down-soft: 
	docker compose -f docker-compose.dev.yaml down

.PHONY: dev-rebuild
dev-rebuild: dev-down dev-up-build

.PHONY: dev-restart
dev-restart: dev-down dev-up

.PHONY: dev-restart-core
dev-restart-core: 
	docker compose -f docker-compose.dev.yaml up -d ingestor adapter-worker auth core

.PHONY: dev-init
dev-init:
	bash init-garage.sh
	$(MAKE) dev-restart-core