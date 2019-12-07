package aws_handler

import (
	//"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	"html/template"
	//"os/exec"
	"../kv_store"
	"regexp"
	"../template_populator"
	//"../logging"
  //"strconv"
  "os"
  "../linux_command_line"
  //"../connectivity_check"
)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func file_exists(filename string) bool {
	file_exists := true
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		file_exists = false
	}
	return file_exists
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func strip_cr_from_string(input_string string) string {

	regex_obj := regexp.MustCompile("\\r")
	new_string := regex_obj.ReplaceAllString(input_string, "")

	return new_string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func strip_cr_cf_from_string(input_string string) string {

	regex_obj := regexp.MustCompile("\\n|\\r")
	new_string := regex_obj.ReplaceAllString(input_string, "")

	return new_string
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AWS_access_handler(w http.ResponseWriter, r *http.Request) {
	
	title := "Configure AWS Access"	 
	 
	retrieved_aws_key_id := ""
	retrieved_aws_secret_key_string := ""

	target_directory := "../aws_bootstrap/"
	linux_command_line.Execute_command_line("mkdir -p " + target_directory)
	
  template_directory := "./templates/cloud_vm_access_and_import/"

	webform_map := kv_store.Create_from_Webform(r)
	
	//###################################################
	 _, route_button_import_keys_from_env := webform_map["aws_key_import"]
	 //###################################################
	 
	if route_button_import_keys_from_env 	{
		
		webform_map["aws_secret_key_id"] = os.Getenv("AWS_ACCESS_KEY_ID")
		webform_map["aws_secret_key"] = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	
	if webform_map["aws_secret_key_id"] != "" || webform_map["aws_secret_key"] != "" {

		template_buffer := template_populator.Load_file(template_directory + "aws_key_id.j3")
	  template_populator.Populate_aws(template_buffer, webform_map, target_directory + "aws_key_id")	  

	  //webform_map["gcp_key_file_string"] = strip_cr_cf_from_string(webform_map["gcp_key_file_string"])

	  template_buffer = template_populator.Load_file(template_directory + "aws_secret_key.j3")
	  template_populator.Populate_aws(template_buffer, webform_map, target_directory + "aws_secret_key")	  

	  template_buffer = template_populator.Load_file(template_directory + "ec2.ini.j3")
	  template_populator.Populate_aws(template_buffer, webform_map, target_directory + "ec2.ini")

	  template_buffer = template_populator.Load_file(template_directory + "aws_credentials.j3")
	  template_populator.Populate_aws(template_buffer, webform_map, target_directory + "aws_credentials")
	}

	linux_command_line.Execute_command_line("cp " + template_directory + "ec2.py.j3 " + target_directory + "ec2.py")
	
	if file_exists(target_directory + "aws_key_id") {
  	retrieved_aws_key_id = template_populator.Load_file_as_string(target_directory + "aws_key_id")	
  }
  
  if file_exists(target_directory + "aws_secret_key") {
  	retrieved_aws_secret_key_string = template_populator.Load_file_as_string(target_directory + "aws_secret_key")
  } 

	p := &Page{Title: title, AWS_key_ID: retrieved_aws_key_id, AWS_secret_key: retrieved_aws_secret_key_string}
	t, _ := template.ParseFiles("templates/html/aws_key_access.html")
  t.Execute(w, p)
}
