# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"
  config.vm.synced_folder ".", "/vagrant", type: "rsync"

  config.vm.define "server" do |server|
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

    server.vm.hostname = "server"
    server.vm.provision "shell", privileged: false, inline: <<-SHELL
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
      nohup ./minitwit > /tmp/out.log 2>&1 &
    SHELL
  end

  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
  SHELL
end