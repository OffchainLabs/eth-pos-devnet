docker ps -a | grep "eth-pos" | awk '{print $1}' | xargs docker stop
docker ps -a | grep "eth-pos" | awk '{print $1}' | xargs docker remove

rm -Rf ./consensus/beacondata ./consensus/validatordata ./consensus/genesis.ssz
rm -Rf ./execution/geth

