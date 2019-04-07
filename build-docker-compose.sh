#!/bin/sh
env GOOS=linux go build -o app
docker-compose up --build