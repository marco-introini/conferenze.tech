# --- Stage di Sviluppo/Build ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
# Installiamo Air per l'hot-reload in dev
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# --- Stage di Produzione (Immagine finale leggera) ---
FROM alpine:latest AS production
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]