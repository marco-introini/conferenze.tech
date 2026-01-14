# Guida Pratica all'Utilizzo dei Log

Questa guida mostra come utilizzare i log migliorati per debugging, monitoring e analisi.

## 🔍 Filtraggio e Ricerca Log

### 1. Trovare Richieste Lente

Identifica richieste che impiegano più di 100ms:

```bash
# Filtra per duration > 100ms
grep "Duration:" backend.log | grep -E "[0-9]{3,}ms"

# Oppure con awk
awk '/Duration:/ && /[0-9]{3,}ms/ {print}' backend.log
```

### 2. Monitorare Errori

Trova tutti gli errori 4xx e 5xx:

```bash
# Tutti gli errori client (4xx)
grep "Status: 4" backend.log

# Tutti gli errori server (5xx)
grep "Status: 5" backend.log

# Errori specifici
grep "Status: 401" backend.log  # Unauthorized
grep "Status: 404" backend.log  # Not Found
grep "Status: 500" backend.log  # Internal Server Error
```

### 3. Tracciare Attività Utente

Segui tutte le azioni di un utente specifico:

```bash
USER_ID="a1b2c3d4-e5f6-7890-abcd-ef1234567890"
grep "User: $USER_ID" backend.log
```

### 4. Analizzare Endpoint Specifici

```bash
# Tutte le chiamate a /api/conferences
grep "/api/conferences" backend.log

# Solo GET
grep "GET /api/conferences" backend.log

# Solo POST con errori
grep "POST /api/conferences" backend.log | grep "Status: [45]"
```

## 📊 Analisi Statistiche

### 1. Conteggio Richieste per Endpoint

```bash
# Conta le richieste per ogni endpoint
awk '{print $3, $4}' backend.log | sort | uniq -c | sort -nr
```

### 2. Statistiche Status Code

```bash
# Conta per status code
grep -o "Status: [0-9]*" backend.log | sort | uniq -c

# Esempio output:
# 1547 Status: 200
#  234 Status: 201
#   45 Status: 401
#   12 Status: 404
#    3 Status: 500
```

### 3. Media Tempi di Risposta

```bash
# Estrae i tempi in ms e calcola media
grep -o "Duration: [0-9]*ms" backend.log | \
  grep -o "[0-9]*" | \
  awk '{sum+=$1; n++} END {print "Media:", sum/n "ms"}'
```

### 4. Top 10 Richieste più Lente

```bash
grep "Duration:" backend.log | \
  sort -t: -k6 -rn | \
  head -10
```

## 🐛 Scenari di Debugging

### Scenario 1: Utente Non Riesce a Fare Login

```bash
# Cerca tentativi di login
grep "POST /api/login" backend.log | tail -20

# Cerca errori 401 (credenziali errate)
grep "POST /api/login" backend.log | grep "Status: 401"

# Cerca da uno specifico IP
grep "192.168.1.100.*POST /api/login" backend.log
```

**Interpretazione:**
- Multiple 401: Password errata o utente inesistente
- 500: Problema server (controllare database)
- Nessun log: Richiesta non arriva al server

### Scenario 2: API Lenta

```bash
# Trova endpoint più lenti
grep "Duration:" backend.log | \
  awk '{print $4, $(NF-4)}' | \
  sort -t: -k2 -rn | \
  head -20
```

**Interpretazione:**
- > 1000ms: Possibile problema database o query N+1
- Aumenta nel tempo: Memory leak o cache piena
- Specifico endpoint: Ottimizza query o logica

### Scenario 3: Accessi Non Autorizzati

```bash
# Trova tentativi di accesso con token invalido
grep "Status: 401" backend.log | grep "User: anonymous"

# Raggruppa per IP sospetti
grep "Status: 401" backend.log | \
  awk '{print $1}' | \
  sort | uniq -c | sort -rn
```

**Azioni:**
- > 10 tentativi da stesso IP: Possibile attacco, considera IP ban
- Pattern regolare: Bot o scraper

### Scenario 4: Monitorare Deploy

```bash
# Dopo un deploy, monitora errori
tail -f backend.log | grep "Status: [45]"

# Confronta tempi prima e dopo deploy
# Prima
grep "Duration:" backend.log.old | \
  awk '{sum+=$(NF-4)} END {print sum/NR}'

# Dopo
grep "Duration:" backend.log | \
  awk '{sum+=$(NF-4)} END {print sum/NR}'
```

## 📈 Dashboard e Monitoring

### Script di Monitoring Real-Time

```bash
#!/bin/bash
# monitor.sh - Dashboard real-time dei log

watch -n 5 '
echo "=== Conferenze.tech - Dashboard ==="
echo ""
echo "Richieste ultima minuto:"
tail -60 backend.log | wc -l
echo ""
echo "Status Codes:"
tail -100 backend.log | grep -o "Status: [0-9]*" | sort | uniq -c
echo ""
echo "Top Endpoints:"
tail -100 backend.log | awk "{print \$4}" | sort | uniq -c | sort -rn | head -5
echo ""
echo "Media Response Time:"
tail -100 backend.log | grep -o "Duration: [0-9]*ms" | \
  grep -o "[0-9]*" | awk "{sum+=\$1; n++} END {print sum/n \"ms\"}"
echo ""
echo "Ultimi Errori:"
tail -100 backend.log | grep "Status: [45]" | tail -3
'
```

### Alert su Errori

```bash
#!/bin/bash
# alert.sh - Invia notifica se troppi errori

THRESHOLD=10
ERRORS=$(tail -100 backend.log | grep "Status: 5" | wc -l)

if [ $ERRORS -gt $THRESHOLD ]; then
    echo "ALERT: $ERRORS errori 5xx negli ultimi 100 log!"
    # Invia email o notifica Slack
fi
```

## 🔧 Integrazione con Strumenti

### 1. Logrotate

Configura rotazione automatica dei log:

```bash
# /etc/logrotate.d/conferenze-tech
/var/log/conferenze-tech/backend.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 conferenze conferenze
    sharedscripts
    postrotate
        systemctl reload conferenze-tech
    endscript
}
```

### 2. Systemd Journal

Redirigi log a systemd:

```bash
# Visualizza log service
journalctl -u conferenze-tech -f

# Filtra per priorità
journalctl -u conferenze-tech -p err

# Ultime 24 ore
journalctl -u conferenze-tech --since "24 hours ago"
```

### 3. ELK Stack (Elasticsearch, Logstash, Kibana)

Pattern Logstash per parsare i log:

```ruby
filter {
  grok {
    match => {
      "message" => "\[%{IPORHOST:client_ip}:%{NUMBER:client_port}\] %{WORD:method} %{URIPATH:path} \| Status: %{NUMBER:status_code} \| Size: %{NUMBER:size_bytes} bytes \| Duration: %{NUMBER:duration_ms}ms \| User: %{DATA:user_id}"
    }
  }
  
  mutate {
    convert => {
      "status_code" => "integer"
      "size_bytes" => "integer"
      "duration_ms" => "integer"
    }
  }
}
```

### 4. Prometheus + Grafana

Esponi metriche da log:

```bash
# Conta richieste per status
curl http://localhost:8080/metrics | grep http_requests_total

# Istogramma latenze
curl http://localhost:8080/metrics | grep http_request_duration_seconds
```

## 📋 Best Practices

### DO ✅

1. **Monitora costantemente i log in produzione**
   ```bash
   tail -f /var/log/conferenze-tech/backend.log
   ```

2. **Imposta alert automatici per errori 5xx**
   - > 10 errori/minuto: Alert WARNING
   - > 50 errori/minuto: Alert CRITICAL

3. **Analizza pattern settimanali**
   ```bash
   # Richieste per giorno della settimana
   grep "2026-01" backend.log | awk '{print substr($1,9,2)}' | sort | uniq -c
   ```

4. **Mantieni storico dei log (almeno 30 giorni)**

5. **Usa grep con -i per ricerche case-insensitive**
   ```bash
   grep -i "error" backend.log
   ```

### DON'T ❌

1. **Non ignorare errori 401/403**
   - Potrebbero indicare tentativi di accesso non autorizzato

2. **Non sottovalutare aumenti graduali di latenza**
   - Spesso segnalano problemi imminenti

3. **Non logare informazioni sensibili**
   - Password, token completi, dati personali

4. **Non dimenticare di ruotare i log**
   - Disco pieno = servizio down

5. **Non fare parsing manuale in produzione**
   - Usa strumenti dedicati (ELK, Splunk, etc.)

## 🎯 Checklist Monitoring Giornaliero

```
[ ] Verifica assenza errori 5xx critici
[ ] Controlla latenza media < 200ms
[ ] Verifica nessun picco anomalo di 401
[ ] Controlla crescita normale del traffico
[ ] Verifica nessun endpoint con latenza > 1s
[ ] Controlla dimensione log file < limite
[ ] Verifica backup log funzionanti
```

## 📞 Troubleshooting Rapido

| Sintomo | Comando | Azione |
|---------|---------|--------|
| API lenta | `grep "Duration: [0-9]{4,}ms"` | Ottimizza query DB |
| Molti 401 | `grep "401" \| wc -l` | Verifica token expiry |
| Errori 500 | `grep "500"` | Check database e logs applicativi |
| Traffic spike | `wc -l backend.log` | Verifica se legittimo o attacco |
| User bloccato | `grep "User: UUID"` | Traccia percorso utente |

## 🚀 Comandi Utili Rapidi

```bash
# Richieste/secondo real-time
tail -f backend.log | pv -l -i 1 -r > /dev/null

# Top 10 IP più attivi
awk '{print $1}' backend.log | cut -d: -f1 | sort | uniq -c | sort -rn | head -10

# Percentuale successi
grep "Status: 2" backend.log | wc -l

# Endpoint mai chiamati (se hai una lista)
comm -23 <(cat endpoints.txt | sort) <(awk '{print $4}' backend.log | sort -u)
```

---

**Nota**: Questi comandi assumono che i log siano scritti in un file `backend.log`. 
Adatta i percorsi in base alla tua configurazione.