#!/usr/bin/env bash
docker compose -f docker-compose-dev.yaml down
touch .env && docker compose -f docker-compose-dev.yaml up -d --build --remove-orphans 
