// client.go

package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Zavr22/world-of-wisdom/internal/pkg/config"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/connection"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/pow"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/protocol"
)

func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	defer conn.Close()

	fmt.Printf("connected to %s\n", address)

	reader := connection.BufioReaderAdapter(bufio.NewReader(conn))
	writer := connection.NetConnAdapter(conn)

	for {
		message, err := HandleActions(ctx, reader, writer)
		if err != nil {
			return fmt.Errorf("failed to handle actions: %w", err)
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func HandleActions(ctx context.Context, reader connection.ConnReader, writer connection.ConnWriter) (string, error) {
	// request challenge
	if err := connection.SendMsg(protocol.Message{Header: protocol.RequestChallenge}, writer); err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// parse response
	msg, err := connection.ReadAndParseMsg(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read and parse msg: %w", err)
	}
	var hashcash pow.Hashcash
	if err := json.Unmarshal([]byte(msg.Payload), &hashcash); err != nil {
		return "", fmt.Errorf("failed to unmarshal hashcash: %w", err)
	}
	fmt.Println("got hashcash:", hashcash)

	// calc hashcash
	conf := ctx.Value("config").(*config.Config)
	hashcash, err = hashcash.CalculateHashcash(conf.HashcashMaxIterations)
	if err != nil {
		return "", fmt.Errorf("failed to compute hashcash: %w", err)
	}
	fmt.Println("hashcash computed:", hashcash)

	// marshall hashcash
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("failed to marshal hashcash: %w", err)
	}

	// send challenge back to server
	if err := connection.SendMsg(protocol.Message{Header: protocol.RequestResource, Payload: string(byteData)}, writer); err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	fmt.Println("challenge sent to server")

	// get results
	msg, err = connection.ReadAndParseMsg(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read and parse msg: %w", err)
	}
	return msg.Payload, nil
}
