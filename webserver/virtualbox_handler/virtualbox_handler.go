// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package virtualbox_handler

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
  Cloud_Platform string
  Instance_construction_type string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DockerPage struct {
	Title string
	Body []byte
	VM_name string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Virtualbox_web_handler(w http.ResponseWriter, r *http.Request) {

	var title = "VirtualBox"
  var show_ssh = "false"
  vm_name := ""
  cloud_platform := "virtualbox"
  instance_construction_type := "Managed"  
  
  webform_map := kv_store.Create_from_Webform(r)
  
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]    
  }
  
  if _, err := os.Stat("./databases/vm_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/vm_store.json")
  }
  
  //###########################  
  _, route_button_delete := webform_map["button_delete"]
  _, route_button_addVM := webform_map["Add VM"]
  _, route_button_ssh := webform_map["button_ssh"]
  _, route_button_sshlaunch := webform_map["button_sshlaunch"]
  _, route_button_docker := webform_map["button_docker"]
  //###########################
  
  if route_button_delete {           
    
    //delete vm from map
    general_utility_web_handler.Referential_integrity_check(top_level_map, vm_name)
    delete(top_level_map, vm_name)
    kv_store.Export_to_json(top_level_map, "./databases/vm_store.json")  
    
  } else if route_button_addVM {
  
    //webform_map := kv_store.Create_from_Webform(r)
    if top_level_map[vm_name]["OS"] == "centos/7" {
      kv_store.Set_kv_entry(webform_map, "bootstrap_script", "../../bootstrap_redhat.sh")
    } else {
      kv_store.Set_kv_entry(webform_map, "bootstrap_script", "../../bootstrap.sh")
    }
    
    top_level_map[vm_name] = webform_map
    top_level_map[vm_name]["public_ip_address"] = top_level_map[vm_name]["ip_address"]
    top_level_map[vm_name]["private_ip_address"] = top_level_map[vm_name]["ip_address"]
    kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", cloud_platform)
		kv_store.Set_kv_entry(top_level_map[vm_name], "instance_construction_type", instance_construction_type)
		
    kv_store.Set_kv_entry(top_level_map[vm_name], "user_id", "vagrant")

    kv_store.Export_to_json(top_level_map, "./databases/vm_store.json")
    
    // security.Add_to_virtualbox_security_map(top_level_map, vm_name)
    
  } else if route_button_ssh {
    // Check connectivity
    //webform_map := kv_store.Create_from_Webform(r)
    
    if vm_name != "" {
      connectivity_check.Check_connectivity(top_level_map, vm_name, "virtualbox")
      kv_store.Export_to_json(top_level_map, "./databases/vm_store.json")
    }
    
  } else if route_button_sshlaunch {
    // Launch VirtualBox VM
    //webform_map := kv_store.Create_from_Webform(r)
    fmt.Println("Launching virtualbox ssh")
    
    if vm_name != "" {
      //linux_command_line.SSH_to_remote_host("vagrant", vm_name, webform_map["ip_address"], "22")
      show_ssh = "true"

      netstat_output, _ := linux_command_line.Execute_command_line("./linux_command_line/netstat_command.sh")
      fmt.Println(netstat_output)

      log_path := "../Logs/ansible/" + vm_name + "_gotty_log.txt"

      gotty_port := 8081
      for (logging.Check_if_regex_string_exists_in_string(netstat_output, strconv.Itoa(gotty_port))) == 1 {

        gotty_port += 1;
      }

      linux_command_line.Run_gotty("vagrant", vm_name, webform_map["public_ip_address"], "22", strconv.Itoa(gotty_port), log_path)
      
      logged_gotty_ip_address := logging.Get_string_from_log(log_path, "(http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+):\\d+")
      logged_gotty_port := logging.Get_string_from_log(log_path, "http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+:(\\d+)")
      
      top_level_map[vm_name]["gotty_ip_address"] = logged_gotty_ip_address
      top_level_map[vm_name]["gotty_port"] = logged_gotty_port
      
      kv_store.Export_to_json(top_level_map, "./databases/vm_store.json")
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
   
  } else if route_button_docker {
    // Open docker list
    //webform_map := kv_store.Create_from_Webform(r)
    
    if vm_name != "" {
      p := &DockerPage{Title: title, VM_name: vm_name}
      t, _ := template.ParseFiles("templates/html/docker_config.html")
      t.Execute(w, p)
    }
  }
  
  //###########################
  
  p := &Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: cloud_platform, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/vm_config.html")
  t.Execute(w, p)

	//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_vm_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/vm_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
