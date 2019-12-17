// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package aws_handler

import (
	"fmt"
	"io/ioutil"
	//"log"
	"encoding/json"
	"net/http"
	"html/template"
	"os/exec"
	"regexp"
	"../kv_store"
	//"./template_populator"
	"path/filepath"
	"../logging"
  "strconv"
  "os"
  "../linux_command_line"
  "../connectivity_check"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_aws_topology () {
  
	cmd := exec.Command("python3", "ec2.py")

	current_directory, _  := os.Getwd()
	cmd.Dir, _ = filepath.Abs(current_directory + "/../aws_bootstrap/")

	out, _ := cmd.Output()
	    
	bytestream := []byte(out)
	ioutil.WriteFile("../aws_bootstrap/aws_output.txt", bytestream, 0644)
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func separate_region_and_az(region_string string, placement_string string) string {

	regex_obj := regexp.MustCompile(region_string + "(.*)")  
  match_array := regex_obj.FindStringSubmatch(placement_string)

  return match_array[1]
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_imported_aws_map (top_level_map map[string]map[string]string) {		

	fmt.Println("\nImporting...")
	
	get_aws_topology()
    
  file_buffer, _ := ioutil.ReadFile("../aws_bootstrap/aws_output.txt")
  
  imported_map_top_level := make(map[string]interface{})
  json.Unmarshal([]byte(file_buffer), &imported_map_top_level)
  
  imported_map := make(map[string]interface{})
  imported_map = imported_map_top_level["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})
  
	for _, instance := range imported_map {
		
		vm_name := instance.(map[string]interface{})["ec2_ip_address"].(string)
		top_level_map[vm_name] = make(map[string]string)

		top_level_map[vm_name]["vm_name"] = vm_name

		top_level_map[vm_name]["public_ip_address"] = instance.(map[string]interface{})["ec2_ip_address"].(string)
		top_level_map[vm_name]["private_ip_address"] = instance.(map[string]interface{})["ec2_private_ip_address"].(string)
		
		region_string := instance.(map[string]interface{})["ec2_region"].(string)
		placement_string := instance.(map[string]interface{})["ec2_placement"].(string)
		top_level_map[vm_name]["region"] = region_string
		zone := separate_region_and_az(region_string, placement_string)
		top_level_map[vm_name]["AZ"] = zone

		top_level_map[vm_name]["public_dns"] = instance.(map[string]interface{})["ec2_public_dns_name"].(string)		

		top_level_map[vm_name]["instance_construction_type"] = "Imported"
		top_level_map[vm_name]["platform_type"] = "aws"
	}

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
  AWS_key_ID string
  AWS_secret_key string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Launch_instance(w http.ResponseWriter, r *http.Request) {
	var title = "Launch AWS Instance"
  
  top_level_map := make(map[string]map[string]string)
  
  if _, err := os.Stat("./databases/aws_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/aws_store.json")
  }
  //////////////////////////////////////////////////////////////////////////////////////////////
  
  vm_map := kv_store.Create_from_Webform(r)
  vm_name := vm_map["vm_name"]
  
  linux_command_line.Execute_command_line("mkdir -p ../Logs")
  linux_command_line.Execute_command_line("mkdir -p ../ansible/ssh_keys/" + vm_name)
      
  Generate_terraform_resources(top_level_map, vm_name)
  // Populate_aws_firewall_default_section(vm_name)
  
	working_directory := "../aws_bootstrap/instances/" + vm_name
	//check(err)
  
	command_string := "bash ./terraform_script.sh"
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
func Terminate_instance(w http.ResponseWriter, r *http.Request) {
	var title = "Terminate Instance"
  
  vm_map := kv_store.Create_from_Webform(r)
  vm_name := vm_map["vm_name"]
  
  linux_command_line.Execute_command_line("mkdir -p ../Logs")
  
	working_directory := "../aws_bootstrap/instances/" + vm_name
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
func AWS_web_handler(w http.ResponseWriter, r *http.Request) {
  
  ami_map := make(map[string]map[string]string)
  ami_map["centos"] = make(map[string]string)
  ami_map["ubuntu"] = make(map[string]string)
  ami_map["centos"]["eu-west-1"] = "ami-3548444c"
  ami_map["centos"]["us-east-1"] = "ami-4bf3d731"
  ami_map["ubuntu"]["eu-west-1"] = "ami-2a7d75c0"
  ami_map["ubuntu"]["us-east-1"] = "ami-f449ac19"
  
	var title = "AWS (Managed VMs)"
  var show_ssh = "false"
  vm_name := ""
  cloud_platform := "aws"
  instance_construction_type := "Managed"
  
  webform_map := kv_store.Create_from_Webform(r)
  
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]
  }
  
  if _, err := os.Stat("./databases/aws_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/aws_store.json")
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
    kv_store.Export_to_json(top_level_map, "./databases/aws_store.json")    
    
  } else if route_button_addInstance {
  
    //webform_map := kv_store.Create_from_Webform(r)
    if webform_map["OS"] == "centos/7" {
    
      kv_store.Set_kv_entry(webform_map, "ami", ami_map["centos"][webform_map["region"]])
      kv_store.Set_kv_entry(webform_map, "package_manager", "yum")

    } else {
    
      kv_store.Set_kv_entry(webform_map, "ami", ami_map["ubuntu"][webform_map["region"]])
      kv_store.Set_kv_entry(webform_map, "package_manager", "apt-get")
    }
    
    top_level_map[vm_name] = webform_map
    kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", "aws")
    kv_store.Set_kv_entry(top_level_map[vm_name], "fqdn", "eu-west-1.compute.amazonaws.com")
    
    //vpc_cidr := Get_vpc_cidr(top_level_map[vm_name]["subnet_address"])
    //subnet_cidr := Get_subnet_cidr(top_level_map[vm_name]["subnet_address"])    
    //kv_store.Set_kv_entry(top_level_map[vm_name], "vpc_cidr", vpc_cidr)
    //kv_store.Set_kv_entry(top_level_map[vm_name], "subnet_cidr", subnet_cidr)
    
    if len(top_level_map) == 1 {      
      kv_store.Set_kv_entry(top_level_map[vm_name], "vm_type", "MAIN")
    }
    
    kv_store.Set_kv_entry(top_level_map[vm_name], "credentials_path", "../../" + "aws_credentials")
    Generate_aws_ssh_keys(top_level_map[vm_name], vm_name)
    
    kv_store.Set_kv_entry(top_level_map[vm_name], "instance_construction_type", instance_construction_type)

    kv_store.Export_to_json(top_level_map, "./databases/aws_store.json")
  
  } else if route_button_ssh {
    // Check connectivity
    //webform_map := kv_store.Create_from_Webform(r)
    
    connectivity_check.Store_aws_vm_resources_in_DB(top_level_map, vm_name)
    connectivity_check.Copy_aws_ssh_keys(vm_name)
    
    if vm_name != "" {
      connectivity_check.Check_connectivity(top_level_map, vm_name, "aws")
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
      
      kv_store.Export_to_json(top_level_map, "./databases/aws_store.json")
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
  }
  
  //###########################
  
  p := &Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: cloud_platform, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/aws_config.html")
  t.Execute(w, p)

  //fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AWS_imported_web_handler(w http.ResponseWriter, r *http.Request) {
  
  ami_map := make(map[string]map[string]string)
  ami_map["centos"] = make(map[string]string)
  ami_map["ubuntu"] = make(map[string]string)
  ami_map["centos"]["eu-west-1"] = "ami-3548444c"
  ami_map["centos"]["us-east-1"] = "ami-4bf3d731"
  ami_map["ubuntu"]["eu-west-1"] = "ami-2a7d75c0"
  ami_map["ubuntu"]["us-east-1"] = "ami-f449ac19"
  
	var title = "AWS (Imported VMs)"
  var show_ssh = "false"
  vm_name := ""
  cloud_platform := "aws"
  instance_construction_type := "Imported"
  
  webform_map := kv_store.Create_from_Webform(r)
  
  top_level_map := make(map[string]map[string]string)
  
  if len(webform_map) != 0 {
    vm_name = webform_map["vm_name"]
  }
  
  if _, err := os.Stat("./databases/aws_imported_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/aws_imported_store.json")
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
  	
  	get_imported_aws_map(top_level_map)
  	kv_store.Export_to_json(top_level_map, "./databases/aws_imported_store.json")

  } else if route_button_set_instance_details {  	  	
  	  	
  	kv_store.Set_kv_entry(top_level_map[vm_name], "ssh_private_key", strip_cr_from_string(webform_map["ssh_private_key"]))  	

  	kv_store.Set_kv_entry(top_level_map[vm_name], "OS", webform_map["OS"])  	
  	kv_store.Set_kv_entry(top_level_map[vm_name], "user_id", webform_map["user_id"])

  	target_directory := "../aws_bootstrap/instances/" + vm_name + "/"
		linux_command_line.Execute_command_line("mkdir -p " + target_directory)
		linux_command_line.Execute_command_line("mkdir -p ../ansible/ssh_keys/" + vm_name)

		Store_aws_ssh_keys_from_imported_VM(top_level_map[vm_name], vm_name)
  	connectivity_check.Copy_aws_ssh_keys(vm_name)

  	kv_store.Export_to_json(top_level_map, "./databases/aws_imported_store.json")

  } else if route_button_delete {

    //delete vm from map
    delete(top_level_map, vm_name)
    kv_store.Export_to_json(top_level_map, "./databases/aws_imported_store.json")    
    
  } else if route_button_addInstance {
  
    //webform_map := kv_store.Create_from_Webform(r)
    if webform_map["OS"] == "centos/7" {
    
      kv_store.Set_kv_entry(webform_map, "ami", ami_map["centos"][webform_map["region"]])
      kv_store.Set_kv_entry(webform_map, "package_manager", "yum")

    } else {
    
      kv_store.Set_kv_entry(webform_map, "ami", ami_map["ubuntu"][webform_map["region"]])
      kv_store.Set_kv_entry(webform_map, "package_manager", "apt-get")
    }
    
    top_level_map[vm_name] = webform_map
    kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", "aws")
    kv_store.Set_kv_entry(top_level_map[vm_name], "fqdn", "eu-west-1.compute.amazonaws.com")    
    
    //vpc_cidr := Get_vpc_cidr(top_level_map[vm_name]["subnet_address"])
    //subnet_cidr := Get_subnet_cidr(top_level_map[vm_name]["subnet_address"])    
    //kv_store.Set_kv_entry(top_level_map[vm_name], "vpc_cidr", vpc_cidr)
    //kv_store.Set_kv_entry(top_level_map[vm_name], "subnet_cidr", subnet_cidr)
    
    if len(top_level_map) == 1 {      
      kv_store.Set_kv_entry(top_level_map[vm_name], "vm_type", "MAIN")
    }
    
    kv_store.Set_kv_entry(top_level_map[vm_name], "credentials_path", "../../" + "aws_credentials")
    Generate_aws_ssh_keys(top_level_map[vm_name], vm_name)
    
    kv_store.Export_to_json(top_level_map, "./databases/aws_imported_store.json")
  
  } else if route_button_ssh {
    // Check connectivity
    //webform_map := kv_store.Create_from_Webform(r)
    
    //connectivity_check.Store_aws_vm_resources_in_DB(top_level_map, vm_name)
    connectivity_check.Copy_aws_ssh_keys(vm_name)
    
    if vm_name != "" {
      connectivity_check.Check_connectivity(top_level_map, vm_name, "aws")
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
      
      kv_store.Export_to_json(top_level_map, "./databases/aws_imported_store.json")
      
      fmt.Println(logged_gotty_ip_address + ":" + logged_gotty_port + "\n")
    }
  }
  
  //###########################
  
  p := &Page{Title: title, Show_SSH: show_ssh, Selected_VM: vm_name, Cloud_Platform: cloud_platform, Instance_construction_type: instance_construction_type}
  t, _ := template.ParseFiles("templates/html/aws_config_imported.html")
  t.Execute(w, p)

  //fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_aws_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/aws_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_aws_imported_map_ajax(w http.ResponseWriter, r *http.Request) {

	top_level_map := kv_store.Create_top_level_map_from_json_file("./databases/aws_imported_store.json")
  
	export_json_string := kv_store.Export_2level_map_to_json_string(top_level_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
