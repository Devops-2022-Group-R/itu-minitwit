# -*- mode: ruby -*-
# vi: set ft=ruby :

$ip_file = "db_ip.txt"

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"
  config.vm.synced_folder ".", "/vagrant", type: "rsync"

  config.vm.define "db" do |server|
    server.vm.provider "virtualbox" do |vb|
      server.vm.network "private_network", ip: "192.168.56.2"
      vb.memory = "1024"
    end

    server.vm.provider :digital_ocean do |provider, override|
      override.vm.box = 'digital_ocean'
      override.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
      override.vm.allowed_synced_folder_types = :rsync
      
      override.ssh.private_key_path = ENV["SSH_KEY"]
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      
      provider.image = 'ubuntu-21-10-x64'
      
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'

      provider.backups_enabled = false
      
      provider.private_networking = true
      provider.ipv6 = false
      provider.monitoring = false
    end

    server.vm.hostname = "dbserver"

    server.trigger.after :up do |trigger|
      trigger.info =  "Writing dbserver's IP to file..."
      trigger.ruby do |env,machine|
        remote_ip = machine.instance_variable_get(:@communicator).instance_variable_get(:@connection_ssh_info)[:host]
        File.write($ip_file, remote_ip)
      end 
    end

    server.vm.provision "shell", inline: <<-SHELL
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
    SHELL
  end

  config.vm.define "server" do |server|
    server.vm.provider "virtualbox" do |vb|
      server.vm.network "private_network", ip: "192.168.56.3"
      vb.memory = "1024"
    end

    server.vm.provider :digital_ocean do |provider, override|
      override.vm.box = 'digital_ocean'
      override.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
      override.vm.allowed_synced_folder_types = :rsync
      
      override.ssh.private_key_path = ENV["SSH_KEY"]
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      
      provider.image = 'ubuntu-21-10-x64'
      
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'

      provider.backups_enabled = false
      
      provider.private_networking = true
      provider.ipv6 = false
      provider.monitoring = false
    end

    server.vm.hostname = "server"

    server.trigger.before :up do |trigger|
      trigger.info =  "Waiting to create server until dbserver's IP is available."
      trigger.ruby do |env,machine|
        while !File.file?($ip_file) do
          sleep(1)
        end
        
        db_ip = File.read($ip_file).strip()
        
        puts "Now, I have it..."
        puts db_ip
      end 
    end

    server.trigger.after :provision do |trigger|
      trigger.ruby do |env,machine|
        File.delete($ip_file) if File.exists? $ip_file
      end 
    end

    server.vm.provision "shell", privileged: false, inline: <<-SHELL
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
    SHELL
  end

  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
  SHELL
end