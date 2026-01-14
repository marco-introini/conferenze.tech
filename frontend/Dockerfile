# Stage 1: Build e Sviluppo
FROM node:22-alpine AS builder
WORKDIR /app

# Copia solo package.json per caching delle dipendenze
COPY package*.json ./

# Installa dipendenze (questo layer viene cachato)
RUN npm ci --prefer-offline --no-audit

# Il resto viene montato come volume in sviluppo
COPY . .

# Stage 2: Final (Produzione con Nginx)
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
