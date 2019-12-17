// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package gcp_handler

import (
	"fmt"
	"io/ioutil"
	//"log"
	"encoding/json"
	"path/filepath"
	"net/http"
	"html/template"
	"os/exec"
	"../kv_store"
	"regexp"
	"strings"
	//"./template_populator"
	"../logging"
  "strconv"
  "os"
  "../linux_command_line"
  "../connectivity_check"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_gcp_topology () {
  
	cmd := exec.Command("python3", "gce.py")

	current_directory, _  := os.Getwd()	
	cmd.Dir, _ = filepath.Abs(current_directory + "/../gcp_bootstrap/")
	
	out, _ := cmd.Output()
	    
	bytestream := []byte(out)
	ioutil.WriteFile("../gcp_bootstrap/gcp_output.txt", bytestream, 0644)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func separate_region_and_az(region_az_string string) (string, string) {

	regex_obj := regexp.MustCompile("(.*)-(.*)")  
  match_array := regex_obj.FindStringSubmatch(region_az_string)

  return match_array[1], match_array[2]
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_imported_gcp_map (top_level_map map[string]map[string]string) {		

	fmt.Println("\nImporting...")

	get_gcp_topology()
    
  file_buffer, _ := ioutil.ReadFile("../gcp_bootstrap/gcp_output.txt")
  
  imported_map_top_level := make(map[string]interface{})
  json.Unmarshal([]byte(file_buffer), &imported_map_top_level)
  
  imported_map := make(map[string]interface{})
  imported_map = imported_map_top_level["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})
  
  project_id := get_gcp_project_id()  

	for _, instance := range imported_map {
		
		vm_name := instance.(map[string]interface{})["gce_name"].(string)
		top_level_map[vm_name] = make(map[string]string)

		top_level_map[vm_name]["vm_name"] = vm_name
		top_level_map[vm_name]["public_ip_address"] = instance.(map[string]interface{})["gce_public_ip"].(string)
		top_level_map[vm_name]["private_ip_address"] = instance.(map[string]interface{})["gce_private_ip"].(string)		

		gce_zone_string := instance.(map[string]interface{})["gce_zone"].(string)
		region, zone := separate_region_and_az(gce_zone_string)
		top_level_map[vm_name]["region"] = region
		top_level_map[vm_name]["AZ"] = zone
		top_level_map[vm_name]["project_id"] = project_id
		top_level_map[vm_name]["instance_construction_type"] = "Imported"
		top_level_map[vm_name]["platform_type"] = "gcp"
	}
  
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_gcp_project_id() string {
	input_stream, _ := ioutil.ReadFile("../gcp_bootstrap/gcp_project_id")	
	return strings.TrimSpace(string(input_stream))
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
      fmt.Println(e.Error() + "\n")
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
  GCP_Project string
  GCP_Json_key string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Launch_instance(w http.ResponseWriter, r *http.Request) {
	var title = "Launch Google Cloud Instance"
  
  top_level_map := make(map[string]map[string]string)
  
  if _, err := os.Stat("./databases/gcp_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/gcp_store.json")
  }
  //////////////////////////////////////////////////////////////////////////////////////////////
  
  vm_map := kv_store.Create_from_Webform(r)
  vm_name := vm_map["vm_name"]
  
  linux_command_line.Execute_command_line("mkdir -p ../Logs")
  linux_command_line.Execute_command_line("mkdir -p ../ansible/ssh_keys/" + vm_name)
  
  Generate_terraform_resources(top_level_map[vm_name])
  // Populate_gcp_firewall_default_section(vm_name)
  
	working_directory := "../gcp_bootstrap/instances/" + vm_name
	//check(err)

	command_string := "bash ./terraform_script.sh"
	log_path := "../Logs/" + vm_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  logpath := "../Logs/" + vm_name + ".log"
  body, err := logging.Get_log(logpath)	
  check(err)	
  p := &Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}

  t, _ := template.ParseFiles("templates/html/log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Terminate_instance(w http.ResponseWriter, r *http.Request) {
	var title = "Terminate Instance"
  
  vm_map := kv_store.Create_from_Webform(r)
  vm_name := vm_map["vm_name"]
  
  linux_command_line.Execute_command_line("mkdir -p ../Logs")
  
	working_directory := "../gcp_bootstrap/instances/" + vm_name
	//check(err)
  
	command_string := "terraform destroy -auto-approve"
	log_path := "../Logs/" + vm_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  logpath := "../Logs/" + vm_name + ".log"
  body, err := logging.Get_log(logpath)	
  check(err)	
  p := &Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}
  //fmt.Println(body)
  t, _ := template.ParseFiles("templates/html/log.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GCP_web_handler(w http.ResponseWriter, r *http.Request) {

  os_image_map := make(map[string]map[string]string)
  os_image_map["centos"] = make(map[string]string)
  os_image_map["ubuntu"] = make(map[string]string)
  
  os_image_map["centos"]["europe-west4"] = "centos-7"
  os_image_map["ubuntu"]["europe-west4"] = "ubuntu-1604-lts"
  
	var title = "Google Cloud (Managed VMs)"
  var show_ssh = "false"
  vm_name := ""
  cloud_platform := "gcp"
  instance_construction_type := "Managed"
  
  webform_map := kv_store.Create_from_Webform(r)
  
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]
  }
  
  if _, err := os.Stat("./databases/gcp_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/gcp_store.json")
  }
  
  //###########################  
  _, route_button_addInstance := webform_map["Add Instance"]
  _, route_button_ssh := webform_map["button_ssh"]
  _, route_button_sshlaunch := webform_map["button_sshlaunch"]
  _, route_button_delete := webform_map["button_delete"]  
  //###########################

  if route_button_delete {

    //delete vm from map	
    delete(top_level_map, vm_name)
    kv_store.Export_to_json(top_level_map, "./databases/gcp_store.json")
    
  } else if route_button_addInstance {
  
    // webform_map := kv_store.Create_from_Webform(r)
    if webform_map["OS"] == "centos/7" {
    
      kv_store.Set_kv_entry(webform_map, "os_image", os_image_map["centos"][webform_map["region"]])
    } else {
    
      kv_store.Set_kv_entry(webform_map, "os_image", os_image_map["ubuntu"][webform_map["region"]])
    }
    
    top_level_map[vm_name] = webform_map
    kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", "gcp")
    
    if len(top_level_map) == 1 {
      kv_store.Set_kv_entry(top_level_map[vm_name], "vm_type", "MAIN")
    }
    
    Generate_gcp_ssh_keys(top_level_map[vm_name])
    
    project_id := get_gcp_project_id()    
    kv_store.Set_kv_entry(top_level_map[vm_name], "project_id", project_id)

    kv_store.Set_kv_entry(top_level_map[vm_name], "instance_construction_type", instance_construction_type)

    kv_store.Export_to_json(top_level_map, "./databases/gcp_store.json")
  
  } else if route_button_ssh {
    // Check connectivity
    //webform_map := kv_store.Create_from_Webform(r)
    
    connectivity_check.Store_gcp_vm_resources_in_DB(top_level_map, vm_name)
    connectivity_check.Copy_gcp_ssh_keys(vm_name)
    
    if vm_name != "" {
      connectivity_check.Check_connectivity(top_level_map, vm_name, "gcp")
    }
    
  } else if route_button_sshlaunch {
    // Launch VirtualBox VM
    //webform_map := kv_store.Create_from_Webform(r)
    
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

      user_id := top_level_map[vm_name]["user_id"]

      linux_command_line.Run_gotty(user_id, vm_name, webform_map["public_ip_address"], "22", strconv.Itoa(gotty_port), log_path)
      
      logged_gotty_ip_address := logging.Get_string_from_log(log_path, "(http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+):\\d+")
      logged_gotty_port := logging.Get_string_from_log(log_path, "http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+:(\\d+)")
      
      top_level_map[vm_name]["gotty_ip_address"] = logged_gotty_ip_address
      top_level_map[vm_name]["gotty_port"] = logged_gotty_port
      
      kv_store.Export_to_json(top_level_map, "./databases/gcp_store.json")
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
  }
  
  //###########################
  
  p := &Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: cloud_platform, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/gcp_config.html")
  t.Execute(w, p)

//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GCP_imported_web_handler(w http.ResponseWriter, r *http.Request) {    

  os_image_map := make(map[string]map[string]string)
  os_image_map["centos"] = make(map[string]string)
  os_image_map["ubuntu"] = make(map[string]string)
  
  os_image_map["centos"]["europe-west4"] = "centos-7"
  os_image_map["ubuntu"]["europe-west4"] = "ubuntu-1604-lts"
  
	var title = "Google Cloud (Imported VMs)"
  var show_ssh = "false"
  vm_name := ""
  cloud_platform := "gcp"
  instance_construction_type := "Imported"
  
  webform_map := kv_store.Create_from_Webform(r)
  
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]
  }
  
  if _, err := os.Stat("./databases/gcp_imported_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/gcp_imported_store.json")
  }
  
  //###########################  
  _, route_button_importInstances := webform_map["Import Instances"]
  _, route_button_set_instance_details := webform_map["SetInstanceDetails"]
  _, route_button_addInstance := webform_map["Add Instance"]
  _, route_button_ssh := webform_map["button_ssh"]
  _, route_button_sshlaunch := webform_map["button_sshlaunch"]
  _, route_button_delete := webform_map["button_delete"]  
  //###########################

  if route_button_importInstances {
  	
  	get_imported_gcp_map(top_level_map)
  	kv_store.Export_to_json(top_level_map, "./databases/gcp_imported_store.json")
  
  } else if route_button_set_instance_details {
  	
  	kv_store.Set_kv_entry(top_level_map[vm_name], "ssh_private_key", strip_cr_from_string(webform_map["ssh_private_key"]))  	

  	kv_store.Set_kv_entry(top_level_map[vm_name], "OS", webform_map["OS"])  	
  	kv_store.Set_kv_entry(top_level_map[vm_name], "user_id", webform_map["user_id"])

  	target_directory := "../gcp_bootstrap/instances/" + vm_name + "/"
		linux_command_line.Execute_command_line("mkdir -p " + target_directory)
		linux_command_line.Execute_command_line("mkdir -p ../ansible/ssh_keys/" + vm_name)

		Store_gcp_ssh_keys_from_imported_VM(top_level_map[vm_name], vm_name)
  	connectivity_check.Copy_gcp_ssh_keys(vm_name)

  	kv_store.Export_to_json(top_level_map, "./databases/gcp_imported_store.json")
    
  } else if route_button_delete {

    //delete vm from map    
    delete(top_level_map, vm_name)
    kv_store.Export_to_json(top_level_map, "./databases/gcp_imported_store.json")
    
  } else if route_button_addInstance {
  
    // webform_map := kv_store.Create_from_Webform(r)
    if webform_map["OS"] == "centos/7" {
    
      kv_store.Set_kv_entry(webform_map, "os_image", os_image_map["centos"][webform_map["region"]])
    } else {
    
      kv_store.Set_kv_entry(webform_map, "os_image", os_image_map["ubuntu"][webform_map["region"]])
    }
    
    top_level_map[vm_name] = webform_map
    kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", "gcp")
    
    if len(top_level_map) == 1 {
      kv_store.Set_kv_entry(top_level_map[vm_name], "vm_type", "MAIN")
    }
    
    Generate_gcp_ssh_keys(top_level_map[vm_name])    
    
    kv_store.Export_to_json(top_level_map, "./databases/gcp_imported_store.json")
  
  } else if route_button_ssh {
    // Check connectivity
    //webform_map := kv_store.Create_from_Webform(r)
    
    //connectivity_check.Store_gcp_vm_resources_in_DB(top_level_map, vm_name)
    connectivity_check.Copy_gcp_ssh_keys(vm_name)
    
    if vm_name != "" {
      connectivity_check.Check_connectivity(top_level_map, vm_name, "gcp")
    }
    
  } else if route_button_sshlaunch {
    // Launch VirtualBox VM
    //webform_map := kv_store.Create_from_Webform(r)
    
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

      user_id := top_level_map[vm_name]["user_id"]

      linux_command_line.Run_gotty(user_id, vm_name, webform_map["public_ip_address"], "22", strconv.Itoa(gotty_port), log_path)
      
      logged_gotty_ip_address := logging.Get_string_from_log(log_path, "(http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+):\\d+")
      logged_gotty_port := logging.Get_string_from_log(log_path, "http:.*?\\d+\\.\\d+\\.\\d+\\.\\d+:(\\d+)")
      
      top_level_map[vm_name]["gotty_ip_address"] = logged_gotty_ip_address
      top_level_map[vm_name]["gotty_port"] = logged_gotty_port
      
      kv_store.Export_to_json(top_level_map, "./databases/gcp_imported_store.json")
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
  }
  
  //###########################
  
  p := &Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: cloud_platform, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/gcp_config_imported.html")
  t.Execute(w, p)

//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_gcp_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/gcp_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_gcp_imported_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/gcp_imported_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
