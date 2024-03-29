package Container

import (
	"FrameworkMultiAgents/Agent"
	"FrameworkMultiAgents/Messages"
	"FrameworkMultiAgents/NetworkService"
	"FrameworkMultiAgents/YellowPage"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Container struct {
	id                  string
	localAdress         string
	agents              map[string]*Agent.Agent
	mainServerAdress    string // null if the Container is the main Container
	mainServerPort      string
	networkService      *NetworkService.NetworkService
	resolveAgentLocally func(agentID string) (string, error)
}

type MainContainer struct {
	Container
	yellowPage YellowPage.YellowPage
}

func NewContainer(mainAddress, localAddress string) *Container {
	networkService := NetworkService.NewNetworkService(mainAddress, localAddress)

	// Prepare the message
	payload := Messages.RegisterContainerPayload{Address: localAddress}
	payloadStr, _ := json.Marshal(payload) // handle error properly
	message := Messages.Message{
		Type:           Messages.RegisterContainer,
		Sender:         localAddress,
		ContentType:    Messages.RegisterContainerContent,
		Content:        string(payloadStr),
		ExpectResponse: true,
	}
	// revoir les notation (content/payload)
	// voir les contexte en go

	// Send the message and wait for a response
	response, err := networkService.SendMessage(message, mainAddress)
	if err != nil {
		log.Fatalf("Failed to register Container: %v", err)
	}
	var answerPayload Messages.RegisterContainerAnswerPayload
	if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
		log.Fatalf("Failed to parse register Container response: %v", err)
	}

	newContainer := &Container{
		id:                  localAddress,
		localAdress:         localAddress,
		agents:              make(map[string]*Agent.Agent),
		mainServerAdress:    mainAddress,
		networkService:      networkService,
		resolveAgentLocally: nil,
	}
	newContainer.networkService.SetContainerOps(newContainer)
	go newContainer.networkService.Start()
	return newContainer
}

func NewMainContainer(mainAdress string) *MainContainer {
	container := Container{
		id:               mainAdress,
		localAdress:      mainAdress,
		agents:           make(map[string]*Agent.Agent),
		mainServerAdress: "",
		networkService:   NetworkService.NewNetworkService(mainAdress, mainAdress),
	}
	go container.networkService.Start()
	mainContainer := MainContainer{
		Container:  container,
		yellowPage: *YellowPage.NewYellowPage(),
	}
	container.networkService.SetContainerOps(&mainContainer)
	mainContainer.yellowPage.RegisterContainer(mainAdress)
	mainContainer.Container.resolveAgentLocally = mainContainer.ResolveAgentAddress
	return &mainContainer
}

func (Container *Container) RegisterContainer(address string) string {
	// No-op for regular containers
	return ""
}

func (Container *Container) RegisterAgent(containerId string) string {
	// No-op for regular containers
	return ""
}

func (MainContainer *MainContainer) RegisterContainer(Address string) string {
	return MainContainer.yellowPage.RegisterContainer(Address)
}

func (MainContainer *MainContainer) RegisterAgent(containerId string) string {
	return MainContainer.yellowPage.RegisterAgent(containerId)
}

func (MainContainer *MainContainer) AddAgent() string {
	agentID := MainContainer.RegisterAgent(MainContainer.id)
	agent := *Agent.NewAgent(agentID, MainContainer.sendMessageToAnotherAgent, MainContainer.GetSyncChannelWithAgent)
	MainContainer.agents[agentID] = &agent
	return agentID
}

func (Container *Container) AddAgent() string {
	var agentID string

	// Prepare the message
	payload := Messages.RegisterAgentPayload{ContainerID: Container.id}
	payloadStr, _ := json.Marshal(payload) // handle error properly
	message := Messages.Message{
		Type:           Messages.RegisterAgent,
		Sender:         Container.localAdress,
		ContentType:    Messages.RegisterAgentContent,
		Content:        string(payloadStr),
		ExpectResponse: true,
	}

	// Send the message and wait for a response
	response, err := Container.networkService.SendMessage(message, Container.mainServerAdress)
	if err != nil {
		log.Fatalf("Failed to register agent: %v", err)
	}

	// Parse the response
	var answerPayload Messages.RegisterAgentAnswerPayload
	if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
		log.Fatalf("Failed to parse register agent response: %v", err)
	}
	agentID = strconv.Itoa(answerPayload.ID)

	// Create the agent
	agent := *Agent.NewAgent(agentID, Container.sendMessageToAnotherAgent, Container.GetSyncChannelWithAgent)
	Container.agents[agentID] = &agent

	return agentID
}

func (Container *Container) PutMessageInMailBox(message Messages.Message, receiverID int) {
	if _, exists := Container.agents[strconv.Itoa(receiverID)]; exists {
		// send the message to the agent
		Container.agents[strconv.Itoa(receiverID)].MailBox <- message
	}
	return
}

func (MainContainer *MainContainer) ResolveAgentAddress(agentID string) (string, error) {
	return MainContainer.yellowPage.ResolveAgentAddress(agentID)
}

func (Container *Container) ResolveAgentAddress(agentID string) (string, error) {
	if Container.mainServerAdress == "" {
		return Container.resolveAgentLocally(agentID)
	} else {
		// send the message to the main Container
		// Prepare the message
		payload := Messages.GetAgentAdressPayload{AgentID: agentID}
		payloadStr, _ := json.Marshal(payload)
		message := Messages.Message{
			Type:           Messages.GetAgentAdress,
			Sender:         Container.localAdress,
			ContentType:    Messages.GetAgentAdressContent,
			Content:        string(payloadStr),
			ExpectResponse: true,
		}
		// Send the message and wait for a response
		response, err := Container.networkService.SendMessage(message, Container.mainServerAdress)
		if err != nil {
			log.Fatalf("Failed to resolve agent address: %v", err)
		}
		// Parse the response
		var answerPayload Messages.GetAgentAdressAnswerPayload
		if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
			log.Fatalf("Failed to parse resolve agent address response: %v", err)
		}
		return answerPayload.Adress, nil
	}

}

func (Container *Container) sendMessageToAnotherAgent(message Messages.Message, receiverId int, agentID int) {
	// SEND ASYNC MESSAGE
	// function to send message to another agent

	// check if the other agent is in the same Container
	if _, exists := Container.agents[strconv.Itoa(receiverId)]; exists {
		// send the message to the agent
		select {
		case Container.agents[strconv.Itoa(receiverId)].MailBox <- message:
		default:
		}
	} else {

		// Resolve the agent address
		receiverIdStr := strconv.Itoa(receiverId)
		receiverAdress, err := Container.ResolveAgentAddress(receiverIdStr)
		if err != nil {
			log.Fatalf("Failed to resolve agent address: %v", err)
		}
		// Send the message
		_, err = Container.networkService.SendMessage(message, receiverAdress)
		if err != nil {
			log.Fatalf("Failed to send message to agent: %v", err)
		}
	}
}

func (Container *Container) GetSyncChannelWithAgent(sourceAgentID, agentId int) (chan Messages.Message, error) {
	// ask agent to return a newly created channel
	// check if the agent is in the same Container
	if _, exists := Container.agents[strconv.Itoa(agentId)]; exists {
		agent := Container.agents[strconv.Itoa(agentId)]
		return agent.GiveNewChannel()
	} else {

		// Resolve the agent address
		agentIdStr := strconv.Itoa(agentId)
		agentAdress, err := Container.ResolveAgentAddress(agentIdStr)
		if err != nil {
			log.Fatalf("Failed to resolve agent address: %v", err)
		}

		// Prepare the message
		payload := Messages.SetSyncCommunicationPayload{AgentID: agentId}
		payloadStr, _ := json.Marshal(payload) // handle error properly
		message := Messages.Message{
			Type:           Messages.SetSyncCommunication,
			Sender:         Container.localAdress,
			ContentType:    Messages.SetSyncCommunicationContent,
			Content:        string(payloadStr),
			ExpectResponse: true,
		}
		// Send the message and wait for a response
		response, err := Container.networkService.SendMessage(message, agentAdress)
		if err != nil {
			log.Fatalf("Failed to get sync channel with agent: %v", err)
		}
		// Parse the response
		var answerPayload Messages.SetSyncCommunicationAnswerPayload
		if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
			log.Fatalf("Failed to parse get sync channel with agent response: %v", err)
		}
		if !answerPayload.Success {
			return nil, fmt.Errorf("Failed to get sync channel with agent")
		}
		channel, err := Container.networkService.CreateSyncChannel(sourceAgentID, agentAdress)
		go Container.networkService.ListenToSyncChannel(channel, agentAdress)
		return channel, nil

	}

}

func (Container *Container) GetAgent(agentID string) *Agent.Agent {
	return Container.agents[agentID]
}

func (Container *Container) Start() {
	for _, agent := range Container.agents {
		go agent.Start()
	}
}

func (Container *Container) UpdateAgentSyncChannel(agentID string, channel chan Messages.Message) {
	Container.agents[agentID].SynchronousChannel = channel
}
