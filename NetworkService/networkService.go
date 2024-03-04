package NetworkService

import (
	"FrameworkMultiAgents/Messages"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type NetworkService struct {
	MainContainerAddress string
	LocalAddress         string
	requestCounter       int64 // For generating unique correlation IDs
	handlerMutex         sync.Mutex
	responseHandlers     map[int64]chan Messages.Message // Map to track response handlers
	connPool             map[string]*websocket.Conn
	connPoolMutex        sync.Mutex
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

func (ns *NetworkService) SendMessage(message Messages.Message, address string) (Messages.Message, error) {
	correlationID := atomic.AddInt64(&ns.requestCounter, 1)
	message.CorrelationID = correlationID

	responseChan := make(chan Messages.Message)
	ns.addHandler(correlationID, responseChan)
	defer ns.removeHandler(correlationID)

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

	select {
	case response := <-responseChan:
		return response, nil
	//TODO: remplacer 30 par une constante
	case <-time.After(30 * time.Second):
		return Messages.Message{}, fmt.Errorf("timeout waiting for response to message with CorrelationID %d", correlationID)
	}
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
	return conn, nil
}

func (ns *NetworkService) addHandler(correlationID int64, ch chan Messages.Message) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	ns.responseHandlers[correlationID] = ch
}

func (ns *NetworkService) Start() {
	for address, conn := range ns.connPool {
		go ns.startListening(address, conn)
	}
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
		fmt.Printf("No handler found for message with CorrelationID %d", message.CorrelationID)
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
