package Agent

type Agent struct {
	ID               int `json:"id"`
	CurrentBehaviour Behaviour
	AgentBehaviours  map[string]Behaviour

	// Comment on met les attributs de l'agent ? Transformer ça en interface ?
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

// rajouter le start
