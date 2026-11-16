#!/usr/bin/env bash
set -euo pipefail

docker compose -f docker-compose-scan.yaml up --build --remove-orphans --pull always
docker compose -f docker-compose-scan.yaml down 