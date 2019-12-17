// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package general_utility_web_handler

import (
	"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	"html/template"
	"regexp"
	//"os/exec"
	"../kv_store"
	//"./map_template"
	"../logging"
  "strconv"
  "os"
  "../linux_command_line"
  //"../connectivity_check"
  //"../security"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Check(e error) {
    if e != nil {
      fmt.Println(e.Error() + "\n")
      //panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Page struct {
	Title string
	Body []byte
	LogPath string
  LogPosition int
  Show_SSH string
  Selected_VM string
  Cloud_Platform string
  Instance_construction_type string
  Docker_container_name string
  Docker string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DockerPage struct {
	Title string
	Body []byte
	Show_SSH string
  Selected_VM string
  Cloud_Platform string
  Instance_construction_type string
  Docker_container_name string
  Docker string
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Referential_integrity_check(top_level_map map[string]map[string]string, vm_name string) {
  
  iptables_security_map := make(map[string]map[string]map[string]string)
  
  if _, err := os.Stat("./databases/iptables_security_store.json"); err == nil {
    iptables_security_map = kv_store.Create_3_level_map_from_json_file("./databases/iptables_security_store.json")
  }
  
  delete(iptables_security_map, vm_name)   
  
  kv_store.Export_3_level_map_to_json(iptables_security_map, "./databases/iptables_security_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generic_security_log_web_handler(w http.ResponseWriter, r *http.Request) {
  
  var title = ""
  log_path := ""
  
  webform_map := kv_store.Create_from_Webform(r)
  
  log_path = webform_map["log_path"]
  fmt.Println(log_path)
  title = log_path
  body, err := logging.Get_log(log_path)
	Check(err)
    
  p := &Page{Title: title, Body: (body), LogPath: log_path, LogPosition: 0}
  
  t, _ := template.ParseFiles("templates/html/ansible_log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_json_store_name(platform_type string, instance_construction_type string) string {

	  json_store_name := ""

	  if instance_construction_type == "Imported" && platform_type == "aws" {
	  	json_store_name = "aws_imported_store.json"	  	
	  } else if instance_construction_type == "Imported" && platform_type == "gcp" {
	  	json_store_name = "gcp_imported_store.json"
	  } else if instance_construction_type == "Managed" && platform_type == "aws" {
	  	json_store_name = "aws_store.json"
	  } else if instance_construction_type == "Managed" && platform_type == "gcp" {
	  	json_store_name = "gcp_store.json"
	  } else {
	  	json_store_name = "vm_store.json"
	  }

	  return json_store_name
	 }
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_concatencated_vm_maps() map[string]map[string]string {

	master_map := make(map[string]map[string]string)
	top_level_map_array := make([](map[string]map[string]string), 5)

	top_level_map_array[0] = kv_store.Create_top_level_map_from_json_file("./databases/vm_store.json")
	top_level_map_array[1] = kv_store.Create_top_level_map_from_json_file("./databases/aws_store.json")
	top_level_map_array[2] = kv_store.Create_top_level_map_from_json_file("./databases/gcp_store.json")
	top_level_map_array[3] = kv_store.Create_top_level_map_from_json_file("./databases/aws_imported_store.json")
	top_level_map_array[4] = kv_store.Create_top_level_map_from_json_file("./databases/gcp_imported_store.json")
	
	for _, top_level_map := range top_level_map_array {
		for map_key, map_value := range top_level_map {
			master_map[map_key] = map_value
		}
	}
	
	return master_map
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_listening_ports(w http.ResponseWriter, r *http.Request) {
	
	vm_name := ""
  platform_type := ""
  instance_construction_type := ""  
  docker_container_name := ""
  ssh_port := ""

	ajax_map := kv_store.Create_from_Webform(r)
  fmt.Println(ajax_map)

	if len(ajax_map) != 0 {
    vm_name = ajax_map["vm_name"]    
    platform_type = ajax_map["platform_type"]
    instance_construction_type = ajax_map["instance_construction_type"]
    docker_container_name = ajax_map["docker_container_name"]
  }

  docker_map := make(map[string]map[string]string)

  if docker_container_name != "" {  	
  	if _, err := os.Stat("./databases/docker_store.json"); err == nil {
	    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")	    
	  }
  }

	/////////////////////////////////////////////////////////////////
	top_level_map := make(map[string]map[string]string)

	json_store_name := Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }  

  if docker_container_name != "" {
  	ssh_port = docker_map[docker_container_name]["external_port_0"]
  } else {
  	ssh_port = "22"
  }

  command_result_string := linux_command_line.Get_remote_listening_ports(top_level_map[vm_name]["user_id"], vm_name, top_level_map[vm_name]["public_ip_address"], ssh_port)

  /////////////////////////////////////////////////////////////////
  regex_obj := regexp.MustCompile("([a-z|A-Z|0-9]+).*?:(\\d+).*LISTEN\\s*(.*?)\n")
	regex_result_array := regex_obj.FindAllStringSubmatch(command_result_string, -1)
	
	command_result_map := make(map[string]map[string]string)	
	
	for row_index, sub_match_array := range regex_result_array {
		
		command_result_map[strconv.Itoa(row_index)] = make(map[string]string)
		
		for col_index, sub_match := range sub_match_array {
			if col_index != 0 {
				command_result_map[strconv.Itoa(row_index)][strconv.Itoa(col_index)] = sub_match
			}
		}
	}  
  /////////////////////////////////////////////////////////////////	

  command_result_map_json := kv_store.Export_2level_map_to_json_string(command_result_map)
	
	fmt.Fprintf(w, "%s", command_result_map_json)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
