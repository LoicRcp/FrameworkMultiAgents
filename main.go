package main

import (
	"FrameworkMultiAgents/Agent"
	"FrameworkMultiAgents/Container"
	"FrameworkMultiAgents/Messages"
	"encoding/json"
	"flag"
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
		err := agent.StartSyncCommunication(1)
		if err != nil {
			fmt.Println(err)
			return
		}

		payload := Messages.InterAgentSyncMessagePayload{
			ReceiverID: 1,
			Content:    "Let's start !",
		}
		payloadStr, _ := json.Marshal(payload)

		agent.SendSyncMessage(Messages.Message{
			Type:        Messages.InterAgentSyncMessage,
			Sender:      fmt.Sprintf("%d", agent.ID),
			ContentType: Messages.InterAgentSyncMessageContent,
			Content:     string(payloadStr),
		})
	}
	fmt.Printf("AGENT %d: Tick %d\n", agent.ID, b.tick)
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
		}, 1)
	} else {
		agent.SendMail(Messages.Message{
			Type:           Messages.InterAgentAsyncMessage,
			Sender:         fmt.Sprintf("%d", agent.ID),
			ContentType:    Messages.InterAgentAsyncMessageContent,
			Content:        "Stop !",
			ExpectResponse: false,
		}, 1)
	}
	fmt.Printf("AGENT %s: Message recu: %s\n", fmt.Sprintf("%d", agent.ID), msg.Content)
	b.tick++
}
func (b *BasicBehaviour1) HandleSyncCommunication(agent *Agent.Agent, msg Messages.Message) {
	if msg.Content == "Ping" && b.tick < 10 {
		agent.SendSyncMessage(Messages.Message{
			Type:        Messages.InterAgentSyncMessage,
			Sender:      fmt.Sprintf("%d", agent.ID),
			ContentType: Messages.InterAgentSyncMessageContent,
			Content:     "Pong",
		})
	} else if b.tick >= 10 {
		agent.SendSyncMessage(Messages.Message{
			Type:        Messages.InterAgentSyncMessage,
			Sender:      fmt.Sprintf("%d", agent.ID),
			ContentType: Messages.InterAgentSyncMessageContent,
			Content:     "Stop !",
		})
	} else {
		agent.StopSynchronousCommunication()
	}
	fmt.Printf("AGENT %d: Message synchrone recu : %s\n", agent.ID, msg.Content)
	b.tick++
}

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
		}, 2)
	}
	fmt.Printf("AGENT %s: Message recu: %s\n", fmt.Sprintf("%d", agent.ID), msg.Content)

}
func (b *BasicBehaviour2) HandleSyncCommunication(agent *Agent.Agent, msg Messages.Message) {
	var payload Messages.InterAgentSyncMessagePayload
	err := json.Unmarshal([]byte(msg.Content), &payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	if payload.Content == "Pong" || payload.Content == "Let's start !" {
		agent.SendSyncMessage(Messages.Message{
			Type:        Messages.InterAgentSyncMessage,
			Sender:      fmt.Sprintf("%d", agent.ID),
			ContentType: Messages.InterAgentSyncMessageContent,
			Content:     "Ping",
		})
	} else if msg.Content == "Stop !" {
		agent.StopSynchronousCommunication()
	}
	fmt.Printf("AGENT %d: Message synchrone recu : %s\n", agent.ID, msg.Content)
}

func main() {

	isMainContainer := flag.Bool("main", false, "Set to true if this process should be the main container")
	port := flag.String("port", "8080", "Set the port number for this container")

	flag.Parse()

	if *isMainContainer {
		fmt.Printf("Starting MainContainer on port %s...\n", *port)
		mainContainer := Container.NewMainContainer("localhost:" + *port)
		agent1 := mainContainer.AddAgent()
		agent1Ref := mainContainer.GetAgent(agent1)
		agent1Ref.RegisterBehaviour("BasicBehaviour", &BasicBehaviour2{})
		agent1Ref.SetBehaviour("BasicBehaviour")
		fmt.Printf("Agent 1: %v\n", agent1)
		time.Sleep(5 * time.Second)
		mainContainer.Start()
		for {
			time.Sleep(1 * time.Second)
		}
	} else {
		//time.Sleep(10 * time.Second)
		fmt.Printf("Starting Container on port %s...\n", *port)
		container := Container.NewContainer("localhost:8080", "localhost:8081")
		agent2 := container.AddAgent()
		agent2Ref := container.GetAgent(agent2)
		agent2Ref.RegisterBehaviour("BasicBehaviour", &BasicBehaviour1{})
		agent2Ref.SetBehaviour("BasicBehaviour")
		fmt.Printf("Agent 2: %v\n", agent2)
		container.Start()
		for {
			time.Sleep(1 * time.Second)
		}
	}

}
