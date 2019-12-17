// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package ansible_handler

import (
	"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	"html/template"
	//"os/exec"
	"../kv_store"
	//"./map_template"
	"../logging"
  "strconv"
  "os"
	"../general_utility_web_handler"
  "../linux_command_line"
  "../connectivity_check"
  //"../security"
)
/*
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
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
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DockerPage struct {
	Title string
	Body []byte
	VM_name string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
*/
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generic_VM_web_handler(w http.ResponseWriter, r *http.Request) {
	
	title := ""
  var show_ssh = "false"
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""
  docker_container_display := ""
  docker_container_name := ""

  webform_map := kv_store.Create_from_Webform(r)
  docker_map := make(map[string]map[string]string)

  top_level_map := make(map[string]map[string]string)
  master_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]    
    platform_type = webform_map["platform_type"]
    instance_construction_type = webform_map["instance_construction_type"]
  }
	
	json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }		

  //###########################      
  _, route_button_sshlaunch := webform_map["button_sshlaunch"]
  _, route_generate_inventory := webform_map["generate_inventory"]
  //###########################
  
  docker_container_name = webform_map["docker_container_name"]

  if docker_container_name != "" {  	
  	if _, err := os.Stat("./databases/docker_store.json"); err == nil {
	    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")	    
	  }
  }

  //#############################################  
  if route_generate_inventory {
  	fmt.Println("Generating inventory only")  	

  	master_map = general_utility_web_handler.Get_concatencated_vm_maps()
  	kv_store.Export_to_json(master_map, "./databases/master_map.json")

  	r.ParseForm()  	
  	process_ansible_inventory(vm_name, top_level_map, webform_map, docker_container_name, webform_map["server_selection"], r.Form["selection_list_of_clients"])
  }    
  //#############################################

  if route_button_sshlaunch {
    
    if vm_name != "" {
      
      show_ssh = "true"

      netstat_output, _ := linux_command_line.Execute_command_line("./linux_command_line/netstat_command.sh")     

      log_path := "../Logs/ansible/" + vm_name + "_gotty_log.txt"

      gotty_port := 8081
      for (logging.Check_if_regex_string_exists_in_string(netstat_output, strconv.Itoa(gotty_port))) == 1 {

        gotty_port += 1;
      }

      user_id := top_level_map[vm_name]["user_id"]

      if docker_container_name != "" {
      	ssh_port := docker_map[docker_container_name]["external_port_0"]
      	linux_command_line.Run_gotty(user_id, vm_name, top_level_map[vm_name]["public_ip_address"], ssh_port, strconv.Itoa(gotty_port), log_path)
      } else {
      	linux_command_line.Run_gotty(user_id, vm_name, top_level_map[vm_name]["public_ip_address"], "22", strconv.Itoa(gotty_port), log_path)      
      }
      
      logged_gotty_ip_address := logging.Get_string_from_log(log_path, "(http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+):\\d+")
      logged_gotty_port := logging.Get_string_from_log(log_path, "http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+:(\\d+)")
      
      top_level_map[vm_name]["gotty_ip_address"] = logged_gotty_ip_address
      top_level_map[vm_name]["gotty_port"] = logged_gotty_port
            
      kv_store.Export_to_json(top_level_map, "./databases/" + json_store_name)
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
   
  }
  
  //###########################
  
  parser := &general_utility_web_handler.Page{Title: title}

  if docker_container_name != "" {

  	docker_container := webform_map["docker_container_name"]
  	docker_container_display = " => Docker Container => " + docker_container    

    parser = &general_utility_web_handler.Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: platform_type, Instance_construction_type: instance_construction_type, Docker_container_name: docker_container_name, Docker: docker_container_display}
  } else {
  	parser = &general_utility_web_handler.Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: platform_type, Instance_construction_type: instance_construction_type, Docker_container_name: docker_container_name, Docker: docker_container_display}
  }

  t, _ := template.ParseFiles("templates/html/vm_details.html")
  t.Execute(w, parser)

	//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func process_ansible_inventory(vm_name string, top_level_map map[string]map[string]string, webform_map map[string]string, docker_container_name string,
																server_selection string, client_list_array []string) string {
		
  _, route_button_shutdown := webform_map["button_shutdown"]
  _, route_button_reboot := webform_map["button_reboot.x"]

  //#######################################################################
  playbook_name := ""
  deployment_list_selection := webform_map["deployment_list_selection"]  

  if deployment_list_selection != "" {
  	playbook_name = string(deployment_list_selection[16:])
  }

  working_directory := "../ansible/" + playbook_name
  //#######################################################################

  if route_button_shutdown {
  
    Generate_inventory(vm_name, top_level_map, "../ansible/vm_delete", docker_container_name, server_selection, client_list_array)
    working_directory = "../ansible/vm_delete/"
    
  } else if route_button_reboot {
        
    Generate_inventory(vm_name, top_level_map, "../ansible/vm_reboot", docker_container_name, server_selection, client_list_array)
    working_directory = "../ansible/vm_reboot/"
  }

	if deployment_list_selection == "deployment_list_basic_config" || deployment_list_selection == "deployment_list_deploy_iptables" {

		ssh_private_key_path := "../ssh_keys/" + vm_name + "/private_key"
		connectivity_check.Run_remote_exec(top_level_map[vm_name]["user_id"], top_level_map[vm_name]["public_ip_address"], ssh_private_key_path)
	}

	if deployment_list_selection == "deployment_list_nexus" {
		kv_store.Create_Nexus_registry_db(vm_name, top_level_map[vm_name]["private_ip_address"], top_level_map[vm_name]["fqdn"])
	}

	//#######################################################################
	Generate_inventory(vm_name, top_level_map, working_directory, docker_container_name, server_selection, client_list_array)
	//#######################################################################
	
  return working_directory
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Ansible_runner(w http.ResponseWriter, r *http.Request) {
  var title = "Ansible Run"
  
  docker_container_name := ""

  top_level_map := make(map[string]map[string]string)
  master_map := make(map[string]map[string]string)
  webform_map := kv_store.Create_from_Webform(r) 
  
  vm_name := webform_map["vm_name"]  
  docker_container_name = webform_map["docker_container_name"]

  var platform_type = webform_map["platform_type"]
  var instance_construction_type = webform_map["instance_construction_type"]
  
  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)  

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }
  
  //////////////////////////////////////////////////////////////////////////////////    
  
  _, route_button_shutdown := webform_map["button_shutdown"]
  _, route_button_reboot := webform_map["button_reboot.x"]     
  _, route_deployment_VM := webform_map["ansible_deploy"]  

  //#############################################    
  
  master_map = general_utility_web_handler.Get_concatencated_vm_maps()
  kv_store.Export_to_json(master_map, "./databases/master_map.json")

  r.ParseForm()  

  working_directory := process_ansible_inventory(vm_name, top_level_map, webform_map, docker_container_name, webform_map["server_selection"], r.Form["selection_list_of_clients"])

  //#############################################  	
  log_path := "../Logs/ansible/ansible_" + vm_name + ".log"
  //#############################################  	

  if docker_container_name != "" {

  	log_path = "../Logs/ansible/ansible_" + vm_name + "_" + webform_map["docker_container_name"] + ".log"
    Ansible_generic_docker_runner(vm_name, webform_map["docker_container_name"], webform_map["ip_address"], working_directory, log_path)

  } else if route_deployment_VM {
  
    Ansible_generic_runner(vm_name, top_level_map[vm_name]["public_ip_address"], working_directory, log_path)
                
  } else if route_button_shutdown {
    
    working_directory = "../ansible/vm_delete/"
    
    kv_store.Set_kv_entry(top_level_map[vm_name], "latest_ping", "Fail")
    kv_store.Set_kv_entry(top_level_map[vm_name], "latest_ssh", "Fail")
    kv_store.Export_to_json(top_level_map, "./databases/" + json_store_name)
    
    Ansible_generic_runner(vm_name, webform_map["ip_address"], working_directory, log_path)     
    
  } else if route_button_reboot {
    
    working_directory = "../ansible/vm_reboot/"
    
    kv_store.Set_kv_entry(top_level_map[vm_name], "latest_ping", "Fail")
    kv_store.Set_kv_entry(top_level_map[vm_name], "latest_ssh", "Fail")
    kv_store.Export_to_json(top_level_map, "./databases/" + json_store_name)
    
    Ansible_generic_runner(vm_name, webform_map["ip_address"], working_directory, log_path)     
  }  
   
	body, err := logging.Get_log(log_path)
	general_utility_web_handler.Check(err)

	p := &general_utility_web_handler.Page{Title: title, Body: (body), LogPath: log_path, LogPosition: 0}
  
  t, _ := template.ParseFiles("templates/html/ansible_log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Ansible_generic_runner(vm_name string, ip_address string, working_directory string, log_path string) {
  
  connectivity_check.Clean_ssh(ip_address, vm_name)    
  
  command_string := "ansible-playbook -i inventory_" + vm_name + ".yml main_deploy.yml"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Ansible_generic_docker_runner(vm_name string, container_name string, ip_address string, working_directory string, log_path string) {
  
  connectivity_check.Clean_ssh(ip_address, vm_name)

  command_string := "ansible-playbook -i inventory_" + vm_name + "_" + container_name + ".yml main_deploy.yml"
  
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)
}
/*
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Generic_security_log_web_handler(w http.ResponseWriter, r *http.Request) {
  
  var title = ""
  log_path := ""
  
  webform_map := kv_store.Create_from_Webform(r)
  
  if webform_map["vm_name"] == "" {
    log_path = "../Logs/ansible/security_" + webform_map["vpc_id"] + ".log"
  } else {
    log_path = "../Logs/ansible/security_" + webform_map["vm_name"] + ".log"
  }
  
  title = log_path
  body, err := logging.Get_log(log_path)
	check(err)
    
  p := &Page{Title: title, Body: (body), LogPath: log_path, LogPosition: 0}
  
  t, _ := template.ParseFiles("templates/html/ansible_log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
*/
