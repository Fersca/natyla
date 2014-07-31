/*
 * Console Module
 *
 * Manage the Natyla administration form console
 *
*/

package natyla

import (
	"net"
	"strings"
	"fmt"
	"strconv"
)

/*
 * Start the command console
 */
func console(){

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleTCPConnection(conn)
	}

}

/*
 * Process each HTTP connection
 */
func handleTCPConnection(conn net.Conn){

	fmt.Println("Connection stablished")

	//Create the array to hold the command
	var command []byte = make([]byte,100)

	for {
		//Read from connection waiting for a command
		cant, err := conn.Read(command)
		if err == nil {

			//read the command and create the string
			var commandStr string = string(command)

			//Exit the connection
			if commandStr[0:4] == "exit" {
				fmt.Println("Cerrando Conexion")
				conn.Close()
				return
			}

			//Get the element
			if commandStr[0:3] == "get" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))
				fmt.Println("Id: ",comandos[2]," - ",len(comandos[2]))

				b,err := getElement(comandos[1],comandos[2])

				if b!=nil {
					conn.Write(b)
					conn.Write([]byte("\n"))
				} else {
					if err==nil{
						conn.Write([]byte("Key not found\n"))
					} else {
						fmt.Println("Error: ", err)
					}
				}
				continue
			}

			//Get the total quantity of elements
			if commandStr[0:8] == "elements" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))

				b, err := getElements(comandos[1])
				if err==nil {
					conn.Write(b)
					conn.Write([]byte("\n"))
				} else {
					fmt.Println("Error: ", err)
				}
				continue
			}

			//return the bytes used
			if commandStr[0:6] == "memory" {

				result := "Uses: "+strconv.FormatInt(memBytes,10)+"bytes, "+ strconv.FormatInt((memBytes/(maxMemBytes/100)),10)+"%\n"
				conn.Write([]byte(result))

				continue
			}


			//POST elements
			if commandStr[0:4] == "post" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))	
				fmt.Println("JSON: ",comandos[2]," - ",len(comandos[2]))

				id,err := createElement(comandos[1],"",comandos[3],true,false)

				var result string
				if err!=nil{
					fmt.Println(err)
				} else {
					//result = "Element Created: "+strconv.Itoa(id)+"\n"
					result = "Element Created: "+id+"\n"
					conn.Write([]byte(result))
				}

				continue
			}

			if commandStr[0:6] == "delete" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				//Get the vale from the cache
				//result := deleteElement(comandos[1],atoi(comandos[2]))
				result := deleteElement(comandos[1],comandos[2])

				if result==false {
					//Return a not-found				
					conn.Write([]byte("Key not found"))
				} else {
					//Return a Ok
					response := "Key: "+comandos[2]+" from: "+comandos[1]+" deleted\n"
					conn.Write([]byte(response))
				}

				continue

			}

			if commandStr[0:6] == "search" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				result, err := search(comandos[1],comandos[2],comandos[3])

				if err!=nil {
					fmt.Println(result)
					conn.Write([]byte("Error searching\n"))
				} else {
					conn.Write([]byte(result))
				}
				continue
			}

			//Exit the connection
			if commandStr[0:4] == "help" {
				result := showHelp()
				conn.Write([]byte(result))
				continue
			}

			//Default Message
			fmt.Println("Comando no definido: ", commandStr)
			conn.Write([]byte("Unknown Command\n"))

		} else {
			fmt.Println("Error: ", err)
		}

	}

}

/*
 * Help
 */
func showHelp() string {

	var help string = "\n\n"

	help += "---------------------------------------------------------------\n"
	help += "Natyla 1.0"
	help += "---------------------------------------------------------------\n\n"

	help += "Telnet Available commands:\n\n"

	help += "- 'exit':                                Close the connection.\n"
	help += "- 'get {collection} {key}':              Get the JSON document from the specified collection.\n"
	help += "- 'elements {collection}':               Get the total elemets from the specified collection.\n"
	help += "- 'memory':                              Get the total ammount of memory used.\n"
	help += "- 'post {collection} {key} {json}':      Save a new JSON document in the specified collection.\n"
	help += "- 'delete {collection} {key}':           Delete the JSON document from the specified collection.\n"
	help += "- 'search {collection} {field} {value}': Search in the specified collection for Jsons with fields in the indicated value.\n"
	help += "\n"
	help += "HTTP Available commands (same as above):\n\n"

	help += "POST/PUT --> localhost:8080/{collection}/{key}    body={json}\n"
	help += "DELETE   --> localhost:8080/{collection}/{key} \n"  
	help += "GET      <-- localhost:8080/{collection}/{key} \n"
	help += "GET      <-- localhost:8080/search?col={collection}&field={field}&value={value}\n"
	help += "\n"
	return help

}


