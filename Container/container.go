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

func (container *container) RegisterContainer(containerID, address string) {
	// No-op for regular containers
	return
}

func (container *container) RegisterAgent(containerId string) (int, error) {
	// No-op for regular containers
	return 0, fmt.Errorf("Not a main container")
}

func (MainContainer *MainContainer) RegisterContainer(containerID, Address string) {
	MainContainer.yellowPage.RegisterContainer(containerID, Address)
}

func (MainContainer *MainContainer) RegisterAgent(containerId string) int {
	return MainContainer.yellowPage.RegisterAgent(containerId)
}

func (container *container) AddAgent() {
	var agentID int
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
		agentID, _ = strconv.Atoi(answerPayload.ContainerID)
	}
	// Create the agent
	agent := *Agent.NewAgent(agentID, func(message Messages.Message, receiverId int) {
		// function to send message to another agent

		// check if the other agent is in the same container
		if _, exists := container.agents[strconv.Itoa(receiverId)]; exists {
			// send the message to the agent
			container.agents[strconv.Itoa(receiverId)].MailBox <- message
		} else {
			// send the message to the main container
			// Prepare the message
			payload := Messages.InterAgentAsyncMessagePayload{ReceiverID: receiverId, Content: message.Content}
			payloadStr, _ := json.Marshal(payload) // handle error properly
			message := Messages.Message{
				Type:           Messages.InterAgentAsyncMessage,
				Sender:         strconv.Itoa(agentID),
				ContentType:    Messages.InterAgentAsyncMessageContent,
				Content:        string(payloadStr),
				ExpectResponse: false,
			}
			// Send the message
			_, err := container.networkService.SendMessage(message, container.mainServerAdress)
			if err != nil {
				log.Fatalf("Failed to send message to agent: %v", err)
			}
		}
	})
	container.agents[strconv.Itoa(agentID)] = agent
}

func (container *container) PutMessageInMailBox(message Messages.Message, receiverID int) {
	if _, exists := container.agents[strconv.Itoa(receiverID)]; exists {
		// send the message to the agent
		container.agents[strconv.Itoa(receiverID)].MailBox <- message
	}
	return
}
