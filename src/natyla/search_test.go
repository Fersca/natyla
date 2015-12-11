package natyla

import (
	"io/ioutil"
	"testing"
)

//Create a "user" resource and after that search for a field, check if this resource is returned
func Test_search_a_resource_based_on_a_field(t *testing.T) {

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")

	//define a json content
	content1 := `{"country":"Argentina","id":2,"name":"Natalia"}`
	content2 := `{"country":"Argentina","id":3,"name":"Agustina"}`

	//create the resource
	post("/users", content1)
	post("/users", content2)

	//search for a resource with equal name
	response := get("/users/search?field=name&equal=Natalia")

	//Check the array with only one resource
	checkContent(t, response, "["+content1+"]")

	//search for a resource that not exists
	response2 := get("/users/search?field=name&equal=Adriana")

	//Check the array with any resource
	checkContent(t, response2, "[]")

	//search for a resource with equal name
	response3 := get("/users/search?field=country&equal=Argentina")

	//Check the array with two resources
	body, _ := ioutil.ReadAll(response3.Body)
	if string(body) != "["+content1+","+content2+"]" {
		//Check the array with two resources in the oder order (travis fails without it)
		if string(body) != "["+content2+","+content1+"]" {
			t.Fatalf("Non-expected content %s, expected %s or %s", string(body), "["+content1+","+content2+"]", "["+content2+","+content1+"]")
		}
	}

	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")
}

//Create resources with numeric fields and performs searches using numbers.
func TestSearchAResourceWithNumericField(t *testing.T) {
	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")

	//define a json content
	content1 := `{"age":24,"country":"Argentina","id":2,"name":"Natalia"}`
	content2 := `{"age":27,"country":"Argentina","id":3,"name":"Adriana"}`

	//create the resource
	post("/users", content1)
	post("/users", content2)

	//search for a resource with age 24
	response := get("/users/search?field=age&equal=24")

	//Check the array with only one resource
	checkContent(t, response, "["+content1+"]")

	//search for a resource with age 27
	response = get("/users/search?field=age&equal=27")

	//Check the array with only one resource
	checkContent(t, response, "["+content2+"]")

	//search for a resource with age 15
	response = get("/users/search?field=age&equal=15")

	//Check the array with no resource
	checkContent(t, response, "[]")

	//search for a resource with age as string
	response = get("/users/search?field=age&equal=two")

	//Check the array with no resource
	checkContent(t, response, "[]")

	//cleanup
	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")
}

func TestAdvancedSearch(t *testing.T) {
	//delete the content from disk if it exists from previous tests
	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")

	//define a json content
	content1 := `{"age":24,"country":"Argentina","id":2,"name":"Natalia"}`
	content2 := `{"age":27,"country":"Argentina","id":3,"name":"Adriana"}`

	//create the resource
	post("/users", content1)
	post("/users", content2)

	//search for a resource with age 24 named Natalia
	response := get("/users?age=24;name=Natalia")

	//Check the array with only one resource
	checkContent(t, response, "["+content1+"]")

	//search for a resource with age 24 named Natalia
	response = get("/users?age=27;name=Natalia")

	//Check the array with no resource
	checkContent(t, response, "[]")

	//search for a resource with country Argentina age 27
	response = get("/users?country=Argentina;age=27")

	//Check the array with no resource
	checkContent(t, response, "["+content2+"]")

	deleteJSONFromDisk("users", "2")
	deleteJSONFromDisk("users", "3")

}