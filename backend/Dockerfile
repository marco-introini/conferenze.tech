# Stage 1: Ambiente di Build e Sviluppo
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Installa Air per hot-reload in sviluppo
RUN go install github.com/air-verse/air@latest

# Copia solo i file di dipendenze per il caching
COPY go.mod go.sum ./
RUN go mod download

# Il resto viene montato come volume in sviluppo
COPY . .

# Stage 2: Immagine Finale per Produzione
FROM alpine:latest AS production
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main ./
CMD ["./main"]
