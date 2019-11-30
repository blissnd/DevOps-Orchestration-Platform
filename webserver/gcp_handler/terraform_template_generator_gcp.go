package gcp_handler

import (
	//"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
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
func Get_vpc_cidr(subnet_address string) string {
	
	subnet_address_array := strings.Split(subnet_address, ".")
	
	return_value := subnet_address_array[0] + "." + subnet_address_array[1] + "." + subnet_address_array[2] + ".0/24"
	
  return return_value
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_subnet_cidr(subnet_address string) string {

  subnet_address_array := strings.Split(subnet_address, ".")
	
	return_value := subnet_address_array[0] + "." + subnet_address_array[1] + "." + subnet_address_array[2] + ".0/24"
	
  return return_value
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
func Store_gcp_credentials(gcp_map map[string]map[string]string, vm_name string) {
  
  credentials_string := ""
  key_id := gcp_map[vm_name]["key_id"]
  secret_key := gcp_map[vm_name]["secret_key"]
  
  current_dir, _ := os.Getwd()
  credentials_path := current_dir + "/../gcp_bootstrap/"
  
  credentials_string += "[default]\ngcp_access_key_id = " + key_id + "\ngcp_secret_access_key = " + secret_key + "\n"
  
  output_file_handle := []byte(credentials_string)
  err := ioutil.WriteFile(credentials_path + "gcp_credentials", output_file_handle, 0644)
  check(err)
  
  kv_store.Set_kv_entry(gcp_map[vm_name], "credentials_path", credentials_path + "gcp_credentials")
}
*/
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_gcp_ssh_keys(gcp_map map[string]string) {
  
  target_directory := "../gcp_bootstrap/instances/" + gcp_map["vm_name"] + "/"
  template_directory := "./templates/terraform/gcp/"
  
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  linux_command_line.Execute_command_line("cp " + template_directory + "terraform_script.sh.j3 " + target_directory + "terraform_script.sh")
  linux_command_line.Execute_command_line("cp " + template_directory + "generate_ssh_keys.sh.j3 " + target_directory + "generate_ssh_keys.sh")
  
  current_dir, _ := os.Getwd()
  os.Chdir(target_directory) 
  linux_command_line.Execute_command_line("bash ./generate_ssh_keys.sh")
  command_output, _ :=  linux_command_line.Execute_command_line("cat ./gcp_key.pub")
  command_output = strings.TrimSuffix(command_output, "\n")
  os.Chdir(current_dir)
  
  kv_store.Set_kv_entry(gcp_map, "key_name", gcp_map["vm_name"] + "_key")
  kv_store.Set_kv_entry(gcp_map, "public_key", command_output)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Store_gcp_ssh_keys_from_imported_VM(gcp_map map[string]string, vm_name string) {
	
	target_directory := "../gcp_bootstrap/instances/" + vm_name + "/"
	
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)

	bytestream := []byte(gcp_map["ssh_private_key"])
	err := ioutil.WriteFile(target_directory + "gcp_key", bytestream, 0600)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_terraform_resources(gcp_map map[string]string) {
  
  target_directory := "../gcp_bootstrap/instances/" + gcp_map["vm_name"] + "/"
  template_directory := "./templates/terraform/gcp/"
  
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  
  template_buffer := template_populator.Load_file(template_directory + "main.tf.j3")
  template_populator.Populate_gcp(template_buffer, gcp_map, target_directory + "main.tf")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_terraform_sec_group_resources(ajax_map map[string]string) {
  
  target_directory := "../gcp_firewall/"
  template_directory := "./templates/terraform/gcp/"
  
  ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  
  linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  linux_command_line.Execute_command_line("cp " + template_directory + "terraform_script.sh.j3 " + target_directory + "terraform_script.sh")

  ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

  file_buffer, err := ioutil.ReadFile(template_directory + "main.tf.firewall.j3")
  check(err)
  file_content_string := string(file_buffer)
  
  main_template_string := template_populator.Get_template_populated_string(file_content_string, ajax_map)
  
  /////////////////////////////////////////////////////////////////////////////////
  
  bytestream := []byte(main_template_string) 
  err = ioutil.WriteFile(target_directory + "main.tf", bytestream, 0644)
  check(err)
  
  //template_populator.Populate_gcp(template_buffer, gcp_map, target_directory + "main.tf")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_gcp_firewall_default_section(vm_name string) {
  
  target_path := "../gcp_bootstrap/instances/" + vm_name + "/main.tf"
  
  default_firewall_section_string := ""      
  default_firewall_section_string += "ports    = [\"22\"]"
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)
  
  bytestream := []byte(file_string)
	err = ioutil.WriteFile(target_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_modified_gcp_firewall_section(hypervisor_port_map map[string]map[string]string, vm_name string) {  
  
  target_path := "../gcp_bootstrap/instances/" + vm_name + "/main.tf"
  
  default_firewall_section_string := ""      
  default_firewall_section_string += "ports    = [\"22\""   
  
  for port_key, port_map  := range hypervisor_port_map {
    
    if port_map["state"] == "OPEN" {
      
      default_firewall_section_string += ", \"" + port_key + "\""       
    }
  }
  
  default_firewall_section_string += "]"
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
    
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)  
  
  bytestream := []byte(file_string)
	err = ioutil.WriteFile(target_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_separate_gcp_firewall_section(hypervisor_port_map map[string]map[string]string) {  
  
  target_path := "../gcp_firewall/main.tf"
  
  ////////////////////// TCP ///////////////////////////

  default_firewall_section_string := ""      
  default_firewall_section_string += "ports    = [\"22\""   
  
  for _, firewall_rule_map  := range hypervisor_port_map {
    
    if firewall_rule_map["state"] == "OPEN" && firewall_rule_map["protocol"] == "TCP" {
      
      default_firewall_section_string += ", \"" + firewall_rule_map["port"] + "\""       
    }
  }
  
  default_firewall_section_string += "]"
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
    
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_tcp_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)  
  
  bytestream := []byte(file_string)
  err = ioutil.WriteFile(target_path, bytestream, 0644)
  check(err)


  ////////////////////// UDP ///////////////////////////

  default_firewall_section_string = ""      
  default_firewall_section_string = "ports    = [\"22\""   
  
  for _, firewall_rule_map  := range hypervisor_port_map {
    
    if firewall_rule_map["state"] == "OPEN" && firewall_rule_map["protocol"] == "UDP" {
      
      default_firewall_section_string += ", \"" + firewall_rule_map["port"] + "\""       
    }
  }
  
  default_firewall_section_string += "]"
  
  file_content, err = ioutil.ReadFile(target_path)
  file_string = string(file_content)
    
  regex_obj = regexp.MustCompile(regexp.QuoteMeta("{<firewall_udp_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)  
  
  bytestream = []byte(file_string)
  err = ioutil.WriteFile(target_path, bytestream, 0644)
  check(err)

}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
