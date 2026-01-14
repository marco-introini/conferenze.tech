# 🚀 API Documentation - Conferenze.tech

## Base URL
```
http://localhost:8080
```

---

## 📋 Indice

- [Autenticazione](#autenticazione)
- [Endpoint Pubblici](#endpoint-pubblici)
- [Endpoint Protetti](#endpoint-protetti)
- [Modelli Dati](#modelli-dati)
- [Codici di Errore](#codici-di-errore)

---

## 🔐 Autenticazione

Gli endpoint protetti richiedono un token nell'header:

```http
Authorization: Bearer your-token-here
```

### Ottenere un token

**Login:**
```bash
POST /api/login
```

**Body:**
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

**Response 200:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    ...
  },
  "token": "your-jwt-token"
}
```

---

## 🌍 Endpoint Pubblici

### 1. Health Check

Verifica stato del server.

```http
GET /health
```

**Response 200:**
```
OK
```

---

### 2. Registrazione Utente

Crea un nuovo account utente.

```http
POST /api/register
```

**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "nickname": "johnd",         // optional
  "city": "Milano",            // optional
  "avatarUrl": "https://...",  // optional
  "bio": "Developer from IT"   // optional
}
```

**Response 201:**
```json
{
  "user": { ... },
  "token": "jwt-token"
}
```

**Errors:**
- `400` - Campi mancanti o email invalida
- `409` - Email già registrata

---

### 3. Login

Autentica un utente esistente.

```http
POST /api/login
```

**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "user": { ... },
  "token": "jwt-token"
}
```

**Errors:**
- `400` - Campi mancanti
- `401` - Credenziali non valide

---

### 4. Lista Conferenze

Ottieni tutte le conferenze disponibili.

```http
GET /api/conferences
```

**Query Parameters:**
- `search` (optional) - Cerca per nome/location
- `location` (optional) - Filtra per città

**Examples:**
```bash
GET /api/conferences
GET /api/conferences?search=golang
GET /api/conferences?location=Milano
```

**Response 200:**
```json
[
  {
    "id": "uuid",
    "title": "GopherCon Italia 2026",
    "date": "2026-09-15T10:00:00Z",
    "location": "Milano",
    "website": "https://gophercon.it",
    "latitude": 45.4642,
    "longitude": 9.1900
  },
  ...
]
```

---

### 5. Dettaglio Conferenza

Ottieni informazioni dettagliate su una conferenza specifica, inclusi i partecipanti.

```http
GET /api/conferences/{conference_id}
```

**Path Parameters:**
- `conference_id` - UUID della conferenza

**Response 200:**
```json
{
  "id": "uuid",
  "title": "GopherCon Italia 2026",
  "date": "2026-09-15T10:00:00Z",
  "location": "Milano",
  "website": "https://gophercon.it",
  "latitude": 45.4642,
  "longitude": 9.1900,
  "attendees": [
    {
      "user": {
        "id": "user-uuid",
        "nickname": "johnd",
        "city": "Milano"
      },
      "needsRide": false,
      "hasCar": true
    },
    ...
  ]
}
```

**Errors:**
- `400` - UUID non valido
- `404` - Conferenza non trovata

---

## 🔒 Endpoint Protetti

**Richiedono header:** `Authorization: Bearer <token>`

### 1. Info Utente Corrente

Ottieni informazioni sull'utente autenticato.

```http
GET /api/me
```

**Response 200:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "nickname": "johnd",
  "city": "Milano",
  "avatarUrl": "https://...",
  "bio": "Developer",
  "createdAt": "2026-01-14T10:00:00Z"
}
```

**Errors:**
- `401` - Token mancante o non valido

---

### 2. Crea Conferenza

Crea una nuova conferenza.

```http
POST /api/conferences
```

**Body:**
```json
{
  "title": "GopherCon Italia 2026",
  "date": "2026-09-15T10:00:00Z",
  "location": "Milano",
  "website": "https://gophercon.it",     // optional
  "latitude": 45.4642,                    // optional
  "longitude": 9.1900                     // optional
}
```

**Response 201:**
```json
{
  "id": "uuid",
  "title": "GopherCon Italia 2026",
  "date": "2026-09-15T10:00:00Z",
  "location": "Milano",
  ...
}
```

**Errors:**
- `400` - Dati non validi
- `401` - Non autenticato

---

### 3. Iscriviti a Conferenza

Registra l'utente corrente a una conferenza.

```http
POST /api/conferences/{conference_id}/register
```

**Path Parameters:**
- `conference_id` - UUID della conferenza

**Body:**
```json
{
  "role": "attendee",           // attendee, speaker, volunteer, organizer
  "notes": "Looking forward!",  // optional
  "needsRide": false,           // default: false
  "hasCar": true                // default: false
}
```

**Response 201:**
```json
{
  "id": "registration-uuid",
  "conferenceId": "conf-uuid",
  "conferenceTitle": "GopherCon Italia 2026",
  "conferenceDate": "2026-09-15T10:00:00Z",
  "conferenceLocation": "Milano",
  "status": "registered",
  "role": "attendee",
  "needsRide": false,
  "hasCar": true,
  "registeredAt": "2026-01-14T10:00:00Z"
}
```

**Errors:**
- `400` - Dati non validi o già iscritto
- `401` - Non autenticato
- `404` - Conferenza non trovata

---

### 4. Le Mie Iscrizioni

Ottieni tutte le conferenze a cui l'utente è iscritto.

```http
GET /api/registrations/{user_id}
```

**Path Parameters:**
- `user_id` - UUID dell'utente (deve corrispondere all'utente autenticato)

**Response 200:**
```json
[
  {
    "id": "registration-uuid",
    "conferenceId": "conf-uuid",
    "conferenceTitle": "GopherCon Italia 2026",
    "conferenceDate": "2026-09-15T10:00:00Z",
    "conferenceLocation": "Milano",
    "status": "registered",
    "role": "attendee",
    "registeredAt": "2026-01-14T10:00:00Z"
  },
  ...
]
```

**Errors:**
- `401` - Non autenticato
- `403` - Non autorizzato (non puoi vedere iscrizioni di altri utenti)

---

### 5. Info Utente (by ID)

Ottieni informazioni su un utente specifico.

```http
GET /api/users/{user_id}
```

**Path Parameters:**
- `user_id` - UUID dell'utente

**Response 200:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  ...
}
```

**Errors:**
- `401` - Non autenticato
- `404` - Utente non trovato

---

### 6. Lista Token

Ottieni tutti i token attivi dell'utente corrente.

```http
GET /api/tokens
```

**Response 200:**
```json
[
  {
    "id": "token-uuid",
    "createdAt": "2026-01-14T10:00:00Z",
    "lastUsedAt": "2026-01-14T12:30:00Z",
    "revoked": false
  },
  ...
]
```

---

### 7. Revoca Token

Revoca un token specifico (logout).

```http
POST /api/token/revoke
```

**Body:**
```json
{
  "tokenId": "uuid"
}
```

**Response 200:**
```json
{
  "message": "Token revoked successfully"
}
```

**Errors:**
- `401` - Non autenticato
- `404` - Token non trovato

---

## 📦 Modelli Dati

### User

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "nickname": "string (optional)",
  "city": "string (optional)",
  "avatarUrl": "string (optional)",
  "bio": "string (optional)",
  "createdAt": "timestamp (RFC3339)"
}
```

### Conference

```json
{
  "id": "uuid",
  "title": "string",
  "date": "timestamp (RFC3339)",
  "location": "string",
  "website": "string (optional)",
  "latitude": "float64 (optional)",
  "longitude": "float64 (optional)"
}
```

### Registration

```json
{
  "id": "uuid",
  "conferenceId": "uuid",
  "conferenceTitle": "string",
  "conferenceDate": "timestamp",
  "conferenceLocation": "string",
  "status": "registered | waitlist | cancelled | attended",
  "role": "attendee | organizer | speaker | volunteer",
  "needsRide": "boolean (optional)",
  "hasCar": "boolean (optional)",
  "registeredAt": "timestamp"
}
```

---

## ❌ Codici di Errore

### Formato Errore Standard

```json
{
  "error": "Human readable error message"
}
```

### Codici HTTP

| Codice | Significato | Quando |
|--------|-------------|--------|
| `200` | OK | Richiesta completata con successo |
| `201` | Created | Risorsa creata con successo |
| `400` | Bad Request | Dati non validi o mancanti |
| `401` | Unauthorized | Token mancante o non valido |
| `403` | Forbidden | Non hai i permessi necessari |
| `404` | Not Found | Risorsa non trovata |
| `409` | Conflict | Risorsa già esistente (es. email duplicata) |
| `500` | Internal Server Error | Errore del server |

---

## 📝 Esempi Completi

### Workflow Completo: Registrazione → Login → Crea Conferenza → Iscriviti

```bash
# 1. Registra utente
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mario@example.com",
    "password": "password123",
    "name": "Mario Rossi",
    "city": "Milano"
  }'
# Response: { "user": {...}, "token": "abc123..." }

# 2. Salva token
TOKEN="abc123..."

# 3. Crea conferenza
CONF_RESPONSE=$(curl -X POST http://localhost:8080/api/conferences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "GopherCon Italia 2026",
    "date": "2026-09-15T10:00:00Z",
    "location": "Milano",
    "website": "https://gophercon.it"
  }')

CONF_ID=$(echo $CONF_RESPONSE | jq -r '.id')

# 4. Iscriviti alla conferenza
curl -X POST http://localhost:8080/api/conferences/$CONF_ID/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "attendee",
    "notes": "Can'\''t wait!",
    "hasCar": true
  }'

# 5. Visualizza conferenza con partecipanti
curl http://localhost:8080/api/conferences/$CONF_ID | jq .

# 6. Le mie iscrizioni
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/me | jq -r '.ID' | \
  xargs -I {} curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/registrations/{}
```

---

## 🔍 Testing

### Con curl
```bash
./backend/curl/api-examples.sh
```

### Con .http files (VS Code REST Client)
```
Apri: backend/api-tests.http
```

### Unit Tests
```bash
make test
```

---

## 🛠️ CORS

Il server ha CORS abilitato con:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

---

## 📌 Note Importanti

1. **Date Format:** Tutte le date sono in formato RFC3339 (`2026-09-15T10:00:00Z`)
2. **UUIDs:** Tutti gli ID sono UUID v4
3. **Null Fields:** I campi opzionali sono stringhe semplici o `null`
4. **Password:** Le password sono hashate con bcrypt server-side
5. **Tokens:** I token sono hash SHA-256 salvati nel database

---

## 🚀 Quick Reference

### Endpoint Pubblici (No Auth)
- ✅ `GET /health`
- ✅ `POST /api/register`
- ✅ `POST /api/login`
- ✅ `GET /api/conferences`
- ✅ `GET /api/conferences/{id}`

### Endpoint Protetti (Require Auth)
- 🔒 `GET /api/me`
- 🔒 `POST /api/conferences`
- 🔒 `POST /api/conferences/{id}/register`
- 🔒 `GET /api/registrations/{user_id}`
- 🔒 `GET /api/users/{user_id}`
- 🔒 `GET /api/tokens`
- 🔒 `POST /api/token/revoke`

---

**Version:** 1.0  
**Last Updated:** 2026-01-14