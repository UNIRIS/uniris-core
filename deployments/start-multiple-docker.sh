#!/bin/sh

echo "Create network"
echo "###################"
docker network create --gateway 172.16.0.1 --subnet 172.16.0.0/24 uniris

echo "Create image"
echo "################"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uniris-miner ../cmd/uniris-miner
docker build -t uniris:latest -f ../build/docker/Dockerfile .
rm uniris-miner

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
for idx in {1..3}
do
    r=$(( $RANDOM % $nb + 1 ))
    seed="172.16.0.1${r}:5000:publickey_${r},"
done
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

docker-compose -f "docker-compose-with-$nb-peers.yml" up --build