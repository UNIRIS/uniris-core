#!/bin/sh

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    *)          machine="UNKNOWN:${unameOut}"
esac

echo "Create network"
echo "###################"
if [ $machine = "Linux" ]
then
    sudo docker network create --gateway 172.16.0.1 --subnet 172.16.0.0/24 uniris
else
    docker network create --gateway 172.16.0.1 --subnet 172.16.0.0/24 uniris
fi

echo "Create image"
echo "################"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uniris-node ../cmd/uniris-node
if [ $machine = "Linux" ]
then
    sudo docker build -t uniris:latest -f ../build/docker/Dockerfile .
else
    docker build -t uniris:latest -f ../build/docker/Dockerfile .
fi

rm uniris-node

nb=$1

{
printf 'version: "2.2"\n'
printf 'networks:\n'
printf '  uniris:\n'
printf '    external: true\n'
printf 'services:\n'
for i in `seq  $nb`
do
printf '  peer%d:\n' "$i"
printf '    image: uniris\n'
printf '    environment:\n'
printf '      - UNIRIS_PUBLIC_KEY=publickey_%d\n' "$i"
printf '      - UNIRIS_PRIVATE_KEY=privatekey_%d\n' "$i"
printf '      - UNIRIS_NETWORK_INTERFACE=eth0\n'
seed=""
if [ $nb -gt 2 ]; then
    for idx in {1..2}
    do
        ok=0
        while [ $ok == 0 ]
        do
            r=$(( $RANDOM % $nb + 1 ))
            if [ $r != $i ];
            then
                seed+="172.16.0.1${r}:5000:publickey_${r};"
                ok=1
            fi
        done
    done
else
    ok=0
    prevR=-1
    while [ $ok == 0 ]
    do
        r=$(( $RANDOM % $nb + 1 ))
        if [ $r != $i ] || [prevR != $r];
        then
            seed="172.16.0.1${r}:5000:publickey_${r};"
            ok=1
        fi
    done
fi

printf '      - UNIRIS_DISCOVERY_SEEDS=%s\n' "$seed"
printf '    ports:\n'
printf '      - 4000\n' "$i"
printf '      - 15672\n' "$i"
printf '    networks:\n'
printf '      uniris:\n'
printf '        ipv4_address: 172.16.0.%d\n' "1$i"
printf '    cpus: 1\n'
printf '    mem_limit: 512M\n'

done

} > "docker-compose-with-$nb-peers.yml"

if [ $machine = "Linux" ]
then
    sudo docker-compose -f "docker-compose-with-$nb-peers.yml" up --build
else
    docker-compose -f "docker-compose-with-$nb-peers.yml" up --build
fi