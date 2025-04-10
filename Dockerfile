# Build stage
FROM node:lts-alpine AS builder
WORKDIR /app
COPY package*.json .
RUN npm install
COPY . .
RUN npm run build

# Production stage with Caddy
FROM caddy:2.9.1-alpine

RUN addgroup -g 1337 -S caddy && \
    adduser -u 1337 -S caddy -G caddy

RUN chown -R caddy:caddy /usr/share/caddy /config /data

COPY --from=builder --chown=caddy:caddy /app/dist /usr/share/caddy
COPY --chown=caddy:caddy etc/caddy/Caddyfile /etc/caddy/Caddyfile

USER caddy

EXPOSE 80 443
