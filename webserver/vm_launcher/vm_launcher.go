package vm_launcher

import (
	"fmt"
  "reflect"
	//"io/ioutil"
	//"strings"
	//"regexp"
  "net/http"
	"os"
	"../linux_command_line"
  "../kv_store"
  "html/template"
  "../template_populator"
  "../logging"
  "../connectivity_check"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Empty struct{}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Page struct {
	Title string
	Body []byte
	LogPath string
  LogPosition int
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
        //panic(e)
        fmt.Println(reflect.TypeOf(Empty{}).PkgPath() + " => " + e.Error() + "\n")
    }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Launch_VM(vm_name string) {
	
  linux_command_line.Execute_command_line("mkdir -p ../Logs")
  linux_command_line.Execute_command_line("mkdir -p ../vm_bootstrap/VMs")
  linux_command_line.Execute_command_line("mkdir -p ../ansible/ssh_keys/" + vm_name)
  
	working_directory := "../vm_bootstrap/VMs/" + vm_name
	//check(err)
	
	command_string := "vagrant up --provision"
	log_path := "../Logs/" + vm_name + ".log"
	linux_command_line.Execute_command_in_background(command_string, working_directory, log_path)		
  
  // Copy private key from vagrant directory
  connectivity_check.Copy_vagrant_private_key(vm_name)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func VM_launch_handler(w http.ResponseWriter, r *http.Request) {
	var title = "Launch"
	top_level_map := make(map[string]map[string]string)
  
  if _, err := os.Stat("./databases/vm_store.json"); err == nil {
    top_level_map = kv_store.Create_top_level_map_from_json_file("./databases/vm_store.json")
  }
  
	vm_map := kv_store.Create_from_Webform(r)
  vm_name := vm_map["vm_name"]
  //top_level_map[vm_name] = vm_map
  kv_store.Set_kv_entry(top_level_map[vm_name], "platform_type", "virtualbox")
  
  // Check connectivity
  connectivity_check.Check_connectivity(top_level_map, vm_name, "virtualbox") 
  
  if top_level_map[vm_name]["OS"] == "centos/7" {
    kv_store.Set_kv_entry(vm_map, "bootstrap_script", "../../bootstrap_redhat.sh")
  } else {
    kv_store.Set_kv_entry(vm_map, "bootstrap_script", "../../bootstrap.sh")
  }
  
  template_buffer := template_populator.Load_file("templates/vagrant/Vagrantfile.j3")
	template_populator.Populate(template_buffer, vm_map, "templates/vagrant/Vagrantfile")
  
	Launch_VM(vm_name)
	logpath := "../Logs/" + vm_name + ".log"
	body, err := logging.Get_log(logpath)
	
	check(err)
	
	p := &Page{Title: title, Body: (body), LogPath: logpath, LogPosition: 0}
	//fmt.Println(body)
	t, _ := template.ParseFiles("templates/html/log.html")
	t.Execute(w, p)
	
  kv_store.Export_to_json(top_level_map, "./databases/vm_store.json")
  
	//fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
