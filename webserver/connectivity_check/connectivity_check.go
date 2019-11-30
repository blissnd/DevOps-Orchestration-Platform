package connectivity_check

import (
	"fmt"
	//"io/ioutil"
	"strings"
	"regexp"
	//"os/exec"
	"os"
	//"../logging"
  "../kv_store"
  "../linux_command_line"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
        fmt.Println(e.Error() + "\n")
    }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Clean_ssh(ip_address string, vm_name string) {
  
  // Copy private key from vagrant directory
  Copy_vagrant_private_key(vm_name)
  
  command_string := "rm ~/.ssh/known_hosts"
	log_path := "../Logs/ansible/ssh_fix_log"
	linux_command_line.Execute_command_in_background(command_string, ".", log_path)
  
  command_string = "ssh-keygen -f /root/.ssh/known_hosts -R " + ip_address
  fmt.Println(command_string)
	log_path = "../Logs/ansible/ssh_fix_log"
  linux_command_line.Execute_command_in_background(command_string, ".", log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Copy_vagrant_private_key(vm_name string) {

  source_path := "../vm_bootstrap/VMs/" + vm_name + "/.vagrant/machines/" + vm_name + "/virtualbox/private_key"
  dest_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  linux_command_line.Execute_command_line("cp " + source_path + " " + dest_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Copy_aws_ssh_keys(vm_name string) {

  source_path := "../aws_bootstrap/instances/" + vm_name + "/aws_key"
  dest_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  linux_command_line.Execute_command_line("cp " + source_path + " " + dest_path)
  
  source_path = "../aws_bootstrap/instances/" + vm_name + "/aws_key.pub"
  dest_path = "../ansible/ssh_keys/" + vm_name + "/aws_key.pub"
  linux_command_line.Execute_command_line("cp " + source_path + " " + dest_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Copy_gcp_ssh_keys(vm_name string) {

  source_path := "../gcp_bootstrap/instances/" + vm_name + "/gcp_key"
  dest_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  linux_command_line.Execute_command_line("cp " + source_path + " " + dest_path)
  
  source_path = "../gcp_bootstrap/instances/" + vm_name + "/gcp_key.pub"
  dest_path = "../ansible/ssh_keys/" + vm_name + "/gcp_key.pub"
  linux_command_line.Execute_command_line("cp " + source_path + " " + dest_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Store_aws_vm_resources_in_DB(aws_map map[string]map[string]string, vm_name string) {
  
  target_directory := "../aws_bootstrap/instances/" + vm_name + "/"
  current_dir, _ := os.Getwd()
  os.Chdir(target_directory) 
  
  command_output, _ :=  linux_command_line.Execute_command_line("terraform output public_dns")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(aws_map[vm_name], "public_dns", command_output)
  
  command_output, _ =  linux_command_line.Execute_command_line("terraform output private_dns")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(aws_map[vm_name], "private_dns", command_output)
  
  command_output, _ =  linux_command_line.Execute_command_line("terraform output public_ip_address")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(aws_map[vm_name], "public_ip_address", command_output)
  
  command_output, _ =  linux_command_line.Execute_command_line("terraform output private_ip_address")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(aws_map[vm_name], "private_ip_address", command_output)
    
  os.Chdir(current_dir)
  
  kv_store.Export_to_json(aws_map, "./databases/aws_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Store_gcp_vm_resources_in_DB(gcp_map map[string]map[string]string, vm_name string) {
  
  target_directory := "../gcp_bootstrap/instances/" + vm_name + "/"
  current_dir, _ := os.Getwd()
  os.Chdir(target_directory) 
  
  command_output, _ :=  linux_command_line.Execute_command_line("terraform output public_ip")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(gcp_map[vm_name], "public_ip_address", command_output)
  
  command_output, _ =  linux_command_line.Execute_command_line("terraform output private_ip")
  command_output = strings.TrimSuffix(command_output, "\n")
  kv_store.Set_kv_entry(gcp_map[vm_name], "private_ip_address", command_output)
    
  os.Chdir(current_dir)
  
  kv_store.Export_to_json(gcp_map, "./databases/gcp_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Run_remote_exec(userid string, ip_address string, ssh_private_key_path string) {

  ssh_command := "ssh" + " -o StrictHostKeyChecking=no " + userid + "@" + ip_address + " -i " + ssh_private_key_path + " sudo apt-get update"  
  linux_command_line.Execute_command_line(ssh_command)

  ssh_command = "ssh" + " -o StrictHostKeyChecking=no " + userid + "@" + ip_address + " -i " + ssh_private_key_path + " sudo yum update"  
  linux_command_line.Execute_command_line(ssh_command)
  
  ssh_command = "ssh" + " -o StrictHostKeyChecking=no " + userid + "@" + ip_address + " -i " + ssh_private_key_path + " sudo apt-get install -y python"  
  linux_command_line.Execute_command_line(ssh_command)
  
  ssh_command = "ssh" + " -o StrictHostKeyChecking=no " + userid + "@" + ip_address + " -i " + ssh_private_key_path + " sudo yum install -y python"  
  linux_command_line.Execute_command_line(ssh_command)
}  
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Check_connectivity(top_level_vm_map map[string]map[string]string, vm_name string, platform_type string) {
  
  ip_address_type := ""
  
  if platform_type == "aws" {
  
    ip_address_type = "public_ip_address"
    
    if top_level_vm_map[vm_name]["OS"] == "centos/7" {
      kv_store.Set_kv_entry(top_level_vm_map[vm_name], "user_id", "centos")
    } else {
      kv_store.Set_kv_entry(top_level_vm_map[vm_name], "user_id", "ubuntu")
    }
     
  } else  if platform_type == "gcp" {
  
    ip_address_type = "public_ip_address"
    
  } else if platform_type == "virtualbox" {
  
    ip_address_type = "ip_address"
    kv_store.Set_kv_entry(top_level_vm_map[vm_name], "user_id", "vagrant")
  }
  
  Clean_ssh(top_level_vm_map[vm_name][ip_address_type], vm_name)
  
  ip_address := top_level_vm_map[vm_name][ip_address_type]
  userid := top_level_vm_map[vm_name]["user_id"]
  ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  
  command_output, err := linux_command_line.Execute_command_line("ping -c 3 " + ip_address)
  check(err)
  
  regex_obj := regexp.MustCompile(".*Command timed out | Host Unreachable.*")
	match_array := regex_obj.FindStringSubmatch(command_output)
  
  if len(match_array) != 0 {
    fmt.Println("Ping failed")
    kv_store.Set_kv_entry(top_level_vm_map[vm_name], "latest_ping", "Fail")
  } else {
    kv_store.Set_kv_entry(top_level_vm_map[vm_name], "latest_ping", "Pass")
  }
  /////////////////////////////////////////////////////////////////////////////////////////////////////////////
  
  ssh_command := "ssh" + " -o StrictHostKeyChecking=no " + userid + "@" + ip_address + " -i " + ssh_private_key_path + " ls /home"
  
  command_output2, err2 := linux_command_line.Execute_command_line(ssh_command)
  check(err2)
  
  fmt.Println(ssh_command + "\n")
  fmt.Println(command_output2 + "\n")
  
  regex_obj = regexp.MustCompile(".*Command timed out.*")
	match_array = regex_obj.FindStringSubmatch(command_output2)
  
  if len(match_array) != 0 || command_output2 == "" {
    fmt.Println("SSH failed")
    kv_store.Set_kv_entry(top_level_vm_map[vm_name], "latest_ssh", "Fail")
  } else {
    kv_store.Set_kv_entry(top_level_vm_map[vm_name], "latest_ssh", "Pass")
  }
  /////////////////////////////////////////////////////////////////////////////////////////////////////////////
  
  if platform_type == "aws" {
    kv_store.Export_to_json(top_level_vm_map, "./databases/aws_store.json")
  } else  if platform_type == "gcp" {
    kv_store.Export_to_json(top_level_vm_map, "./databases/gcp_store.json")
  } else if platform_type == "virtualbox" {
    kv_store.Export_to_json(top_level_vm_map, "./databases/vm_store.json")
  }
  
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Check_docker_container_connectivity(docker_map map[string]map[string]string, vm_name string, container_name string, vm_ip_address string) {
  
  external_ssh_port := docker_map[container_name]["external_port_0"]
  
  Clean_ssh(vm_ip_address, vm_name)
  
  ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"

  ssh_command := "ssh" + " -o StrictHostKeyChecking=no vagrant@" + vm_ip_address + " -i " + ssh_private_key_path + " -p " + external_ssh_port +" ls /home"
  
  command_output2, err2 := linux_command_line.Execute_command_line(ssh_command)
  check(err2)
  
  fmt.Println(ssh_command + "\n")
  fmt.Println(command_output2 + "\n")
  
  regex_obj := regexp.MustCompile(".*Command timed out.*")
	match_array := regex_obj.FindStringSubmatch(command_output2)
  
  if len(match_array) != 0 || command_output2 == "" {
    fmt.Println("SSH failed")
    kv_store.Set_kv_entry(docker_map[container_name], "latest_ssh", "Fail")
  } else {
    kv_store.Set_kv_entry(docker_map[container_name], "latest_ssh", "Pass")
  }
  /////////////////////////////////////////////////////////////////////////////////////////////////////////////
  
  kv_store.Export_to_json(docker_map, "./databases/docker_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
