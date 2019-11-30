package main
 
import (
    //"fmt"
    "encoding/json"
    "os"
    //"io"
    "io/ioutil"
    "reflect"
    "strconv"
)

func recurse_map(sub_map map[string]interface{}, indentation string, new_sub_map map[string]interface{}) {

	//println("[[[MAP]]]")

	for index, elem := range sub_map {
		//print("<" + reflect.TypeOf(elem).String() + ">")

		new_sub_map[index] = elem

		print(indentation + index + " => ")

		if elem == nil {
			continue;
		}

		interface_type := reflect.TypeOf(elem)

		if interface_type.Kind() == reflect.Array || interface_type.Kind() == reflect.Slice {

			array_length := len(elem.([]interface{}))

			println()

			new_sub_map[index] = make([]interface{}, array_length)
			new_sub_map[index] = new_sub_map[index].([]interface{})

			recurse_array(elem.([]interface{}), indentation + " ", new_sub_map[index].([]interface{}))

		} else if interface_type.Kind() == reflect.Map {			
			
			println()

			new_sub_map[index] = make(map[string]interface{})	

			recurse_map(elem.(map[string]interface{}), indentation + " ", new_sub_map[index].(map[string]interface{}))

		} else if interface_type.Kind() == reflect.String {

			new_sub_map[index] = elem.(string)
			println(indentation + elem.(string))

		} else if interface_type.Kind() == reflect.Int {

			println(indentation + string(elem.(int)))

		} else if interface_type.Kind() == reflect.Float64 {

			println(indentation + strconv.FormatFloat(elem.(float64), 'f', 2, 64))

		}

	}
}

/////////////////////////////////////////////////////////////////////////////////////////////

func recurse_array(sub_map_array []interface{}, indentation string, new_sub_map_array []interface{}) {

	//println("[[[ARRAY]]]")

	for index, elem := range sub_map_array {
		//print("<" + reflect.TypeOf(elem).String() + ">")

		new_sub_map_array[index] = elem

		print(indentation + strconv.Itoa(index) + " => ")

		if elem == nil {
			continue;
		}
		
		interface_type := reflect.TypeOf(elem)

		if interface_type.Kind() == reflect.Array || interface_type.Kind() == reflect.Slice {

			array_length := len(elem.([]interface{}))

			println()

			new_sub_map_array[index] = make([]interface{}, array_length)
			new_sub_map_array[index] = new_sub_map_array[index].([]interface{})

			recurse_array(elem.([]interface{}), indentation + " ", new_sub_map_array[index].([]interface{}))

		} else if interface_type.Kind() == reflect.Map {				
				
			new_sub_map_array[index] = make(map[string]interface{})	

			recurse_map(elem.(map[string]interface{}), indentation + " ", new_sub_map_array[index].(map[string]interface{}))

		} else if interface_type.Kind() == reflect.String {

			new_sub_map_array[index] = elem.(string)
			println(indentation + elem.(string))

		} else if interface_type.Kind() == reflect.Int {

			println(indentation + string(elem.(int)))

		} else if interface_type.Kind() == reflect.Float64 {

			println(indentation + strconv.FormatFloat(elem.(float64), 'f', 2, 64))

		}

	}
}

///////////////////////////////////////////////////////////

func main() {

	filename := os.Args[1]

	file_buffer, _ := ioutil.ReadFile(filename)
	
	imported_map := make(map[string]interface{})
	new_map := make(map[string]interface{})

	json.Unmarshal([]byte(file_buffer), &imported_map)

	recurse_map(imported_map, "", new_map)

	println()

	/////////////////

}

