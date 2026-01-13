# Stage 1: Ambiente di Build e Sviluppo
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main .

# Stage 2: Immagine Finale
FROM alpine:latest AS production
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]