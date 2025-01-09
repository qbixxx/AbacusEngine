package controller

import (
	"abacus_engine/internal/interpreter"
	"abacus_engine/internal/state"
	"abacus_engine/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type AppController struct {
	ui           *ui.UI
	stateManager *state.StateManager
	interpreter  *interpreter.Interpreter
}

// NewAppController crea una nueva instancia del controlador de la aplicación.
func NewAppController(rows int) *AppController {
	ui := ui.NewUI(rows)
	stateManager := state.NewStateManager(ui.UpdateStateInfo)
	interpreter := interpreter.NewInterpreter(ui.MainPage.Table) // Conexión directa a MemoryTable

	// Configurar el callback para manejar initAddress desde la UI
	ui.SetInitAddressCallback(func(value string) {

		var initAddress int
		numericData64, _ := strconv.ParseInt(value, 16, 0)
		initAddress = int(numericData64)
		interpreter.SetInitAddress(initAddress)

		// Actualizar la información del intérprete en la UI
		rip, ac, init, enabled := interpreter.GetState()
		ui.UpdateInterpreterInfo(rip, ac, init, enabled)

	})

	return &AppController{
		ui:           ui,
		stateManager: stateManager,
		interpreter:  interpreter,
	}
}

// Run inicia la aplicación.
func (ac *AppController) Run() error {
	app := tview.NewApplication()

	// Configurar manejo de teclas
	app.SetInputCapture(ac.HandleKeyEvent)

	// Configurar la interfaz principal
	app.SetRoot(ac.ui.Pages, true)
	ac.updateInterpreterInfo()
	// Ejecutar la aplicación
	return app.Run()
}

// HandleKeyEvent maneja los eventos de teclado.
func (ac *AppController) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {

	case tcell.KeyCtrlE: // Edit
		ac.stateManager.SetState(state.Edit)
		ac.interpreter.Clean()
		ac.updateInterpreterInfo()

	case tcell.KeyCtrlD: // Debug
		if ac.interpreter.IsRunnable() {
			ac.stateManager.SetState(state.Debug)
			ac.interpreter.SetForDebug()
			ac.updateInterpreterInfo()

		}

	case tcell.KeyCtrlK: // Reset
		if ac.stateManager.GetCurrentState() == state.Edit {
			ac.interpreter.Reset()
			ac.updateInterpreterInfo()
		}

	case tcell.KeyCtrlR: // Run
		ac.stateManager.SetState(state.Run)
		ac.interpreter.SetForDebug()
		for ac.interpreter.IsRunnable() {
			ac.interpreter.Step()
			ac.updateInterpreterInfo()
		}
		ac.stateManager.SetState(state.Edit)

	case tcell.KeyCtrlI: // Input address form
		if ac.stateManager.GetCurrentState() == state.Edit {
			ac.setInitAddress()
		}

	case tcell.KeyEnter:
		if ac.stateManager.GetCurrentState() == state.Debug {
			ac.interpreter.Step()
			if !ac.interpreter.IsRunnable() {
				ac.stateManager.SetState(state.Edit)
			}
			ac.updateInterpreterInfo()
		}
		if ac.stateManager.GetCurrentState() == state.Edit {

		}

	// Ignorar otras teclas en modos específicos
	default:
		if ac.stateManager.GetCurrentState() == state.Debug {

			return nil
		}
	}

	return event
}
func (ac *AppController) setInitAddress() {
	ac.ui.ShowInputField()
}

// updateInterpreterInfo actualiza la sección de información del intérprete en la UI.
func (ac *AppController) updateInterpreterInfo() {
	instructionPointer, accumulator, initAddress, runnable := ac.interpreter.GetState()
	ac.ui.UpdateInterpreterInfo(instructionPointer, accumulator, initAddress, runnable)
}
