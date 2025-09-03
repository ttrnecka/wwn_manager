FROM golang:1.23-alpine AS backend-build

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

ENV CGO_ENABLED=0
ENV GOPROXY=direct

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates
RUN  cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      update-ca-certificates;

RUN apk update && apk add bash dos2unix git 2>/dev/null

WORKDIR /app

COPY go.work* .
COPY webapi webapi
RUN go mod download
RUN go build -o server ./webapi

# final stage
FROM alpine 
WORKDIR /app
COPY --from=backend-build /app/server /app/server

EXPOSE 8888
ENTRYPOINT ["/app/server"]