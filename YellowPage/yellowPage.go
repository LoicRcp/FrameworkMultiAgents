package YellowPage

import (
	"strconv"
	"sync"
)

type YellowPage struct {
	AgentRegistry     map[string]string
	ContainerRegistry map[string]string
	mutex             sync.Mutex
	maxID             int
	availableIDQueue  []int
}

func NewYellowPage() *YellowPage {
	return &YellowPage{
		AgentRegistry:     make(map[string]string),
		ContainerRegistry: make(map[string]string),
		mutex:             sync.Mutex{},
		maxID:             0,
		availableIDQueue:  make([]int, 0),
	}
}

func (yellowPage *YellowPage) getAvailableID() int {
	yellowPage.mutex.Lock()
	defer yellowPage.mutex.Unlock()
	if len(yellowPage.availableIDQueue) == 0 {
		yellowPage.maxID++
		return yellowPage.maxID
	} else {
		id := yellowPage.availableIDQueue[0]
		yellowPage.availableIDQueue = yellowPage.availableIDQueue[1:]
		return id
	}
}

// adress = ip:port
func (yellowPage *YellowPage) RegisterContainer(containerID string, adress string) {
	yellowPage.ContainerRegistry[containerID] = adress

}

func (yellowPage *YellowPage) RegisterAgent(containerID string) int {
	id := yellowPage.getAvailableID()
	yellowPage.AgentRegistry[strconv.Itoa(id)] = containerID
	return id
}
