package natyla

import (
    "testing"
    "io/ioutil"
)

//Create a "user" resource and after that search for a filed, check if this resource is retourned
func Test_search_a_resource_based_on_a_field(t *testing.T) {

	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "2")
	deleteJsonFromDisk("users", "3")
	
	//define a json content
	content1:= "{\"country\":\"Argentina\",\"id\":2,\"name\":\"Natalia\"}"
	content2:= "{\"country\":\"Argentina\",\"id\":3,\"name\":\"Agustina\"}"

	//create the resource
	post("/users", content1)
	post("/users", content2)
	
	//search for a resource with equal name
	response := get("/search?resource=users&field=name&value=Natalia")
		
	//Check the array with only one resource
	checkContent(t,response,"["+content1+"]")
	
	//search for a resource that not exists
	response2 := get("/search?resource=users&field=name&value=Adriana")
		
	//Check the array with any resource
	checkContent(t,response2,"[]")

	//search for a resource with equal name
	response3 := get("/search?resource=users&field=country&value=Argentina")
		
	//Check the array with two resources
	body, _ := ioutil.ReadAll(response3.Body)
	if string(body) != "["+content1+","+content2+"]" {
		//Check the array with two resources in the oder order (travis fails without it)	
		if string(body) != "["+content2+","+content1+"]" {
			t.Fatalf("Non-expected content %s, expected %s or %s", string(body), "["+content1+","+content2+"]", "["+content2+","+content1+"]")
		}
	}
			
	//delete the content from disk if it exists from previous tests
	deleteJsonFromDisk("users", "2")
	deleteJsonFromDisk("users", "3")	
}

