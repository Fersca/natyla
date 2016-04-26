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

/*
 * Init for testing
 */
func restAPI() {
	//Create the webserver
	http.Handle("/", http.HandlerFunc(processRequest))
	err := http.ListenAndServe("0.0.0.0:"+config["api_port"].(string), nil)
	if err != nil {
		fmt.Print("Natyla ListenAndServe Error", err)
		fmt.Println("popi")

	}
}

/*
 * Process the commands recived from internet
 */
func processRequest(w http.ResponseWriter, req *http.Request) {

	//If favicon.ico then return nothing by now.... :TODO
	if req.URL.Path == "/favicon.ico" {
		return
	}

	//Get the headers map
	headerMap := w.Header()
	//Add the new headers
	headerMap.Add("System", "Natyla 1.0")
	//Print Information
	printRequest(req)

	//get the resources from url
	comandos := strings.Split(req.URL.Path[1:], "/")

	//check if the request is on the root, in this case return 400 - Bad request
	if comandos[0] == "" {
		w.WriteHeader(400)
		w.Write([]byte("Need to specify the resource. Eg: '/users/1' for GET or '/users/' with content for POST"))
		return
	}

	//Get user from acces token
	user, scope := authToken(req.FormValue("access_token"))

	if user == "" {
		//If token is invalid return Unauthorized response.
		headerMap.Add("Unauthorized", "You need to have a valid token")
		w.WriteHeader(401)
		return
	}

	//Performs action based on the request Method
	switch req.Method {

	case "GET":

		//Search advanced. /cars?color:red;owner[]=Natalia;owner[]=Adriana
		if len(comandos) == 1 {
			fmt.Println("Search advanced for", comandos, "with", req.URL.Query())
			collection := comandos[0]

			err := req.ParseForm()

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}

			//query is in req.Form after ParseForm is executed
			result, err := advancedSearch(collection, req.Form)

			if err != nil {
				fmt.Println(result)
				w.WriteHeader(500)
				return
			}

			render(result, w, req)
			return
		}

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
			render(result, w, req)
			return
		}

		//Get the value from the cache
		element, err := getElement(comandos[0], comandos[1])

		if element != nil {
			//Write the response to the client
			render(element, w, req)
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

		if scope != readWrite {
			//If token is invalid return Unauthorized response.
			headerMap.Add("Unauthorized", "You need to have a read/write token")
			w.WriteHeader(401)
			return
		}

		//Create the array to hold the body
		p := make([]byte, req.ContentLength)

		//Reads the body content
		req.Body.Read(p)

		//Check if it is a token creation request
		if comandos[0] == "token" {
			token, tokenErr := createToken(string(p))
			if tokenErr != nil {
				fmt.Println(tokenErr)
				headerMap.Add("Error", "Error creating token for "+user)
				w.WriteHeader(500)
			}
			p = token
		}

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
			headerMap.Add("location", comandos[0]+"/"+id)
			//Response the 201 - created to the client
			w.WriteHeader(201)
		}

	case "DELETE":

		if scope != readWrite {
			//If token is invalid return Unauthorized response.
			headerMap.Add("Unauthorized", "You need to have a read/write token")
			w.WriteHeader(401)
			return
		}

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
 * Verify token is needed and authenticate
 */
func authToken(token string) (string, string) {

	if config["admin_token"] == nil || config["admin_token"] == "" || strings.ToLower(token) == strings.ToLower(config["admin_token"].(string)) {
		return "admin", readWrite
	}

	if token == "" {
		return "", ""
	}

	e, _ := getElement("token", token)

	if e == nil {
		return "", ""
	}

	m, _ := convertJSONToMap(string(e))

	//TODO: in the future, return the user name, instead of "admin" when natyla hables multiple users
	return "admin", m["scope"].(string)

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
		fmt.Println("Params: ", req.RequestURI)
		fmt.Println("Headers: ", req.Header)
		fmt.Println("Accept HTML:", acceptHtml(req))

	}
}

/*
 * Render the json output based on the accept header
 */
func render(element []byte, w http.ResponseWriter, req *http.Request) {

	//Get the headers map
	headerMap := w.Header()

	if acceptHtml(req) {
		prettyContent := readPrettyTemplate()
		headerMap.Add("Content-Type", "text/html")
		w.Write([]byte(strings.Replace(string(prettyContent), "##ELEMENT##", string(element), -1)))
	} else {
		//Add the new headers
		headerMap.Add("Content-Type", "application/json")
		w.Write([]byte(element))
	}

}

/*
 * Check if the request accept html as return type
 */
func acceptHtml(req *http.Request) bool {

	if req.Header["Accept"] != nil {
		return contains(strings.Split(req.Header["Accept"][0], ","), "text/html")
	}

	return false
}

/*
 *Check is the slide contains a text
 */
func contains(s []string, e string) bool {

	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
