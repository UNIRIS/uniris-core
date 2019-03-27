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


privateKeys[0]="0066fc58626e8e245f58c9c705609f697071bde3968d33ab1022ea4488832f3f5aee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"
publicKeys[0]="00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"
privateKeys[1]="0077ff2d9233bdad72fb5c7a178003fb8c7fb3a69625075a8abc3064e62e714b66a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"
publicKeys[1]="00a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"
privateKeys[2]="00b274c45f99fe42de79433e8c8b71edcf4e61b6b8e55406883ba2c423785705af29dc9568f67727da6b0d93ca8538fe171a6bcd239dfd08dea748986848f24466"
publicKeys[2]="0029dc9568f67727da6b0d93ca8538fe171a6bcd239dfd08dea748986848f24466"
privateKeys[3]="00b84ac63695ef472b9a23dab78e388cb6ba2a7a0578b51c3ba4e4d7e5b23ee28a947b6ca58d87a126fc410858046c3edaea3dd1d570275502e7b45331c47fb655"
publicKeys[3]="00947b6ca58d87a126fc410858046c3edaea3dd1d570275502e7b45331c47fb655"
privateKeys[4]="00292fab83ef20397b07d2f9f0adaab2ff093d269111c70d3f80d1b7cc9ee6c1b7ca8a194ecb4ecc61287124a7f5d4db80a1a3d203ed43ace2e4fabf5e78a6bc83"
publicKeys[4]="00ca8a194ecb4ecc61287124a7f5d4db80a1a3d203ed43ace2e4fabf5e78a6bc83"
privateKeys[5]="003dfcbdfb38b042a9c319677a93999c5e3125a64e1441c7409d4c7db5434021295475d994fecb8492dfe05df3025c7fdfacdec323970b1c96c2d75f7ac5ef0fc3"
publicKeys[5]="005475d994fecb8492dfe05df3025c7fdfacdec323970b1c96c2d75f7ac5ef0fc3"
privateKeys[6]="00a5805009c794e7a90203417405c2ace5189af3a9fd68bc4ccfb5f1ccf4a0e00cb43683f3f473af5719f294f93519c579fbcf1f603080df8b0ccfe042aa5896b4"
publicKeys[6]="00b43683f3f473af5719f294f93519c579fbcf1f603080df8b0ccfe042aa5896b4"
privateKeys[7]="0076b7c35321d6f4e126bd6efaff8e33114f56a309756ef492d206b38e5a70df25254bd5c54d29afc54156042f50f0b3424c4f1a60882a3e2ea71d4e803ae301c5"
publicKeys[7]="00254bd5c54d29afc54156042f50f0b3424c4f1a60882a3e2ea71d4e803ae301c5"

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

ok=0
while [ $ok == 0 ]
do
    r=$(( $RANDOM % $nb ))
    rKey=${publicKeys[$r]}

    found=-1
    for j in "${prevKeys[@]}"
    do
        if [ $j == $rKey ]
        then
            found=1
        fi
    done

    if [ $found == -1 ]
    then
        printf '      - UNIRIS_PUBLIC_KEY=%s\n' ${publicKeys[$r]}
        printf '      - UNIRIS_PRIVATE_KEY=%s\n' ${privateKeys[$r]}
        prevKeys[$i]=${publicKeys[$r]}
        ok=1
    fi
done

printf '      - UNIRIS_NETWORK_INTERFACE=eth0\n'
seed=""
if [ $nb -gt 2 ]; then
    for idx in {1..2}
    do
        ok=0
        while [ $ok == 0 ]
        do
            r=$(( $RANDOM % $nb ))
            if [ $((r+1)) != $i ];
            then
                seed+="172.16.0.1$((r+1)):5000:${publicKeys[$r]};"
                ok=1
                prevSeed[$i]=${r}
            fi
        done
    done
else
    r=$(( $RANDOM % $nb ))
    seed="172.16.0.1$((r+1)):5000:${publicKeys[$r]};"
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