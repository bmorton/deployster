# -*- mode: ruby -*-
# vi: set ft=ruby :

# This is a simplified version of the CoreOS Vagrant repo available at:
# https://github.com/coreos/coreos-vagrant
#
# It is meant to help you quickly get setup with Deployster so you can easily
# test out the functionality and explore its features.

VAGRANTFILE_API_VERSION = "2"
Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.box = "coreos-alpha"
  config.vm.box_version = ">= 550.0.0"
  config.vm.box_url = "http://alpha.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json"

  ["vmware_fusion", "vmware_workstation"].each do |vmware|
    config.vm.provider vmware do |v, override|
      override.vm.box_url = "http://alpha.release.core-os.net/amd64-usr/current/coreos_production_vagrant_vmware_fusion.json"
      v.gui = false
      v.vmx['memsize'] = 1024
      v.vmx['numvcpus'] = 1
    end
  end

  config.vm.provider :virtualbox do |v|
    # On VirtualBox, we don't have guest additions or a functional vboxsf
    # in CoreOS, so tell Vagrant that so it can be smarter.
    v.check_guest_additions = false
    v.functional_vboxsf     = false

    v.gui = false
    v.memory = 1024
    v.cpus = 1
  end

  if Vagrant.has_plugin?("vagrant-vbguest") then
    config.vbguest.auto_update = false
  end

  config.vm.hostname = "deployster"
  config.vm.network "forwarded_port", guest: 2375, host: 2375, auto_correct: true
  config.vm.network "forwarded_port", guest: 80, host: 8080, auto_correct: true
  config.vm.network "forwarded_port", guest: 3000, host: 3000, auto_correct: true
  config.vm.network :private_network, ip: "172.17.8.100"

  config.vm.synced_folder ".", "/home/core/share", id: "core", nfs: true, mount_options: ['nolock,vers=3,udp']

  config.vm.provision :file,  source: "vagrant/cloud-config.yml", destination: "/tmp/vagrantfile-user-data"
  config.vm.provision :shell, inline: "docker pull mailgun/vulcand:v0.7.0"
  config.vm.provision :shell, inline: "docker pull bmorton/deployster:latest"
  config.vm.provision :shell, inline: "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant/", privileged: true
end
