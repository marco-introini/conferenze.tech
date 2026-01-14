.PHONY: dev dev-fast up down restart clean logs seed migrate build-prod deploy prune test test-api test-verbose test-coverage

# Avvio rapido senza rebuild (usa cache)
dev-fast:
	docker compose up

# Primo avvio o dopo modifiche ai Dockerfile
dev:
	docker compose up --build

# Avvia in background
up:
	docker compose up -d

# Ferma i container
down:
	docker compose down

# Restart veloce
restart:
	docker compose restart

# Pulizia completa (container, network, volumes anonimi)
clean:
	docker compose down --remove-orphans

# Pulizia profonda (include volumi named)
clean-all:
	docker compose down -v --remove-orphans
	docker system prune -f

# Visualizza logs
logs:
	docker compose logs -f

# Logs specifici
logs-backend:
	docker compose logs -f backend

logs-frontend:
	docker compose logs -f frontend

logs-db:
	docker compose logs -f db

# Seeding database
seed:
	docker compose exec backend go run ./cmd/seeder

# Migrate database
migrate:
	docker compose exec -T db psql -U user -d conferenzetech -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	docker compose exec -T db psql -U user -d conferenzetech < backend/schema.sql

# Build produzione
build-prod:
	docker build --target production -t my-app-backend ./backend
	docker build --target production -t my-app-frontend ./frontend

# Rimuove immagini e cache inutilizzate
prune:
	docker system prune -af --volumes

# Deploy
deploy:
	ssh user@minipc-ip "cd /app && git pull && docker compose up -d --build"

# Shell nei container
shell-backend:
	docker compose exec backend sh

shell-frontend:
	docker compose exec frontend sh

shell-db:
	docker compose exec db psql -U user -d conferenzetech

# Testing
test:
	cd backend && go test -v ./...

test-verbose:
	cd backend && go test -v -race -count=1 ./...

test-coverage:
	cd backend && go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

test-api:
	./backend/curl/api-examples.sh

test-http:
	@echo "Apri backend/api-tests.http in VS Code con REST Client extension"
	@echo "Oppure usa IntelliJ IDEA HTTP Client"

# Run unit tests nel container
test-docker:
	docker compose exec backend go test -v ./...
