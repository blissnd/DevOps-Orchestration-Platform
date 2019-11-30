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
  //"strconv"
  "os"
	"../general_utility_web_handler"
  "../linux_command_line"
  "../connectivity_check"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Docker_web_handler(w http.ResponseWriter, r *http.Request) {

	var title = "Docker"
	var show_ssh = "false"
  docker_container_name := ""
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""
  webform_map := kv_store.Create_from_Webform(r)
  working_directory := ""

  docker_map := make(map[string]map[string]string)
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
  	vm_name = webform_map["vm_name"]
    docker_container_name = webform_map["docker_container_name"]
    platform_type = webform_map["platform_type"]
    instance_construction_type = webform_map["instance_construction_type"]
  }

  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }
  
  if _, err := os.Stat("./databases/docker_store.json"); err == nil {
    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  }
  
  //###########################
  _, route_button_addDocker := webform_map["button_add_docker"]
  _, route_button_delete := webform_map["button_delete"]  
  //###########################
  
  if route_button_addDocker {
  
    docker_map[docker_container_name] = webform_map
    docker_map[docker_container_name]["ip_address"] = top_level_map[vm_name]["public_ip_address"]
    docker_map[docker_container_name]["user_id"] = top_level_map[vm_name]["user_id"]
    kv_store.Export_to_json(docker_map, "./databases/docker_store.json")
    
  } else if route_button_delete {
    //delete docker container from VM
    working_directory = "../ansible/docker_container_delete/"
    Ansible_generic_docker_runner(vm_name, docker_container_name, top_level_map[vm_name]["ip_address"], working_directory, "../Logs/ansible/ansible_" + vm_name + "_" + docker_container_name + ".log")     
    
    //delete vm from map
    delete(docker_map, webform_map["docker_container_name"])
    kv_store.Export_to_json(docker_map, "./databases/docker_store.json")      
   
  }
  
  //###########################
  
  p := &general_utility_web_handler.DockerPage{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: platform_type, Instance_construction_type: instance_construction_type}

  t, _ := template.ParseFiles("templates/html/docker_config.html")
  t.Execute(w, p)

	//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Ansible_docker_runner(w http.ResponseWriter, r *http.Request) {  // POSSIBLY DEPRACATED!
  var title = "Ansible Run"
  
  webform_map := kv_store.Create_from_Webform(r)
  fmt.Println(webform_map)
  
  vm_name := webform_map["vm_name"]
  
  working_directory := "../ansible/basic_config/"

  connectivity_check.Clean_ssh(webform_map["ip_address"], vm_name)

  command_string := "ansible-playbook -i inventory_" + vm_name + ".yml main_deploy.yml"
	log_path := "../../Logs/ansible/ansible_" + vm_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  log_path = "../Logs/ansible/ansible_" + vm_name + ".log"
	body, err := logging.Get_log(log_path)
	general_utility_web_handler.Check(err)
	
	p := &general_utility_web_handler.Page{Title: title, Body: (body), LogPath: log_path, LogPosition: 0}
  
  t, _ := template.ParseFiles("templates/html/ansible_log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_docker_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Docker_launch_handler(w http.ResponseWriter, r *http.Request) {

  var title = "Launch Docker Container"
  docker_container_name := ""
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""

  webform_map := kv_store.Create_from_Webform(r)  
  docker_map := make(map[string]map[string]string)
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
  	vm_name = webform_map["vm_name"]
    docker_container_name = webform_map["docker_container_name"]
    platform_type = webform_map["platform_type"]
    instance_construction_type = webform_map["instance_construction_type"]
  }

  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }    
  
  if _, err := os.Stat("./databases/docker_store.json"); err == nil {
    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  }
  
  Generate_docker_resources(docker_map[docker_container_name], vm_name, docker_container_name)

  Generate_inventory_for_docker_container_launch(top_level_map[vm_name], top_level_map[vm_name]["public_ip_address"], docker_map[docker_container_name], "../ansible/docker_config")
  Generate_inventory_for_docker_container_launch(top_level_map[vm_name], top_level_map[vm_name]["public_ip_address"], docker_map[docker_container_name], "../ansible/docker_container_delete")    
  
  working_directory := "../ansible/docker_config/"
  connectivity_check.Clean_ssh(top_level_map[vm_name]["ip_address"], vm_name)
  command_string := "ansible-playbook -i inventory_" + vm_name + "_" + docker_container_name + ".yml main_deploy.yml"
	log_path := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  //###########################
  
  logpath := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	body, err := logging.Get_log(logpath)
	
	general_utility_web_handler.Check(err)
	
	p := &general_utility_web_handler.Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}
	//fmt.Println(body)
	t, _ := template.ParseFiles("templates/html/log.html")
	t.Execute(w, p)
  
  //kv_store.Export_to_json(docker_map, "./databases/docker_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Docker_image_launch_handler(w http.ResponseWriter, r *http.Request) {

  var title = "Launch Docker Container from Nexus Image Registry"
  docker_container_name := ""
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""

  webform_map := kv_store.Create_from_Webform(r)
  
  if len(webform_map) != 0 {
  	vm_name = webform_map["vm_name"]
    docker_container_name = webform_map["docker_container_name"]
    platform_type = webform_map["platform_type"]
    instance_construction_type = webform_map["instance_construction_type"]
  }

  docker_map := make(map[string]map[string]string)
  top_level_map := make(map[string]map[string]string)
  
  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }
  
  if _, err := os.Stat("./databases/docker_store.json"); err == nil {
    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  }
  
  Generate_docker_resources_for_image(docker_map[docker_container_name], vm_name, docker_container_name)

  Generate_inventory_for_docker_container_launch(top_level_map[vm_name], top_level_map[vm_name]["public_ip_address"], docker_map[docker_container_name], "../ansible/docker_launch_from_registry_image")
  Generate_inventory_for_docker_container_launch(top_level_map[vm_name], top_level_map[vm_name]["public_ip_address"], docker_map[docker_container_name], "../ansible/docker_container_delete")    
  
  working_directory := "../ansible/docker_launch_from_registry_image/"
  connectivity_check.Clean_ssh(top_level_map[vm_name]["ip_address"], vm_name)
  command_string := "ansible-playbook -i inventory_" + vm_name + "_" + docker_container_name + ".yml main_deploy.yml"
	log_path := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  //###########################
  
  logpath := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	body, err := logging.Get_log(logpath)
	
	general_utility_web_handler.Check(err)
	
	p := &general_utility_web_handler.Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}
	//fmt.Println(body)
	t, _ := template.ParseFiles("templates/html/log.html")
	t.Execute(w, p)
  
  //kv_store.Export_to_json(docker_map, "./databases/docker_store.json")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Docker_registry_commit_handler(w http.ResponseWriter, r *http.Request) {

  var title = "Commit to Docker Registry"
  docker_container_name := ""
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""
  
  webform_map := kv_store.Create_from_Webform(r)
  
  if len(webform_map) != 0 {
  	vm_name = webform_map["vm_name"]
    docker_container_name = webform_map["docker_container_name"]
    platform_type = webform_map["platform_type"]
    instance_construction_type = webform_map["instance_construction_type"]
  }

  docker_map := make(map[string]map[string]string)
  top_level_map := make(map[string]map[string]string)
  
  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/" + json_store_name); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }  
  
  if _, err := os.Stat("./databases/docker_store.json"); err == nil {
    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")
  }       

  Generate_docker_resources_for_image(docker_map[docker_container_name], vm_name, docker_container_name)
  
  Generate_inventory_for_docker_image_commit(top_level_map[vm_name], top_level_map[vm_name]["public_ip_address"], docker_map[docker_container_name], "../ansible/docker_commit")

  working_directory := "../ansible/docker_commit/"
  connectivity_check.Clean_ssh(top_level_map[vm_name]["ip_address"], vm_name)
  command_string := "ansible-playbook -i inventory_" + vm_name + "_" + docker_container_name + ".yml main_deploy.yml"
	log_path := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  //###########################
  
  logpath := "../Logs/docker/" + vm_name + "_" + docker_container_name + ".log"
	body, err := logging.Get_log(logpath)
	
	general_utility_web_handler.Check(err)
	
	p := &general_utility_web_handler.Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}
	//fmt.Println(body)
	t, _ := template.ParseFiles("templates/html/log.html")
	t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
