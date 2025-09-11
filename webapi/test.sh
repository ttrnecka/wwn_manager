#!/usr/bin/env bash
systemctl --user start podman.socket
export DOCKER_HOST=unix://`podman info --format '{{.Host.RemoteSocket.Path}}'`
go test ./...