package controller

import (
	"abacus_engine/internal/interpreter"
	"abacus_engine/internal/state"
	"abacus_engine/internal/ui"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
		//fmt.Printf("Setting initAddress to %s\n", value)
		// Convertir el valor a entero
		var initAddress int
		fmt.Sscanf(value, "%d", &initAddress)
		interpreter.SetInitAddress(initAddress)

		// Actualizar la información del intérprete en la UI
		rip, ac, init, enabled := interpreter.GetState()
		ui.UpdateInterpreterInfo(fmt.Sprintf("RIP: %03X\nAccumulator: %d\nInit Address: %d\nEnabled: %v", rip, ac, init, enabled))

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
	case tcell.KeyCtrlE:
		ac.stateManager.SetState(state.Edit)
	case tcell.KeyCtrlD:
		if ac.interpreter.IsRunnable(){
			ac.stateManager.SetState(state.Debug)
			ac.interpreter.SetForDebug()
			ac.updateInterpreterInfo()

		}

	//reset
	case tcell.KeyCtrlK:
		if ac.stateManager.GetCurrentState() == state.Edit{
			ac.interpreter.Reset()
			ac.updateInterpreterInfo()
		}
		

	case tcell.KeyCtrlR:
		ac.stateManager.SetState(state.Run)
	case tcell.KeyCtrlI:
		if ac.stateManager.GetCurrentState() == state.Edit{
			ac.setInitAddress()
		}
	// Ejecutar instrucción en modo Debug
	case tcell.KeyEnter:
		if ac.stateManager.GetCurrentState() == state.Debug {
			ac.interpreter.Step()
			ac.updateInterpreterInfo()
		}

	// Ignorar otras teclas en modos específicos
	default:
		if ac.stateManager.GetCurrentState() == state.Run {
			// No permitir interacciones adicionales en modo Run
			return nil
		}
	}

	return event
}
func (ac *AppController) setInitAddress() {
	ac.ui.SwitchToPage("input")
	//ac.ui.SwitchToPage("main")

}

// updateInterpreterInfo actualiza la sección de información del intérprete en la UI.
func (ac *AppController) updateInterpreterInfo() {
	instructionPointer, accumulator, initAddress, runnable := ac.interpreter.GetState()
	ac.ui.UpdateInterpreterInfo(
		fmt.Sprintf("RIP: %03X\nAccumulator: %d\nInit Address: %d\nEnabled: %v", instructionPointer, accumulator, initAddress, runnable),
	)
}
