#!/bin/sh

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    *)          machine="UNKNOWN:${unameOut}"
esac

echo "Build UNIRIS binary"
echo "###################"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uniris-node ../../cmd/uniris-node


echo "Build UNIRIS image"
echo "###################"

if [ $machine = "Linux" ]
then
     sudo docker build -t uniris:latest .
else
     docker build -t uniris:latest .
fi

rm uniris-node