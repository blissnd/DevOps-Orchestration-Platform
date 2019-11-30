package linux_command_line

import (
	"fmt"
	//"io/ioutil"
	"strings"
	//"regexp"
	"os/exec"
  "context"
  "time"
	//"os"
	"../logging"	
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
      fmt.Println(e.Error() + "\n")
      //panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Copy_vagrant_private_key(vm_name string) {

  source_path := "../vm_bootstrap/VMs/" + vm_name + "/.vagrant/machines/" + vm_name + "/virtualbox/private_key"
  dest_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  Execute_command_line("cp " + source_path + " " + dest_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Clean_ssh_specify_port(ip_address string, vm_name string, port string) {
  
  // Copy private key from vagrant directory
  Copy_vagrant_private_key(vm_name)
  
  command_string := "rm ~/.ssh/known_hosts"
	log_path := "../Logs/ansible/ssh_fix_log"
	Execute_command_in_background(command_string, ".", log_path)
  
  command_string = "ssh-keygen -f /root/.ssh/known_hosts -R [" + ip_address + "]:" + port
  fmt.Println(command_string)
	log_path = "../Logs/ansible/ssh_fix_log"
  Execute_command_in_background(command_string, ".", log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Clean_ssh(ip_address string, vm_name string) {
  
  // Copy private key from vagrant directory
  Copy_vagrant_private_key(vm_name)
  
  command_string := "rm ~/.ssh/known_hosts"
	log_path := "../Logs/ansible/ssh_fix_log"
	Execute_command_in_background(command_string, ".", log_path)
  
  command_string = "ssh-keygen -f /root/.ssh/known_hosts -R " + ip_address
  fmt.Println(command_string)
	log_path = "../Logs/ansible/ssh_fix_log"
  Execute_command_in_background(command_string, ".", log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Execute_command_line(command_string string) (string, error) {
	command_array := strings.Split(command_string, " ")
  
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel() // The cancel should be deferred so resources are cleaned up
  
	linux_command := exec.CommandContext(ctx, command_array[0])
  //logging.Store_command_log(linux_command, "TEST_LOG")
	array_length := len(command_array)
	command_array = command_array[1: array_length]
	
	for current_arg := range command_array {
		linux_command.Args = append(linux_command.Args, command_array[current_arg])
	}

	//linux_command.Dir = working_directory
	
  //command_output := logging.Get_command_log(linux_command)
  
  //err := linux_command.Start()
  command_output, err := linux_command.Output()
  check(err)
  
  if ctx.Err() == context.DeadlineExceeded {
		command_output = []byte("Command timed out")
  }
  
  return string(command_output), err
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Execute_command_in_background(command_string string, working_directory string, log_path string) {
	
	command_array := strings.Split(command_string, " ")
	
	linux_command := exec.Command(command_array[0])
	array_length := len(command_array)
	command_array = command_array[1: array_length]
	
	for current_arg := range command_array {
		linux_command.Args = append(linux_command.Args, command_array[current_arg])
	}

	linux_command.Dir = working_directory

	logging.Store_command_log_and_run(linux_command, log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Execute_command_in_background_and_wait(command_string string, working_directory string, log_path string) {

	command_array := strings.Split(command_string, " ")
	
	linux_command := exec.Command(command_array[0])
	array_length := len(command_array)
	command_array = command_array[1: array_length]
	
	for current_arg := range command_array {
		linux_command.Args = append(linux_command.Args, command_array[current_arg])
	}

	linux_command.Dir = working_directory
	
	logging.Store_command_log_and_wait(linux_command, log_path)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Run_gotty(user_id string, vm_name string, ip_address string, ssh_port string, gotty_port string, log_path string) {

	Clean_ssh_specify_port(ip_address, vm_name, ssh_port)

	ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"

	command_line_string := "../software/gotty -w --timeout 5 -config ./linux_command_line/gotty_config --port " + gotty_port + " --close-timeout 5 ssh " + user_id + "@" + ip_address + " -p " + ssh_port + " -i " + ssh_private_key_path + " -o StrictHostKeyChecking=no "
	fmt.Println(command_line_string + "\n")

	Execute_command_in_background(command_line_string, ".", log_path)

	//out, err := Execute_command_line(command_line_string)
	//check(err)
	//fmt.Println(out + "\n")

	/*
	linux_command := exec.Command("../software/gotty")
	linux_command.Args = append(linux_command.Args, "ssh")
	linux_command.Args = append(linux_command.Args, user_id + "@" + ip_address)
  
  linux_command.Args = append(linux_command.Args, "-p")
  linux_command.Args = append(linux_command.Args, port)
  
  linux_command.Args = append(linux_command.Args, "-i")
  linux_command.Args = append(linux_command.Args, ssh_private_key_path)
	
	//var err = linux_command.Run()
  //check(err)
  */
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_remote_listening_ports(user_id string, vm_name string, ip_address string, ssh_port string) string {

	ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"

	output, _ := Execute_command_line("ssh -o StrictHostKeyChecking=no " + user_id + "@" + ip_address + " -p " + ssh_port + " -i " + ssh_private_key_path + " sudo netstat -tuplna | grep -i listen")

	return output
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Get_kubernetes_dashboard_token(user_id string, vm_name string, ip_address string, ssh_port string) string {

	ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"

	linux_command := "ssh -o StrictHostKeyChecking=no " + user_id + "@" + ip_address + " -p " + ssh_port + " -i " + ssh_private_key_path
	linux_command += " kubectl -n kube-system describe secrets $(kubectl -n kube-system get secret | grep \"service-controller-token-*\" | awk '{print $1'})"

	output, _ := Execute_command_line(linux_command)

	return output
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SSH_to_remote_host(user_id string, vm_name string, ip_address string, port string) {
  
  Clean_ssh_specify_port(ip_address, vm_name, port)
  
  ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  
  linux_command := exec.Command("gnome-terminal")
	
	linux_command.Args = append(linux_command.Args, "-x")
	linux_command.Args = append(linux_command.Args, "ssh")
	linux_command.Args = append(linux_command.Args, user_id + "@" + ip_address)
  
  linux_command.Args = append(linux_command.Args, "-p")
  linux_command.Args = append(linux_command.Args, port)
  
  linux_command.Args = append(linux_command.Args, "-i")
  linux_command.Args = append(linux_command.Args, ssh_private_key_path)
	
	var err = linux_command.Run()
  check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SSH_to_remote_aws_host(user_id string, vm_name string, ip_address string, port string) {

  ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  
  linux_command := exec.Command("gnome-terminal")
	
	linux_command.Args = append(linux_command.Args, "-x")
	linux_command.Args = append(linux_command.Args, "ssh")
	linux_command.Args = append(linux_command.Args, user_id + "@" + ip_address)
  
  linux_command.Args = append(linux_command.Args, "-p")
  linux_command.Args = append(linux_command.Args, port)
  
  linux_command.Args = append(linux_command.Args, "-i")
  linux_command.Args = append(linux_command.Args, ssh_private_key_path)
	
	var err = linux_command.Run()
  check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SSH_to_remote_gcp_host(user_id string, vm_name string, ip_address string, port string) {

  ssh_private_key_path := "../ansible/ssh_keys/" + vm_name + "/private_key"
  
  linux_command := exec.Command("gnome-terminal")
	
	linux_command.Args = append(linux_command.Args, "-x")
	linux_command.Args = append(linux_command.Args, "ssh")
	linux_command.Args = append(linux_command.Args, user_id + "@" + ip_address)
  
  linux_command.Args = append(linux_command.Args, "-p")
  linux_command.Args = append(linux_command.Args, port)
  
  linux_command.Args = append(linux_command.Args, "-i")
  linux_command.Args = append(linux_command.Args, ssh_private_key_path)
	
	var err = linux_command.Run()
  check(err)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
