package kubernetes_handler

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
  //"../connectivity_check"
  //"../security"
)
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
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Kubernetes_web_handler(w http.ResponseWriter, r *http.Request) {
	
	title := "Kubernetes Cluster"
  var show_ssh = "false"
  vm_name := ""
  platform_type := ""
  instance_construction_type := ""
  docker_container_display := ""
  docker_container_name := ""

  webform_map := kv_store.Create_from_Webform(r)
  docker_map := make(map[string]map[string]string)

  top_level_map := make(map[string]map[string]string)
  //master_map := make(map[string]map[string]string)
  
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
  //###########################
  
  docker_container_name = webform_map["docker_container_name"]

  if docker_container_name != "" {  	
  	if _, err := os.Stat("./databases/docker_store.json"); err == nil {
	    docker_map = kv_store.Create_top_level_map_from_json_file("./databases/docker_store.json")	    
	  }
  }  

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

  t, _ := template.ParseFiles("templates/html/kubernetes.html")
  t.Execute(w, parser)

	//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_master_map_ajax(w http.ResponseWriter, r *http.Request) {

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

  export_json_string := kv_store.Export_2level_map_to_json_string(master_map)
  fmt.Fprintf(w, "%s", export_json_string)
  //fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_kubernetes_dashboard_token_ajax(w http.ResponseWriter, r *http.Request) {
	
	vm_name := ""  
  ssh_port := ""
  master_map := make(map[string]map[string]string)

	ajax_map := kv_store.Create_from_Webform(r)
  fmt.Println(ajax_map)

	if len(ajax_map) != 0 {
    vm_name = ajax_map["vm_name"]
  } 
  
  master_map = general_utility_web_handler.Get_concatencated_vm_maps();
  
  ssh_port = "22"
  
  command_result_string := linux_command_line.Get_kubernetes_dashboard_token(master_map[vm_name]["user_id"], vm_name, master_map[vm_name]["public_ip_address"], ssh_port)
  
  /////////////////////////////////////////////////////////////////
	
	command_result_map := make(map[string]string)
	command_result_map["kubernetes_dashboard_token"] = command_result_string
	
  /////////////////////////////////////////////////////////////////

  command_result_map_json := kv_store.Export_1level_map_to_json_string(command_result_map)
	
	fmt.Fprintf(w, "%s", command_result_map_json)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
