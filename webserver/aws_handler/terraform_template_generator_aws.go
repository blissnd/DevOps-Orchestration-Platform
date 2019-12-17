// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package aws_handler

import (
	//"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
	//"path/filepath"
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
func Store_aws_credentials(aws_map map[string]map[string]string, vm_name string) { 	//////// Deprecated
  
  credentials_string := ""
  key_id := aws_map[vm_name]["key_id"]
  secret_key := aws_map[vm_name]["secret_key"]
  
  current_dir, _ := os.Getwd()
  credentials_path := current_dir + "/../aws_bootstrap/"
  
  credentials_string += "[default]\naws_access_key_id = " + key_id + "\naws_secret_access_key = " + secret_key + "\n"
  
  output_file_handle := []byte(credentials_string)
  err := ioutil.WriteFile(credentials_path + "aws_credentials", output_file_handle, 0644)
  check(err)    
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Store_aws_ssh_keys_from_imported_VM(aws_map map[string]string, vm_name string) {
	
	target_directory := "../aws_bootstrap/instances/" + vm_name + "/"	
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)

	bytestream := []byte(aws_map["ssh_private_key"])
	err := ioutil.WriteFile(target_directory + "aws_key", bytestream, 0600)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_aws_ssh_keys(aws_map map[string]string, vm_name string) {
  
  target_directory := "../aws_bootstrap/instances/" + vm_name + "/"
  template_directory := "./templates/terraform/aws/"
  
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  linux_command_line.Execute_command_line("cp " + template_directory + "terraform_script.sh.j3 " + target_directory + "terraform_script.sh")
  linux_command_line.Execute_command_line("cp " + template_directory + "generate_ssh_keys.sh.j3 " + target_directory + "generate_ssh_keys.sh")
  
  current_dir, _ := os.Getwd()
  os.Chdir(target_directory) 
  linux_command_line.Execute_command_line("bash ./generate_ssh_keys.sh")
  command_output, _ :=  linux_command_line.Execute_command_line("cat ./aws_key.pub")
  command_output = strings.TrimSuffix(command_output, "\n")
  os.Chdir(current_dir)
  
  kv_store.Set_kv_entry(aws_map, "key_name", aws_map["vm_name"] + "_key")
  kv_store.Set_kv_entry(aws_map, "public_key", command_output)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_terraform_resources(aws_map map[string]map[string]string, vm_name string) {
  
  target_directory := "../aws_bootstrap/instances/" + vm_name + "/"
  template_directory := "./templates/terraform/aws/"
  
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  
  file_buffer, err := ioutil.ReadFile(template_directory + "main.tf.j3")
	check(err)
	file_content_string := string(file_buffer)
  
  main_template_string := template_populator.Get_template_populated_string(file_content_string, aws_map[vm_name])
  
  /////////////////////////////////////////////////////////////////////////////////
  
  file_buffer2, err := ioutil.ReadFile(template_directory + "instance.j3")
	check(err)
	file_content_string2 := string(file_buffer2)
	instance_template_string := ""
	
  instance_template_string += template_populator.Get_template_populated_string(file_content_string2, aws_map[vm_name])
  
  final_file_string := template_populator.Get_regex_populated_string(main_template_string, "{<instance_section>}", instance_template_string)
  
  bytestream := []byte(final_file_string)	
	err = ioutil.WriteFile(target_directory + "main.tf", bytestream, 0644)
	check(err)
	
  //template_populator.Populate_aws(template_buffer, aws_map, target_directory + "main.tf")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generate_terraform_sec_group_resources(ajax_map map[string]string) {
  
  target_directory := "../aws_firewall/"
  template_directory := "./templates/terraform/aws/"
  
  ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  
  linux_command_line.Execute_command_line("mkdir -p " + target_directory)
  linux_command_line.Execute_command_line("cp " + template_directory + "terraform_script.sh.j3 " + target_directory + "terraform_script.sh")

  ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

  file_buffer, err := ioutil.ReadFile(template_directory + "security_groups.tf.j3")
  check(err)
  file_content_string := string(file_buffer)
  
  main_template_string := template_populator.Get_template_populated_string(file_content_string, ajax_map)
  
  /////////////////////////////////////////////////////////////////////////////////
  
  bytestream := []byte(main_template_string) 
  err = ioutil.WriteFile(target_directory + "main.tf", bytestream, 0644)
  check(err)
  
  //template_populator.Populate_aws(template_buffer, aws_map, target_directory + "main.tf")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_aws_firewall_default_section(vm_name string) {
  
  target_path := "../aws_bootstrap/instances/" + vm_name + "/main.tf"
  
  default_firewall_section_string := ""    
  
  default_firewall_section_string += "ingress { \n"
  default_firewall_section_string += "      from_port = \"22\"\n"
  default_firewall_section_string += "      to_port   = \"22\"\n"
  default_firewall_section_string += "      protocol  = \"tcp\"\n"
  default_firewall_section_string += "      cidr_blocks = [\"0.0.0.0/0\"]\n"
  default_firewall_section_string += "    }\n"  
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)
  
  bytestream := []byte(file_string)
	err = ioutil.WriteFile(target_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_modified_aws_firewall_section(hypervisor_port_map map[string]map[string]string, vm_name string) {  
  
  target_path := "../aws_bootstrap/instances/" + vm_name + "/main.tf"
  
  default_firewall_section_string := ""    
  
  //default_firewall_section_string += "    {\n"
  //default_firewall_section_string += "      from_port = \"22\"\n"
  //default_firewall_section_string += "      to_port   = \"22\"\n"
  //default_firewall_section_string += "      protocol  = \"tcp\"\n"
  //default_firewall_section_string += "      cidr_blocks = [\"0.0.0.0/0\"]\n"
  //default_firewall_section_string += "    }"  
  
  for port_key, port_map  := range hypervisor_port_map {
    
    if port_map["state"] == "OPEN" {
            
      default_firewall_section_string += "ingress { \n"
      default_firewall_section_string += "      from_port = \"" + port_key + "\"\n"
      default_firewall_section_string += "      to_port   = \"" + port_key + "\"\n"
      default_firewall_section_string += "      protocol  = \"tcp\"\n"
      default_firewall_section_string += "      cidr_blocks = [\"" + port_map["source_cidr"] + "\"]\n"
      default_firewall_section_string += "    }"      
    }  
  }    
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
    
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)  
  
  bytestream := []byte(file_string)
	err = ioutil.WriteFile(target_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Populate_separate_aws_firewall_section(hypervisor_port_map map[string]map[string]string) {  
  
  target_path := "../aws_firewall/main.tf"
  
  default_firewall_section_string := ""  
  default_firewall_section_string += "ingress { \n"
  
  default_firewall_section_string += "      from_port = \"22\"\n"
  default_firewall_section_string += "      to_port   = \"22\"\n"
  default_firewall_section_string += "      protocol  = \"tcp\"\n"
  default_firewall_section_string += "      cidr_blocks = [\"0.0.0.0/0\"]\n"
  default_firewall_section_string += "    }\n"  
  
  for _, firewall_rule_map  := range hypervisor_port_map {
    
    if firewall_rule_map["state"] == "OPEN" {
      
      default_firewall_section_string += "    ingress {\n"      
      default_firewall_section_string += "      from_port = \"" + firewall_rule_map["port"] + "\"\n"
      default_firewall_section_string += "      to_port   = \"" + firewall_rule_map["port"] + "\"\n"
      default_firewall_section_string += "      protocol  = \"" + firewall_rule_map["protocol"] + "\"\n"
      default_firewall_section_string += "      cidr_blocks = [\"" + firewall_rule_map["source_cidr"] + "\"]\n"
      default_firewall_section_string += "    }\n"      
    }  
  }
  
  file_content, err := ioutil.ReadFile(target_path)
  file_string := string(file_content)
    
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<firewall_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, default_firewall_section_string)  
  
  bytestream := []byte(file_string)
  err = ioutil.WriteFile(target_path, bytestream, 0644)
  check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
