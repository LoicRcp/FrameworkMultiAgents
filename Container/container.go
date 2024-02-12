package Container

import "FrameworkMultiAgents/Agent"

type container struct {
	id           string
	agents       map[string]Agent.Agent // map
	serverAdress string
}
