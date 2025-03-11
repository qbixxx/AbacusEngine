package ui

import (
	"github.com/rivo/tview"
)


// UI encapsula la estructura gráfica del programa.
type UI struct {

	Tview	*tview.Application
	PageHolder *tview.Pages
	MainPage      *MainPage
	onInitAddress func(string) // Callback para manejar el initAddress
}


func NewUI(rows int) *UI {

	ui := &UI{
		MainPage:	newMainPage(rows),
		PageHolder: tview.NewPages(),
		Tview:		tview.NewApplication().EnableMouse(true),
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

	ui.PageHolder.AddPage("main", ui.MainPage.RootGrid, true, true)
	ui.PageHolder.AddPage("input", modal(inputField, 42, 3), true, true)
	ui.PageHolder.SwitchToPage("main")

	return ui
}






