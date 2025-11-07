#!/bin/bash
docker compose --file deploy/docker-compose.yaml --project-directory ./ up --force-recreate --build