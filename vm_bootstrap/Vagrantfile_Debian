Vagrant.configure(2) do |config|
	
	config.vm.define "VM-1" do |node|
	
		node.vm.box = "ubuntu/xenial64"
		node.vm.hostname = "VM-1"
		node.vm.network "private_network", ip: "192.168.2.10"
		
		node.vm.provider "virtualbox" do |vb|
			 # Display the VirtualBox GUI when booting the machine
				vb.gui = true
			
			# Customize the amount of memory on the VM:
			vb.memory = "1524"
		end
		
		node.vm.provision :shell, path: "bootstrap.sh"
	end
	
	config.vm.define "VM-2" do |node|

		node.vm.box = "ubuntu/xenial64"
		node.vm.hostname = "VM-2"
		node.vm.network "private_network", ip: "192.168.2.20"
		
		node.vm.provider "virtualbox" do |vb|
			 # Display the VirtualBox GUI when booting the machine
				vb.gui = true
			
			# Customize the amount of memory on the VM:
			vb.memory = "1524"
		end
		
		node.vm.provision :shell, path: "bootstrap.sh"
	end
	
end
