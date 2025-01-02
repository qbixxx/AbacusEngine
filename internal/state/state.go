package state

type ProgramState int

const (
	Edit ProgramState = iota
	Debug
	Run
)

func (ps ProgramState) String() string {
	switch ps {
	case Edit:
		return "Edit"
	case Debug:
		return "Debug"
	case Run:
		return "Run"
	default:
		return "Unknown"
	}
}

type StateManager struct {
	currentState  ProgramState
	OnStateChange func(state string)
}

func NewStateManager(onStateChange func(state string)) *StateManager {
	sm := StateManager{
		currentState:  Edit,
		OnStateChange: onStateChange,
	}
	sm.SetState(Edit)
	return &sm
}

func (sm *StateManager) SetState(newState ProgramState) {
	sm.currentState = newState
	if sm.OnStateChange != nil {
		sm.OnStateChange(newState.String())
	}
}

func (sm *StateManager) GetCurrentState() ProgramState {
	return sm.currentState
}
