package gcp_handler

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
  "path/filepath"
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
func get_gcp_service_account_email(gcp_key_file_string string) string {	

	match_result := ""

	regex_string := "(.|\n)*?client_email\":.*?\"(.*?)\""
	
  regex_obj := regexp.MustCompile(regex_string)
  match_result = regex_obj.FindString(gcp_key_file_string)      
  
  if match_result == "" {
    return ""
  } else {  	
  	match_array := regex_obj.FindStringSubmatch(gcp_key_file_string)
    return match_array[2]
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GCP_access_handler(w http.ResponseWriter, r *http.Request) {

	title := "Configure GCP Access"
	
	retrieved_gcp_project_id := ""
	retrieved_gcp_json_key_file_string := ""

	target_directory := "../gcp_bootstrap/"
	target_directory_absolute_path, _ := filepath.Abs(target_directory)
  template_directory := "./templates/cloud_vm_access_and_import/"

	webform_map := kv_store.Create_from_Webform(r)
	
	if webform_map["gcp_project_id"] != "" || webform_map["gcp_key_file_string"] != "" {

		service_account_email := get_gcp_service_account_email(webform_map["gcp_key_file_string"])
		webform_map["gcp_service_account_email"] = service_account_email

		template_buffer := template_populator.Load_file(template_directory + "cloud_account_params.py.j3")
	  template_populator.Populate_gcp(template_buffer, webform_map, target_directory + "secrets.py")	  

	  webform_map["gcp_key_file_string"] = strip_cr_cf_from_string(webform_map["gcp_key_file_string"])

	  template_buffer = template_populator.Load_file(template_directory + "gcp_credentials.json.j3")
	  template_populator.Populate_gcp(template_buffer, webform_map, target_directory + "gcp_credentials.json")

	  template_buffer = template_populator.Load_file(template_directory + "gcp_project_id.j3")
	  template_populator.Populate_gcp(template_buffer, webform_map, target_directory + "gcp_project_id")

	  webform_map["cloud_account_params_path"] = target_directory_absolute_path + "/secrets.py"

	  template_buffer = template_populator.Load_file(template_directory + "gce.ini.j3")
	  template_populator.Populate_gcp(template_buffer, webform_map, target_directory + "gce.ini")
	}

  linux_command_line.Execute_command_line("cp " + template_directory + "gce.py.j3 " + target_directory + "gce.py")

  if file_exists(target_directory + "gcp_project_id") {
  	retrieved_gcp_project_id = template_populator.Load_file_as_string(target_directory + "gcp_project_id")	
  }
  
  if file_exists(target_directory + "gcp_credentials.json") {
  	retrieved_gcp_json_key_file_string = template_populator.Load_file_as_string(target_directory + "gcp_credentials.json")
  }  

	p := &Page{Title: title, GCP_Project: retrieved_gcp_project_id, GCP_Json_key: retrieved_gcp_json_key_file_string}
	t, _ := template.ParseFiles("templates/html/gcp_key_access.html")
  t.Execute(w, p)
}
