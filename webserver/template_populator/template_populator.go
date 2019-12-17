// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package template_populator

import (
	//"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	//"os/exec"
	"os"
	"../linux_command_line"
	"../kv_store"
)
/////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
        panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////
func Load_file(pathname string) []string {
	file_content, err := ioutil.ReadFile(pathname)
	check(err)
	file_buffer := strings.Split(string(file_content), "\n")
	return file_buffer
}
/////////////////////////////////////////////////////////////////////////////////
func Load_file_as_string(pathname string) string {
	file_content, err := ioutil.ReadFile(pathname)
	check(err)
	file_buffer := string(file_content)
	return file_buffer
}
/////////////////////////////////////////////////////////////////////////////////
func Populate(file_buffer []string, vm_map map[string]string, output_file string) {
	output_string := ""
	
	for line := range file_buffer {
	
		for map_key, map_value  := range vm_map {
			regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
			file_buffer[line] = regex_obj.ReplaceAllString(file_buffer[line], map_value)
		}
		output_string += file_buffer[line] + "\n"
	}
	
	bytestream := []byte(output_string)
	
	dir_name := "../vm_bootstrap/VMs/" + vm_map["vm_name"] + "/"
	linux_command_line.Execute_command_line("mkdir -p " + dir_name)
	
	err := ioutil.WriteFile(dir_name + "Vagrantfile", bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_aws(file_buffer []string, vm_map map[string]string, output_path string) {
	output_string := ""
	
	for line := range file_buffer {
	
		for map_key, map_value  := range vm_map {
			regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
			file_buffer[line] = regex_obj.ReplaceAllString(file_buffer[line], map_value)
		}
		output_string += file_buffer[line] + "\n"
	}
	
	bytestream := []byte(output_string)
	
	err := ioutil.WriteFile(output_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_gcp(file_buffer []string, vm_map map[string]string, output_path string) {
	output_string := ""
	
	for line := range file_buffer {
	
		for map_key, map_value  := range vm_map {
			regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
			file_buffer[line] = regex_obj.ReplaceAllString(file_buffer[line], map_value)
		}
		output_string += file_buffer[line] + "\n"
	}
	
	bytestream := []byte(output_string)
	
	err := ioutil.WriteFile(output_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_iptables(file_buffer []string, vm_map map[string]string, output_path string) {
	output_string := ""
	
	for line := range file_buffer {
	
		for map_key, map_value  := range vm_map {
			regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
			file_buffer[line] = regex_obj.ReplaceAllString(file_buffer[line], map_value)
		}
		output_string += file_buffer[line] + "\n"
	}
	
	bytestream := []byte(output_string)
	
	err := ioutil.WriteFile(output_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_docker(file_buffer []string, docker_submap map[string]string, output_path string) {
	output_string := ""
	
	for line := range file_buffer {
	
		for map_key, map_value  := range docker_submap {
			regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
			file_buffer[line] = regex_obj.ReplaceAllString(file_buffer[line], map_value)
		}
		output_string += file_buffer[line] + "\n"
	}
	
	bytestream := []byte(output_string)
	
	err := ioutil.WriteFile(output_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_docker_port_section(docker_run_dest_path string, docker_submap map[string]string) {
  
  port_section_string := ""
  file_buffer, err := ioutil.ReadFile(docker_run_dest_path)
  file_string := string(file_buffer)
  
  port_section_string += "-p "
  port_section_string += docker_submap["external_port_0"]
  port_section_string += ":"
  port_section_string += "22"
  
  for map_key, map_value := range docker_submap {
  
    if strings.HasPrefix(map_key, "external_port_") && map_key != "external_port_0" {
    
      port_number := map_key[14:len(map_key)]
      port_section_string += " -p "
      port_section_string += map_value
      port_section_string += ":"
      port_section_string += docker_submap["internal_port_" + port_number]
    }
  }
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<port_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, port_section_string) 
  
  bytestream := []byte(file_string)
	
	err = ioutil.WriteFile(docker_run_dest_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Populate_docker_registry_section(docker_run_dest_path string, docker_submap map[string]string, nexus_admin_password string) {
  
  /// Fetch Nexus Registry Map ///
  nexus_registry_map := make(map[string]string)
  
  if _, err := os.Stat("./databases/nexus_registry.json"); err == nil {
    nexus_registry_map = kv_store.Create_single_level_map_from_json_file("./databases/nexus_registry.json")
  }
  ////////////////////////////////
  
  nexus_registry := nexus_registry_map["vm_name"] + "." + nexus_registry_map["fqdn"]
  image_name := docker_submap["image_name"]
  
  docker_login_string := "--username admin --password " + nexus_admin_password
  docker_registry_string := nexus_registry + ":" + "18444/" + image_name
  
  file_buffer, err := ioutil.ReadFile(docker_run_dest_path)
  file_string := string(file_buffer)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<docker_login_string>}"))
  file_string = regex_obj.ReplaceAllString(file_string, docker_login_string) 
  
  regex_obj = regexp.MustCompile(regexp.QuoteMeta("{<nexus_registry_section>}"))
  file_string = regex_obj.ReplaceAllString(file_string, docker_registry_string)
  
  bytestream := []byte(file_string)
	
	err = ioutil.WriteFile(docker_run_dest_path, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Get_template_populated_string(file_content_string string, input_map map[string]string) string {
  	
	for map_key, map_value  := range input_map {
	
	  regex_obj := regexp.MustCompile(regexp.QuoteMeta("{<" + map_key + ">}"))
		file_content_string = regex_obj.ReplaceAllString(file_content_string, map_value)
	}
	
	return file_content_string
}
/////////////////////////////////////////////////////////////////////////////////
func Get_regex_populated_string(file_content_string string, regex_string string, replace_string string) string {
	
	regex_obj := regexp.MustCompile(regexp.QuoteMeta(regex_string))
	file_content_string = regex_obj.ReplaceAllString(file_content_string, replace_string)
  
	return file_content_string
}
/////////////////////////////////////////////////////////////////////////////////

