/*
 * Search Module
 *
 * Manage the search engine on Natyle
 *
 */
package natyla

import (
	"encoding/json"
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


func advancedSearch(collection string, query map[string][]string) ([]byte, error) 	{
	arr := make([]interface{}, 0)
	cc := collections[collection]

	//Get each value from collection
	for _, valueNode := range cc.Mapa {
		//Compare fields that are inside query
		docHasQueryField := false
		wasMatched := true
		for toMatchKey, toMatchValues := range query {
			//With fields in the saved value if they exist
			if suspect, ok := valueNode.Value.(node).V[toMatchKey]; ok {
				docHasQueryField = true
				//If value is json.Number compare as string
				if reflect.TypeOf(suspect).String() == "json.Number" {
					//For now compare only the first value from the Match Values array
					if toMatchValues[0] != string(suspect.(json.Number)) {
						wasMatched = false
						break
					}
				} else if suspect != toMatchValues[0] {
					wasMatched = false
					break
				}
			}
		}

		if docHasQueryField && wasMatched {
			arr = append(arr, valueNode.Value.(node).V)
		}
	}

	b, err := json.Marshal(arr)

	return b, err
}