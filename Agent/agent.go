package Agent

import (
	"FrameworkMultiAgents/Messages"
	"fmt"
	"strconv"
	"time"
)

type Agent struct {
	ID                      int `json:"id"`
	CurrentBehaviour        Behaviour
	AgentBehaviours         map[string]Behaviour
	MailBox                 chan Messages.Message
	SendAsyncMessageToAgent func(message Messages.Message, receiverId int, agentId int)
	GetSyncChannelWithAgent func(agentId int) (chan Messages.Message, error)
	SynchronousChannel      chan Messages.Message
}

func (agent *Agent) Perceive() {
	if agent.CurrentBehaviour == nil {
		fmt.Println("No behaviour set for agent")
		return
	}
	agent.CurrentBehaviour.Perceive(agent)
}

func (agent *Agent) Decide() {
	if agent.CurrentBehaviour == nil {
		fmt.Println("No behaviour set for agent")
		return
	}
	agent.CurrentBehaviour.Decide(agent)
}

func (agent *Agent) Act() {
	if agent.CurrentBehaviour == nil {
		fmt.Println("No behaviour set for agent")
		return
	}
	agent.CurrentBehaviour.Act(agent)
}

type Behaviour interface {
	Perceive(agent *Agent, params ...interface{})
	Decide(agent *Agent, params ...interface{})
	Act(agent *Agent, params ...interface{})
	HandleMailboxMessage(agent *Agent, message Messages.Message)
	HandleSyncCommunication(agent *Agent, message Messages.Message)
}

func (agent *Agent) takeMail(message Messages.Message) {
	agent.MailBox <- message
}

func (agent *Agent) SendMail(message Messages.Message, receiverId int) {
	agent.SendAsyncMessageToAgent(message, receiverId, agent.ID)
}

func NewAgent(id string, sendMessageToContainer func(message Messages.Message, receiverId, agentId int), GetSyncChannelWithAgent func(agentId int) (chan Messages.Message, error)) *Agent {
	idInt, _ := strconv.Atoi(id)
	return &Agent{
		ID:                      idInt,
		CurrentBehaviour:        nil,
		AgentBehaviours:         make(map[string]Behaviour),
		MailBox:                 make(chan Messages.Message, 50),
		SendAsyncMessageToAgent: sendMessageToContainer,
		GetSyncChannelWithAgent: GetSyncChannelWithAgent,
		SynchronousChannel:      nil,
	}
}

func (agent *Agent) StartSyncCommunication(receiverId int) error {
	if agent.SynchronousChannel != nil {
		fmt.Errorf("The agent already has a synchronous communication")
	}
	channel, err := agent.GetSyncChannelWithAgent(receiverId)
	if err != nil {
		return err
	}
	agent.SynchronousChannel = channel
	return nil
}

func (agent *Agent) SendSyncMessage(message Messages.Message) error {
	if agent.SynchronousChannel == nil {
		fmt.Errorf("The agent does not have a synchronous communication")
	}
	agent.SynchronousChannel <- message
	return nil
}

// function "giveNewChannel", which is used to give the container a synchronous channel to communicate with another agent.
// The agent will create a new channel, return it to the container.
func (agent *Agent) GiveNewChannel() (chan Messages.Message, error) {
	if agent.SynchronousChannel == nil {
		agent.SynchronousChannel = make(chan Messages.Message)
	}
	return nil, fmt.Errorf("The agent already has a synchronous communication")
}

func (agent *Agent) getChannel() chan Messages.Message {
	return agent.SynchronousChannel
}

func (agent *Agent) StopSynchronousCommunication() {
	close(agent.SynchronousChannel)
	agent.SynchronousChannel = nil
}

func (agent *Agent) handleSyncCommunication(message Messages.Message) {
	agent.CurrentBehaviour.HandleSyncCommunication(agent, message)
}

func (agent *Agent) RegisterBehaviour(name string, behaviour Behaviour) {
	agent.AgentBehaviours[name] = behaviour
}

func (agent *Agent) SetBehaviour(name string) {
	agent.CurrentBehaviour = agent.AgentBehaviours[name]
}

func (agent *Agent) RemoveBehaviour(name string) {
	delete(agent.AgentBehaviours, name)
}

func (agent *Agent) Start() {
	for {
		select {
		case message := <-agent.MailBox:
			if message.Type == Messages.Death {
				// handle death
			} else {
				agent.CurrentBehaviour.HandleMailboxMessage(agent, message)
			}
		case message := <-agent.SynchronousChannel:
			agent.handleSyncCommunication(message)
		default:
			agent.Perceive()
			agent.Decide()
			agent.Act()
		}
		time.Sleep(1 * time.Second)
	}
}
