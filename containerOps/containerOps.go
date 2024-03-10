package containerOps

import "FrameworkMultiAgents/Messages"

type ContainerOps interface {
	RegisterContainer(address string) string
	RegisterAgent(agentID string) (int, error)
	PutMessageInMailBox(message Messages.Message, receiverID int)
	ResolveAgentAddress(agentID string) (string, error)
}
