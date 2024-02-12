package Container

import (
	"FrameworkMultiAgents/Agent"
)

type container struct {
	// id = ipAdress
	id               string
	port             string
	agents           map[string]Agent.Agent
	mainServerAdress string // null if the container is the main container
	mainServerPort   string
}

type MainContainer struct {
	container
	yellowPage YellowPage.YellowPage
}

func NewContainer(ID, mainAdress, mainPort, port string) *container {
	return &container{
		id:               ID,
		port:             port,
		agents:           make(map[string]Agent.Agent),
		mainServerAdress: mainAdress,
		mainServerPort:   mainPort,
	}
}

func NewMainContainer(ID, port string) *MainContainer {
	return &MainContainer{
		container: container{
			id:               ID,
			port:             port,
			agents:           make(map[string]Agent.Agent),
			mainServerAdress: "",
			mainServerPort:   "",
		},
		yellowPage: *YellowPage.NewYellowPage(),
	}
}
func (MainContainer *MainContainer) RegisterContainer(address, port string) {
	MainContainer.yellowPage.RegisterContainer(address, port)
}

func (MainContainer *MainContainer) RegisterAgent(containerId string) int {
	return MainContainer.yellowPage.RegisterAgent(containerId)

}

func (container *container) AddAgent(agent Agent.Agent) {
	// ask the yellowPage for an ID

	// add the agent to the map
	container.agents[agent.ID] = agent
}
