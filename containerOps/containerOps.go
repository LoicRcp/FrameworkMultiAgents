package containerOps

import "FrameworkMultiAgents/Messages"

type ContainerOps interface {
	RegisterContainer(address string) string
	RegisterAgent(agentID string) string
	PutMessageInMailBox(message Messages.Message, receiverID int)
	ResolveAgentAddress(agentID string) (string, error)
	UpdateAgentSyncChannel(agentID string, channel chan Messages.Message)
}
