#!/bin/sh


echo "Build UNIRIS binary"
echo "###################"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uniris ../cmd/*

echo "Build UNIRIS image"
echo "###################"
docker build -t uniris:latest .

rm uniris