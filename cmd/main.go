package main

import (
	"abacus_engine/internal/state"
	"abacus_engine/internal/ui"
	//"abacus_engine/internal/memory"
	"github.com/rivo/tview"
)

func main() {
	rows := 500

	// Inicializar componentes
	app, err := initializeApplication(rows)
	if err != nil {
		panic(err)
	}

	// Ejecutar la aplicación
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// initializeApplication configura la aplicación con todos los componentes.
func initializeApplication(rows int) (*tview.Application, error) {
	// Crear modelos de datos
	//memoryModel := memory.NewMemory(rows)

	// Crear UI
	uiManager := ui.NewUI(rows)//, memoryModel)

	// Crear gestor de estados
	stateManager := state.NewStateManager(uiManager.UpdateStateInfo)
	stateManager.SetState(state.Edit)

	// Crear aplicación
	app := tview.NewApplication()
	app.SetInputCapture(stateManager.HandleKeyEvent)

	// Configurar la UI en la raíz
	app.SetRoot(uiManager.RootGrid, true)
	

	return app, nil
}
