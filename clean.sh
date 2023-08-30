#! /bin/bash

docker rm -f $(docker ps -a -q)

sudo rm -rf ./consensus/beacondata
sudo rm -rf ./consensus/validatordata
sudo rm -rf ./execution/geth

