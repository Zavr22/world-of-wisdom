package protocol

import (
	"fmt"
	"strings"
)

// Header of TCP-message in protocol, means type of message
type Header string

const (
	Quit              Header = "Quit"              // close connection
	RequestChallenge         = "RequestChallenge"  // request new challenge from server
	ResponseChallenge        = "ResponseChallenge" // message with challenge for client
	RequestResource          = "RequestResource"   // message with solved challenge
	ResponseResource         = "ResponseResource"  // message with useful info is solution is correct, or with error if not
)

// Message - message struct for both server and client
type Message struct {
	Header  Header //type of message
	Payload string //payload, could be json, quote or be empty
}

// Stringify - stringify message to send it by tcp-connection
// divider between header and payload is |
func (m *Message) Stringify() string {
	return fmt.Sprintf("%s|%s", m.Header, m.Payload)
}

// Parse - parses Message from str, checks header and payload
func Parse(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	var msgType Header
	// message has view as 1|payload (payload is optional)
	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	// try to parse header
	msgType = Header(parts[0])
	msg := Message{
		Header: msgType,
	}
	// last part after | is payload
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}
