package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"os/exec"
	"./kv_store"
	//"./map_template"
	"./vm_launcher"
	"./logging"
  "strconv"
  "./ansible_handler"
  "./kubernetes_handler"
	"./general_utility_web_handler"
  "./virtualbox_handler"
  "./aws_handler"
  "./gcp_handler"
  "./security"
)
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
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func loadPage(title string) (*Page, error) {
		
	linux_command := exec.Command("bash")
	linux_command.Args = append(linux_command.Args, "./replace_escapes.sh")
	linux_command.Args = append(linux_command.Args, "../vm_bootstrap/Ansible/ubuntu-xenial-16.04-cloudimg-console.log")
	linux_command.Args = append(linux_command.Args, "./logs/test_log.log")
	
	err := linux_command.Run()
	
	if err != nil {
			panic(err)
	}
	
	filename := "./logs/test_log.log"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
			return nil, err
	}
	
	return &Page{Title: title, Body: (body)}, nil
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func index_handler(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("./index.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_nav_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("./nav_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_title_bar_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("./title_bar.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_operations_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("./operations_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_ssh_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("templates/html/ssh_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_ssh_popup_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("templates/html/ssh_popup_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_listening_port_popup_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("templates/html/listening_port_popup_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handler_kubernetes_dashboard_token_popup_frame(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("templates/html/kubernetes_dashboard_token_popup_frame.html")  // Parse template file.
    t.Execute(w, nil)  // merge.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func log_capture_handler(w http.ResponseWriter, r *http.Request) {
	var title = "Log"

	p, err := loadPage(title)

	if err != nil {
			p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("templates/html/log.html")
	t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ssh_test_handler(w http.ResponseWriter, r *http.Request) {
	//var title = "Test"
	
	//linux_command := exec.Command("gnome-terminal -x bash -c 'echo \"hello\"; ssh localhost'")
	
	linux_command := exec.Command("gnome-terminal")
	
	linux_command.Args = append(linux_command.Args, "-x")
	linux_command.Args = append(linux_command.Args, "ssh")
	linux_command.Args = append(linux_command.Args, "blissnd@localhost")
	
	var err = linux_command.Run()
	
	if err != nil {
		panic(err)
	}
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func configuration_handler(w http.ResponseWriter, r *http.Request) {
	var title = "Configure"
	
	p := &Page{Title: title}

	t, _ := template.ParseFiles("templates/html/vm_config.html")
	t.Execute(w, p)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func log_refresh_handler(w http.ResponseWriter, r *http.Request) {

	ajax_map := kv_store.Create_from_Webform(r)
	logpath := kv_store.Get_kv_entry(ajax_map, "log_path")
  log_position := kv_store.Get_kv_entry(ajax_map, "log_position")
  log_position_int, _ := strconv.Atoi(log_position)
  
  log_content, new_log_position_int, err := logging.Get_remaining_log(logpath, log_position_int)
	check(err)

  kv_store.Set_kv_entry(ajax_map, "log_content", string(log_content))
  kv_store.Set_kv_entry(ajax_map, "LogPosition", strconv.Itoa(new_log_position_int))
    
	export_json_string := kv_store.Export_1level_map_to_json_string(ajax_map)
	fmt.Fprintf(w, "%s", export_json_string)
	//fmt.Println(ajax_map)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
  //http.Handle("/", http.FileServer(http.Dir(".")))
  
  http.HandleFunc("/", index_handler)

  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
  http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))  
  
  //////////////
  http.HandleFunc("/nav_frame/", handler_nav_frame)
  http.HandleFunc("/title_bar_frame/", handler_title_bar_frame)
  http.HandleFunc("/operations_frame/", handler_operations_frame)
  http.HandleFunc("/ssh_frame/", handler_ssh_frame)
  http.HandleFunc("/ssh_popup_frame/", handler_ssh_popup_frame)
  http.HandleFunc("/listening_port_popup_frame/", handler_listening_port_popup_frame)
  //////////////

  http.HandleFunc("/virtualbox/", virtualbox_handler.Virtualbox_web_handler)
  
	http.HandleFunc("/launch/", ssh_test_handler)

	http.HandleFunc("/configure/create_vm/", vm_launcher.VM_launch_handler)
  http.HandleFunc("/configure/launch_aws_instance/", aws_handler.Launch_instance)
  http.HandleFunc("/configure/launch_gcp_instance/", gcp_handler.Launch_instance)
  
  http.HandleFunc("/aws/terminate_aws_instance/", aws_handler.Terminate_instance)
  http.HandleFunc("/gcp/terminate_gcp_instance", gcp_handler.Terminate_instance)
  
	http.HandleFunc("/configure/create_docker/", ansible_handler.Docker_launch_handler)
	http.HandleFunc("/configure/create_docker_from_image/", ansible_handler.Docker_image_launch_handler)
	http.HandleFunc("/configure/docker_commit/", ansible_handler.Docker_registry_commit_handler)
  
  http.HandleFunc("/configure/", configuration_handler)
	http.HandleFunc("/get_log/", log_refresh_handler)  
  http.HandleFunc("/ansible/run_ansible/", ansible_handler.Ansible_runner)
  
  http.HandleFunc("/get_vm_map/", virtualbox_handler.Get_vm_map_ajax)
  http.HandleFunc("/get_aws_map/", aws_handler.Get_aws_map_ajax)
  http.HandleFunc("/get_gcp_map/", gcp_handler.Get_gcp_map_ajax)
  http.HandleFunc("/get_aws_imported_map/", aws_handler.Get_aws_imported_map_ajax)
  http.HandleFunc("/get_gcp_imported_map/", gcp_handler.Get_gcp_imported_map_ajax)
  
  http.HandleFunc("/get_iptables_security_map_ajax/", security.Get_iptables_security_map_ajax)
  http.HandleFunc("/get_aws_firewall_map_ajax/", security.Get_aws_firewall_map_ajax)
  http.HandleFunc("/get_gcp_firewall_map_ajax/", security.Get_gcp_firewall_map_ajax)
  
	http.HandleFunc("/specific_vm/", ansible_handler.Generic_VM_web_handler)

	http.HandleFunc("/kubernetes/", kubernetes_handler.Kubernetes_web_handler)
	http.HandleFunc("/get_kubernetes_dashboard_token_ajax/", kubernetes_handler.Get_kubernetes_dashboard_token_ajax)
	http.HandleFunc("/get_master_map_ajax/", kubernetes_handler.Get_master_map_ajax)
	http.HandleFunc("/kubernetes_dashboard_token_popup_frame/", handler_kubernetes_dashboard_token_popup_frame)

  http.HandleFunc("/docker/", ansible_handler.Docker_web_handler)
  http.HandleFunc("/get_docker_map/", ansible_handler.Get_docker_map_ajax)
  
  http.HandleFunc("/aws/", aws_handler.AWS_web_handler)
  http.HandleFunc("/gcp/", gcp_handler.GCP_web_handler)
  http.HandleFunc("/aws_imported/", aws_handler.AWS_imported_web_handler)
  http.HandleFunc("/gcp_imported/", gcp_handler.GCP_imported_web_handler)
  http.HandleFunc("/aws_key_access/", aws_handler.AWS_access_handler)
  http.HandleFunc("/gcp_key_access/", gcp_handler.GCP_access_handler)
  
  http.HandleFunc("/iptables_security/", security.IP_tables_security_web_handler)
  http.HandleFunc("/aws_hypervisor_security/", security.AWS_hypervisor_security_web_handler)
  http.HandleFunc("/gcp_hypervisor_security/", security.GCP_hypervisor_security_web_handler)
  
  http.HandleFunc("/ajax_modify_iptables_firewall/", security.AJAX_modify_iptables_firewall)
  http.HandleFunc("/ajax_aws_modify_hypervisor_firewall/", security.AJAX_AWS_modify_hypervisor_firewall)
  http.HandleFunc("/ajax_gcp_modify_hypervisor_firewall/", security.AJAX_GCP_modify_hypervisor_firewall)
  
  http.HandleFunc("/get_generic_security_log/", general_utility_web_handler.Generic_security_log_web_handler)
  http.HandleFunc("/ajax_get_listening_ports/", general_utility_web_handler.Get_listening_ports)
  
	log.Fatal(http.ListenAndServe(":6543", nil))
}
