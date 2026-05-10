FROM docker.io/library/golang:1.25.9-alpine AS backend-build

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

ENV CGO_ENABLED=0
ENV GOPROXY=direct

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates
RUN  cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      update-ca-certificates;

RUN apk update && apk --no-cache add bash dos2unix git 2>/dev/null

WORKDIR /app

COPY go.work* .
COPY webapi webapi
RUN go mod download
RUN go build -o server ./webapi

# final stage
FROM alpine 
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=backend-build --chown=appuser:appuser /app/server /app/server

USER appuser
EXPOSE 8888
ENTRYPOINT ["/app/server"]