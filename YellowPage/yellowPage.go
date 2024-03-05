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
	maxID uint64
}

func NewYellowPage() *YellowPage {
	return &YellowPage{
		AgentRegistry:     make(map[string]string),
		ContainerRegistry: make(map[string]string),
		mutex:             sync.Mutex{},
		maxID:             0,
	}
}

func (yellowPage *YellowPage) getAvailableID() uint64 {
	yellowPage.mutex.Lock()
	defer yellowPage.mutex.Unlock()
	yellowPage.maxID++
	return yellowPage.maxID
}

// adress = ip:port
func (yellowPage *YellowPage) RegisterContainer(containerID string, adress string) {
	yellowPage.ContainerRegistry[containerID] = adress
}
func (yellowPage *YellowPage) RegisterAgent(containerID string) string {
	id := strconv.FormatUint(yellowPage.getAvailableID(), 10)
	yellowPage.AgentRegistry[id] = containerID
	return id
}
