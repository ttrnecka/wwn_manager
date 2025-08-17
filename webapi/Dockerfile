FROM golang:1.23-alpine AS build
ARG IMPORT_ZSCALER_CERT=false

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

ENV CGO_ENABLED=0

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates

RUN if [ "$IMPORT_ZSCALER_CERT" = "true" ]; then \
      cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      update-ca-certificates; \
    else \
      echo "Skipping Zscaler cert install"; \
    fi

RUN apk update && apk add bash dos2unix 2>/dev/null

WORKDIR /app

COPY go.work* .
COPY webapi webapi
RUN go mod download
RUN go build -o /webapi ./webapi

# final stage
FROM alpine 
WORKDIR /
COPY --from=build /webapi /webapi
EXPOSE 8888
ENTRYPOINT ["/webapi"]