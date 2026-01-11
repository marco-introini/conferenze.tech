# Stage 1: Ambiente di Build e Sviluppo
FROM golang:1.25-alpine AS builder
WORKDIR /app
# Copiamo solo i file dei moduli per sfruttare la cache di Docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compiliamo il binario per la produzione
RUN go build -o /main .

# Stage 2: Immagine Finale Piuma per il Mini PC
FROM alpine:latest AS production
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]