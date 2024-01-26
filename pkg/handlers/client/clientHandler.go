package client

import (
	"bufio"
	"fmt"
	"log"
	chat "oliva-back-chat/internal/gen"
	"os"
	"strings"
)

type ClientHandle struct {
	Stream     chat.Services_ChatServiceClient
	clientName string
}

func (ch *ClientHandle) ClientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)

	}
	ch.clientName = strings.TrimRight(msg, "\r\n")
}
func (ch *ClientHandle) SendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		clientMessage = strings.TrimRight(clientMessage, "\r\n")
		if err != nil {
			log.Printf("Failed to read from console :: %v", err)
			continue
		}

		clientMessageBox := &chat.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.Stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending to server :: %v", err)
		}

	}

}

func (ch *ClientHandle) ReceiveMessage() {

	for {
		resp, err := ch.Stream.Recv()
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		log.Printf("%s : %s", resp.Name, resp.Body)
	}
}
