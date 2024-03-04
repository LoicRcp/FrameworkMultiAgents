package containerOps

type ContainerOps interface {
	RegisterContainer(containerID, address string)
	RegisterAgent(agentID string) (int, error)
}
