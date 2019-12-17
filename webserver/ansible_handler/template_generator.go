// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package ansible_handler

import (
	//"fmt"
	"io/ioutil"
	"strings"
  "strconv"
	//"regexp"
	"os"
	"../general_utility_web_handler"
	"../linux_command_line"
  "../template_populator"
  "../kv_store"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func indent(indentation_level int) string {
  string_to_indent := ""
  
  for current_position := 0;  current_position < indentation_level; current_position++ {
    string_to_indent += "  "
  }
  return string_to_indent
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_inventory(vm_name string, top_level_vm_map map[string]map[string]string, output_dir string, docker_container_name string,
													server_selection string, client_list_array []string) {   
  
  //// Copy ansible config file across
  ansible_config_template_path := "./templates/ansible/ansible.cfg.j3"
  target_path := output_dir + "/ansible.cfg"
  linux_command_line.Execute_command_line("cp " + ansible_config_template_path + " " + target_path)   
  
  /// Fetch docker map ///
  docker_map := make(map[string]map[string]string)
  
  if _, err := os.Stat("./databases/docker_store.json"); err == nil {
    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  }

  /// Fetch master/concatenated map ///
  master_map := make(map[string]map[string]string)
  
  if _, err := os.Stat("./databases/master_map.json"); err == nil {
    master_map = kv_store.Create_top_level_map_from_json_file("./databases/master_map.json")
  }

  ///////////////////////////////////////////////////    

  /// Fetch Nexus Registry Map ///
  nexus_registry_map := make(map[string]string)
  
  if _, err := os.Stat("./databases/nexus_registry.json"); err == nil {
    nexus_registry_map = kv_store.Create_single_level_map_from_json_file("./databases/nexus_registry.json")
  }
  ////////////////////////////////////

  inventory_string := ""

  var vm_name_array []string
  
  for vm_name, _ := range top_level_vm_map {
    vm_name_array = append(vm_name_array, vm_name)
  }    

  inventory_string = ""
  indentation_level := 0
  
  inventory_string += indent(indentation_level) + "masters:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "hosts:\n"
  indentation_level++
  
  if docker_container_name != "" {
    inventory_string += indent(indentation_level) + top_level_vm_map[vm_name]["public_ip_address"] + ":" + docker_map[docker_container_name]["external_port_0"] + ":\n"
  } else {
    inventory_string += indent(indentation_level) + top_level_vm_map[vm_name]["public_ip_address"] + ":22:\n"
  }
  
  indentation_level++
  inventory_string += indent(indentation_level) + "ansible_ssh_private_key_file:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  indentation_level -= 2
  inventory_string += indent(indentation_level) + "vars:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "home_dir:  /home/" + top_level_vm_map[vm_name]["user_id"] + "\n"
  inventory_string += indent(indentation_level) + "username:  " + top_level_vm_map[vm_name]["user_id"] + "\n"    
  inventory_string += indent(indentation_level) + "vm_name:  " + vm_name + "\n"
  inventory_string += indent(indentation_level) + "ssh_private_key:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  inventory_string +=  indent(indentation_level) + "ssh_other_keys:\n"
  indentation_level++
  
  for vm := range vm_name_array {
    inventory_string +=  indent(indentation_level) + "- { src:  ../ssh_keys/" + vm_name_array[vm] + "/private_key, dest: " + vm_name_array[vm] + ".key }\n"
  }
  
  indentation_level --
  inventory_string += indent(indentation_level) + "os_name:  "
  
  os_name := ""

  if docker_container_name != "" {
    os_name = docker_map[docker_container_name]["OS"]
  } else {
    os_name = top_level_vm_map[vm_name]["OS"]
  }

  if os_name == "centos/7" {
    inventory_string += "centos\n"
  } else {
    inventory_string += "ubuntu\n"
  }
  
  inventory_string += indent(indentation_level) + "ip_address:  " + top_level_vm_map[vm_name]["private_ip_address"] + "\n"
  ip_address_array := strings.Split(top_level_vm_map[vm_name]["private_ip_address"], ".")
  
  if len(ip_address_array) == 4 {
    inventory_string += indent(indentation_level) + "ip_address_part_1: " + ip_address_array[0] + "\n"
    inventory_string += indent(indentation_level) + "ip_address_part_2: " + ip_address_array[1] + "\n"
    inventory_string += indent(indentation_level) + "ip_address_part_3: " + ip_address_array[2] + "\n"
    inventory_string += indent(indentation_level) + "ip_address_part_4: " + ip_address_array[3] + "\n"
  }
  
  server_ip_address := master_map[server_selection]["private_ip_address"]

  inventory_string += indent(indentation_level) + "platform_type: " + top_level_vm_map[vm_name]["platform_type"] + "\n"
  inventory_string += indent(indentation_level) + "fqdn: " + top_level_vm_map[vm_name]["fqdn"] + "\n"
  inventory_string += indent(indentation_level) + "server: " + server_ip_address + "\n"
  inventory_string += indent(indentation_level) + "nexus_registry_dns:  " + nexus_registry_map["vm_name"] + "." + nexus_registry_map["fqdn"] + "\n"
  inventory_string += indent(indentation_level) + "container_name: " + docker_container_name + "\n"
  inventory_string += indent(indentation_level) + "image_name: " + docker_map[docker_container_name]["image_name"] + "\n"
  
  ldap_dn := ""
  
  fqdn_array := strings.Split(top_level_vm_map[vm_name]["fqdn"], ".")
  
  if len(fqdn_array) == 0 {
    ldap_dn += "dc=" + top_level_vm_map[vm_name]["fqdn"]
  } else {
    for current_index := 0; current_index < len(fqdn_array); current_index++ {
    
      inventory_string += indent(indentation_level) + "fqdn_part_" + strconv.Itoa(current_index) + ":  " + fqdn_array[current_index] + "\n"
      ldap_dn += "dc=" + fqdn_array[current_index]
      
      if current_index < len(fqdn_array) - 1 {
        ldap_dn += ","
      }
    }
  }
  
  inventory_string += indent(indentation_level) + "ldap_dn: " + ldap_dn + "\n"
  
  inventory_string +=  indent(indentation_level) + "other_hosts:\n"
  indentation_level++
  
  for _, client_vm_name := range client_list_array {      
    
    host_ip_address_array:= strings.Split(master_map[client_vm_name]["private_ip_address"], ".")      
  
    inventory_string +=  indent(indentation_level) + "- { ip_address: " + master_map[client_vm_name]["private_ip_address"] + ", "
    
    if len(host_ip_address_array) == 4 {
      inventory_string +=  indent(indentation_level) + "ip_address_part_1: " +  host_ip_address_array[0] + ", "
      inventory_string +=  indent(indentation_level) + "ip_address_part_2: " +  host_ip_address_array[1] + ", "
      inventory_string +=  indent(indentation_level) + "ip_address_part_3: " +  host_ip_address_array[2] + ", "
      inventory_string +=  indent(indentation_level) + "ip_address_part_4: " +  host_ip_address_array[3] + ", "
    }
    
    inventory_string +=  indent(indentation_level) + "fqdn: " +  master_map[client_vm_name]["fqdn"] + ", "
    inventory_string +=  indent(indentation_level) + "vm_name: " +  master_map[client_vm_name]["vm_name"]
    
    inventory_string += " }\n"
  }
  
  indentation_level--
  
  //#####################################################
  output_file_handle := []byte(inventory_string)

  inventory_file_name := ""

  if docker_container_name != "" {
    inventory_file_name = "/inventory_" + vm_name + "_" + docker_container_name + ".yml"
  } else {
    inventory_file_name = "/inventory_" + vm_name + ".yml"
  }

  err := ioutil.WriteFile(output_dir + inventory_file_name, output_file_handle, 0644)
  general_utility_web_handler.Check(err)
 
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_inventory_for_docker_image_commit(vm_map map[string]string, ip_address string, docker_submap map[string]string, output_dir string) {    
  
	input_stream, _ := ioutil.ReadFile("../ansible/nexus/nexus_admin_password")	
	nexus_admin_password := strings.TrimSpace(string(input_stream))

  //// Copy ansible config file across
  ansible_config_template_path := "./templates/ansible/ansible.cfg.j3"
  target_path := output_dir + "/ansible.cfg"
  linux_command_line.Execute_command_line("cp " + ansible_config_template_path + " " + target_path)
    
  /// Fetch Nexus Registry Map ///
  nexus_registry_map := make(map[string]string)
  
  if _, err := os.Stat("./databases/nexus_registry.json"); err == nil {
    nexus_registry_map = kv_store.Create_single_level_map_from_json_file("./databases/nexus_registry.json")
  }
  ////////////////////////////////////
  
  vm_name := docker_submap["vm_name"]
  container_name := docker_submap["docker_container_name"]
  image_name := docker_submap["image_name"]
    
  inventory_string := ""

  inventory_string = ""
  indentation_level := 0
  
  inventory_string += indent(indentation_level) + "masters:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "hosts:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + ip_address + ":22:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "ansible_ssh_private_key_file:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  indentation_level -= 2
  inventory_string += indent(indentation_level) + "vars:\n"
  indentation_level++  
  inventory_string += indent(indentation_level) + "home_dir:  /home/" + vm_map["user_id"] + "\n"
  inventory_string += indent(indentation_level) + "username:  " + vm_map["user_id"] + "\n"
  inventory_string += indent(indentation_level) + "vm_name:  " + vm_name + "\n"
  inventory_string += indent(indentation_level) + "fqdn:  " + vm_map["fqdn"] + "\n"
  inventory_string += indent(indentation_level) + "nexus_registry:  " + nexus_registry_map["ip_address"] + "\n"
  inventory_string += indent(indentation_level) + "nexus_registry_dns:  " + nexus_registry_map["vm_name"] + "." + nexus_registry_map["fqdn"] + "\n"
  inventory_string += indent(indentation_level) + "nexus_admin_password:  " + nexus_admin_password + "\n"
  inventory_string += indent(indentation_level) + "container_name:  " + container_name + "\n"
  inventory_string += indent(indentation_level) + "image_name:  " + image_name + "\n"
  inventory_string += indent(indentation_level) + "ssh_private_key:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  
  inventory_string += indent(indentation_level) + "os_name:  "
  
  if docker_submap["OS"] == "centos/7" {
    inventory_string += "centos\n"
  } else {
    inventory_string += "ubuntu\n"
  }
  
  output_file_handle := []byte(inventory_string)

  err := ioutil.WriteFile(output_dir + "/inventory_" + vm_name + "_" + container_name + ".yml", output_file_handle, 0644)
  general_utility_web_handler.Check(err)
    
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_inventory_for_docker_container_launch(vm_map map[string]string, ip_address string, docker_submap map[string]string, output_dir string) {
  
  input_stream, _ := ioutil.ReadFile("../ansible/nexus/nexus_admin_password")	
	nexus_admin_password := strings.TrimSpace(string(input_stream))

  //// Copy ansible config file across
  ansible_config_template_path := "./templates/ansible/ansible.cfg.j3"
  target_path := output_dir + "/ansible.cfg"
  linux_command_line.Execute_command_line("cp " + ansible_config_template_path + " " + target_path)
  
  /// Fetch Nexus Registry Map ///
  nexus_registry_map := make(map[string]string)
  
  if _, err := os.Stat("./databases/nexus_registry.json"); err == nil {
    nexus_registry_map = kv_store.Create_single_level_map_from_json_file("./databases/nexus_registry.json")
  }
  ////////////////////////////////////

  vm_name := docker_submap["vm_name"]
  docker_container_name := docker_submap["docker_container_name"]
  image_name := docker_submap["image_name"]
   
  inventory_string := ""

  inventory_string = ""
  indentation_level := 0
  
  inventory_string += indent(indentation_level) + "masters:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "hosts:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + ip_address + ":22:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "ansible_ssh_private_key_file:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  indentation_level -= 2
  inventory_string += indent(indentation_level) + "vars:\n"
  indentation_level++
  inventory_string += indent(indentation_level) + "home_dir:  /home/" + vm_map["user_id"] + "\n"
  inventory_string += indent(indentation_level) + "username:  " + vm_map["user_id"] + "\n"    
  inventory_string += indent(indentation_level) + "vm_name:  " + vm_name + "\n"
  inventory_string += indent(indentation_level) + "nexus_registry:  " + nexus_registry_map["ip_address"] + "\n"
  inventory_string += indent(indentation_level) + "nexus_registry_dns:  " + nexus_registry_map["vm_name"] + "." + nexus_registry_map["fqdn"] + "\n"
  inventory_string += indent(indentation_level) + "nexus_admin_password:  " + nexus_admin_password + "\n"
  inventory_string += indent(indentation_level) + "container_name:  " + docker_container_name + "\n"
  inventory_string += indent(indentation_level) + "image_name:  " + image_name + "\n"
  inventory_string += indent(indentation_level) + "ssh_private_key:  "
  inventory_string += "../ssh_keys/" + vm_name + "/private_key\n"
  
  inventory_string += indent(indentation_level) + "os_name:  "
  
  if docker_submap["OS"] == "centos/7" {
    inventory_string += "centos\n"
  } else {
    inventory_string += "ubuntu\n"
  }
  
  output_file_handle := []byte(inventory_string)

  err := ioutil.WriteFile(output_dir + "/inventory_" + vm_name + "_" + docker_container_name + ".yml", output_file_handle, 0644)
  general_utility_web_handler.Check(err)
    
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_docker_resources(docker_submap map[string]string, vm_name string, container_name string) {

	linux_command_line.Execute_command_line("mkdir -p ../Logs/docker")

  dockerfile_template_path := ""
  public_key_src_path := ""
  dockerfile_dest_path := ""
  docker_bootstrap_script_src_path := ""
  docker_resource_dest_dir := ""
  docker_build_src_path := ""
  docker_build_dest_path := ""
  docker_run_src_path := ""
  docker_run_dest_path := ""
  sudo_override_src_path := ""
  
  public_key_src_path = "../ansible/ssh_keys/" + vm_name + "/authorized_keys"
  docker_build_src_path = "./templates/docker/docker_build.sh.j3"
  docker_run_src_path = "./templates/docker/docker_run.sh.j3"
  sudo_override_src_path = "./templates/docker/myOverrides"
  
  if docker_submap["OS"] == "centos/7" {
  
    linux_command_line.Execute_command_line("mkdir -p ../docker/centos")
    dockerfile_template_path = "./templates/docker/dockerfile_centos.j3"
    docker_bootstrap_script_src_path = "./templates/docker/docker_bootstrap_script_centos.sh.j3"
    dockerfile_dest_path = "../docker/centos/dockerfile"
    docker_build_dest_path = "../docker/centos/docker_build.sh"
    docker_run_dest_path = "../docker/centos/docker_run.sh"
    docker_resource_dest_dir = "../docker/centos/"
    
  } else if docker_submap["OS"] == "ubuntu/xenial64" {
    
    linux_command_line.Execute_command_line("mkdir -p ../docker/ubuntu")
    dockerfile_template_path = "./templates/docker/dockerfile_ubuntu.j3"
    docker_bootstrap_script_src_path = "./templates/docker/docker_bootstrap_script_ubuntu.sh.j3"
    dockerfile_dest_path = "../docker/ubuntu/dockerfile"
    docker_build_dest_path = "../docker/ubuntu/docker_build.sh"
    docker_run_dest_path = "../docker/ubuntu/docker_run.sh"
    docker_resource_dest_dir = "../docker/ubuntu/"
  }
  
  template_buffer := template_populator.Load_file(dockerfile_template_path)
  template_populator.Populate_docker(template_buffer, docker_submap, dockerfile_dest_path)

  template_myoverrides_buffer := template_populator.Load_file(sudo_override_src_path)
  template_populator.Populate_docker(template_myoverrides_buffer, docker_submap, docker_resource_dest_dir + "myOverrides")
  
  docker_build_buffer := template_populator.Load_file(docker_build_src_path)
  template_populator.Populate_docker(docker_build_buffer, docker_submap, docker_build_dest_path)
  
  docker_run_buffer := template_populator.Load_file(docker_run_src_path)
  template_populator.Populate_docker(docker_run_buffer, docker_submap, docker_run_dest_path)
  
  template_populator.Populate_docker_port_section(docker_run_dest_path, docker_submap)
  
  linux_command_line.Execute_command_line("cp " + public_key_src_path + " " + docker_resource_dest_dir)
  linux_command_line.Execute_command_line("cp " + docker_bootstrap_script_src_path + " " + docker_resource_dest_dir + "docker_bootstrap_script.sh")  
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_docker_resources_for_image(docker_submap map[string]string, vm_name string, container_name string) {

  input_stream, _ := ioutil.ReadFile("../ansible/nexus/nexus_admin_password") 
  nexus_admin_password := strings.TrimSpace(string(input_stream))

  public_key_src_path := ""
  docker_bootstrap_script_src_path := ""
  docker_resource_dest_dir := ""
  docker_run_src_path := ""
  docker_run_dest_path := ""      
  
  public_key_src_path = "../ansible/ssh_keys/" + vm_name + "/authorized_keys"
  docker_run_src_path = "./templates/docker/docker_run_from_image.sh.j3"
  
  if docker_submap["OS"] == "centos/7" {
  
    linux_command_line.Execute_command_line("mkdir -p ../docker/centos")
    docker_bootstrap_script_src_path = "./templates/docker/docker_bootstrap_script_centos.sh.j3"
    docker_run_dest_path = "../docker/centos/docker_run.sh"
    docker_resource_dest_dir = "../docker/centos/"
    
  } else if docker_submap["OS"] == "ubuntu/xenial64" {
    
    linux_command_line.Execute_command_line("mkdir -p ../docker/ubuntu")
    docker_bootstrap_script_src_path = "./templates/docker/docker_bootstrap_script_ubuntu.sh.j3"
    docker_run_dest_path = "../docker/ubuntu/docker_run.sh"
    docker_resource_dest_dir = "../docker/ubuntu/"
  }
    
  docker_run_buffer := template_populator.Load_file(docker_run_src_path)
  template_populator.Populate_docker(docker_run_buffer, docker_submap, docker_run_dest_path)
  
  template_populator.Populate_docker_port_section(docker_run_dest_path, docker_submap)
  template_populator.Populate_docker_registry_section(docker_run_dest_path, docker_submap, nexus_admin_password)
  
  linux_command_line.Execute_command_line("cp " + public_key_src_path + " " + docker_resource_dest_dir)
  linux_command_line.Execute_command_line("cp " + docker_bootstrap_script_src_path + " " + docker_resource_dest_dir + "docker_bootstrap_script.sh")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

