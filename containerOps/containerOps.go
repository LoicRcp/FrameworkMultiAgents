package containerOps

import "FrameworkMultiAgents/Messages"

type ContainerOps interface {
	RegisterContainer(containerID, address string)
	RegisterAgent(agentID string) (int, error)
	PutMessageInMailBox(message Messages.Message, receiverID int)
}
