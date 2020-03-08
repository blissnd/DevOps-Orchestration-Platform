# DevOps-Orchestration-Platform
Bootstrap an entire IT infrastructure and deploy &amp; configure applications, from a GoLang web application

For documentation & examples, see:

  [DevOps Orchestration Platform](https://software-automation.com/?page_id=132)

  [YouTube Video](https://www.youtube.com/watch?v=yuCfdvAKZrw)


Prereqs:
=========
GoLang go version go1.10.4 linux/amd64 must be installed with the GoLang paths set-up correctly
User running the application must have sudo rights (added to sudoers with NOPASSWD)

Usage:
=======
Navigate to <project_root>/webserver, then run:
  > ./build_and_run.sh
  
Go to http://localhost:6543 with either Firefox or Chrome

---
IMPORTANT NOTES ABOUT VIRTUALBOX
=========================
1. When using VirtualBox, the GoLang web server must be run as root:
	
	> sudo ./build_and_run.sh
	
2. Ensure that the IP address of the VirtualBox host network adapter (usually named 'vboxnet0') is not used as an IP address for any of the VMs and also ensure any chosen IP addresses for VirtualBox VMs are on the SAME SUBNET as vboxnet0.

3. As VirtualBox requires two network adapters (e.g. NAT adapter on 10.0.2.15 & host adapter on 192.168.3/24), the following file currently needs to be manually modified by ssh'ing into all Kubernetes nodes including the master:

/etc/systemd/system/kubelet.service.d/10-kubeadm.conf

And include (example of virtualbox host adapter IP 192.168.3.5):

Environment="KUBELET_EXTRA_ARGS=--node-ip=192.168.3.5
