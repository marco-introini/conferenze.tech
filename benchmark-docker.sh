#!/bin/bash

# Script per benchmarking startup Docker

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Docker Startup Benchmark${NC}\n"

# Funzione per misurare il tempo
measure_time() {
    local description=$1
    local command=$2

    echo -e "${YELLOW}Test: ${description}${NC}"

    # Cleanup prima del test
    docker compose down -v --remove-orphans > /dev/null 2>&1
    sleep 2

    # Misura tempo
    start=$(date +%s)
    eval "$command"
    end=$(date +%s)

    duration=$((end - start))
    echo -e "${GREEN}✓ Completato in ${duration}s${NC}\n"

    # Cleanup dopo il test
    docker compose down > /dev/null 2>&1
    sleep 2

    return $duration
}

echo "Questo script confronterà i tempi di startup."
echo "Ogni test include build e startup completo."
echo ""
read -p "Premere INVIO per continuare..."

# Test 1: Primo build (cold start)
echo -e "\n${BLUE}═══════════════════════════════════════${NC}"
echo -e "${BLUE}Test 1: Primo Build (Cold Start)${NC}"
echo -e "${BLUE}═══════════════════════════════════════${NC}\n"

docker system prune -af --volumes > /dev/null 2>&1
measure_time "Build completo da zero" "docker compose up --build -d && docker compose logs --tail=20"
cold_time=$?

# Test 2: Rebuild con cache
echo -e "\n${BLUE}═══════════════════════════════════════${NC}"
echo -e "${BLUE}Test 2: Rebuild con Cache${NC}"
echo -e "${BLUE}═══════════════════════════════════════${NC}\n"

measure_time "Rebuild con dipendenze in cache" "docker compose up --build -d && docker compose logs --tail=20"
rebuild_time=$?

# Test 3: Startup senza rebuild
echo -e "\n${BLUE}═══════════════════════════════════════${NC}"
echo -e "${BLUE}Test 3: Startup Veloce (dev-fast)${NC}"
echo -e "${BLUE}═══════════════════════════════════════${NC}\n"

# Build una volta
docker compose up --build -d > /dev/null 2>&1
docker compose down > /dev/null 2>&1
sleep 2

# Ora misura solo startup
measure_time "Startup senza rebuild" "docker compose up -d && docker compose logs --tail=20"
fast_time=$?

# Risultati
echo -e "\n${BLUE}═══════════════════════════════════════${NC}"
echo -e "${BLUE}📊 RISULTATI BENCHMARK${NC}"
echo -e "${BLUE}═══════════════════════════════════════${NC}\n"

echo "┌─────────────────────────────────────┬──────────┐"
echo "│ Test                                │ Tempo    │"
echo "├─────────────────────────────────────┼──────────┤"
printf "│ %-35s │ %6ss │\n" "Primo Build (Cold)" "$cold_time"
printf "│ %-35s │ %6ss │\n" "Rebuild con Cache" "$rebuild_time"
printf "│ %-35s │ %6ss │\n" "Startup Veloce (dev-fast)" "$fast_time"
echo "└─────────────────────────────────────┴──────────┘"

# Calcola miglioramento
if [ $cold_time -gt 0 ]; then
    cache_improvement=$(( (cold_time - rebuild_time) * 100 / cold_time ))
    fast_improvement=$(( (cold_time - fast_time) * 100 / cold_time ))

    echo -e "\n${GREEN}🎯 Miglioramenti:${NC}"
    echo "  • Cache dipendenze: ${cache_improvement}% più veloce"
    echo "  • Startup senza rebuild: ${fast_improvement}% più veloce"
fi

echo -e "\n${YELLOW}💡 Raccomandazioni:${NC}"
echo "  • Primo avvio: make dev (1 volta)"
echo "  • Uso quotidiano: make dev-fast (molto più veloce)"
echo "  • Dopo cambio dipendenze: make dev"

# Cleanup finale
echo -e "\n${BLUE}Pulizia finale...${NC}"
docker compose down -v --remove-orphans > /dev/null 2>&1

echo -e "\n${GREEN}✓ Benchmark completato!${NC}\n"
