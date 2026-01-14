#!/bin/bash

# Conferenze.tech API - Test con curl
# Uso: ./api-examples.sh
# Richiede: jq (brew install jq)

set -e  # Exit on error (disabilitato per test)
set +e  # Permetti errori per gestirli

BASE_URL="http://localhost:8080"
TOKEN=""

# Verifica jq installato
if ! command -v jq &> /dev/null; then
    echo -e "\033[0;31m✗ jq non trovato. Installalo con: brew install jq\033[0m"
    exit 1
fi

# Verifica server sia up
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "\033[0;31m✗ Server non raggiungibile su $BASE_URL\033[0m"
    echo -e "\033[1;33mAvvia il server con: make dev-fast\033[0m"
    exit 1
fi

# Colori per output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funzione per stampare sezioni
section() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Funzione per stampare comandi
cmd() {
    echo -e "${YELLOW}▶ $1${NC}"
}

# Funzione per successo
success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Funzione per errore
error() {
    echo -e "${RED}✗ $1${NC}"
}

################################################################################
# AUTH & REGISTRATION
################################################################################

section "1. AUTENTICAZIONE E REGISTRAZIONE"

cmd "Registra nuovo utente"
echo 'curl -X POST '"$BASE_URL"'/api/register \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "name": "Mario Rossi",
    "email": "mario.rossi@example.com",
    "password": "securepass123"
  }'"'"

echo -e "\n${GREEN}Prova:${NC}"
curl -s -X POST "$BASE_URL/api/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test'$(date +%s)'@example.com",
    "password": "password123"
  }' | jq .

echo ""
cmd "Login e ottieni token"
echo 'curl -X POST '"$BASE_URL"'/api/login \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "email": "marco@marcointroini.it",
    "password": "password"
  }'"'"

echo -e "\n${GREEN}Prova (salva il token!):${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "marco@marcointroini.it",
    "password": "password"
  }')

# Verifica se la risposta è JSON valido
if echo "$RESPONSE" | jq . > /dev/null 2>&1; then
    echo "$RESPONSE" | jq .
    # Estrai token
    TOKEN=$(echo "$RESPONSE" | jq -r '.token // empty')

    if [ -n "$TOKEN" ]; then
        success "Token ottenuto: ${TOKEN:0:20}..."
        export TOKEN
    else
        error "Login fallito - token non trovato nella response"
    fi
else
    error "Login fallito - response non è JSON valido:"
    echo "$RESPONSE"
    echo ""
    echo -e "${YELLOW}Possibili cause:${NC}"
    echo "  1. Server non è in esecuzione (make dev-fast)"
    echo "  2. Database non è pronto (make logs-db)"
    echo "  3. Tabella user_tokens mancante (make migrate)"
    echo "  4. Credenziali errate"
    exit 1
fi

echo ""
cmd "Verifica token (GET /api/me)"
echo 'curl '"$BASE_URL"'/api/me \
  -H "Authorization: Bearer $TOKEN"'

if [ -n "$TOKEN" ]; then
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s "$BASE_URL/api/me" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

################################################################################
# CONFERENCES
################################################################################

section "2. GESTIONE CONFERENZE"

cmd "Lista tutte le conferenze (pubblico)"
echo "curl $BASE_URL/api/conferences"
echo -e "\n${GREEN}Prova:${NC}"
curl -s "$BASE_URL/api/conferences" | jq .

echo ""
cmd "Cerca conferenze per nome"
echo 'curl "'"$BASE_URL"'/api/conferences?search=golang"'
echo -e "\n${GREEN}Prova:${NC}"
curl -s "$BASE_URL/api/conferences?search=go" | jq .

echo ""
cmd "Crea nuova conferenza (richiede auth)"
echo 'curl -X POST '"$BASE_URL"'/api/conferences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "title": "GopherCon Italia 2026",
    "date": "2026-09-15T10:00:00Z",
    "location": "Milano",
    "website": "https://gophercon.it",
    "latitude": 45.4642,
    "longitude": 9.1900
  }'"'"

if [ -n "$TOKEN" ]; then
    echo -e "\n${GREEN}Prova:${NC}"
    CONF_RESPONSE=$(curl -s -X POST "$BASE_URL/api/conferences" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Test Conference '$(date +%H%M%S)'",
        "date": "2026-12-01T10:00:00Z",
        "location": "Online"
      }')

    if echo "$CONF_RESPONSE" | jq . > /dev/null 2>&1; then
        echo "$CONF_RESPONSE" | jq .
        CONF_ID=$(echo "$CONF_RESPONSE" | jq -r '.id // empty')
        if [ -n "$CONF_ID" ]; then
            success "Conferenza creata con ID: $CONF_ID"
            export CONF_ID
        fi
    else
        error "Errore nella creazione conferenza:"
        echo "$CONF_RESPONSE"
    fi
fi

echo ""
cmd "Dettaglio conferenza"
if [ -n "$CONF_ID" ]; then
    echo "curl $BASE_URL/api/conferences/$CONF_ID"
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s "$BASE_URL/api/conferences/$CONF_ID" | jq .
else
    echo "curl $BASE_URL/api/conferences/{conference_id}"
fi

echo ""
cmd "Aggiorna conferenza (richiede ownership)"
if [ -n "$CONF_ID" ] && [ -n "$TOKEN" ]; then
    echo 'curl -X PUT '"$BASE_URL"'/api/conferences/'"$CONF_ID"' \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "title": "Test Conference AGGIORNATO",
    "date": "2026-12-01T10:00:00Z",
    "location": "Milano",
    "website": "https://example.com"
  }'"'"

    echo -e "\n${GREEN}Prova:${NC}"
    curl -s -X PUT "$BASE_URL/api/conferences/$CONF_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Test Conference UPDATED",
        "date": "2026-12-01T10:00:00Z",
        "location": "Milano"
      }' | jq .
else
    echo 'curl -X PUT '"$BASE_URL"'/api/conferences/{id} ...'
fi

echo ""
cmd "Le mie conferenze"
if [ -n "$TOKEN" ]; then
    echo "curl $BASE_URL/api/my-conferences \\"
    echo '  -H "Authorization: Bearer $TOKEN"'
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s "$BASE_URL/api/my-conferences" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

################################################################################
# REGISTRATIONS
################################################################################

section "3. ISCRIZIONI CONFERENZE"

cmd "Registrati a una conferenza"
if [ -n "$CONF_ID" ] && [ -n "$TOKEN" ]; then
    echo 'curl -X POST '"$BASE_URL"'/api/conferences/'"$CONF_ID"'/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "notes": "Interessato ai workshop"
  }'"'"

    echo -e "\n${GREEN}Prova:${NC}"
    curl -s -X POST "$BASE_URL/api/conferences/$CONF_ID/register" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "notes": "Test registration"
      }' | jq .
fi

echo ""
cmd "Lista partecipanti (organizzatore only)"
if [ -n "$CONF_ID" ] && [ -n "$TOKEN" ]; then
    echo "curl $BASE_URL/api/conferences/$CONF_ID/registrations \\"
    echo '  -H "Authorization: Bearer $TOKEN"'
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s "$BASE_URL/api/conferences/$CONF_ID/registrations" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

echo ""
cmd "Le mie iscrizioni"
if [ -n "$TOKEN" ]; then
    echo "curl $BASE_URL/api/my-registrations \\"
    echo '  -H "Authorization: Bearer $TOKEN"'
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s "$BASE_URL/api/my-registrations" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

echo ""
cmd "Cancella iscrizione"
if [ -n "$CONF_ID" ] && [ -n "$TOKEN" ]; then
    echo "curl -X DELETE $BASE_URL/api/conferences/$CONF_ID/register \\"
    echo '  -H "Authorization: Bearer $TOKEN"'
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s -X DELETE "$BASE_URL/api/conferences/$CONF_ID/register" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

################################################################################
# ERROR CASES
################################################################################

section "4. TEST CASI D'ERRORE"

cmd "404 - Conferenza non esistente"
echo "curl $BASE_URL/api/conferences/00000000-0000-0000-0000-000000000000"
echo -e "\n${GREEN}Prova:${NC}"
curl -s "$BASE_URL/api/conferences/00000000-0000-0000-0000-000000000000" | jq .

echo ""
cmd "401 - Richiesta senza autenticazione"
echo 'curl -X POST '"$BASE_URL"'/api/conferences \
  -H "Content-Type: application/json" \
  -d '"'"'{"name":"Test"}'"'"
echo -e "\n${GREEN}Prova:${NC}"
curl -s -X POST "$BASE_URL/api/conferences" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test"}' | jq .

echo ""
cmd "401 - Token invalido"
echo "curl $BASE_URL/api/me \\"
echo '  -H "Authorization: Bearer invalid-token"'
echo -e "\n${GREEN}Prova:${NC}"
curl -s "$BASE_URL/api/me" \
  -H "Authorization: Bearer invalid-token-12345" | jq .

echo ""
cmd "400 - Dati invalidi"
cmd "Dati invalidi (400)"
echo 'curl -X POST '"$BASE_URL"'/api/conferences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "title": "",
    "date": "invalid-date"
  }'"'"

if [ -n "$TOKEN" ]; then
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s -X POST "$BASE_URL/api/conferences" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "",
        "date": "invalid-date"
      }' | jq .
fi

################################################################################
# CLEANUP
################################################################################

section "5. CLEANUP"

cmd "Elimina conferenza di test"
if [ -n "$CONF_ID" ] && [ -n "$TOKEN" ]; then
    echo "curl -X DELETE $BASE_URL/api/conferences/$CONF_ID \\"
    echo '  -H "Authorization: Bearer $TOKEN"'
    echo -e "\n${GREEN}Prova:${NC}"
    curl -s -X DELETE "$BASE_URL/api/conferences/$CONF_ID" \
      -H "Authorization: Bearer $TOKEN"
    success "Conferenza eliminata"
fi

################################################################################
# RIEPILOGO
################################################################################

section "RIEPILOGO"

echo -e "${GREEN}Token salvato in variabile:${NC} \$TOKEN"
echo -e "${GREEN}Conference ID:${NC} \$CONF_ID"
echo ""
echo -e "${YELLOW}Per usare il token in altri comandi:${NC}"
echo "export TOKEN='$TOKEN'"
echo ""
echo -e "${YELLOW}Esempi rapidi:${NC}"
echo "curl $BASE_URL/api/conferences | jq ."
echo 'curl -H "Authorization: Bearer $TOKEN" '"$BASE_URL"'/api/me | jq .'
echo ""
echo -e "${BLUE}Note:${NC}"
echo "• Installa 'jq' per formattare JSON: brew install jq (macOS) o apt install jq (Linux)"
echo "• Verifica server sia running: make dev-fast"
echo "• Se errori 'Failed to save token': make migrate"
echo "• Usa -v per vedere headers: curl -v ..."
echo "• Usa -i per vedere status: curl -i ..."
