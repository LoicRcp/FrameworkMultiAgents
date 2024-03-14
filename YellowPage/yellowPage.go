package YellowPage

import (
	"strconv"
	"sync"
)

type YellowPage struct {
	AgentRegistry     map[string]string
	ContainerRegistry map[string]string

	// eviter un maximum les mutex, donc utiliser un/des channels dans le networkService avec un select pour éviter la concurrence
	mutex sync.Mutex

	// enlever la queue et utiliser un int64 directementZ
	maxIDAgent     uint64
	maxIDContainer uint64
}

func NewYellowPage() *YellowPage {
	return &YellowPage{
		AgentRegistry:     make(map[string]string),
		ContainerRegistry: make(map[string]string),
		mutex:             sync.Mutex{},
		maxIDAgent:        0,
		maxIDContainer:    0,
	}
}

func (yellowPage *YellowPage) getAvailableID() uint64 {
	yellowPage.mutex.Lock()
	defer yellowPage.mutex.Unlock()
	yellowPage.maxIDAgent++
	return yellowPage.maxIDAgent
}

// adress = ip:port
func (yellowPage *YellowPage) RegisterContainer(adress string) string {
	yellowPage.mutex.Lock()
	defer yellowPage.mutex.Unlock()
	yellowPage.maxIDContainer++
	id := strconv.FormatUint(yellowPage.maxIDContainer, 10)
	yellowPage.ContainerRegistry[id] = adress
	return id
}
func (yellowPage *YellowPage) RegisterAgent(containerID string) string {
	id := strconv.FormatUint(yellowPage.getAvailableID(), 10)
	yellowPage.AgentRegistry[id] = containerID
	return id
}

func (yellowPage *YellowPage) ResolveAgentAddress(agentID string) (string, error) {
	containerID, ok := yellowPage.AgentRegistry[agentID]
	if !ok {
		return "", nil
	}
	return containerID, nil
	adress, ok := yellowPage.ContainerRegistry[containerID]
	if !ok {
		return "", nil
	}
	return adress, nil
}
