package main

import (
	"abacus_engine/internal/ui"
	"github.com/rivo/tview"
)

func main() {
	// Crear UI
	rows := 500
	myUI := ui.NewUI(rows)

	// Crear aplicación
	app := tview.NewApplication()

	// Establecer la grid principal en la aplicación
	if err := app.SetRoot(myUI.Grid, true).Run(); err != nil {
		panic(err)
	}
}
