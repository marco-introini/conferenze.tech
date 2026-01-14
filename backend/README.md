# Backend per conferenze.tech

## 🧪 Testing API

Il progetto supporta **4 approcci** per testare le API:

### 1. Test Automatici (Raccomandato)
```bash
make test              # Unit tests con httptest
make test-verbose      # Con race detector
make test-coverage     # HTML coverage report
```

### 2. File .http (VS Code REST Client)
```bash
# Apri backend/api-tests.http in VS Code
# Clicca "Send Request" su ogni chiamata
```

### 3. Script curl
```bash
./backend/curl/api-examples.sh  # Script completo interattivo
```

### 4. curl Manuale
```bash
# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"marco@marcointroini.it","password":"password"}'

# Usa il token per chiamate protette
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer <token>"
```

📖 **Guide complete**: 
- [TESTING.md](./TESTING.md) - Testing guide
- [API.md](./API.md) - API documentation

## 📚 Documentazione API

Vedi [API.md](./API.md) per:
- Lista completa endpoint
- Endpoint pubblici vs protetti
- Modelli dati
- Esempi di utilizzo
- Codici di errore

### Quick Reference

**Endpoint Pubblici (No Auth):**
- ✅ `GET /api/conferences` - Lista conferenze
- ✅ `GET /api/conferences/{id}` - Dettaglio conferenza
- ✅ `POST /api/register` - Registrazione
- ✅ `POST /api/login` - Login

**Endpoint Protetti (Auth Required):**
- 🔒 `POST /api/conferences/create` - Crea conferenza
- 🔒 `POST /api/conferences/{id}/register` - Iscriviti
- 🔒 `GET /api/me` - Info utente
- 🔒 `GET /api/registrations/{user_id}` - Le mie iscrizioni

## Struttura

```
backend/
├── db/
│   ├── models.go    # Struct per User, Conference, ConferenceRegistration
│   ├── querier.go   # Interfaccia Querier con parametri per le query
│   ├── queries.go   # Implementazione delle query
│   └── db.go        # Connessione DB e gestione transazioni
```

## Setup

1. Crea un progetto Go:

```bash
mkdir backend
cd backend
go mod init conferenze.tech/backend
```

2. Aggiungi le dipendenze:

```bash
go get github.com/lib/pq
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt
```

3. Copia i file dalla directory `frontend-conferenze/backend/db/` nella directory `backend/db/`

4. Configura le variabili d'ambiente o crea un file `.env`:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=conferenze
export DB_SSLMODE=disable
```

5. Esegui le migration:

```bash
psql -h localhost -U postgres -d conferenze -f schema.sql
```

## Utilizzo

```go
package main

import (
    "log"
    "fmt"
    "os"
    "conferenze.tech/backend/db"
)

func main() {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"),
    )

    database, err := db.New(dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    // Esempio: lista conferenze
    conferences, err := database.ListUpcomingConferences(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for _, c := range conferences {
        fmt.Printf("%s - %s\n", c.Title, c.Location)
    }
}
```

## Query disponibili

### Users
- `CreateUser` - Crea un nuovo utente
- `GetUserByID` - Ottieni utente per ID
- `GetUserByEmail` - Ottieni utente per email
- `UpdateUser` - Aggiorna profilo utente
- `UpdateUserPassword` - Aggiorna password

### Conferences
- `CreateConference` - Crea una conferenza
- `GetConferenceByID` - Ottieni conferenza per ID
- `ListConferences` - Lista tutte le conferenze
- `ListUpcomingConferences` - Lista conferenze future
- `ListConferencesByLocation` - Filtra per città
- `UpdateConference` - Aggiorna conferenza
- `DeleteConference` - Elimina conferenza

### Registrations
- `RegisterUserToConference` - Registra utente a conferenza
- `GetRegistration` - Ottieni registrazione specifica
- `GetRegistrationsByConference` - Lista registrazioni per conferenza
- `GetRegistrationsByUser` - Lista registrazioni per utente
- `UpdateRegistrationStatus` - Aggiorna stato registrazione
- `CancelRegistration` - Cancella registrazione
- `DeleteRegistration` - Elimina registrazione

### Carpooling
- `ListUsersNeedingRide` - Lista utenti che cercano passaggio
- `ListUsersOfferingRide` - Lista utenti che offrono passaggio

### Stats
- `GetConferenceStats` - Statistiche conferenza
