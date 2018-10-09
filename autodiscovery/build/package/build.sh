#!/bin/sh

echo "Create network"
echo "###################"
docker network create --gateway 172.16.0.1 --subnet 172.16.0.0/24 uniris

echo "Build UNIRIS image"
echo "###################"
docker build -t uniris:latest uniris

echo "Build Discovery image"
echo "###################"
docker build -t uniris-discovery:latest -f Dockerfile ../..
