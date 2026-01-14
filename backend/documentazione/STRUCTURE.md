# Backend Structure

Questo documento descrive la struttura del backend dell'applicazione conferenze.tech.

## Organizzazione dei File

Il codice è stato organizzato in moduli separati per migliorare la manutenibilità e la leggibilità:

### File Principali

- **`main.go`** - Entry point dell'applicazione, gestisce la connessione al database e l'avvio del server

- **`server.go`** - Definizione della struct `Server` e configurazione delle route HTTP

- **`types.go`** - Definizione di tutti i tipi di dati, request e response utilizzati nell'applicazione

- **`utils.go`** - Funzioni di utility per la conversione tra tipi SQL nullable e puntatori Go

### Moduli Funzionali

#### Autenticazione
- **`auth.go`** - Gestione dell'autenticazione, generazione e hashing dei token, middleware di autenticazione

#### Handlers HTTP

- **`handlers_user.go`** - Operazioni sugli utenti:
  - `Register` - Registrazione nuovo utente
  - `Login` - Autenticazione utente
  - `GetMe` - Recupero informazioni utente per ID
  - `GetMeFromToken` - Recupero informazioni utente dal token

- **`handlers_conference.go`** - Operazioni sulle conferenze:
  - `ListConferences` - Elenco di tutte le conferenze
  - `GetConference` - Dettagli di una conferenza con partecipanti
  - `CreateConference` - Creazione nuova conferenza

- **`handlers_registration.go`** - Gestione iscrizioni:
  - `RegisterToConference` - Iscrizione utente a una conferenza
  - `GetUserRegistrations` - Elenco iscrizioni di un utente

- **`handlers_token.go`** - Gestione token:
  - `GetTokens` - Elenco token di un utente
  - `RevokeToken` - Revoca di un token

#### Middleware
- **`logging.go`** - Middleware di logging:
  - `loggingMiddleware` - Logging dettagliato delle richieste HTTP
    - Metodo HTTP e percorso
    - Status code della risposta
    - Dimensione della risposta in bytes
    - Durata della richiesta
    - Utente autenticato (se presente)
- **`security.go`** - Middleware di sicurezza:
  - `corsMiddleware` - Gestione CORS

## Flusso delle Richieste

1. **Richiesta HTTP** → `loggingMiddleware` → `corsMiddleware`
2. **Route pubbliche** (`/api/register`, `/api/login`) → Handler diretto
3. **Route protette** (`/api/*`) → `authMiddleware` → Handler specifico

## Dipendenze

- `github.com/google/uuid` - Gestione UUID
- `github.com/jackc/pgx/v5` - Driver PostgreSQL
- `github.com/lib/pq` - Driver PostgreSQL aggiuntivo
- Database queries generate con `sqlc` in `./db`

## Compilazione ed Esecuzione

```bash
# Compilazione
go build

# Esecuzione
./backend

# Variabili d'ambiente
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
PORT=8080
```

## Logging

Il middleware di logging cattura e registra informazioni dettagliate per ogni richiesta:

```
[127.0.0.1:54321] GET /api/conferences | Status: 200 | Size: 1234 bytes | Duration: 45ms | User: a1b2c3d4-e5f6-7890-abcd-ef1234567890
[127.0.0.1:54322] POST /api/login | Status: 401 | Size: 28 bytes | Duration: 12ms | User: anonymous
[127.0.0.1:54323] POST /api/conferences/create | Status: 201 | Size: 256 bytes | Duration: 78ms | User: a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

- **Request info**: Indirizzo client, metodo HTTP, percorso
- **Response info**: Status code HTTP, dimensione risposta
- **Performance**: Durata dell'elaborazione della richiesta
- **Security**: ID utente autenticato o "anonymous" per richieste pubbliche

## Note

- Tutti gli handler utilizzano un timeout di 5 secondi per le operazioni sul database
- Le password sono hashate usando bcrypt
- I token sono memorizzati come hash SHA-256 nel database
- Le risposte sono sempre in formato JSON
- Il logging cattura automaticamente tutte le richieste con dettagli completi