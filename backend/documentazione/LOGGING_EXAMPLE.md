# Logging Examples

Questo documento mostra esempi di output del sistema di logging migliorato.

## Formato Log

Ogni richiesta HTTP viene loggata con il seguente formato:

```
[IP:PORT] METHOD PATH | Status: CODE | Size: BYTES bytes | Duration: TIME | User: USER_ID
```

## Esempi di Log

### Richieste Pubbliche (Senza Autenticazione)

```
[127.0.0.1:54321] POST /api/register | Status: 201 | Size: 456 bytes | Duration: 123ms | User: anonymous
[127.0.0.1:54322] POST /api/login | Status: 200 | Size: 389 bytes | Duration: 89ms | User: anonymous
[192.168.1.100:60123] POST /api/login | Status: 401 | Size: 28 bytes | Duration: 12ms | User: anonymous
[127.0.0.1:54323] GET /health | Status: 200 | Size: 2 bytes | Duration: 1ms | User: anonymous
```

### Richieste Autenticate

```
[127.0.0.1:54324] GET /api/me | Status: 200 | Size: 234 bytes | Duration: 45ms | User: a1b2c3d4-e5f6-7890-abcd-ef1234567890
[127.0.0.1:54325] GET /api/conferences | Status: 200 | Size: 3456 bytes | Duration: 67ms | User: a1b2c3d4-e5f6-7890-abcd-ef1234567890
[127.0.0.1:54326] POST /api/conferences/create | Status: 201 | Size: 278 bytes | Duration: 156ms | User: b2c3d4e5-f6a7-8901-bcde-f12345678901
[192.168.1.50:61234] GET /api/registrations/abc123 | Status: 200 | Size: 1024 bytes | Duration: 34ms | User: c3d4e5f6-a7b8-9012-cdef-123456789012
```

### Errori e Casi Speciali

```
[127.0.0.1:54327] GET /api/conferences/invalid-id | Status: 400 | Size: 28 bytes | Duration: 5ms | User: d4e5f6a7-b8c9-0123-def1-234567890123
[127.0.0.1:54328] GET /api/conferences/99999999-9999-9999-9999-999999999999 | Status: 404 | Size: 22 bytes | Duration: 23ms | User: e5f6a7b8-c9d0-1234-ef12-345678901234
[192.168.1.75:62345] POST /api/conferences/create | Status: 500 | Size: 28 bytes | Duration: 234ms | User: f6a7b8c9-d0e1-2345-f123-456789012345
[127.0.0.1:54329] GET /api/me | Status: 401 | Size: 27 bytes | Duration: 3ms | User: anonymous
```

### Richieste CORS Preflight

```
[127.0.0.1:54330] OPTIONS /api/conferences | Status: 200 | Size: 0 bytes | Duration: 1ms | User: anonymous
[192.168.1.100:63456] OPTIONS /api/me | Status: 200 | Size: 0 bytes | Duration: 0ms | User: anonymous
```

## Interpretazione dei Dati

### Status Codes
- **2xx**: Successo (200 OK, 201 Created)
- **4xx**: Errori client (400 Bad Request, 401 Unauthorized, 404 Not Found)
- **5xx**: Errori server (500 Internal Server Error)

### Size
- Dimensione in bytes della risposta HTTP
- Utile per identificare risposte pesanti o anomale
- `0 bytes` tipicamente per OPTIONS o errori senza body

### Duration
- Tempo totale di elaborazione della richiesta
- Include tempo di database, processing e serializzazione
- Utile per identificare endpoint lenti o problemi di performance

### User
- `anonymous`: Richieste pubbliche o non autenticate
- `UUID`: Identificativo univoco dell'utente autenticato
- Utile per audit, debug e tracking delle attività utente

## Best Practices

1. **Monitoring Performance**: Monitora le richieste con duration > 500ms
2. **Security Audit**: Analizza pattern di `401` per tentativi di accesso non autorizzato
3. **Traffic Analysis**: Usa gli IP per identificare pattern di utilizzo
4. **User Activity**: Traccia le azioni degli utenti tramite il campo User
5. **Error Detection**: Filtra per status >= 500 per identificare problemi server