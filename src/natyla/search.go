/*
 * Search Module
 *
 * Manage the search engine on Natyle
 *
 */
package natyla

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/*
 * Search the jsons that has the key with the specified value
 */
func search(col, key, value string) ([]byte, error) {

	arr := make([]interface{}, 0)
	cc := collections[col]

	//Search the Map for the value
	for _, v := range cc.Mapa {
		//TODO: This is absolutely inefficient, I'm creating a new array for each iteration. Fix this.
		//Is this possible to have something like java ArrayLists  ?
		nod := v.Value.(node)

		//Only check if field exists in document
		if nodeValue, ok := nod.V[key]; ok {
			//In case field is json.Number conver to string, otherwise check directly.
			if reflect.TypeOf(nodeValue).String() == "json.Number" {
				if value == string(nodeValue.(json.Number)) {
					arr = append(arr, nod.V)
				} 
			} else if nodeValue == value {
				arr = append(arr, nod.V)
			}
		}		
	}

	//Create the Json object
	b, err := json.Marshal(arr)

	return b, err
}
