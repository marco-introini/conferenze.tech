.PHONY: dev build deploy clean

# Ferma e rimuove container, reti e immagini orfane
clean:
	docker compose down --remove-orphans

dev:
	docker compose down
	docker compose up --build

build-prod:
	docker build --target production -t my-app-backend ./backend

deploy:
	ssh user@minipc-ip "cd /app && git pull && docker compose up -d --build"