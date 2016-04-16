package natyla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"container/list"
	"time"
)

/*
RETST API TESTS
*/

//Test the init of the REST API and the telnet admin
func Test_Start_REST_API_and_Telnet(t *testing.T) {
	//Only for code coverage because it is not necessary
	go Start()

	//Set the default token as empty
	config["admin_token"] = ""

}

//Test if not existing resource return "404 - Not Found"
func Test_NotFoundGET(t *testing.T) {

	//create the request
	response := get("/pipi/1")

	//check if natyla responds with 404
	checkStatus(t, response, 404)

}

//Test if not existing resource return "404 - Not Found"
func Test_Head_not_supported(t *testing.T) {

	//create the request
	response := head("/pipi/1")

	//check if natyla responds with 404
	checkStatus(t, response, 405)

}

//Test if request to "/" return a "400 - bad request"
func Test_BadRequestOnRoot(t *testing.T) {

	//create the request
	response := get("/")

	//check if natyla responds with 400
	checkStatus(t, response, 400)

}

//Test if request to "/favicon.ico" return a "200 - OK with empty body"
func Test_OKonFavicon(t *testing.T) {

	//create the request
	response := get("/favicon.ico")

	//check if natyla responds with 200
	checkStatus(t, response, 200)

}

//test if we can create a resource without an ID
func Test_CreateASimpleResourceWithOutId(t *testing.T) {

	//create the request
	response := post("/users", "{\"name\":\"Fernando\"}")

	//check if natyla responds with 400 created
	checkStatus(t, response, 400)

}

//tests if we can create a resource with an invalid token
func Test_Try_To_Create_With_Invalid_Token(t *testing.T) {

	//Set the new value for the token
	config["admin_token"] = "newTokenExample"

	//create the request
	response := post("/users", "{\"name\":\"Fernando\"}")

	//check if natyla responds with 401 401 Unauthorized
	checkStatus(t, response, 401)

	//Restore the precious value
	config["admin_token"] = ""

}

//test if we can create a simple JSON resource and get it
func Test_CreateASimpleResourceAndGetIt(t *testing.T) {

	//Clean the list
	lista = list.New()

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "1")

	//define a JSON content
	content := "{\"id\":1,\"name\":\"Valeria\"}"

	fmt.Println("membytes: ", memBytes, maxMemBytes)
	//create the resource
	responsePost := post("/users", content)

	//sleep for 1 second in order to let the content be saved to disk
	sleep()

	printList()
	
	//check if natyla responds with 201 created
	checkStatus(t, responsePost, 201)

	//checks is the location header is correct
	checkHeader(t, responsePost, "Location", "users/1")

	//check if the resource is in the memory map
	cc := collections["users"]
	element := cc.Mapa["1"]
	if element == nil {
		t.Fatalf("Memory element does not exists")
	}
	
	//check the different element flags
	if element.Value.(*node).Deleted == true {
		t.Fatalf("The element is marked as deleted and it shouldnt")
	}
	if element.Value.(*node).Swap == true {
		t.Fatalf("The element is marked as swaped and it shouldnt")
	}
	
	//check the element content in memory
	JSON, _ := json.Marshal(element.Value.(*node).V)
	if string(JSON) != content {
		t.Fatalf("Non-expected memory content %s, expected %s", string(JSON), content)
	}

	//check if the element is in the LRU and if its the same as the eleent in the cache
	lisChan <- 1
	firstElement := lista.Front()
	<-lisChan
	if firstElement == nil {
		t.Fatalf("LRU element does not exists")
	}
	if firstElement != element {
		t.Fatalf("The element is not the same as the LRU element")
	}

	//get the resource
	response := get("/users/1")

	//check if natyla responds with 200
	checkStatus(t, response, 200)

	//check the body content
	checkContent(t, response, content)

	//check the disk file
	diskContent, _ := readJSONFromDisK("users", "1")

	//check if the content is the same on the disk
	if string(diskContent) != content {
		t.Fatalf("Non-expected disk content %s, expected %s", string(diskContent), content)
	}

	/*
		// TODO: Do a PUT and check if the content changes in memory, in disk and if the LRU element go to the front
	*/

	//Delete the resource
	deleteReq("/users/1")

	//check the status code
	checkStatus(t, response, 200)

	//sleep for 1 second in order to let the gorutine finish
	sleep()

	//check if it was marked as deleted in memory
	delElement := cc.Mapa["1"]
	if delElement == nil {
		t.Fatalf("Memory element does not exists")
	}
	//check the different element flags
	if delElement.Value.(*node).Deleted == false {
		t.Fatalf("The element is not marked as deleted and it shouldnt")
	}
	//check the nil value
	if delElement.Value.(*node).V != nil {
		t.Fatalf("The element Value is not nit and it should")
	}

	//check if it is not in the LRU any more
	lisChan <- 1
	firstElement = lista.Front()
	<-lisChan
	if firstElement == delElement {
		t.Fatalf("The element is the same as the first, its wrong")
	}

	//check if it is not in the disk any more
	_, err := readJSONFromDisK("users", "1")
	if err == nil {
		t.Fatalf("the JSON exists in disk and it shouldnt")
	}

}

func Test_Try_To_Create_With_Valid_Token(t *testing.T) {

	//Set the new value for the token
	config["admin_token"] = "test"

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "10")

	//define a JSON content
	content := "{\"id\":10,\"name\":\"Gilda\"}"

	//create the resource
	responsePost := post("/users?access_token=test", content)

	//sleep for 1 second in order to let the content be saved to disk
	sleep()

	//check if natyla responds with 201 created
	checkStatus(t, responsePost, 201)

	//Delete the resource without token
	response := deleteReq("/users/10")

	//check the status code, invalid token
	checkStatus(t, response, 401)

	//Delete the resource
	response = deleteReq("/users/10?access_token=test")

	//check the status code, delete OK
	checkStatus(t, response, 200)

	//sleep for 1 second in order to let the content be deleted from disk
	sleep()

	//Restore the precious value
	config["admin_token"] = ""

}

func printList(){
	
	fmt.Println("******************************************************************************")
	fmt.Println("********************   LISTA                  ********************************")
	fmt.Println("******************************************************************************")
	f := lista.Front()
	b := lista.Back()
	fmt.Println("Frente: ", f)
	fmt.Println("Back  : ", b)
	e := lista.Front()
	
	for;e!=nil; {
		n := e.Value.(*node)
		prev := e.Prev()
		next := e.Next()
		fmt.Println("nodo: ",n, "prev: ",prev, ", next: ", next)
		e = e.Next()		
	}		
	fmt.Println("******************************************************************************")
	
}
func Test_the_swap_functionality(t *testing.T) {

	//Clean the list
	lista = list.New()
	
	//printList()
	
	//delete the content used in the test
	deleteJSONFromDisk("sequence", "1")
	deleteJSONFromDisk("sequence", "2")
	deleteJSONFromDisk("sequence", "3")

	sleep()

	//create an element
	post("/sequence", "{\"id\":1,\"name\":\"First\"}")

	//create another element (with put that by now is the same as POST)
	post("/sequence", "{\"id\":2,\"name\":\"Second\"}")

	//printList()

	//check the last element (should be the first)
	lisChan <- 1
	lastElement := lista.Back()
	<-lisChan

	n := lastElement.Value.(*node) 
	if n.V["name"] != "First" {
		fmt.Println("name: ",n.V["name"])
		t.Fatalf("The last element is not the first element addded")
	}
	
	//check the memory and put the max amount to that value, so the next element creation will purge the LRU (the first element)
	fmt.Println("Memory 1: ", memBytes)

	//whait for the swap gorutine to finish
	sleep()
	
	//store the previous value
	tempMemory := memBytes
	memBytes = maxMemBytes

	fmt.Println("Memory 2: ", memBytes)
	//create the third element
	post("/sequence", "{\"id\":3,\"name\":\"Third\"}")

	//whait for the swap gorutine to finish
	sleep()
	
	//printList()
	
	//check the last element in the LRU (it should be the second, no the first)
	lisChan <- 1
	lastElement2 := lista.Back()
	<-lisChan
	
	n = lastElement2.Value.(*node) 
	if n.V["name"] != "Second" {
		fmt.Println("name: ",n.V["name"])
		//printList()
		t.Fatalf("The last element is not the second element added")
	}

	//find the "first" element in the map, it should be marked as swapped and the content should be empty
	cc := collections["sequence"]
	firstElement := cc.Mapa["1"]
	fmt.Println("firstElement: ", firstElement, "address:",&firstElement)
	if firstElement.Value.(*node).Swap == false {
		t.Fatalf("The node should be marked as swapped")
	}

	if firstElement.Value.(*node).V != nil {
		t.Fatalf("The node content should be empty")
	}

	//get the "first" element, and check the value, it should have been taken from disk
	response := get("/sequence/1")

	//Check the array with only one resource
	checkContent(t, response, "{\"id\":1,\"name\":\"First\"}")

	//delete the content used in the test
	deleteJSONFromDisk("sequence", "1")
	deleteJSONFromDisk("sequence", "2")

	//delete the last element and check if its removed from memory then the not found cache is not enabled
	cacheNotFound = false
	deleteElement("sequence", "3")
	deletedElement := cc.Mapa["3"]
	if deletedElement != nil {
		t.Fatalf("The node exists and it shouldnt")
	}

	cacheNotFound = true
	memBytes = tempMemory
	
}

//Get an element that is not in the cache but it is in the disk
func Test_get_element_that_is_only_in_disk(t *testing.T) {

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "4")

	//define a JSON content
	content1 := "{\"id\":4,\"name\":\"Jimena\"}"

	//create the file
	saveJSONToDisk(true, "users", "4", content1)

	//search for a resource with equal name
	response := get("/users/4")

	//Check the array with only one resource
	checkContent(t, response, content1)

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "4")

	//check if now it is in the memory and in the LRU
	cc := collections["users"]
	element := cc.Mapa["4"]

	if element == nil {
		t.Fatalf("Element does not exist in memory")
	}
	lisChan <- 1
	firstElement := lista.Front()
	<-lisChan
	if firstElement == nil {
		t.Fatalf("LRU element does not exists")
	}
	if firstElement != element {
		t.Fatalf("The element is not the same as the LRU element")
	}

	deleteElement("users", "4")
	deletedElement := cc.Mapa["4"]
	if deletedElement.Value.(*node).Deleted != true {
		t.Fatalf("The element has not been marked as deleted")
	}

}

//Get an element that is not in the cache but it is in the disk
func Test_delete_an_element_that_is_not_in_the_memory(t *testing.T) {

	//search for a resource with equal name
	response := deleteReq("/users/5")

	//check if its returns a not found
	checkStatus(t, response, 404)

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "5")

	//define a JSON content
	content1 := "{\"id\":5,\"name\":\"Sabrina\"}"

	//create the file
	saveJSONToDisk(true, "users", "5", content1)

	//search for the resource
	response = deleteReq("/users/5")

	//check if its returns a 404 not found (because the not found is cached but the file is in disk)
	checkStatus(t, response, 404)

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "5")

	//define a JSON content
	content1 = "{\"id\":6,\"name\":\"Alejandra\"}"

	//create the file
	saveJSONToDisk(true, "users", "6", content1)

	//delete the resource that is not in memory but is in the disk
	response = deleteReq("/users/6")

	//check if the code is 200 because it was in the disk
	checkStatus(t, response, 200)

	//check if it is not in the disk any more
	_, err := readJSONFromDisK("users", "6")
	if err == nil {
		t.Fatalf("the JSON exists in disk and it shouldnt")
	}

}

////////////////////// Utility Functions //////////////////////

//Sleep for some seconds
func sleep() {
	time.Sleep(500 * time.Millisecond)
}

//check the body content of a response
func checkContent(t *testing.T, response *httptest.ResponseRecorder, content string) {

	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != content {
		t.Fatalf("Non-expected content %s, expected %s", string(body), content)
	}
}

//Check the status code
func checkStatus(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	if response.Code != expected {
		t.Fatalf("Non-expected status code %v :\n\tbody: %v", expected, response.Code)
	}
}

//Check the specific header value
func checkHeader(t *testing.T, response *httptest.ResponseRecorder, header string, value string) {
	if response.Header().Get(header) != value {
		t.Fatalf("Header: %s, get:%s expected:%s", header, response.Header().Get(header), value)
	}
}

//Simulate a request and returs a response type
func get(url string) *httptest.ResponseRecorder {
	//create the request
	request, _ := http.NewRequest("GET", url, nil)
	response := httptest.NewRecorder()

	//execute the request
	processRequest(response, request)
	return response
}

//Simulate a request and returs a response type
func head(url string) *httptest.ResponseRecorder {
	//create the request
	request, _ := http.NewRequest("HEAD", url, nil)
	response := httptest.NewRecorder()

	//execute the request
	processRequest(response, request)
	return response
}

//create a resource
func post(url string, content string) *httptest.ResponseRecorder {
	//create the post request
	request, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(content)))
	//add the JSON header
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	//execute the request
	processRequest(response, request)
	return response
}

//delete a resource
func deleteReq(url string) *httptest.ResponseRecorder {
	//create the post request
	request, _ := http.NewRequest("DELETE", url, nil)
	response := httptest.NewRecorder()

	//execute the request
	processRequest(response, request)
	return response
}
