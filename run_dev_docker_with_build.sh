#!/usr/bin/env bash
touch .env && docker compose -f docker-compose-dev.yaml up -d --build
