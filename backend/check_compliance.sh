#!/bin/bash

echo "🔍 Go Compliance Check - conferenze.tech backend"
echo "=================================================="
echo ""

cd "$(dirname "$0")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

score=0
total=0

# 1. go fmt
echo "1️⃣  Checking go fmt..."
total=$((total+1))
if [ -z "$(gofmt -l .)" ]; then
    echo -e "${GREEN}✅ All files properly formatted${NC}"
    score=$((score+1))
else
    echo -e "${RED}❌ Files need formatting:${NC}"
    gofmt -l .
fi
echo ""

# 2. go vet
echo "2️⃣  Running go vet..."
total=$((total+1))
if go vet ./... 2>&1 | grep -q "^"; then
    echo -e "${RED}❌ go vet found issues${NC}"
    go vet ./...
else
    echo -e "${GREEN}✅ go vet passed${NC}"
    score=$((score+1))
fi
echo ""

# 3. staticcheck
echo "3️⃣  Running staticcheck..."
total=$((total+1))
if command -v ~/go/bin/staticcheck &> /dev/null; then
    if ~/go/bin/staticcheck ./... 2>&1 | grep -v "^$" | grep -q "."; then
        echo -e "${YELLOW}⚠️  staticcheck found issues:${NC}"
        ~/go/bin/staticcheck ./... | head -10
    else
        echo -e "${GREEN}✅ staticcheck passed${NC}"
        score=$((score+1))
    fi
else
    echo -e "${YELLOW}⚠️  staticcheck not installed (skipping)${NC}"
fi
echo ""

# 4. tests
echo "4️⃣  Running tests..."
total=$((total+1))
if go test ./... -v 2>&1 | grep -q "PASS"; then
    echo -e "${GREEN}✅ Tests passed${NC}"
    score=$((score+1))
else
    echo -e "${RED}❌ Tests failed${NC}"
    go test ./...
fi
echo ""

# 5. build
echo "5️⃣  Building..."
total=$((total+1))
if go build -o /dev/null ./... 2>&1; then
    echo -e "${GREEN}✅ Build successful${NC}"
    score=$((score+1))
else
    echo -e "${RED}❌ Build failed${NC}"
    go build ./...
fi
echo ""

# Results
echo "=================================================="
echo -e "📊 Score: ${GREEN}$score${NC}/$total"
echo ""

if [ $score -eq $total ]; then
    echo -e "${GREEN}🎉 All checks passed! Code is compliant!${NC}"
    exit 0
elif [ $score -ge $((total-1)) ]; then
    echo -e "${YELLOW}⚠️  Almost there! Minor issues to fix.${NC}"
    exit 1
else
    echo -e "${RED}❌ Several issues found. Check output above.${NC}"
    exit 1
fi
