package natyla

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
RETST API TESTS
*/

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
	lastElement := lista.Back()
	if lastElement == nil {
		t.Fatalf("LRU element does not exists")
	}
	if lastElement != element {
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
	lastElement = lista.Back()
	if lastElement != nil {
		t.Fatalf("LRU element exists and it shouldnt")
	}

	//check if it is not in the disk any more
	_, err := readJsonFromDisK("users", "1")
	if err == nil {
		t.Fatalf("the json exists in disk and it shouldnt")
	}

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

	//TODO: check if now it is in the memory and in the LRU
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
