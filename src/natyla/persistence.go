/*
 * Persistence Module
 *
 * Manage the Disck access to persistence, delete and write the JSON objects
 *
 */
package natyla

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
 * Save the Json to disk
 */
func saveJsonToDisk(createDir bool, col, id, valor string) {

	if createDir {
		os.Mkdir(config["data_dir"].(string)+"/"+col, 0777)
	}

	err := ioutil.WriteFile(config["data_dir"].(string)+"/"+col+"/"+id+".json", []byte(valor), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

/*
 * Delete the Json from disk
 */
func deleteJsonFromDisk(col, clave string) error {
	return os.Remove(config["data_dir"].(string) + "/" + col + "/" + clave + ".json")
}

/*
 * Read the Json from disk
 */
func readJsonFromDisK(col, clave string) ([]byte, error) {
	fmt.Println("Read from disk: ", col, " - ", clave)
	content, err := ioutil.ReadFile(config["data_dir"].(string) + "/" + col + "/" + clave + ".json")
	if err != nil {
		fmt.Println(err)
	}

	return content, err
}

/*
 * Create the data directory
 */
func createDataDir() {
	//create the data directory, if it already exist, do nothing
	os.Mkdir(config["data_dir"].(string), 0777)
}

/*
 * Read the config file
 */
func readConfig() {

	//read the config file
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Can't found 'config.json' using default parameters")
		config = make(map[string]interface{})
		config["token"] = "adminToken"
		var maxMemdefault json.Number = json.Number("10485760")
		config["memory"] = maxMemdefault
		config["data_dir"] = "data"
		config["api_port"] = "8080"
		config["telnet_port"] = "8081"
	} else {
		config, err = convertJsonToMap(string(content))
	}

	fmt.Println("Using Config:", config)

}

//Hold the html template
var template []byte

/*
 * Read pretty print html from disk
 */
func readPrettyTemplate() []byte {
	if template != nil {
		return template
	} else {
		//get the template from disk
		content, err := ioutil.ReadFile("pretty.html")
		if err != nil {
			//in case of error return a simple template
			template = []byte("<html><body><b>##ELEMENT##</b></body></html>")
		} else {
			template = content
		}
		return content
	}
}
