//
// Persistence Module
//
// Manage the Disck access to persistence, delete and write the JSON objects
//
package natyla

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
 * Save the Json to disk
 */
func saveJSONToDisk(createDir bool, col, id, valor string) {

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
func deleteJSONFromDisk(col, clave string) error {
	return os.Remove(config["data_dir"].(string) + "/" + col + "/" + clave + ".json")
}

/*
 * Read the Json from disk
 */
func readJSONFromDisK(col, clave string) ([]byte, error) {
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
		config["admin_token"] = "adminToken"
		maxMemdefault := json.Number("10485760")
		config["memory"] = maxMemdefault
		config["data_dir"] = "data"
		config["api_port"] = "8080"
		config["telnet_port"] = "8081"
	} else {
		config, _ = convertJSONToMap(string(content))
	}

	fmt.Println("Using Config:", config)

}

/*
 * Read all files so they are cached.
 */
func readAllFromDisk() (nRead uint64) {
	nRead = 0
	//Walk through data directory.
	filepath.Walk(config["data_dir"].(string), func(path string, f os.FileInfo, err error) error {
		//When we encounter .json file run the algorithm
		if filepath.Ext(path) == ".json" {
			//Replace paths so we can work with windows and unix paths
			splitPath := strings.Split(strings.Replace(path, `\`, `/`, -1), `/`)
			//Only work with paths conforming to data_dir/collection/id.json
			if len(splitPath) == 3 {
				col := splitPath[1]
				id := splitPath[2]

				content, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Println(err)
				} else {
					_, err := createElement(col, id, string(content), false, false) // set to not save to disk
					if err != nil {
						fmt.Println("Invalid Disk JSON")
					} else {
						nRead++
					}
				}
			}
		}
		return nil
	})

	return nRead
}

//Hold the html template
var template []byte

/*
 * Read pretty print html from disk
 */
func readPrettyTemplate() []byte {
	if template != nil {
		return template
	}
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
