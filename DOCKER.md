# 🚀 Guida Docker Ottimizzata

## Setup Iniziale (Prima Volta)

```bash
# Build delle immagini con cache
make dev
```

## Sviluppo Quotidiano

```bash
# Avvio veloce (usa cache, molto più rapido!)
make dev-fast

# O in background
make up
```

## Ottimizzazioni Implementate

### ✅ Backend (Go)
- **Air hot-reload**: modifica il codice e si ricarica automaticamente
- **Cache dipendenze Go**: `go mod download` eseguito solo se cambiano `go.mod`/`go.sum`
- **Volume persistente**: `/go/pkg/mod` condiviso tra rebuild

### ✅ Frontend (React/Vite)
- **Cache node_modules**: volume Docker persistente
- **HMR Vite**: hot module replacement nativo
- **npm ci**: installazione più veloce e deterministica

### ✅ Database
- **Health check**: backend parte solo quando PostgreSQL è pronto
- **Volume persistente**: dati non persi tra restart

### ✅ Build Cache
- **Layer caching**: dipendenze separate dal codice
- **Multi-stage builds**: immagini leggere per produzione
- **`.dockerignore`**: file inutili esclusi dal context

## Comandi Utili

```bash
# Sviluppo
make dev-fast          # Avvio rapido (raccomandato)
make dev               # Rebuild completo
make restart           # Restart veloce

# Gestione
make down              # Ferma container
make logs              # Vedi tutti i logs
make logs-backend      # Logs solo backend
make clean             # Pulizia leggera
make clean-all         # Pulizia completa (include volumi)

# Database
make migrate           # Applica schema
make seed              # Popola dati test
make shell-db          # Shell PostgreSQL

# Debug
make shell-backend     # Shell nel container backend
make shell-frontend    # Shell nel container frontend
```

## Tempi di Avvio

| Comando | Prima Volta | Successive |
|---------|-------------|------------|
| `make dev` | ~2-3 min | ~30-60s |
| `make dev-fast` | N/A | ~5-10s |

## Troubleshooting

### Build lento dopo modifiche
```bash
# Pulisci cache e ricostruisci
make clean-all
make dev
```

### Hot-reload non funziona (Backend)
Verifica che Air sia attivo:
```bash
make logs-backend
# Dovresti vedere: "Air is watching..."
```

### Hot-reload non funziona (Frontend)
Controlla che il volume node_modules sia montato:
```bash
docker compose ps
```

### Errori di dipendenze
```bash
# Backend: forza re-download
docker compose exec backend go mod download

# Frontend: reinstalla
docker compose exec frontend npm ci
```

### Database non si connette
```bash
# Verifica health check
docker compose ps
# Status deve essere "healthy" per db
```

## File Modificati

- ✅ `backend/Dockerfile` - Air hot-reload + cache
- ✅ `frontend/Dockerfile` - npm ci + cache
- ✅ `docker-compose.yml` - Volumi + health checks
- ✅ `backend/.air.toml` - Configurazione Air
- ✅ `backend/.dockerignore` - Ottimizza context
- ✅ `frontend/.dockerignore` - Ottimizza context
- ✅ `Makefile` - Comandi semplificati

## Note Sviluppo

- **Non serve ricostruire** le immagini dopo ogni modifica al codice
- **Air** rileva modifiche `.go` e ricompila automaticamente
- **Vite HMR** aggiorna il browser senza refresh
- I volumi persistenti velocizzano molto i riavvii