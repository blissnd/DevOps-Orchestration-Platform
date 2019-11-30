package logging

import (
	"fmt"
	//"io/ioutil"
	//"strings"
	//"regexp"
	"os"
	"os/exec"
	"io/ioutil"
  "bufio"
  "regexp"
  "bytes"
  "reflect"
  "time"
)
/////////////////////////////////////////////////////////////////////////////////
type Empty struct{}
/////////////////////////////////////////////////////////////////////////////////
func check(e error) {
    if e != nil {
        fmt.Println(reflect.TypeOf(Empty{}).PkgPath() + " => " + e.Error() + "\n")
        //panic(e)
    }
}
/////////////////////////////////////////////////////////////////////////////////
func Store_command_log(linux_command *exec.Cmd, log_path string) {
	outfile, err := os.Create(log_path)
	check(err)
	defer outfile.Close()
	linux_command.Stdout = outfile
	linux_command.Stderr = outfile
}
/////////////////////////////////////////////////////////////////////////////////
func Store_command_log_and_run(linux_command *exec.Cmd, log_path string) {
  //dir, _ := os.Getwd()

  outfile, err := os.Create(log_path)
	check(err)
	defer outfile.Close()
	linux_command.Stdout = outfile
	linux_command.Stderr = outfile
	
	err = linux_command.Start()
	check(err)
}
/////////////////////////////////////////////////////////////////////////////////
func Store_command_log_and_wait(linux_command *exec.Cmd, log_path string) {
  //dir, _ := os.Getwd()

  outfile, err := os.Create(log_path)
	check(err)
	defer outfile.Close()
	linux_command.Stdout = outfile
	linux_command.Stderr = outfile
	
	err = linux_command.Start()
	check(err)
	linux_command.Wait()
}
/////////////////////////////////////////////////////////////////////////////////
func Get_command_log(linux_command *exec.Cmd) string {

  var output_buffer bytes.Buffer
  linux_command.Stdout = &output_buffer
  linux_command.Stderr = &output_buffer

	err := linux_command.Start()
	check(err)
  return output_buffer.String()
}
/////////////////////////////////////////////////////////////////////////////////
func Get_log(path string) ([]byte, error) {
	body, err := ioutil.ReadFile(path)
	return body, err
}
/////////////////////////////////////////////////////////////////////////////////
func Get_remaining_log(logpath string, log_position int) ([]byte, int, error) {
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta("\x1b"))
  regex_obj2 := regexp.MustCompile(regexp.QuoteMeta("\x0D"))
  
  log_file_handle, err := os.Open(logpath) 
  check(err)
  if err != nil {
    os.Create(logpath)
  }
  
  scanner := bufio.NewScanner(log_file_handle)
  var remaining_string string = ""
    
  current_position := 0

  for scanner.Scan() {    
    
    if current_position >= log_position {
      scanned_text := scanner.Text()
      scanned_text = regex_obj.ReplaceAllString(scanned_text, "")
      scanned_text = regex_obj2.ReplaceAllString(scanned_text, "")
      
      remaining_string += scanned_text + "\n"
    }
    current_position += 1
	}
  
  err = scanner.Err()
  
	return []byte(remaining_string), current_position, err
}
/////////////////////////////////////////////////////////////////////////////////
func Get_string_from_log(log_path string, regex_string string) string {
  
  start_time := time.Now()

  match_result := ""
  wait_for_log_entry := false
  string_log := ""
  elapsed := 0
  
  for wait_for_log_entry == false && elapsed < 10 {
    
    end_time := time.Now()
    elapsed = int(end_time.Sub(start_time)) / 1000000000
    
    body, _ := Get_log(log_path)
    string_log = string(body)    

    regex_obj := regexp.MustCompile(regex_string)
    match_result = regex_obj.FindString(string_log)      

    if match_result != "" {
       wait_for_log_entry = true
    }
  }
  
  fmt.Println(match_result)

  regex_obj := regexp.MustCompile(regex_string)
  
  match_array := regex_obj.FindStringSubmatch(string_log)
    
  if match_result == "" {
    return ""
  } else {    
    return match_array[1]
  }
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Check_if_regex_string_exists_in_log(log_path string, regex_string string) int {
    
  os.Chdir("../webserver/")
  
  body, _ := Get_log(log_path)
  string_log := string(body)
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta(regex_string))
  match_result := regex_obj.FindString(string_log)
  
  if match_result == "" {
    return 0
  } else {    
    return 1
  }

}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Check_if_regex_string_exists_in_string(string_log string, regex_string string) int {
    
  os.Chdir("../webserver/")    
  
  regex_obj := regexp.MustCompile(regexp.QuoteMeta(regex_string))
  match_result := regex_obj.FindString(string_log)
  
  if match_result == "" {
    return 0
  } else {    
    return 1
  }

}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
