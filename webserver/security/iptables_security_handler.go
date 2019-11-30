package security

import (
	"fmt"
	"time"
	//"encoding/json"
	//"io/ioutil"
	"regexp"
	//"log"
	"net/http"
	"html/template"
	//"os/exec"
	"../kv_store"
	"../general_utility_web_handler"
	//"./template_populator"
	"../logging"
  //"strconv"
  "os"
  //"../linux_command_line"
  //"../connectivity_check"
  "../ansible_handler"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func IP_tables_security_web_handler(w http.ResponseWriter, r *http.Request) {
  
  //////////////////////////////////////////////////////////////////////////////
    
  iptables_security_map := make(map[string]map[string]map[string]string)
  
  if _, err := os.Stat("./databases/iptables_security_store.json"); err == nil {
    iptables_security_map = kv_store.Create_3_level_map_from_json_file("./databases/iptables_security_store.json")
  }
  
  //////////////////////////////////////////////////////////////////////////////
  
  webform_map := kv_store.Create_from_Webform(r)    

  vm_name := webform_map["vm_name"]
  platform_type := webform_map["platform_type"]
  instance_construction_type := webform_map["instance_construction_type"]
  firewall_rule_map_key := webform_map["protocol"] + "-" + webform_map["port"]

  title := vm_name + " / " + platform_type + " / " + instance_construction_type  
    
  //##########################################
  _, route_button_add_Rule := webform_map["Add Rule"]
  _, route_button_delete := webform_map["button_delete"]  
  //##########################################
  log_path_prefix := "../Logs/ansible/security"
  //##########################################    
  
  if route_button_add_Rule && webform_map[vm_name] != "null" {            
           
    if _, exist := iptables_security_map[vm_name]; exist == false {    	
    	iptables_security_map[vm_name] = make(map[string]map[string]string)
    }

    if _, exist := iptables_security_map[vm_name][firewall_rule_map_key]; exist == false {    	
    	iptables_security_map[vm_name][firewall_rule_map_key] = make(map[string]string)
    }

    iptables_security_map[vm_name][firewall_rule_map_key] = webform_map
    
    kv_store.Export_3_level_map_to_json(iptables_security_map, "./databases/iptables_security_store.json")
    
  } else if route_button_delete {
    
    vm_name = webform_map["vm_name"]    
    firewall_rule_map_key := webform_map["protocol"] + "-" + webform_map["port"]
    
    delete(iptables_security_map[vm_name], firewall_rule_map_key)
    
    kv_store.Export_3_level_map_to_json(iptables_security_map, "./databases/iptables_security_store.json")
  }
  
  //###########################
    
  p := &Page{Title: title, LogPath: log_path_prefix, Selected_VM: vm_name, Cloud_Platform: platform_type, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/iptables_security.html")
  t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_iptables_security_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := make(map[string]map[string]map[string]string)

	if _, err := os.Stat("./databases/iptables_security_store.json"); err == nil {
    top_level_map = kv_store.Create_3_level_map_from_json_file("./databases/iptables_security_store.json")
  }

	export_json_string := kv_store.Export_3level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AJAX_modify_iptables_firewall(w http.ResponseWriter, r *http.Request) {
  
  ip_address := ""
  success_or_failure := 0
  
  //////////////////////////////////////////////////////////////////////////////
  iptables_security_map := make(map[string]map[string]map[string]string)
  
  if _, err := os.Stat("./databases/iptables_security_store.json"); err == nil {
    iptables_security_map = kv_store.Create_3_level_map_from_json_file("./databases/iptables_security_store.json")
  }  
  //////////////////////////////////////////////////////////////////////////////
  ajax_map := kv_store.Create_from_Webform(r)
  //////////////////////////////////////////////////////////////////////////////      
  vm_name := ajax_map["vm_name"]
  firewall_rule_map_key := ajax_map["protocol"] + "-" + ajax_map["port"]

  fmt.Println(vm_name)
  
  platform_type := ajax_map["platform_type"]
  instance_construction_type := ajax_map["instance_construction_type"]
  
  //////////////////////////////////////////////////////////////////////////////
  top_level_map := make(map[string]map[string]string)
  
  json_store_name := general_utility_web_handler.Get_json_store_name(platform_type, instance_construction_type)

  if _, err := os.Stat("./databases/vm_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/" + json_store_name)
  }
  //////////////////////////////////////////////////////////////////////////////

  ip_address = top_level_map[vm_name]["public_ip_address"]
  
  //##########################################
  log_path := ajax_map["log_path"]
  //##########################################   
  
  port_number := ajax_map["port"]
  old_state := ajax_map["port_state"]
  protocol := ajax_map["protocol"]
  
  fmt.Println(ajax_map)
  fmt.Println(top_level_map)
  
  if old_state == "CLOSED" {
    // Open Port
    Generate_firewall_port_open_rule(iptables_security_map, vm_name, firewall_rule_map_key, protocol, port_number)

    ansible_handler.Generate_inventory(vm_name, top_level_map, "../ansible/open_firewall_port", "", "", nil)
    
    working_directory := "../ansible/open_firewall_port/"
    ansible_handler.Ansible_generic_runner(vm_name, ip_address, working_directory, log_path)
    
    iptables_security_map[vm_name][firewall_rule_map_key]["state"] = "OPEN"
    
  } else {
    // Close Port
    Generate_firewall_port_close_rule(iptables_security_map, vm_name, firewall_rule_map_key, protocol, port_number)
    
    ansible_handler.Generate_inventory(vm_name, top_level_map, "../ansible/close_firewall_port", "", "", nil)

    working_directory := "../ansible/close_firewall_port/"    
    ansible_handler.Ansible_generic_runner(vm_name, ip_address, working_directory, log_path)
    
    iptables_security_map[vm_name][firewall_rule_map_key]["state"] = "CLOSED"
  }
  
  success_or_failure = Check_log(log_path)
    
  if success_or_failure == 0 {
    if iptables_security_map[vm_name][firewall_rule_map_key]["state"] == "CLOSED" {

      iptables_security_map[vm_name][firewall_rule_map_key]["state"] = "OPEN"
      
    } else {
      iptables_security_map[vm_name][firewall_rule_map_key]["state"] = "CLOSED"
    }
  }
    
  kv_store.Export_3_level_map_to_json(iptables_security_map, "./databases/iptables_security_store.json")
  
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
func Check_log(log_path string) int {    
  
  start_time := time.Now()

  match_result := ""
  wait_for_log_entry := false
  string_log := ""
  elapsed := 0
  
  for wait_for_log_entry == false && elapsed < 20 {
    
    end_time := time.Now()
    elapsed = int(end_time.Sub(start_time)) / 1000000000
    
    body, _ := logging.Get_log(log_path)
    string_log = string(body)    

    regex_obj := regexp.MustCompile("ok=.*changed=.*unreachable=.*failed=.")
    match_result = regex_obj.FindString(string_log)
    
    if match_result != "" {
       wait_for_log_entry = true
    }
  }
    
  regex_obj := regexp.MustCompile("ok=.*changed=.*unreachable=0.*failed=0.")
  
  match_result = regex_obj.FindString(string_log)
    
  if match_result == "" {
    return 0
  } else {    
    return 1
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
