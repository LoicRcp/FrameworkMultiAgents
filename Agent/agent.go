package Agent

import (
	"FrameworkMultiAgents/Messages"
	"strconv"
)

type Agent struct {
	ID                      int `json:"id"`
	CurrentBehaviour        Behaviour
	AgentBehaviours         map[string]Behaviour
	MailBox                 chan Messages.Message
	SendAsyncMessageToAgent func(message Messages.Message, receiverId int)
	GetSyncChannelWithAgent func(agentId int) chan Messages.Message
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
}

func (agent *Agent) takeMail(message Messages.Message) {
	agent.MailBox <- message
}

func (agent *Agent) sendMail(message Messages.Message, receiverId int) {
	agent.SendAsyncMessageToAgent(message, receiverId)
}

func NewAgent(id string, sendMessageToContainer func(message Messages.Message, receiverId int)) *Agent {
	idInt, _ := strconv.Atoi(id)
	return &Agent{
		ID:                      idInt,
		AgentBehaviours:         make(map[string]Behaviour),
		MailBox:                 make(chan Messages.Message),
		SendAsyncMessageToAgent: sendMessageToContainer,
	}
}
