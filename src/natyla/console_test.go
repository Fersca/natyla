package natyla

import (
    "testing"
    "net"
    //"bufio"
    //"fmt"
    "time"
)

type DummyConn struct {
	Command int
}
func (d DummyConn) Read(b []byte) (n int, err error) {

	if d.Command == 0 {
		b[0] = 'h'
		b[1] = 'e'
		b[2] = 'l'
		b[3] = 'p'
		d.Command = 1
	    return 4,nil
    }

	/*    
	if d.Command == 1 {
		b[0] = 'e'
		b[1] = 'x'
		b[2] = 'i'
		b[3] = 't'
		d.Command = 2
	    return 5,nil
    }
    */
    
    return 4,nil;
}

func (d DummyConn) Write(b []byte) (n int, err error) {
    return 50,nil
}
func (d DummyConn) Close() error {
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

func Test_first_connection_to_telnet_console(t *testing.T) {

 	conn := DummyConn{}
 	
	handleTCPConnection(conn)
	
}

