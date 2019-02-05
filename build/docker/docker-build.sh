#!/bin/sh

# echo "Create network"
# echo "###################"
# docker network create --gateway 172.16.0.1 --subnet 172.16.0.0/24 uniris

echo "Build UNIRIS binary"
echo "###################"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uniris ../cmd/*

echo "Build UNIRIS image"
echo "###################"
docker build -t uniris:latest .

rm uniris