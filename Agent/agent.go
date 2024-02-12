package Agent

type Agent struct {
	ID               int `json:"id"`
	CurrentBehaviour Behaviour
	AgentBehaviours  map[string]Behaviour
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
