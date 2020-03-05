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
IMPORTANT NOTE ABOUT VIRTUALBOX: Ensure that the IP address of the VirtualBox host network adapter (usually named 'vboxnet0') 
 is not used as an IP address for any of the VMs and also ensure any chosen IP addresses for VirtualBox VMs are on the SAME SUBNET as vboxnet0.
 