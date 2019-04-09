#!/bin/sh
env GOOS=linux go build -o app
docker build -f Dockerfile-dev . -t linxianer12/product-golang 
docker run -p 3000:3000 linxianer12/product-golang 