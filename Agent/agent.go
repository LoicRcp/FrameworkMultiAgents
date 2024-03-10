package Agent

import (
	"FrameworkMultiAgents/Messages"
	"fmt"
	"strconv"
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
	agent.CurrentBehaviour.Perceive()
}

func (agent *Agent) Decide() {
	agent.CurrentBehaviour.Decide()
}

func (agent *Agent) Act() {
	agent.CurrentBehaviour.Act()
}

type Behaviour interface {
	Perceive(params ...interface{})
	Decide(params ...interface{})
	Act(params ...interface{})
	HandleMailboxMessage(message Messages.Message)
	HandleSyncCommunication(message Messages.Message)
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
		MailBox:                 make(chan Messages.Message),
		SendAsyncMessageToAgent: sendMessageToContainer,
		GetSyncChannelWithAgent: GetSyncChannelWithAgent,
		SynchronousChannel:      nil,
	}
}

func (Agent *Agent) StartSyncCommunication(receiverId int) error {
	if Agent.SynchronousChannel != nil {
		fmt.Errorf("The agent already has a synchronous communication")
	}
	channel, err := Agent.GetSyncChannelWithAgent(receiverId)
	if err != nil {
		return err
	}
	Agent.SynchronousChannel = channel
	return nil
}

func (Agent *Agent) SendSyncMessage(message Messages.Message) error {
	if Agent.SynchronousChannel == nil {
		fmt.Errorf("The agent does not have a synchronous communication")
	}
	Agent.SynchronousChannel <- message
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
	agent.CurrentBehaviour.HandleSyncCommunication(message)
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
				agent.CurrentBehaviour.HandleMailboxMessage(message)
			}
		case message := <-agent.SynchronousChannel:
			agent.handleSyncCommunication(message)
		default:
			agent.Perceive()
			agent.Decide()
			agent.Act()
		}
	}
}
