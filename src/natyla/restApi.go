/*
 * API Module
 *
 * Manage the REST API Access to Natyla
 *
 */
package natyla

import (
	"fmt"
	"net/http"
	"strings"
)

func restAPI() {
	//Create the webserver
	http.Handle("/", http.HandlerFunc(processRequest))
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Printf("Natyla ListenAndServe Error", err)
	}
}

/*
 * Process the commands recived from internet
 */
func processRequest(w http.ResponseWriter, req *http.Request){

	//If favicon.ico then return nothing by now.... :TODO
	if req.URL.Path =="/favicon.ico" {
		return 
	}

	//Get the headers map	
	headerMap := w.Header()
	//Add the new headers
	headerMap.Add("System", "Natyla 1.0")
	//PrintInformation
	printRequest(req)

	//get the resources from url	
	comandos := strings.Split(req.URL.Path[1:],"/")

	//check if the request is on the root, in this case return 400 - Bad request
	if comandos[0] == "" {
		w.WriteHeader(400)
		w.Write([]byte("Need to specify the resource. Eg: '/users/1' for GET or '/users/' with content for POST"))
		return
	} 
	
	//Performs action based on the request Method
	switch req.Method {

	case "GET":

		//Serch for the specific field in the collection		
		if comandos[1] == "search" {
			col := comandos[0]
			key := req.FormValue("field")
			value := req.FormValue("equal")
			fmt.Println("Searching for:", col, key, value)
			result, err := search(col, key, value)
			if err != nil {
				fmt.Println(result)
				w.WriteHeader(500)
				return
			}
			w.Write(result)
			return		
		}

		//Get the vale from the cache
		element, err := getElement(comandos[0], comandos[1])

		if element != nil {
			//Write the response to the client
			headerMap.Add("Content-Type", "application/json")
			w.Write([]byte(element))
		} else {
			if err == nil {
				//Return a not-found
				w.WriteHeader(404)
			} else {
				headerMap.Add("Error", "Invalid JSON Disk")
				w.WriteHeader(500)
			}
		}

	case "PUT":
		fallthrough
	case "POST":
		//Create the array to hold the body
		var p []byte = make([]byte, req.ContentLength)

		//Reads the body content
		req.Body.Read(p)

		//Save the element in the cache
		id, err := createElement(comandos[0], "", string(p), true, false)

		if err != nil {
			fmt.Println("Error code:", err.Error())
			if err.Error() == "invalid_id" {
				headerMap.Add("Error", "Invalid ID field")
				w.WriteHeader(400)
			} else {
				fmt.Println(err)
				w.WriteHeader(500)
			}
		} else {
			//headerMap.Add("element_id",strconv.Itoa(id))
			headerMap.Add("location", comandos[0]+"/"+id)
			//Response the 201 - created to the client
			w.WriteHeader(201)
		}

	case "DELETE":
		//Get the vale from the cache
		//result := deleteElement(comandos[0],atoi(comandos[1]))
		result := deleteElement(comandos[0], comandos[1])
		if result == false {
			//Return a not-found
			w.WriteHeader(404)
		} else {
			//Return a Ok
			w.WriteHeader(200)
		}

	default:
		if enablePrint {
			fmt.Println("Not Supported: ", req.Method)
		}
		//Method Not Allowed
		w.WriteHeader(405)
	}

}

/*
 * Print the request information
 */
func printRequest(req *http.Request) {

	//Print request information
	if enablePrint {
		fmt.Println("-------------------")
		fmt.Println("Method: ", req.Method)
		fmt.Println("URL: ", req.URL)
		fmt.Println("Headers: ", req.Header)
	}
}
