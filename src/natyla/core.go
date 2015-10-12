/*
 * Natyla - FullStack API/Cache/Store
 *
 * 2014 - Fernando Scasserra - twitter: @fersca.
 *
 * Natyla is a persistance cache system written in golang that performs in constant time.
 * It keeps a MAP to store the object internally, and a Double Linked list to purge the LRU elements.
 *
 * LRU updates are done in backgrounds gorutines.
 * LRU and MAP modifications are performed through channels in order to keep them synchronized.
 * Bytes stored are counted in order to limit the amount of memory used by the application.
 *
 * Core Module
 * Manage the internal Memory Access, LRU, concurrency and Swapping
 */

package natyla

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

//Create the list to support the LRU List
var lista *list.List

//Max byte in memory (Key + Data), today set to 100KB
var maxMemBytes int64
var memBytes int64 = 0

//const pointerLen int = 4+8 //Bytes of pointer in 32bits machines plus int64 for the key of element in hashmemBytes
var cacheNotFound bool = true

//Channes to sync the List, map
var lisChan chan int

//chennel to acces to the collection map
var collectionChan chan int

//Print information
const enablePrint bool = true

//Create the map that stores the list of collectionsge
var collections map[string]collectionChannel
var config map[string]interface{}

/*
 * Init the system variables
 */
func init() {

	//Welcome Message
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Starting Natyla...")
	fmt.Println("Version: 1.02")

	//Set the thread quantity based on the number of CPU's
	coreNum := runtime.NumCPU()

	fmt.Println("Number of cores: ", coreNum)

	//read the config file
	readConfig()

	//create the data directory
	createDataDir()

	//set max memory form config
	maxMemBytes, _ = config["memory"].(json.Number).Int64()
	fmt.Println("Max memory defined as: ", maxMemBytes/1024/1024, " Mbytes")

	runtime.GOMAXPROCS(coreNum)

	//Create a new doble-linked list to act as LRU
	lista = list.New()

	//Create the channels
	lisChan = make(chan int, 1)
	collectionChan = make(chan int, 1)

	collections = make(map[string]collectionChannel)

	//Read collections from disk
	nRead := readAllFromDisk()
	fmt.Println("Read", nRead, "entries from disk")

	fmt.Println("Ready, API Listening on http://localhost:8080, Telnet on port 8081")
	fmt.Println("------------------------------------------------------------------")
}

/*
 * Create the server
 */
func Start() {
	//Start the console
	go console()

	//Start the rest API
	restAPI()

}

/*
 * Convert a Json string to a map
 */
func convertJsonToMap(valor string) (map[string]interface{}, error) {

	//Create the Json element
	d := json.NewDecoder(strings.NewReader(valor))
	d.UseNumber()
	var f interface{}
	err := d.Decode(&f)

	if err != nil {
		return nil, err
	}

	//transform it to a map
	m := f.(map[string]interface{})
	return m, nil

}

/*
 * Create the element in the collection
 */
func createElement(col string, id string, valor string, saveToDisk bool, deleted bool) (string, error) {

	//create the list element
	var elemento *list.Element
	b := []byte(valor)

	if deleted == false {

		//Create the Json element
		d := json.NewDecoder(strings.NewReader(valor))
		d.UseNumber()
		var f interface{}
		err := d.Decode(&f)

		if err != nil {
			return "", err
		}

		//transform it to a map
		m := f.(map[string]interface{})

		//Checks the data tye of the ID field
		switch m["id"].(type) {
		case json.Number:
			//id = strconv.FormatFloat(m["id"].(float64),'f',-1,64)
			id = m["id"].(json.Number).String()
		case string:
			id = m["id"].(string)
		default:
			return "", errors.New("invalid_id")
		}

		//Add the value to the list and get the pointer to the node
		n := node{m, false, false}

		lisChan <- 1
		elemento = lista.PushFront(n)
		<-lisChan

	} else {

		//if not found cache is disabled
		if cacheNotFound == false {
			return id, nil
		}

		fmt.Println("Creating node as deleted: ", col, id)
		//create the node as deleted
		var n node
		n.V = nil
		n.Deleted = true

		elemento = &list.Element{Value: n}

	}

	//get the collection-channel relation
	cc := collections[col]
	var createDir bool = false

	if cc.Mapa == nil {

		fmt.Println("Creating new collection: ", col)
		//Create the new map and the new channel
		var newMapa map[string]*list.Element
		var newMapChann chan int
		newMapa = make(map[string]*list.Element)
		newMapChann = make(chan int, 1)

		newCC := collectionChannel{newMapa, newMapChann}
		newCC.Mapa[id] = elemento

		//The collection doesn't exist, create one
		collectionChan <- 1
		collections[col] = newCC
		<-collectionChan
		createDir = true

	} else {
		fmt.Println("Using collection: ", col)
		//Save the node in the map
		cc.Canal <- 1
		cc.Mapa[id] = elemento
		<-cc.Canal
	}

	//if we are creating a deleted node, do not save it to disk
	if deleted == false {

		//Increase the memory counter in a diffetet gorutinie, save to disk and purge LRU
		go func() {
			//Increments the memory counter (Key + Value in LRU + len of col name, + Key in MAP)
			memBytes += int64(len(b))

			if enablePrint {
				fmt.Println("Inc Bytes: ", memBytes)
			}

			//Save the Json to disk, if it is not already on disk
			if saveToDisk == true {
				saveJsonToDisk(createDir, col, id, valor)
			}

			//Purge de LRU
			purgeLRU()
		}()
	}

	return id, nil
}

/*
 * Get the element from the Map and push the element to the first position of the LRU-List
 */
func getElement(col string, id string) ([]byte, error) {

	cc := collections[col]

	//Get the element from the map
	elemento := cc.Mapa[id]

	//checks if the element exists in the cache
	if elemento == nil {
		fmt.Println("Elemento not in memory, reading disk, ID: ", id)

		//read the disk
		content, er := readJsonFromDisK(col, id)

		//if file doesnt exists cache the not found and return nil
		if er != nil {
			//create the element and set it as deleted
			createElement(col, id, "", false, true) // set as deleted and do not save to disk
		} else {
			//Create the element from the disk content
			_, err := createElement(col, id, string(content), false, false) // set to not save to disk
			if err != nil {
				return nil, errors.New("Invalid Disk JSON")
			}
		}

		//call get element again (recursively)
		return getElement(col, id)
	}

	//If the Not-found is cached, return false directely
	if elemento.Value.(node).Deleted == true {
		fmt.Println("Not-Found cached detected on getting, ID: ", id)
		return nil, nil
	}

	//Move the element to the front of the LRU-List using a gorutine
	go moveFront(elemento)

	//Check if the element is mark as swapped
	if elemento.Value.(node).Swap == true {

		//Read the swapped json from disk
		b, _ := readJsonFromDisK(col, id)

		//TODO: read if there was an error and do something...

		m, err := convertJsonToMap(string(b))

		if err != nil {
			return nil, err
		}

		//save the map in the node, mark it as un-swapped
		var unswappedNode node
		unswappedNode.V = m
		unswappedNode.Swap = false
		elemento.Value = unswappedNode

		//increase de memory counter
		memBytes += int64(len(b))

		//as we have load content from disk, we have to purge LRU
		go purgeLRU()
	}

	//Return the element
	b, err := json.Marshal(elemento.Value.(node).V)
	return b, err

}

/*
 * Get the number of elements
 */
func getElements(col string) ([]byte, error) {
	cc := collections[col]
	b, err := json.Marshal(len(cc.Mapa))

	return b, err
}

/*
 * Purge the LRU List deleting the last element
 */
func purgeLRU() {

	//Checks the memory limit and decrease it if it's necessary
	for memBytes > maxMemBytes {

		//sync this procedure
		lisChan <- 1

		//Print Message
		fmt.Println(memBytes, " - ", maxMemBytes)
		fmt.Println("Max memory reached! swapping", memBytes)
		fmt.Println("LRU Elements: ", lista.Len())

		//Get the last element and remove it. Sync is not needed because nothing
		//happens if the element is moved in the middle of this rutine, at last it will be removed
		var lastElement *list.Element = lista.Back()
		if lastElement == nil {
			fmt.Println("Empty LRU")
			//unsync
			<-lisChan
			return
		}

		//Remove the element from the LRU
		deleteElementFromLRU(lastElement)

		var swappedNode node
		swappedNode.V = nil
		swappedNode.Swap = true
		swappedNode.Deleted = false

		(*lastElement).Value = swappedNode

		/*

			//Save the collection and the key in two variables (to use later to update the map
			col := lastElement.Value.(node).col
			key := lastElement.Value.(node).key

			//Create a new element as "S"wapped node
			var swappedNode node
			swappedNode.V = nil
			swappedNode.Swap = true
			swappedNode.col = col
			swappedNode.key = key
			swappedNode.Deleted = false

			//Replace de MAP content with the new swapped node
			cc := collections[col]
			var mapElement *list.Element = cc.Mapa[key]

			(*mapElement).Value = swappedNode

		*/

		//Print a purge
		if enablePrint {
			fmt.Println("Purge Done: ", memBytes)
		}

		//unsync
		<-lisChan
	}

}

/*
 * Move the element to the front of the LRU, because it was readed or updated
 */
func moveFront(elemento *list.Element) {
	//Move the element
	lisChan <- 1
	lista.MoveToFront(elemento)
	<-lisChan
	if enablePrint {
		fmt.Println("LRU Updated")
	}
}

/*
 * Delete the element from the disk, and if its enable, cache the not-found
 */
func deleteElement(col string, clave string) bool {

	//Get the element collection
	cc := collections[col]

	//Get the element from the map
	elemento := cc.Mapa[clave]

	//checks if the element exists in the cache
	if elemento != nil {

		//if it is marked as deleted, return a not-found directly without checking the disk
		if elemento.Value.(node).Deleted == true {
			fmt.Println("Not-Found cached detected on deleting, ID: ", clave)
			return false
		}

		//the node was not previously deleted....so exists in the disk

		//if not-found cache is enabled, mark the element as deleted
		if cacheNotFound == true {

			//created a new node and asign it to the element
			var deletedNode node
			deletedNode.V = nil
			deletedNode.Deleted = true
			elemento.Value = deletedNode
			fmt.Println("Caching Not-found for, ID: ", clave)

		} else {
			//if it is not enabled, delete the element from the memory
			cc.Canal <- 1
			delete(cc.Mapa, clave)
			<-cc.Canal
		}

		//In both cases, remove the element from the list and from disk in a separated gorutine
		go func() {

			lisChan <- 1
			deleteElementFromLRU(elemento)
			<-lisChan

			deleteJsonFromDisk(col, clave)

			//Print message
			if enablePrint {
				fmt.Println("Delete successfull, ID: ", clave)
			}
		}()

	} else {

		fmt.Println("Delete element not in memory, ID: ", clave)

		//Create a new element with the key in the cache, to save a not-found if it is enable
		createElement(col, clave, "", false, true)

		//Check is the element exist in the disk
		err := deleteJsonFromDisk(col, clave)

		//if exists, direcly remove it and return true
		//if it not exist return false (because it was not found)
		if err == nil {
			return true
		} else {
			return false
		}

	}

	return true

}

/*
 * Delete the element from de LRU and decrement the counters
 */
func deleteElementFromLRU(elemento *list.Element) {

	//Decrement the byte counter, decrease the Key * 2 + Value
	var n node = (*elemento).Value.(node)

	b, _ := json.Marshal(n.V)
	memBytes -= int64(len(b))

	//Delete the element in the LRU List
	lista.Remove(elemento)

	fmt.Println("Dec Bytes: ", len(b))

}
