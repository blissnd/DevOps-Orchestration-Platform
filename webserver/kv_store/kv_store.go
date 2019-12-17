// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

package kv_store

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
  //"fmt"
)
/////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
        panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////
func Create_map() map[string]string {
	new_map := make(map[string]string)
	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Set_kv_entry(map_name map[string]string, key string, val string) map[string]string {
	map_name[key] = val
	return map_name
}
/////////////////////////////////////////////////////////////////////////////////
func Get_kv_entry(map_name map[string]string, key string) string {
	return map_name[key]
}
/////////////////////////////////////////////////////////////////////////////////
func Get_map_size(map_name map[string]string) int {
	return len(map_name)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_to_json(map_name map[string]map[string]string, pathname string) {
	json_string_to_export, _ := json.Marshal(map_name)
	
	bytestream := []byte(json_string_to_export)
	err := ioutil.WriteFile(pathname, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_3_level_map_to_json(map_name map[string]map[string]map[string]string, pathname string) {
	json_string_to_export, _ := json.Marshal(map_name)
	
	bytestream := []byte(json_string_to_export)
	err := ioutil.WriteFile(pathname, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_4_level_map_to_json(map_name map[string]map[string]map[string]map[string]string, pathname string) {
	json_string_to_export, _ := json.Marshal(map_name)
	
	bytestream := []byte(json_string_to_export)
	err := ioutil.WriteFile(pathname, bytestream, 0644)
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_1level_map_to_json_string(map_name map[string]string) string {
	json_string_to_export, _ := json.Marshal(map_name)
	
	return string(json_string_to_export)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_2level_map_to_json_string(map_name map[string]map[string]string) string {
	json_string_to_export, _ := json.Marshal(map_name)
	
	return string(json_string_to_export)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_3level_map_to_json_string(map_name map[string]map[string]map[string]string) string {
	json_string_to_export, _ := json.Marshal(map_name)
	
	return string(json_string_to_export)
}
/////////////////////////////////////////////////////////////////////////////////
func Export_4level_map_to_json_string(map_name map[string]map[string]map[string]map[string]string) string {
	json_string_to_export, _ := json.Marshal(map_name)
	
	return string(json_string_to_export)
}
/////////////////////////////////////////////////////////////////////////////////
func Import_from_json(pathname string) map[string]map[string]string {
	input_stream, err := ioutil.ReadFile(pathname)
	check(err)
	
	new_map := make(map[string]map[string]string)
	json.Unmarshal([]byte(input_stream), &new_map)
	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Import_from_3_level_json(pathname string) map[string]map[string]map[string]string {
	input_stream, err := ioutil.ReadFile(pathname)
	check(err)
	
	new_map := make(map[string]map[string]map[string]string)
	json.Unmarshal([]byte(input_stream), &new_map)
	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Import_from_4_level_json(pathname string) map[string]map[string]map[string]map[string]string {
	input_stream, err := ioutil.ReadFile(pathname)
	check(err)
	
	new_map := make(map[string]map[string]map[string]map[string]string)
	json.Unmarshal([]byte(input_stream), &new_map)
	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Import_single_level_map_from_json(pathname string) map[string]string {
	input_stream, err := ioutil.ReadFile(pathname)
	check(err)
	
	new_map := make(map[string]string)
	json.Unmarshal([]byte(input_stream), &new_map)
	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Create_from_Webform(r *http.Request) map[string]string {
	new_map := Create_map()
	r.ParseForm()
	
	for key, val := range r.Form {
		new_map = Set_kv_entry(new_map, key, string(val[0]))
	}

	return new_map
}
/////////////////////////////////////////////////////////////////////////////////
func Create_top_level_map_from_json_file(pathname string) map[string]map[string]string {

  top_level_map := Import_from_json(pathname)
  
  return top_level_map
}
/////////////////////////////////////////////////////////////////////////////////
func Create_3_level_map_from_json_file(pathname string) map[string]map[string]map[string]string {

  three_level_map := Import_from_3_level_json(pathname)
  
  return three_level_map
}
/////////////////////////////////////////////////////////////////////////////////
func Create_4_level_map_from_json_file(pathname string) map[string]map[string]map[string]map[string]string {

  four_level_map := Import_from_4_level_json(pathname)
  
  return four_level_map
}
////////////////////////////////////////////////////////////////////////////////
func Create_single_level_map_from_json_file(pathname string) map[string]string {

  single_level_map := Import_single_level_map_from_json(pathname)
  
  return single_level_map
}
/////////////////////////////////////////////////////////////////////////////////
func Create_Nexus_registry_db(vm_name string, ip_address string, fqdn string) {

	nexus_registry_map := make(map[string]string)
  
  nexus_registry_map["vm_name"] = vm_name
  nexus_registry_map["ip_address"] = ip_address
  nexus_registry_map["fqdn"] = fqdn
  
  json_string_to_export, _ := json.Marshal(nexus_registry_map)
  
  bytestream := []byte(json_string_to_export)
		err := ioutil.WriteFile("./databases/nexus_registry.json", bytestream, 0644)
		check(err)
}
/////////////////////////////////////////////////////////////////////////////////
