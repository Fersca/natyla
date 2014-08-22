package natyla

import (
	"net"
	"testing"
	"fmt"
	"time"
)

type DummyConn struct {
	T *testing.T
}

var ready bool = false //the boolean is because an /n is sent at the end, so write is executed twice
func (d DummyConn) Write(b []byte) (n int, err error) {

	//test the "help" command and check if the repsonse is correct
	if commandNumber == 1 {
		response := string(b)
		if response != showHelp() {
			d.T.Fatalf("Response: '%s' different from expected: '%s'",response, showHelp())
		}	
	}

	//test if the response was "unknown command"
	if commandNumber == 2 {
		response := string(b)
		if response != "Unknown Command\n" {
			d.T.Fatalf("Response: '%s' different from expected: '%s'",response, "Unknown Command\n")
		}	
	}

	//test the "elements" command and check if the repsonse is correct
	if commandNumber == 3 && ready==false {
		fmt.Println("Entra una vez")
		response := string(b)
		if response != "1" {
			d.T.Fatalf("Response elements: '%s' different from expected: '%s'",response, "1")
		} else {
			ready = true
			fmt.Println("Pone true")
		}
	}

	return 50, nil
}
func (d DummyConn) Close() error {	
	//check if the exit was in the specific command
	if commandNumber != 4 {
		d.T.Fatalf("Close connection in an invalid command: %n %s",commandNumber, "expected: 3")
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
var commandNumber int = 0

func (d DummyConn) Read(b []byte) (n int, err error) {

	//execute the command in the following sequense....
	
	//increase the sequence
	commandNumber++
	fmt.Println("d.Command: ",commandNumber)
	
	//send the "help" command to ses the help screen
	if commandNumber == 1 { 	
		return sendCommand("help",b), nil	
	}
	
	//test an unknown command
	if commandNumber == 2 {
		return sendCommand("pipi",b), nil	
	}

	//test the "elements" command
	if commandNumber == 3 {
	
		//define a json content
		content := "{\"id\":1,\"name\":\"Grande\"}"

		//delete the previous disk content
		deleteJsonFromDisk("casa", "1")
		
		//create the resource
		responsePost:=post("/casa", content)
	
		//check the response conde
		checkStatus(d.T, responsePost, 201)
 		
		return sendCommand("elements casa",b), nil	
	}

	//exit the telnet 	
	return sendCommand("exit",b), nil	
	
}

func sendCommand(content string, b []byte) int {
	byts := []byte(content)
	for pos:=0; pos<len(byts);pos++ {
		b[pos] = byts[pos]
	} 	
	return len(byts)+2
}

//Connects agains the telnet service and send the command
func Test_first_connection_to_telnet_console(t *testing.T) {

	//Create the dummy connection
	conn := DummyConn{t}

	//process the connection
	handleTCPConnection(conn)

}
