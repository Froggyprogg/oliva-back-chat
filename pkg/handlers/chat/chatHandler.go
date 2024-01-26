package chat

import (
	"google.golang.org/grpc"
	"log"
	"math/rand"
	chat "oliva-back-chat/internal/gen"
	"sync"
	"time"
)

type ChatServer struct {
	chat.UnimplementedServicesServer
}

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageQue struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageQueObject = messageQue{}

func Register(gRPCServer *grpc.Server) {
	chat.RegisterServicesServer(gRPCServer, &ChatServer{})
}
func (is *ChatServer) ChatService(csi chat.Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e3)

	go recieveFromStream(csi, clientUniqueCode)

	errch := make(chan error)
	go sendToStream(csi, clientUniqueCode, errch)

	return <-errch
}
func recieveFromStream(csi_ chat.Services_ChatServiceServer, clientUniqueCode_ int) {

	for {
		req, err := csi_.Recv()
		if err != nil {
			log.Printf("Error reciving request from client :: %v", err)
			break

		} else {
			messageQueObject.mu.Lock()
			messageQueObject.MQue = append(messageQueObject.MQue, messageUnit{ClientName: req.Name, MessageBody: req.Body, MessageUniqueCode: rand.Intn(1e8), ClientUniqueCode: clientUniqueCode_})
			messageQueObject.mu.Unlock()
			log.Printf("%v", messageQueObject.MQue[len(messageQueObject.MQue)-1])
		}

	}

}
func sendToStream(csi_ chat.Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {

	for {

		for {
			time.Sleep(500 * time.Millisecond)
			messageQueObject.mu.Lock()
			if len(messageQueObject.MQue) == 0 {
				messageQueObject.mu.Unlock()
				break
			}
			senderUniqueCode := messageQueObject.MQue[0].ClientUniqueCode
			senderName4client := messageQueObject.MQue[0].ClientName
			message4client := messageQueObject.MQue[0].MessageBody
			messageQueObject.mu.Unlock()
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&chat.FromServer{Name: senderName4client, Body: message4client})

				if err != nil {
					errch_ <- err
				}
				messageQueObject.mu.Lock()
				if len(messageQueObject.MQue) >= 2 {
					messageQueObject.MQue = messageQueObject.MQue[1:] // if send success > delete message
				} else {
					messageQueObject.MQue = []messageUnit{}
				}
				messageQueObject.mu.Unlock()

			}

		}

		time.Sleep(1 * time.Second)

	}

}
