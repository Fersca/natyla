package natyla

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"time"
	"encoding/json"
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
	content:= "{\"id\":1,\"name\":\"Valeria\"}"

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
	if element==nil {
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
	if lastElement==nil {
		t.Fatalf("LRU element does not exists")
	}
	if lastElement!=element {
		t.Fatalf("The element is not the same as the LRU element")
	}
	
	//get the resource
	response := get("/users/1")

	//check if natyla responds with 200
	checkStatus(t, response, 200)
	
	//check the body content
	checkContent(t,response,content)

	//check the disk file
	diskContent, _ := readJsonFromDisK("users","1")
	
	//check if the content is the same on the disk
	if string(diskContent) != content {
		t.Fatalf("Non-expected disk content %s, expected %s", string(diskContent), content)
	}
	
	//Delete the resource
	deleteReq("/users/1")
	
	//check the status code
	checkStatus(t, response, 200)
	
	//sleep for 1 second in order to let the gorutine finish
	sleep()
	
	//check if it was marked as deleted in memory
	delElement := cc.Mapa["1"]
	if delElement==nil {
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
	if lastElement!=nil {
		t.Fatalf("LRU element exists and it shouldnt")
	}
	
	//check if it is not in the disk any more
	_, err := readJsonFromDisK("users", "1")
	if err==nil{
		t.Fatalf("the json exists in disk and it shouldnt")
	}
		
}

////////////////////// Utility Functions //////////////////////

//Sleep for some seconds
func sleep(){
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
