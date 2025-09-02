# build stage
# AS frontend-build
FROM node:20-alpine
ARG IMPORT_ZSCALER_CERT=false

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates

RUN if [ "$IMPORT_ZSCALER_CERT" = "true" ]; then \
      cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      npm config set cafile /usr/local/share/ca-certificates/zscaler-root-ca.crt; \
    else \
      echo "Skipping Zscaler cert install"; \
    fi

WORKDIR /app
# ENV NODE_EXTRA_CA_CERTS=/usr/local/share/ca-certificates/zscaler-root-ca.crt
COPY frontend/package*.json ./
RUN npm ci  --force --loglevel verbose
COPY frontend/ .
RUN npm run build

# # final stage
# FROM alpine 
# WORKDIR /app
# COPY --from=frontend-build /app/dist /app/dist
