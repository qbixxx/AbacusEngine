// internal/ui/ui.go
package ui

import (
	"github.com/rivo/tview"
	//"github.com/gdamore/tcell/v2"
)

// UI encapsula la estructura gráfica del programa.
type UI struct {
	Tview         *tview.Application
	PageHolder    *tview.Pages
	MainPage      *MainPage
	onInitAddress func(string) // Callback para manejar el initAddress
}

func NewUI(rows int) *UI {
	ui := &UI{
		MainPage:   newMainPage(rows),
		PageHolder: tview.NewPages(),
		Tview:      tview.NewApplication().EnableMouse(true),
	}

	inputField := newInputField(ui)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	// lógica para mostrar modal de error
	showModal := func(message string) {
		modal := tview.NewModal().
			SetText(message).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				ui.PageHolder.SwitchToPage("main")
			})
		ui.PageHolder.AddPage("error-modal", modal, true, true)
	}

	switchToMain := func() {
		ui.PageHolder.SwitchToPage("main")
	}

	fileManager := NewFileManagerPage(ui.MainPage.Table.GetTable(), switchToMain, showModal)

	ui.PageHolder.AddPage("main", ui.MainPage.RootGrid, true, true)
	ui.PageHolder.AddPage("input", modal(inputField, 42, 3), true, true)
	ui.PageHolder.AddPage("file-manager", fileManager, true, false)
	ui.PageHolder.SwitchToPage("main")

	return ui
}
