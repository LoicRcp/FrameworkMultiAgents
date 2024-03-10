package NetworkService

import (
	"FrameworkMultiAgents/Messages"
	"FrameworkMultiAgents/containerOps"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type SyncCommunication struct {
	receiverID     int
	receiverAdress string
	syncChannel    chan Messages.Message
}

type NetworkService struct {
	MainContainerAddress string
	LocalAddress         string
	requestCounter       int64 // For generating unique correlation IDs
	handlerMutex         sync.Mutex
	responseHandlers     map[int64]chan Messages.Message // Map to track response handlers
	connPool             map[string]*websocket.Conn
	connPoolMutex        sync.Mutex
	containerOps         containerOps.ContainerOps
	syncChannels         map[int]SyncCommunication
}

func NewNetworkService(mainContainerAddress, localAddress string) *NetworkService {
	ns := &NetworkService{
		MainContainerAddress: mainContainerAddress,
		LocalAddress:         localAddress,
		responseHandlers:     make(map[int64]chan Messages.Message),
		connPool:             make(map[string]*websocket.Conn),
	}
	return ns
}

func (ns *NetworkService) SetContainerOps(ops containerOps.ContainerOps) {
	ns.containerOps = ops
}

func (ns *NetworkService) SendMessage(message Messages.Message, address string) (Messages.Message, error) {
	correlationID := atomic.AddInt64(&ns.requestCounter, 1)
	message.CorrelationID = correlationID

	var responseChan chan Messages.Message
	if message.ExpectResponse {
		responseChan = make(chan Messages.Message)
		ns.addHandler(correlationID, responseChan)
		defer ns.removeHandler(correlationID)
	}

	conn, err := ns.getConnection(address)
	if err != nil {
		return Messages.Message{}, fmt.Errorf("error getting connection: %w", err)
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return Messages.Message{}, fmt.Errorf("error marshaling message: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return Messages.Message{}, fmt.Errorf("WriteMessage error: %w", err)
	}

	if message.ExpectResponse {
		select {
		case response := <-responseChan:
			return response, nil
		case <-time.After(30 * time.Second): // Consider making this timeout configurable
			return Messages.Message{}, fmt.Errorf("timeout waiting for response to message with CorrelationID %d", correlationID)
		}
	}

	// If no response is expected, return immediately
	return Messages.Message{}, nil
}

func (ns *NetworkService) getConnection(address string) (*websocket.Conn, error) {
	ns.connPoolMutex.Lock()
	defer ns.connPoolMutex.Unlock()

	// Use existing connection if available
	if conn, exists := ns.connPool[address]; exists && conn != nil {
		return conn, nil
	}

	// Create a new connection and add it to the pool
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s", address), nil)
	if err != nil {
		return nil, fmt.Errorf("WebSocket Dial Error: %w", err)
	}
	ns.connPool[address] = conn
	go ns.startListening(address, conn)
	return conn, nil
}

func (ns *NetworkService) addHandler(correlationID int64, ch chan Messages.Message) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	ns.responseHandlers[correlationID] = ch
}

func (ns *NetworkService) startListening(address string, conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v", err)
			return
		}
		var message Messages.Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			fmt.Printf("Error unmarshaling message: %v", err)
			return
		}
		ns.processIncomingMessage(message)
	}
}

func (ns *NetworkService) processIncomingMessage(message Messages.Message) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	if ch, exists := ns.responseHandlers[message.CorrelationID]; exists {
		ch <- message
	} else {
		if message.Type == Messages.RegisterContainer {
			var payload Messages.RegisterContainerPayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling RegisterContainerPayload: %v", err)
				return
			}
			id := ns.containerOps.RegisterContainer(payload.Address)
			// send response without expecting a reply
			payload2 := Messages.RegisterContainerAnswerPayload{
				ContainerID: id,
			}
			payloadStr, _ := json.Marshal(payload2)
			response := Messages.Message{
				Type:           Messages.RegisterContainerAnswer,
				Sender:         ns.LocalAddress,
				ContentType:    Messages.RegisterContainerAnswerContent,
				Content:        string(payloadStr),
				CorrelationID:  message.CorrelationID,
				ExpectResponse: false,
			}
			ns.SendMessage(response, message.Sender)
		} else if message.Type == Messages.RegisterAgent {
			var payload Messages.RegisterAgentPayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling RegisterAgentPayload: %v", err)
				return
			}
			id, err := ns.containerOps.RegisterAgent(payload.ContainerID)
			if err != nil {
				fmt.Printf("Error registering agent: %v", err)
				return
			}
			payload2 := Messages.RegisterAgentAnswerPayload{
				ID: id,
			}
			payloadStr, _ := json.Marshal(payload2)
			response := Messages.Message{
				Type:           Messages.RegisterAgentAnswer,
				Sender:         ns.LocalAddress,
				ContentType:    Messages.RegisterAgentAnswerContent,
				Content:        string(payloadStr),
				CorrelationID:  message.CorrelationID,
				ExpectResponse: false,
			}
			ns.SendMessage(response, message.Sender)
		} else if message.Type == Messages.InterAgentAsyncMessage {
			var payload Messages.InterAgentAsyncMessagePayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling InterAgentMessagePayload: %v", err)
				return
			}
			ns.containerOps.PutMessageInMailBox(message, payload.ReceiverID)
		} else if message.Type == Messages.GetAgentAdress {
			var payload Messages.GetAgentAdressPayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling GetAgentAdressPayload: %v", err)
				return
			}
			address, err := ns.containerOps.ResolveAgentAddress(payload.AgentID)
			if err != nil {
				fmt.Printf("Error resolving agent address: %v", err)
				return
			}
			payload2 := Messages.GetAgentAdressAnswerPayload{
				Adress: address,
			}
			payloadStr, _ := json.Marshal(payload2)
			response := Messages.Message{
				Type: Messages.GetAgentAdressAnswer,

				Sender:         ns.LocalAddress,
				ContentType:    Messages.GetAgentAdressAnswerContent,
				Content:        string(payloadStr),
				CorrelationID:  message.CorrelationID,
				ExpectResponse: false,
			}
			ns.SendMessage(response, message.Sender)
		} else if message.Type == Messages.SetSyncCommunication {
			var payload Messages.SetSyncCommunicationPayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling SetSyncCommunicationPayload: %v", err)
				return
			}
			if _, exists := ns.syncChannels[payload.AgentID]; exists {
				payload2 := Messages.SetSyncCommunicationAnswerPayload{
					Success: false,
				}
				payloadStr, _ := json.Marshal(payload2)
				response := Messages.Message{
					Type:           Messages.SetSyncCommunicationAnswer,
					Sender:         ns.LocalAddress,
					ContentType:    Messages.SetSyncCommunicationAnswerContent,
					Content:        string(payloadStr),
					CorrelationID:  message.CorrelationID,
					ExpectResponse: false,
				}
				ns.SendMessage(response, message.Sender)
			} else {
				ns.syncChannels[payload.AgentID] = SyncCommunication{
					receiverID:     payload.AgentID,
					receiverAdress: message.Sender,
					syncChannel:    make(chan Messages.Message),
				}
				payload2 := Messages.SetSyncCommunicationAnswerPayload{
					Success: true,
				}
				payloadStr, _ := json.Marshal(payload2)
				response := Messages.Message{
					Type:           Messages.SetSyncCommunicationAnswer,
					Sender:         ns.LocalAddress,
					ContentType:    Messages.SetSyncCommunicationAnswerContent,
					Content:        string(payloadStr),
					CorrelationID:  message.CorrelationID,
					ExpectResponse: false,
				}
				ns.SendMessage(response, message.Sender)
			}
		} else if message.Type == Messages.InterAgentSyncMessage {
			var payload Messages.InterAgentSyncMessagePayload
			if err := json.Unmarshal([]byte(message.Content), &payload); err != nil {
				fmt.Printf("Error unmarshaling InterAgentSyncMessagePayload: %v", err)
				return
			}
			if _, exists := ns.syncChannels[payload.ReceiverID]; exists {
				ns.syncChannels[payload.ReceiverID].syncChannel <- message
			} else {
				fmt.Printf("No synchronous channel found for agent with ID %d", payload.ReceiverID)
			}

		} else {
			fmt.Printf("No handler found for message with CorrelationID %d", message.CorrelationID)
		}
		return
	}
}

func (ns *NetworkService) removeHandler(correlationID int64) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	if ch, exists := ns.responseHandlers[correlationID]; exists {
		close(ch)
		delete(ns.responseHandlers, correlationID)
	}
}

func (ns *NetworkService) GetSyncChannelWithAgent(SenderID, ReceiverID int) (chan Messages.Message, error) {
	if _, exists := ns.syncChannels[SenderID]; exists {
		return nil, fmt.Errorf("Agent already has a synchronous communication")
	}
	var payload Messages.SetSyncCommunicationPayload
	payload.AgentID = ReceiverID
	payloadStr, _ := json.Marshal(payload)
	message := Messages.Message{
		Type:           Messages.SetSyncCommunication,
		Sender:         ns.LocalAddress,
		ContentType:    Messages.SetSyncCommunicationContent,
		Content:        string(payloadStr),
		CorrelationID:  0,
		ExpectResponse: false,
	}
	receiverAddress, err := ns.containerOps.ResolveAgentAddress(strconv.Itoa(ReceiverID))

	response, err := ns.SendMessage(message, receiverAddress)
	if err != nil {
		return nil, fmt.Errorf("error sending SetSyncCommunication message: %w", err)
	}
	if response.Type != Messages.SetSyncCommunicationAnswer {
		return nil, fmt.Errorf("unexpected response type: %v", response.Type)
	}
	var payload2 Messages.SetSyncCommunicationAnswerPayload
	if err := json.Unmarshal([]byte(response.Content), &payload2); err != nil {
		return nil, fmt.Errorf("error unmarshaling SetSyncCommunicationAnswerPayload: %w", err)
	}

	ns.syncChannels[SenderID] = SyncCommunication{
		receiverID:     ReceiverID,
		receiverAdress: receiverAddress,
		syncChannel:    make(chan Messages.Message),
	}

	go ns.listenToSyncChannel(ns.syncChannels[SenderID].syncChannel, ns.syncChannels[SenderID].receiverAdress)

	return ns.syncChannels[SenderID].syncChannel, nil

}

func (ns *NetworkService) listenToSyncChannel(ch chan Messages.Message, address string) {
	for {
		message := <-ch
		ns.SendMessage(message, address)
	}
}

func (ns *NetworkService) GetSyncChannel(agentID int) (chan Messages.Message, error) {
	if syncComm, exists := ns.syncChannels[agentID]; exists {
		return syncComm.syncChannel, nil
	}
	return nil, fmt.Errorf("no synchronous channel found for agent with ID %d", agentID)
}
