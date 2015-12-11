package natyla

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type DummyConn struct {
	T *testing.T
}

func (d DummyConn) Write(b []byte) (n int, err error) {

	//test the "help" command and check if the repsonse is correct
	checkCommand(1, b, showHelp(), d.T)

	//test if the response was "unknown command" for an unknown command
	checkCommand(2, b, "Unknown Command\n", d.T)

	//test if the response was "Element Created" for a post command
	checkCommand(3, b, "Element Created: 1\n", d.T)

	//test the "post" command and check if the repsonse is correct with the content
	checkCommand(4, b, "{\"id\":1,\"name\":\"Grande\"}\n", d.T)

	//test the "post" command and check if the repsonse is correct with the content
	checkCommand(5, b, "1\n", d.T)

	//test the "search" command and check if the repsonse is correct with the content
	checkCommand(6, b, "[{\"id\":1,\"name\":\"Grande\"}]\n", d.T)

	//test the "delete" command response
	checkCommand(7, b, "Key: 1 from: casa deleted\n", d.T)

	//test the "memory" command response
	if commandNumber == 8 {
		response := string(b)
		if response[:4] != "Uses" {
			d.T.Fatalf("Response: '%s' different from expected: '%s'", response[:4], "Uses")
		}
	}

	//test the get and delete to unknown keys
	checkCommand(9, b, "Key not found\n", d.T)
	checkCommand(10, b, "Key not found\n", d.T)

	return 50, nil
}
func (d DummyConn) Close() error {
	//check if the exit was in the specific command
	if commandNumber != 11 {
		d.T.Fatalf("Close connection in an invalid command: %n %s", commandNumber, "expected: 6")
	}
	return nil
}

func (d DummyConn) LocalAddr() net.Addr {
	return nil
}
func (d DummyConn) RemoteAddr() net.Addr {
	return nil
}
func (d DummyConn) SetDeadline(t time.Time) error {
	return nil
}
func (d DummyConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (d DummyConn) SetWriteDeadline(t time.Time) error {
	return nil
}

//variable to hold the sequence of the commands
var commandNumber int

func (d DummyConn) Read(b []byte) (n int, err error) {

	//execute the command in the following sequense....

	//increase the sequence
	commandNumber++
	fmt.Println("d.Command: ", commandNumber)

	//send the "help" command to ses the help screen
	if commandNumber == 1 {
		return sendCommand("help", b), nil
	}

	//test an unknown command
	if commandNumber == 2 {
		return sendCommand("pipi", b), nil
	}

	//test the "post" command
	if commandNumber == 3 {

		//define a json content
		content := "{\"id\":1,\"name\":\"Grande\"}"

		//delete the previous disk content
		deleteJSONFromDisk("casa", "1")

		command := "post casa " + content

		return sendCommand(command, b), nil
	}

	//test the "get" command
	if commandNumber == 4 {
		return sendCommand("get casa 1", b), nil
	}

	//test the "elements" command
	if commandNumber == 5 {
		//sends the elements command to the server
		return sendCommand("elements casa", b), nil
	}

	//test the "get" command
	if commandNumber == 6 {
		return sendCommand("search casa name Grande", b), nil
	}

	//test the "get" command
	if commandNumber == 7 {
		return sendCommand("delete casa 1", b), nil
	}

	//test the "get" command
	if commandNumber == 8 {
		return sendCommand("memory", b), nil
	}

	//try to get and delete unknown keys
	if commandNumber == 9 {
		return sendCommand("get puf 1", b), nil
	}

	//try to get and delete unknown keys
	if commandNumber == 10 {
		return sendCommand("delete purr 1", b), nil
	}

	//exit the telnet
	return sendCommand("exit", b), nil

}

//Check the command number, with the apropiated response content
func checkCommand(number int, b []byte, content string, t *testing.T) {

	//test if the response was "Element Created"
	if commandNumber == number {
		response := string(b)
		if response != content {
			t.Fatalf("Response: '%s' different from expected: '%s'", response, content)
		}
	}
}

//send the command to the telnet service
func sendCommand(content string, b []byte) int {
	byts := []byte(content)
	for pos := 0; pos < len(byts); pos++ {
		b[pos] = byts[pos]
	}
	return len(byts) + 2
}

//Connects agains the telnet service and send the command
func Test_first_connection_to_telnet_console(t *testing.T) {

	//Create the dummy connection
	conn := DummyConn{t}

	//process the connection
	handleTCPConnection(conn)

}
