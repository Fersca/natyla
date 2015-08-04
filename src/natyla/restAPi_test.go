package natyla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	config["token"] = ""

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
	config["token"] = "newTokenExample"

	//create the request
	response := post("/users", "{\"name\":\"Fernando\"}")

	//check if natyla responds with 401 401 Unauthorized
	checkStatus(t, response, 401)

	//Restore the precious value
	config["token"] = ""

}

//test if we can create a simple json resource and get it
func Test_CreateASimpleResourceAndGetIt(t *testing.T) {

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "1")

	//define a json content
	content := "{\"id\":1,\"name\":\"Valeria\"}"

	//create the resource
	responsePost := post("/users", content)

	//sleep for 1 second in order to let the content be saved to disk
	sleep()

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
	if element.Value.(node).Deleted == true {
		t.Fatalf("The element is marked as deleted and it shouldnt")
	}
	if element.Value.(node).Swap == true {
		t.Fatalf("The element is marked as swaped and it shouldnt")
	}
	//check the element content in memory
	json, _ := json.Marshal(element.Value.(node).V)
	if string(json) != content {
		t.Fatalf("Non-expected memory content %s, expected %s", string(json), content)
	}

	//check if the element is in the LRU and if its the same as the eleent in the cache
	firstElement := lista.Front()
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
	diskContent, _ := readJsonFromDisK("users", "1")

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
	if delElement.Value.(node).Deleted == false {
		t.Fatalf("The element is not marked as deleted and it shouldnt")
	}
	//check the nil value
	if delElement.Value.(node).V != nil {
		t.Fatalf("The element Value is not nit and it should")
	}

	//check if it is not in the LRU any more
	firstElement = lista.Front()
	if firstElement == delElement {
		t.Fatalf("The element is the same as the first, its wrong")
	}

	//check if it is not in the disk any more
	_, err := readJsonFromDisK("users", "1")
	if err == nil {
		t.Fatalf("the json exists in disk and it shouldnt")
	}

}

func Test_Try_To_Create_With_Valid_Token(t *testing.T) {

	//Set the new value for the token
	config["token"] = "test"

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "1")

	//define a json content
	content := "{\"id\":1,\"name\":\"Valeria\"}"

	//create the resource
	responsePost := post("/users?access_token=test", content)

	//sleep for 1 second in order to let the content be saved to disk
	sleep()

	//check if natyla responds with 201 created
	checkStatus(t, responsePost, 201)

	//Delete the resource
	deleteReq("/users/1")

	//Restore the precious value
	config["token"] = ""

}

func Test_the_swap_functionality(t *testing.T) {

	//delete the content used in the test
	deleteJsonFromDisk("sequence", "1")
	deleteJsonFromDisk("sequence", "2")
	deleteJsonFromDisk("sequence", "3")

	//create an element
	post("/sequence", "{\"id\":1,\"name\":\"First\"}")

	//create another element (with put that by now is the same as POST)
	post("/sequence", "{\"id\":2,\"name\":\"Second\"}")

	//check the last element (should be the first)
	lastElement := lista.Back()

	//check the memory and put the max amount to that value, so the next element creation will purge the LRU (the first element)
	fmt.Println("Memory: ", memBytes)
	//store the previous value
	tempMemory := maxMemBytes
	maxMemBytes = memBytes + 1

	//create the third element
	post("/sequence", "{\"id\":3,\"name\":\"Third\"}")

	//whait for the swap gorutine to finish
	sleep()

	//check the last element in the LRU (it should be the second, no the first)
	lastElement2 := lista.Back()
	if lastElement == lastElement2 {
		t.Fatalf("The last element is the same and it shouldnt")
	}

	//find the "first" element in the map, it should be marked as swapped and the content should be empty
	cc := collections["sequence"]
	firstElement := cc.Mapa["1"]
	if firstElement.Value.(node).Swap == false {
		t.Fatalf("The node should be marked as swapped")
	}

	if firstElement.Value.(node).V != nil {
		t.Fatalf("The node content should be empty")
	}

	//get the "first" element, and check the value, it should have been taken from disk
	response := get("/sequence/1")

	//Check the array with only one resource
	checkContent(t, response, "{\"id\":1,\"name\":\"First\"}")

	//delete the content used in the test
	deleteJsonFromDisk("sequence", "1")
	deleteJsonFromDisk("sequence", "2")

	//delete the last element and check if its removed from memory then the not found cache is not enabled
	cacheNotFound = false
	deleteElement("sequence", "3")
	deletedElement := cc.Mapa["3"]
	if deletedElement != nil {
		t.Fatalf("The node exists and it shouldnt")
	}

	cacheNotFound = true
	maxMemBytes = tempMemory

}

//Get an element that is not in the cache but it is in the disk
func Test_get_element_that_is_only_in_disk(t *testing.T) {

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "4")

	//define a json content
	content1 := "{\"id\":4,\"name\":\"Jimena\"}"

	//create the file
	saveJsonToDisk(true, "users", "4", content1)

	//search for a resource with equal name
	response := get("/users/4")

	//Check the array with only one resource
	checkContent(t, response, content1)

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "4")

	//check if now it is in the memory and in the LRU
	cc := collections["users"]
	element := cc.Mapa["4"]

	if element == nil {
		t.Fatalf("Element does not exist in memory")
	}

	firstElement := lista.Front()

	if firstElement == nil {
		t.Fatalf("LRU element does not exists")
	}
	if firstElement != element {
		t.Fatalf("The element is not the same as the LRU element")
	}

	deleteElement("users", "4")
	deletedElement := cc.Mapa["4"]
	if deletedElement.Value.(node).Deleted != true {
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
	deleteJsonFromDisk("users", "5")

	//define a json content
	content1 := "{\"id\":5,\"name\":\"Sabrina\"}"

	//create the file
	saveJsonToDisk(true, "users", "5", content1)

	//search for the resource
	response = deleteReq("/users/5")

	//check if its returns a 404 not found (because the not found is cached but the file is in disk)
	checkStatus(t, response, 404)

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "5")

	//define a json content
	content1 = "{\"id\":6,\"name\":\"Alejandra\"}"

	//create the file
	saveJsonToDisk(true, "users", "6", content1)

	//delete the resource that is not in memory but is in the disk
	response = deleteReq("/users/6")

	//check if the code is 200 because it was in the disk
	checkStatus(t, response, 200)

	//check if it is not in the disk any more
	_, err := readJsonFromDisK("users", "6")
	if err == nil {
		t.Fatalf("the json exists in disk and it shouldnt")
	}

}

////////////////////// Utility Functions //////////////////////

//Sleep for some seconds
func sleep() {
	time.Sleep(1 * time.Second)
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
	//add the json header
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
