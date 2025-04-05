package controller

import (
	"abacus_engine/internal/filemanager"
	"abacus_engine/internal/interpreter"
	"abacus_engine/internal/state"
	"abacus_engine/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type AppController struct {
	App          *tview.Application
	Ui           *ui.UI
	StateManager *state.StateManager
	Interpreter  *interpreter.Interpreter
	FileManager  *filemanager.FileManager
	Adapter      *ui.MemoryTableAdapter
}

func NewAppController(rows int) *AppController {
	app := tview.NewApplication().EnableMouse(true)
	uiInstance := ui.NewUI(rows)
	adapter := ui.NewMemoryTableAdapter(uiInstance.MainPage.GetTable())
	stateManager := state.NewStateManager(uiInstance.UpdateStateInfo)
	interpreterInstance := interpreter.NewInterpreter(uiInstance.MainPage.Table)

	uiInstance.SetInitAddressCallback(func(value string) {
		if numericData64, err := strconv.ParseInt(value, 16, 0); err == nil {
			interpreterInstance.SetInitAddress(int(numericData64))
		}
		rip, ac, init, enabled := interpreterInstance.GetState()
		uiInstance.UpdateInterpreterInfo(rip, ac, init, enabled)
	})

    fileManagerPage := ui.NewFileManagerPage(adapter, func() {
        uiInstance.PageHolder.SwitchToPage("main")
    }, func(msg string) {
        uiInstance.ShowModal(msg)
    })
    
    uiInstance.PageHolder.AddPage("file-manager", fileManagerPage, true, false)

	return &AppController{
		App:          app,
		Ui:           uiInstance,
		StateManager: stateManager,
		Interpreter:  interpreterInstance,
		Adapter:      adapter,
	}
}

func (ac *AppController) Run() error {
	ac.App.SetInputCapture(ac.HandleKeyEvent)
	ac.App.SetRoot(ac.Ui.PageHolder, true)
	ac.InitializeHeap(4001, 4096)
	ac.updateInterpreterInfo()
	return ac.App.Run()
}

func (ac *AppController) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlE:
		ac.enterEditMode()
	case tcell.KeyCtrlD:
		ac.enterDebugMode()
	case tcell.KeyCtrlK:
		if ac.StateManager.GetCurrentState() == state.Edit {
			ac.Interpreter.Reset()
			ac.updateInterpreterInfo()
		}
	case tcell.KeyCtrlR:
		ac.runInterpreter()
	case tcell.KeyCtrlO:
		ac.Ui.PageHolder.SwitchToPage("file-manager")
	case tcell.KeyCtrlI:
		if ac.StateManager.GetCurrentState() == state.Edit {
			ac.setInitAddress()
		}
	case tcell.KeyCtrlH:
		ac.ToggleHeap()
	case tcell.KeyEnter:
		ac.stepDebug()
    
    default:
        if ac.StateManager.GetCurrentState() == state.Debug {
            return nil
        }
	}


	return event
}

func (ac *AppController) updateInterpreterInfo() {
	rip, acVal, init, enabled := ac.Interpreter.GetState()
	ac.Ui.UpdateInterpreterInfo(rip, acVal, init, enabled)
}

func (ac *AppController) enterEditMode() {
	ac.StateManager.SetState(state.Edit)
	ac.Interpreter.Clean()
	ac.updateInterpreterInfo()
}

func (ac *AppController) enterDebugMode() {
	if ac.Interpreter.IsRunnable() {
		ac.StateManager.SetState(state.Debug)
		ac.Interpreter.SetForDebug()
		ac.updateInterpreterInfo()
	}
}

func (ac *AppController) runInterpreter() {
	ac.StateManager.SetState(state.Run)
	ac.Interpreter.SetForDebug()
	for ac.Interpreter.IsRunnable() {
		ac.Interpreter.Step()
		ac.updateInterpreterInfo()
	}
	ac.StateManager.SetState(state.Edit)
}

func (ac *AppController) stepDebug() {
	if ac.StateManager.GetCurrentState() == state.Debug {
		ac.Interpreter.Step()
		if !ac.Interpreter.IsRunnable() {
			ac.StateManager.SetState(state.Edit)
		}
		ac.updateInterpreterInfo()
	}
}

func (ac *AppController) setInitAddress() {
	ac.Ui.ShowInitAddressInput()
}

func (ac *AppController) ToggleHeap() {
	ac.Ui.ToggleHeapPage()
}

func (ac *AppController) InitializeHeap(from, to int) {
	ac.Ui.CreateHeapView(from, to)
    ac.Ui.MainPage.Table.OnHeapUpdate = func(row int, value string) {
		ac.Ui.MainPage.UpdateHeap(row, value, from, to)
	}
}
