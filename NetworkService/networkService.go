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
}

func NewNetworkService(mainContainerAddress, localAddress string) *NetworkService {
	ns := &NetworkService{
		MainContainerAddress: mainContainerAddress,
		LocalAddress:         localAddress,
		responseHandlers:     make(map[int64]chan Messages.Message),
	}
	return ns
}

// SendMessage dynamically connects to the specified WebSocket server and sends a message.
func (ns *NetworkService) SendMessage(message Messages.Message, address string) (Messages.Message, error) {
	correlationID := atomic.AddInt64(&ns.requestCounter, 1)
	message.CorrelationID = correlationID

	responseChan := make(chan Messages.Message)
	ns.addHandler(correlationID, responseChan)
	defer ns.removeHandler(correlationID)

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s", address), nil)
	if err != nil {
		return Messages.Message{}, fmt.Errorf("WebSocket Dial Error: %w", err)
	}
	defer conn.Close()

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
	case <-time.After(30 * time.Second):
		return Messages.Message{}, fmt.Errorf("timeout waiting for response to message with CorrelationID %d", correlationID)
	}
}

func (ns *NetworkService) addHandler(correlationID int64, ch chan Messages.Message) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	ns.responseHandlers[correlationID] = ch
}

func (ns *NetworkService) removeHandler(correlationID int64) {
	ns.handlerMutex.Lock()
	defer ns.handlerMutex.Unlock()
	if ch, exists := ns.responseHandlers[correlationID]; exists {
		close(ch)
		delete(ns.responseHandlers, correlationID)
	}
}
