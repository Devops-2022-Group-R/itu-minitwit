defaults = './.infrastructure/Vagrantfile.defaults'
load defaults if File.exists?(defaults)


$ip_file = "db_ip.txt"

Vagrant.configure("2") do |config|
  config.vm.define "client" do |client|
    client.vm.hostname = "client"

    virtual_box_network(client, "192.168.56.4", 80, 8080)

    client.vm.provision "shell", path: ".infrastructure/build/build_frontend.sh"
  end

  config.vm.define "db" do |db|
    db.vm.hostname = "dbserver"

    vbDbIp = "192.168.56.2"
    virtual_box_network(db, vbDbIp)

    db.vm.provision "shell", path: ".infrastructure/build/build_db.sh"

    db.vm.provider "virtualbox" do |vb, override|
      override.trigger.after :up do |trigger|
        trigger.info =  "Writing dbserver's IP to file with virtualbox..."

        trigger.ruby do |env, machine|
          File.write($ip_file, vbDbIp)
        end
      end
    end

    db.vm.provider :digital_ocean do |digital, override|
      override.trigger.after :up do |trigger|
        trigger.info =  "Writing dbserver's IP to file with digitalocean..."

        trigger.ruby do |env,machine|
          p "Writing with digitalocean"
          remote_ip = machine.instance_variable_get(:@communicator).instance_variable_get(:@connection_ssh_info)[:host]
          File.write($ip_file, remote_ip)
        end 
      end
    end
  end

  config.vm.define "server" do |server|
    config.vm.hostname = "server"

    virtual_box_network(server, "192.168.56.3", 80, 80)
    
    server.vm.provision "shell", path: ".infrastructure/build/build_backend.sh"

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
  end
end
