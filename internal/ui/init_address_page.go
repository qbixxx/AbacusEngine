package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var inputField *tview.InputField
var initAddress string

func newInputField(ui *UI) *tview.InputField {
	// Configurar el campo de entrada
	inputField = tview.NewInputField().
		SetLabel("Enter Init Address: ").
		SetFieldWidth(20).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Recuperar el valor ingresado
				initAddress = inputField.GetText()
				if initAddress == "" {
					initAddress = "-1"
				}

				if ui.onInitAddress != nil {
					ui.onInitAddress(initAddress)
				}
				// Volver a la página principal
				ui.PageHolder.SwitchToPage("main")

			}
		})

	inputField.SetFieldBackgroundColor(tcell.ColorLightGreen).SetFieldTextColor(tcell.ColorBlack)
	inputField.Box.SetBorder(true)

	return inputField
}

func (ui *UI) ShowInputField() {
	ui.PageHolder.ShowPage("input")

}

// SetInitAddressCallback configura la función callback.
func (ui *UI) SetInitAddressCallback(callback func(string)) {
	ui.onInitAddress = callback
}
