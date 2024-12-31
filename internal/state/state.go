
package state

import (
	//"fmt"
	"github.com/gdamore/tcell/v2"
)

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
	currentState   ProgramState
	OnStateChange  func(state string)
}

func NewStateManager(onStateChange func(state string)) *StateManager {
	return &StateManager{
		currentState:  Edit,
		OnStateChange: onStateChange,
	}
}

func (sm *StateManager) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch sm.currentState {
	case Edit:
		return sm.handleEditKeys(event)
	case Debug:
		return sm.handleDebugKeys(event)
	case Run:
		return sm.handleRunKeys(event)
	default:
		return nil
	}
}

// handleEditKeys maneja las teclas en modo Edit.
func (sm *StateManager) handleEditKeys(event *tcell.EventKey) *tcell.EventKey {
	// F1, F2, F3 siempre están permitidas
	switch event.Key() {
	case tcell.KeyF1:
		sm.SetState(Edit)
	case tcell.KeyF2:
		sm.SetState(Debug)
	case tcell.KeyF3:
		sm.SetState(Run)
	default:
		// Permitir todas las demás teclas
		return event
	}
	return nil
}

// handleDebugKeys maneja las teclas en modo Debug.
func (sm *StateManager) handleDebugKeys(event *tcell.EventKey) *tcell.EventKey {
	// F1, F2, F3 siempre están permitidas
	switch event.Key() {
	case tcell.KeyF1:
		sm.SetState(Edit)
	case tcell.KeyF2:
		sm.SetState(Debug)
	case tcell.KeyF3:
		sm.SetState(Run)
	case tcell.KeyDown: // Solo permitir flecha hacia abajo
		return event
	}
	// Ignorar cualquier otra tecla
	return nil
}

// handleRunKeys maneja las teclas en modo Run.
func (sm *StateManager) handleRunKeys(event *tcell.EventKey) *tcell.EventKey {
	// F1, F2, F3 siempre están permitidas
	switch event.Key() {
	case tcell.KeyF1:
		sm.SetState(Edit)
	case tcell.KeyF2:
		sm.SetState(Debug)
	case tcell.KeyF3:
		sm.SetState(Run)
	}
	// Ignorar cualquier otra tecla
	return nil
}

func (sm *StateManager) SetState(newState ProgramState) {
	sm.currentState = newState
	if sm.OnStateChange != nil {
		sm.OnStateChange(newState.String())
	}
}