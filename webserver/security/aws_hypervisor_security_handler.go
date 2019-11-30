package security

import (
	"fmt"
	//"encoding/json"
	//"io/ioutil"
	"regexp"
	//"log"
	"net/http"
	"html/template"
	//"os/exec"
	"../kv_store"
	//"./template_populator"
	"../logging"
  //"strconv"
  "os"
  "../linux_command_line"
  //"../connectivity_check"
  //"../ansible_handler"
  "../aws_handler"
)
/*
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
      fmt.Println(e.Error() + "\n")
      panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Page struct {
	Title string
	Body []byte
	LogPath string
  LogPosition int
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type webform_struct struct {
  vpc_id string
  ip_address string
  source_cidr string
  port string
  log_path string
  firewall_type string
  platform string
}
*/
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func synchronise_aws_firewall_map(top_level_map map[string]map[string]string, aws_firewall_map map[string]map[string]map[string]string) { // DEPRECATED
  
  if len(top_level_map) > 0 {
  
    for _, vm_map := range top_level_map {    
      vpc_id := vm_map["vpc_id"]      
      
      if _, exist := aws_firewall_map[vpc_id]; exist == false {
        aws_firewall_map[vpc_id] = make(map[string]map[string]string)
      }
    }
    
    kv_store.Export_3_level_map_to_json(aws_firewall_map, "./databases/aws_firewall_store.json")
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AWS_hypervisor_security_web_handler(w http.ResponseWriter, r *http.Request) {

	var title = "AWS Security Group Configuration"  
  //top_level_map := make(map[string]map[string]string)
  
  //if _, err := os.Stat("./databases/aws_store.json"); err == nil {
  //  top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/aws_store.json")
  //}
  
  /////////////////////////////////////////////////////////////////////////////////////
  webform_map := kv_store.Create_from_Webform(r)
  /////////////////////////////////////////////////////////////////////////////////////
  aws_firewall_map := make(map[string]map[string]map[string]string)
  
  if _, err := os.Stat("./databases/aws_firewall_store.json"); err == nil {
    aws_firewall_map = kv_store.Create_3_level_map_from_json_file("./databases/aws_firewall_store.json")
  }
  
  //synchronise_aws_firewall_map(top_level_map, aws_firewall_map)
  
  /////////////////////////////////////////////////////////////////////////////////////
  
  firewall_id := webform_map["firewall_id"]
  firewall_rule_map_key := webform_map["protocol"] + "-" + webform_map["port"]
  
  //##########################################
  _, route_button_add_Rule := webform_map["Add Rule"]
  _, route_button_delete := webform_map["button_delete"]  
  //##########################################
  log_path_prefix := "../Logs/ansible/security"
  //##########################################
  
  if route_button_add_Rule && webform_map[firewall_id] != "null" {
    
    if _, exist := aws_firewall_map[firewall_id]; exist == false {    	
    	aws_firewall_map[firewall_id] = make(map[string]map[string]string)
    }

    if _, exist := aws_firewall_map[firewall_id][firewall_rule_map_key]; exist == false {    	
    	aws_firewall_map[firewall_id][firewall_rule_map_key] = make(map[string]string)
    }

    aws_firewall_map[firewall_id][firewall_rule_map_key] = webform_map        
        
    kv_store.Export_3_level_map_to_json(aws_firewall_map, "./databases/aws_firewall_store.json")          
    
  } else if route_button_delete {
    
    firewall_id := webform_map["firewall_id"]
  	firewall_rule_map_key := webform_map["protocol"] + "-" + webform_map["port"]

    delete(aws_firewall_map[firewall_id], firewall_rule_map_key)
    
    kv_store.Export_3_level_map_to_json(aws_firewall_map, "./databases/aws_firewall_store.json")
  }
  
  //###########################
  p := &Page{Title: title, LogPath: log_path_prefix}
  t, _ := template.ParseFiles("templates/html/aws_hypervisor_security.html")
  t.Execute(w, p)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_aws_firewall_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := make(map[string]map[string]map[string]string)

	if _, err := os.Stat("./databases/aws_firewall_store.json"); err == nil {
    top_level_map = kv_store.Create_3_level_map_from_json_file("./databases/aws_firewall_store.json")
  }

	export_json_string := kv_store.Export_3level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_main_aws_vm_name(top_level_map map[string]map[string]string, vpc_id string) string {
  
  vm_name := ""
  
  for _, vm_map := range top_level_map {
    
    if vm_map["vm_type"] == "MAIN" {
      vm_name = vm_map["vm_name"]
    }
  }
  
  return vm_name
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AJAX_AWS_modify_hypervisor_firewall(w http.ResponseWriter, r *http.Request) {
  
  success_or_failure := 0  
  
  //////////////////////////////////////////////////////////////////////////////
  ajax_map := kv_store.Create_from_Webform(r)
  //////////////////////////////////////////////////////////////////////////////
  
  aws_firewall_map := make(map[string]map[string]map[string]string)
    
  firewall_id := ajax_map["firewall_id"]
  firewall_rule_map_key := ajax_map["protocol"] + "-" + ajax_map["port"]

  if _, err := os.Stat("./databases/aws_firewall_store.json"); err == nil {
    aws_firewall_map = kv_store.Create_3_level_map_from_json_file("./databases/aws_firewall_store.json")
  }
  
  //##########################################
  log_path := ajax_map["log_path"]
  //##########################################   
    
  old_state := ajax_map["port_state"]
  
  if old_state == "CLOSED" {
  
    aws_firewall_map[firewall_id][firewall_rule_map_key]["state"] = "OPEN"
    
  } else {
  
    aws_firewall_map[firewall_id][firewall_rule_map_key]["state"] = "CLOSED"
  }        
 
  //aws_handler.Generate_terraform_resources(top_level_map, vm_name)
  //aws_handler.Populate_modified_aws_firewall_section(aws_firewall_map[firewall_id], vm_name)

  aws_handler.Generate_terraform_sec_group_resources(aws_firewall_map[firewall_id][firewall_rule_map_key])
  aws_handler.Populate_separate_aws_firewall_section(aws_firewall_map[firewall_id])

  success_or_failure = Run_aws_sec_group_terraform(log_path)
  
  if success_or_failure == 0 {
    if aws_firewall_map[firewall_id][firewall_rule_map_key]["state"] == "CLOSED" {

      aws_firewall_map[firewall_id][firewall_rule_map_key]["state"] = "OPEN"
      
    } else {
      aws_firewall_map[firewall_id][firewall_rule_map_key]["state"] = "CLOSED"
    }
  }
  
  kv_store.Export_3_level_map_to_json(aws_firewall_map, "./databases/aws_firewall_store.json")
  
  
  ////////////////////////////////////////////////////////////////////////////// 
  
  json_return_string := ""
  
  if success_or_failure == 1 {
    json_return_string = "{\"result\": \"pass\"}"
  } else {
    json_return_string = "{\"result\": \"fail\"}"
  }
  
	fmt.Fprintf(w, "%s", json_return_string)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Run_aws_terraform(vm_name string, log_path string) int {    
  
	working_directory := "../aws_bootstrap/instances/" + vm_name
  
	command_string := "bash ./terraform_script.sh"
	log_path_2 := "../../" + log_path
	linux_command_line.Execute_command_in_background_and_wait(command_string, working_directory, log_path_2)		
  
  body, _ := logging.Get_log(log_path)
  string_log := string(body)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("Apply complete!"))
  match_result := regex_obj.FindString(string_log)
  
  if match_result == "" {
    return 0
  } else {    
    return 1
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Run_aws_sec_group_terraform(log_path string) int {    
  
  working_directory := "../aws_firewall"
  
  command_string := "bash ./terraform_script.sh"
  linux_command_line.Execute_command_in_background_and_wait(command_string, working_directory, log_path)    
  
  body, _ := logging.Get_log(log_path)
  string_log := string(body)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("Apply complete!"))
  match_result := regex_obj.FindString(string_log)
  
  if match_result == "" {
    return 0
  } else {    
    return 1
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
