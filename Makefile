.PHONY: dev build deploy clean seed

# Ferma e rimuove container, reti e immagini orfane
clean:
	docker compose down --remove-orphans

dev:
	docker compose down
	docker compose up --build

seed:
	docker compose exec backend go run ./cmd/seeder

migrate:
	docker compose exec -T db psql -U user -d conferenzetech < backend/schema.sql

build-prod:
	docker build --target production -t my-app-backend ./backend

deploy:
	ssh user@minipc-ip "cd /app && git pull && docker compose up -d --build"