#!/bin/bash

sudo apt-get update
sudo apt-get install -y postgresql postgresql-contrib

db="vagrant"
user="vagrant"
password="vagrant"

sudo -u postgres psql -c "CREATE USER $user WITH ENCRYPTED PASSWORD '$password';"
sudo -u postgres psql -c "CREATE DATABASE $db;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $db to $user;"

export configPath=$(sudo -u postgres psql -c 'SHOW config_file' 2>/dev/null | grep /etc)
export hbaPath=$(sudo -u postgres psql -c 'SHOW hba_file' 2>/dev/null | grep /etc)

localNetwork=$(ip route | awk '/.* eth1 .*/ {print $1}')
ip=$(ip route | awk '/.* eth1 .*/ {print $NF}')

echo "listen_addresses = '$ip'" | sudo tee -a $configPath
# Type database user address method
echo "host $db $user $localNetwork md5" | sudo tee -a $hbaPath

sudo service postgresql restart