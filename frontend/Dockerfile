FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

# Vite inlines VITE_* at build time, so this is an arg, not a runtime env var.
ARG VITE_BACKEND_API_URL
ENV VITE_BACKEND_API_URL=$VITE_BACKEND_API_URL

COPY . .
RUN npm run build


FROM nginxinc/nginx-unprivileged:1.27-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# The base image runs as uid 101 and cannot bind a privileged port.
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/ >/dev/null || exit 1
