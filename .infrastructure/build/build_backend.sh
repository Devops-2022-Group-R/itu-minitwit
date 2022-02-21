#!/bin/bash

export DB_IP=`cat /vagrant/db_ip.txt`
export ENVIRONMENT="PRODUCTION"
export CONNECTION_STRING="host=$DB_IP user=vagrant password=vagrant dbname=vagrant port=5432 sslmode=disable" 

sudo apt-get update
sudo apt-get install -y gcc

echo "Downloading go"
wget -q https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
echo "Go downloaded"

sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin
export PORT=80

mkdir $HOME/bin
mkdir $HOME/minitwit
cp -r /vagrant/* $HOME/minitwit

cd $HOME/minitwit
go build -o $HOME/bin/minitwit ./src

cd $HOME/bin
./minitwit initDb
nohup ./minitwit > /tmp/out.log 2>&1 &