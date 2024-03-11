package main

import (
	"FrameworkMultiAgents/Container"
	"FrameworkMultiAgents/Messages"
	"fmt"
)

type BasicBehaviour struct{}

func (b *BasicBehaviour) Perceive() {}
func (b *BasicBehaviour) Decide()   {}
func (b *BasicBehaviour) Act()      {}
func (b *BasicBehaviour) HandleMailboxMessage(msg Messages.Message) {
	fmt.Println("Message reçu dans le conteneur principal:", msg.Content)
}
func (b *BasicBehaviour) HandleSyncCommunication(msg Messages.Message) {}

func main() {

	mainContainer := Container.NewMainContainer("localhost:8080")

	mainContainer.AddAgent()

	mainContainer.Start()
}
