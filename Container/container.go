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

func NewContainer(ID, mainAdress, mainPort, port string) *container {
	return &container{
		id:               ID,
		port:             port,
		agents:           make(map[string]Agent.Agent),
		mainServerAdress: mainAdress,
		mainServerPort:   mainPort,
	}
}

func (container *container) AddAgent(agent Agent.Agent) {
	// ask the yellowPage for an ID

	// add the agent to the map
	container.agents[agent.ID] = agent
}
