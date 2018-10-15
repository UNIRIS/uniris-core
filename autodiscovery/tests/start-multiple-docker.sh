#!/bin/sh

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
printf '    build:\n'
printf '      dockerfile: build/package/Dockerfile\n'
printf '      context: ..\n'
printf '    environment:\n'
printf '      - UNIRIS_PUBLICKEY=publickey_%d\n' "$i"
printf '      - UNIRIS_NETWORK_INTERFACE=eth0\n'
if [ $i -eq 1 ]
then
    printf '      - UNIRIS_DISCOVERY_SEEDS=127.0.0.1:3545:publickey_1\n'
else
    seed=""
    for idx in {1..3}
    do
        r=$(( $RANDOM % $nb + 1 ))
        seed="${seed}peer${r}:3545:publickey_${r},"
    done
    printf '      - UNIRIS_DISCOVERY_SEEDS=%s\n' "$seed"
fi
printf '    networks:\n'
printf '      uniris:\n'
printf '        aliases:\n'
printf '          - peer%d\n' "$i"
printf '        ipv4_address: 172.16.0.%d\n' "1$i"
printf '    cpus: 1\n'
printf '    mem_limit: 512M\n'

done

} > bench.docker-compose.yml

docker-compose -f bench.docker-compose.yml up --build