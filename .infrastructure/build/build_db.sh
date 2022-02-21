#!/bin/bash

sudo apt-get update
sudo apt-get install -y postgresql postgresql-contrib

sudo -u postgres psql -c "CREATE USER vagrant WITH ENCRYPTED PASSWORD 'vagrant';"
sudo -u postgres psql -c "CREATE DATABASE vagrant;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE vagrant to vagrant;"

export configPath=$(sudo -u postgres psql -c 'SHOW config_file' 2>/dev/null | grep /etc)
export hbaPath=$(sudo -u postgres psql -c 'SHOW hba_file' 2>/dev/null | grep /etc)

echo "listen_addresses = '*'" | sudo tee -a $configPath
echo "host samerole all 0.0.0.0/0 md5" | sudo tee -a $hbaPath

sudo service postgresql restart