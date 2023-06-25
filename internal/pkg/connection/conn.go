package connection

import (
	"bufio"
	"fmt"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/protocol"
	"net"
)

// ConnReader defines an interface for reading from a connection.
type ConnReader interface {
	ReadString(delim byte) (string, error)
}

// ConnWriter defines an interface for writing to a connection.
type ConnWriter interface {
	Write(p []byte) (int, error)
}

// ReadAndParseMsg reads a message from a connection and parses it using the protocol package.
func ReadAndParseMsg(reader ConnReader) (*protocol.Message, error) {
	msgStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read msg: %w", err)
	}
	msg, err := protocol.Parse(msgStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse msg: %w", err)
	}
	return msg, nil
}

// SendMsg formats a message using the protocol package and sends it over a connection.
func SendMsg(msg protocol.Message, writer ConnWriter) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := writer.Write([]byte(msgStr))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// BufioReaderAdapter wraps a bufio.Reader in a ConnReader interface.
func BufioReaderAdapter(reader *bufio.Reader) ConnReader {
	return &bufioReader{reader}
}

// bufioReader is a wrapper around bufio.Reader that implements the ConnReader interface.
type bufioReader struct {
	*bufio.Reader
}

// ReadString implements the ReadString method of the ConnReader interface.
func (b *bufioReader) ReadString(delim byte) (string, error) {
	return b.Reader.ReadString(delim)
}

// NetConnAdapter wraps a net.Conn in a ConnWriter interface.
func NetConnAdapter(conn net.Conn) ConnWriter {
	return &netConn{conn}
}

// netConn is a wrapper around net.Conn that implements the ConnWriter interface.
type netConn struct {
	net.Conn
}

// Write implements the Write method of the ConnWriter interface.
func (n *netConn) Write(p []byte) (int, error) {
	return n.Conn.Write(p)
}
