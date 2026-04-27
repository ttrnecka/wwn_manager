# build stage
FROM docker.io/library/node:20-alpine AS frontend-build

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates
RUN cat /tmp/zscaler-root-ca.crt >> /etc/ssl/certs/ca-certificates.crt
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci  --force --loglevel verbose
COPY frontend/ .
RUN npm run build

# final stage
FROM alpine 
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
COPY --from=frontend-build --chown=appuser:appgroup /app/dist /app/dist

USER appuser
