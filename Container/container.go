package Container

import (
	"FrameworkMultiAgents/Agent"
	"FrameworkMultiAgents/Messages"
	"FrameworkMultiAgents/NetworkService"
	"FrameworkMultiAgents/YellowPage"
	"encoding/json"
	"log"
)

type container struct {
	id               string
	localAdress      string
	agents           map[string]Agent.Agent
	mainServerAdress string // null if the container is the main container
	mainServerPort   string
	networkService   *NetworkService.NetworkService
}

type MainContainer struct {
	container
	yellowPage YellowPage.YellowPage
}

func NewContainer(mainAddress, localAddress string) *container {
	networkService := NetworkService.NewNetworkService(mainAddress, localAddress)

	// Prepare the message
	payload := Messages.RegisterContainerPayload{Address: localAddress}
	payloadStr, _ := json.Marshal(payload) // handle error properly
	message := Messages.Message{
		Type:        Messages.RegisterContainer,
		Sender:      localAddress,
		ContentType: Messages.RegisterContainerContent,
		Content:     string(payloadStr),
	}

	// Send the message and wait for a response
	response, err := networkService.SendMessage(message, mainAddress)
	if err != nil {
		log.Fatalf("Failed to register container: %v", err)
	}
	var answerPayload Messages.RegisterContainerAnswerPayload
	if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
		log.Fatalf("Failed to parse register container response: %v", err)
	}

	return &container{
		id:               answerPayload.ID, // Use the ID returned from the main container
		localAdress:      localAddress,
		agents:           make(map[string]Agent.Agent),
		mainServerAdress: mainAddress,
		networkService:   networkService,
	}
}

func NewMainContainer(ID, mainAdress string) *MainContainer {
	return &MainContainer{
		container: container{
			id:               "0",
			localAdress:      mainAdress,
			agents:           make(map[string]Agent.Agent),
			mainServerAdress: mainAdress,
		},
		yellowPage: *YellowPage.NewYellowPage(),
	}
}
func (MainContainer *MainContainer) RegisterContainer(containerID, Address string) {
	MainContainer.yellowPage.RegisterContainer(containerID, Address)
}

func (MainContainer *MainContainer) RegisterAgent(containerId string) int {
	return MainContainer.yellowPage.RegisterAgent(containerId)
}

func (container *container) AddAgent(agent Agent.Agent) {

}
