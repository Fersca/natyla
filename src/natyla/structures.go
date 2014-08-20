/*
 * Structures Module
 *
 * Structures that Natyla uses internally
 *
 */
package natyla

import (
	"container/list"
)

//Struct to hold the value and the key in the LRU
type node struct {
	V       map[string]interface{}
	Swap    bool
	Deleted bool
	col     string
	key     string
}

//Struct to hold the value and the key in the LRU
type searchNode struct {
	Id       string
	Document map[string]interface{}
}

//Holds the relation between the diferent collections of element with the corresponding channel to write it
type collectionChannel struct {
	Mapa  map[string]*list.Element
	Canal chan int
}
