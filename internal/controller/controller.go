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
	Ui           *ui.UI
	StateManager *state.StateManager
	Interpreter  *interpreter.Interpreter
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
		Ui:           ui,
		StateManager: stateManager,
		Interpreter:  interpreter,
	}
}

// Run inicia la aplicación.
func (ac *AppController) Run() error {
	app := tview.NewApplication().EnableMouse(true)

	// Configurar manejo de teclas
	app.SetInputCapture(ac.HandleKeyEvent)
	ac.InitializeHeap(4001, 4096)
	// Configurar la interfaz principal
	app.SetRoot(ac.Ui.PageHolder, true)

	ac.updateInterpreterInfo()

	// Ejecutar la aplicación
	return app.Run()
}

// HandleKeyEvent maneja los eventos de teclado.
func (ac *AppController) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {

	case tcell.KeyCtrlE: // Edit
		ac.StateManager.SetState(state.Edit)
		ac.Interpreter.Clean()
		ac.updateInterpreterInfo()

	case tcell.KeyCtrlD: // Debug
		if ac.Interpreter.IsRunnable() {
			ac.StateManager.SetState(state.Debug)
			ac.Interpreter.SetForDebug()
			ac.updateInterpreterInfo()

		}

	case tcell.KeyCtrlK: // Reset
		if ac.StateManager.GetCurrentState() == state.Edit {
			ac.Interpreter.Reset()
			ac.updateInterpreterInfo()
		}

	case tcell.KeyCtrlR: // Run
		ac.StateManager.SetState(state.Run)
		ac.Interpreter.SetForDebug()
		for ac.Interpreter.IsRunnable() {
			ac.Interpreter.Step()
			ac.updateInterpreterInfo()
		}
		ac.StateManager.SetState(state.Edit)

	case tcell.KeyCtrlI: // Input address form
		if ac.StateManager.GetCurrentState() == state.Edit {
			ac.setInitAddress()
		}
	case tcell.KeyCtrlH:
		ac.ToggleHeap()

	case tcell.KeyEnter:
		if ac.StateManager.GetCurrentState() == state.Debug {
			ac.Interpreter.Step()
			if !ac.Interpreter.IsRunnable() {
				ac.StateManager.SetState(state.Edit)
			}
			ac.updateInterpreterInfo()
		}
		if ac.StateManager.GetCurrentState() == state.Edit {

		}

	// Ignorar otras teclas en modos específicos
	default:
		if ac.StateManager.GetCurrentState() == state.Debug {

			return nil
		}
	}

	return event
}

func (ac *AppController) InitializeHeap(heapStart, heapEnd int) {
	ui := ac.Ui
	ui.CreateHeapView(heapStart, heapEnd)

	// Callback para actualizaciones
	ui.MainPage.Table.OnHeapUpdate = func(row int, value string) {
		ui.MainPage.UpdateHeap(row, value, heapStart, heapEnd)
	}
}

func (ac *AppController) ToggleHeap() {
	mp := ac.Ui.MainPage

	mp.IsHeapVisible = !mp.IsHeapVisible // Alternar visibilidad
	mp.UpdateLayout()                    // Actualizar el diseño de la página principal
}

func (ac *AppController) setInitAddress() {
	ac.Ui.ShowInputField()
}

// updateInterpreterInfo actualiza la sección de información del intérprete en la UI.
func (ac *AppController) updateInterpreterInfo() {
	instructionPointer, accumulator, initAddress, runnable := ac.Interpreter.GetState()
	ac.Ui.UpdateInterpreterInfo(instructionPointer, accumulator, initAddress, runnable)
}
