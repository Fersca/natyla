package natyla

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
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

//test if we can create a simple json resource
func Test_CreateASimpleResource(t *testing.T) {

	//create the request
	response := post("/users", "{\"id\":1,\"name\":\"Fernando\"}")

	//check if natyla responds with 201 created
	checkStatus(t, response, 201)

	//checks is the location header is correct
	checkHeader(t, response, "Location", "users/1")

}

//test if we can create a resource without an ID
func Test_CreateASimpleResourceWithOutId(t *testing.T) {

	//create the request
	response := post("/users", "{\"name\":\"Fernando\"}")

	//check if natyla responds with 400 created
	checkStatus(t, response, 400)

}

////////////////////// Utility Functions ////////////////////// 

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
