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
		Type:           Messages.RegisterContainer,
		Sender:         localAddress,
		ContentType:    Messages.RegisterContainerContent,
		Content:        string(payloadStr),
		ExpectResponse: true,
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

	container := &container{
		id:               localAddress,
		localAdress:      localAddress,
		agents:           make(map[string]Agent.Agent),
		mainServerAdress: mainAddress,
		networkService:   networkService,
	}
	container.networkService.SetContainerOps(container)
	return container
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

func (container *container) RegisterContainer(address string) string {
	// No-op for regular containers
	return ""
}

func (container *container) RegisterAgent(containerId string) (int, error) {
	// No-op for regular containers
	return 0, fmt.Errorf("Not a main container")
}

func (MainContainer *MainContainer) RegisterContainer(Address string) string {
	return MainContainer.yellowPage.RegisterContainer(Address)
}

func (MainContainer *MainContainer) RegisterAgent(containerId string) string {
	return MainContainer.yellowPage.RegisterAgent(containerId)
}

func (container *container) AddAgent() {
	var agentID string
	if mainContainer, isMain := interface{}(container).(*MainContainer); isMain {
		agentID = mainContainer.RegisterAgent(container.id)
	} else {

		// Prepare the message
		payload := Messages.RegisterAgentPayload{ContainerID: container.id}
		payloadStr, _ := json.Marshal(payload) // handle error properly
		message := Messages.Message{
			Type:           Messages.RegisterAgent,
			Sender:         container.localAdress,
			ContentType:    Messages.RegisterAgentContent,
			Content:        string(payloadStr),
			ExpectResponse: true,
		}

		// Send the message and wait for a response
		response, err := container.networkService.SendMessage(message, container.mainServerAdress)
		if err != nil {
			log.Fatalf("Failed to register agent: %v", err)
		}

		// Parse the response
		var answerPayload Messages.RegisterAgentPayload
		if err := json.Unmarshal([]byte(response.Content), &answerPayload); err != nil {
			log.Fatalf("Failed to parse register agent response: %v", err)
		}
		agentID = answerPayload.ContainerID
	}
	// Create the agent
	agent := *Agent.NewAgent(agentID, func(message Messages.Message, receiverId int) {
		// SEND ASYNC MESSAGE
		// function to send message to another agent

		// check if the other agent is in the same container
		if _, exists := container.agents[strconv.Itoa(receiverId)]; exists {
			// send the message to the agent
			container.agents[strconv.Itoa(receiverId)].MailBox <- message
		} else {

			// Resolve the agent address
			receiverIdStr := strconv.Itoa(receiverId)
			receiverAdress, err := container.ResolveAgentAddress(receiverIdStr)
			if err != nil {
				log.Fatalf("Failed to resolve agent address: %v", err)
			}

			// send the message to the main container
			// Prepare the message
			payload := Messages.InterAgentAsyncMessagePayload{ReceiverID: receiverId, Content: message.Content}
			payloadStr, _ := json.Marshal(payload) // handle error properly
			message := Messages.Message{
				Type:           Messages.InterAgentAsyncMessage,
				Sender:         agentID,
				ContentType:    Messages.InterAgentAsyncMessageContent,
				Content:        string(payloadStr),
				ExpectResponse: false,
			}
			// Send the message
			_, err = container.networkService.SendMessage(message, receiverAdress)
			if err != nil {
				log.Fatalf("Failed to send message to agent: %v", err)
			}
		}
	})
	container.agents[agentID] = agent
}

func (container *container) PutMessageInMailBox(message Messages.Message, receiverID int) {
	if _, exists := container.agents[strconv.Itoa(receiverID)]; exists {
		// send the message to the agent
		container.agents[strconv.Itoa(receiverID)].MailBox <- message
	}
	return
}

func (container *container) ResolveAgentAddress(agentID string) (string, error) {
	// check if the agent is in the same container
	if mainContainer, isMain := interface{}(container).(*MainContainer); isMain {
		return mainContainer.yellowPage.ResolveAgentAddress(agentID)
	}
	// send the message to the main container
	// Prepare the message
	payload := Messages.GetAgentAdressPayload{AgentID: agentID}
	payloadStr, _ := json.Marshal(payload)
	message := Messages.Message{
		Type:           Messages.GetAgentAdress,
		Sender:         container.localAdress,
		ContentType:    Messages.GetAgentAdressContent,
		Content:        string(payloadStr),
		ExpectResponse: true,
	}
	// Send the message and wait for a response
	response, err := container.networkService.SendMessage(message, container.mainServerAdress)
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
