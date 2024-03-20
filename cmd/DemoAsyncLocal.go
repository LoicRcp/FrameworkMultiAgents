package main

import (
	"FrameworkMultiAgents/Agent"
	"FrameworkMultiAgents/Container"
	"FrameworkMultiAgents/Messages"
	"fmt"
	"time"
)

type BasicBehaviour1 struct {
	tick int
}

func (b *BasicBehaviour1) Perceive(agent *Agent.Agent, params ...interface{}) {}
func (b *BasicBehaviour1) Decide(agent *Agent.Agent, params ...interface{})   {}
func (b *BasicBehaviour1) Act(agent *Agent.Agent, params ...interface{}) {
	if b.tick == 0 {
		agent.SendMail(Messages.Message{
			Type:           Messages.InterAgentAsyncMessage,
			Sender:         fmt.Sprintf("%d", agent.ID),
			ContentType:    Messages.InterAgentAsyncMessageContent,
			Content:        "Start !",
			ExpectResponse: false,
		}, 2)
	}
	b.tick++
}
func (b *BasicBehaviour1) HandleMailboxMessage(agent *Agent.Agent, msg Messages.Message) {
	if b.tick < 10 {
		agent.SendMail(Messages.Message{
			Type:           Messages.InterAgentAsyncMessage,
			Sender:         fmt.Sprintf("%d", agent.ID),
			ContentType:    Messages.InterAgentAsyncMessageContent,
			Content:        "Pong",
			ExpectResponse: false,
		}, 2)
	} else {
		agent.SendMail(Messages.Message{
			Type:           Messages.InterAgentAsyncMessage,
			Sender:         fmt.Sprintf("%d", agent.ID),
			ContentType:    Messages.InterAgentAsyncMessageContent,
			Content:        "Stop !",
			ExpectResponse: false,
		}, 2)
	}
	fmt.Printf("AGENT %s: Message recu: %s\n", fmt.Sprintf("%d", agent.ID), msg.Content)
	b.tick++
}
func (b *BasicBehaviour1) HandleSyncCommunication(agent *Agent.Agent, msg Messages.Message) {}

type BasicBehaviour2 struct{}

func (b *BasicBehaviour2) Perceive(agent *Agent.Agent, params ...interface{}) {}
func (b *BasicBehaviour2) Decide(agent *Agent.Agent, params ...interface{})   {}
func (b *BasicBehaviour2) Act(agent *Agent.Agent, params ...interface{}) {
}
func (b *BasicBehaviour2) HandleMailboxMessage(agent *Agent.Agent, msg Messages.Message) {
	if msg.Content != "Stop !" {
		agent.SendMail(Messages.Message{
			Type:           Messages.InterAgentAsyncMessage,
			Sender:         fmt.Sprintf("%d", agent.ID),
			ContentType:    Messages.InterAgentAsyncMessageContent,
			Content:        "Ping",
			ExpectResponse: false,
		}, 1)
	}
	fmt.Printf("AGENT %s: Message recu: %s\n", fmt.Sprintf("%d", agent.ID), msg.Content)

}
func (b *BasicBehaviour2) HandleSyncCommunication(agent *Agent.Agent, msg Messages.Message) {}

func main() {

	mainContainer := Container.NewMainContainer("localhost:8080")

	agent1 := mainContainer.AddAgent()
	agent2 := mainContainer.AddAgent()

	agent1Ref := mainContainer.GetAgent(agent1)
	agent2Ref := mainContainer.GetAgent(agent2)

	agent1Ref.RegisterBehaviour("BasicBehaviour", &BasicBehaviour1{})
	agent2Ref.RegisterBehaviour("BasicBehaviour", &BasicBehaviour2{})

	agent1Ref.SetBehaviour("BasicBehaviour")
	agent2Ref.SetBehaviour("BasicBehaviour")

	fmt.Printf("Agent 1: %v\n", agent1)
	fmt.Printf("Agent 2: %v\n", agent2)

	mainContainer.Start()

	for {
		time.Sleep(1 * time.Second)
	}

}