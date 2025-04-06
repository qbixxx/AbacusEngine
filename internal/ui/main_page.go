package ui

import (
	"abacus_engine/internal/memory"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const asciiTitle = "[lightgreen]" +
	"_____________                             \n" +
	"___    |__  /_______ __________  _________\n" +
	"__  /| |_  __ \\  __ `/  ___/  / / /_  ___/\n" +
	"_  ___ |  /_/ / /_/ // /__ / /_/ /_(__  ) \n" +
	"/_/  |_/_.___/\\__,_/ \\___/ \\__,_/ /____/  \n" +
	"                                         \n" +
	"__________              _____             \n" +
	"___  ____/_____________ ___(_)___________ \n" +
	"__  __/  __  __ \\_  __ `/_  /__  __ \\  _ \\\n" +
	"_  /___  _  / / /  /_/ /_  / _  / / /  __/\n" +
	"/_____/  /_/ /_/_\\__, / /_/  /_/ /_/\\___/ \n" +
	"                /____/                    \n"

type MainPage struct {
	RootGrid        *tview.Grid
	Table           *memory.MemoryTable
	HeapTable       *tview.Table // Tabla del heap
	MenuGrid        *tview.Grid
	Title           *tview.TextView
	InfoState       *tview.TextView
	InfoInterpreter *tview.TextView
	Footer          *tview.TextView
	IsHeapVisible   bool // Indica si la tabla del heap está visible
}

func newMainPage(rows int) *MainPage {
	// Inicializar MainPage
	mainPage := &MainPage{
		Table:           memory.NewMemoryTable(rows),
		Title:           tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText(asciiTitle),
		InfoState:       tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter),
		InfoInterpreter: tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter),
		Footer:          tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("[black:green]^E[white:black] Edit	[black:green]^D[white:black] Debug	[black:green]^R[white:black] Run	[black:green]^I[white:black] Set Init Address	[black:green]^K[white:black] Reset	[black:green]^H[white:black] Show Heap	[black:green]^O[white:black] Open File"),
	}

	// Configurar MenuGrid
	mainPage.MenuGrid = tview.NewGrid().
		SetRows(16, 1, -1).
		AddItem(mainPage.Title, 0, 0, 1, 1, 0, 0, false).
		AddItem(mainPage.InfoState, 1, 0, 1, 1, 0, 0, true).
		AddItem(mainPage.InfoInterpreter, 2, 0, 1, 1, 0, 0, true)

	// Configurar RootGrid
	mainPage.RootGrid = tview.NewGrid().
		SetRows(-1, 1).
		SetColumns(44, 0).
		AddItem(mainPage.MenuGrid, 0, 0, 1, 1, 0, 0, false).
		AddItem(mainPage.Table.GetTable(), 0, 1, 1, 1, 0, 0, true).
		AddItem(mainPage.Footer, 1, 0, 1, 2, 0, 0, false)

	return mainPage
}

func (ui *UI) UpdateStateInfo(state string) {

	var color string
	switch state {

	case "Edit":
		color = "[yellow]"
	case "Debug":
		color = "[turquoise]"
	case "Run":
		color = "[white]"
	}

	ui.MainPage.InfoState.SetText(fmt.Sprintf("[green]Mode: "+color+"%s", state))
}

// UpdateStateTitle actualiza la sección Interpreter Info con el estado actual del interprete.
func (ui *UI) UpdateInterpreterInfo(rip, acc, initAdr int, runnable bool) {

	var colorEnable string

	if runnable {
		colorEnable = "[green]"
	} else {
		colorEnable = "[red]"
	}

	if rip == -1 {
		ui.MainPage.InfoInterpreter.SetText(fmt.Sprintf("RIP: [red]undefined[white]\nAccumulator: %03X\nInit Address: "+colorEnable+"%03X\n"+"[white]Enabled: "+colorEnable+"%v", acc, initAdr, runnable))
	} else {
		ui.MainPage.InfoInterpreter.SetText(fmt.Sprintf("RIP: %03X\nAccumulator: %03X\nInit Address: "+colorEnable+"%03X\n"+"[white]Enabled: "+colorEnable+"%v", rip, acc, initAdr, runnable))
	}

}

func (mp *MainPage) GetTable() *tview.Table {
	return mp.Table.GetTable()
}



func (mp *MainPage) UpdateLayout() {
	mp.RootGrid.Clear() // Limpia la disposición actual

	if mp.IsHeapVisible {
		// Layout con la tabla del heap visible
		mp.RootGrid.SetRows(-1, 1).
			SetColumns(44, -1, 27). // Se añade una columna para el heap
			AddItem(mp.MenuGrid, 0, 0, 1, 1, 0, 0, false).
			AddItem(mp.Table.GetTable(), 0, 1, 1, 1, 0, 0, true).
			AddItem(mp.HeapTable, 0, 2, 1, 1, 0, 0, false). // Añadir el heap
			AddItem(mp.Footer, 1, 0, 1, 3, 0, 0, false)
	} else {
		// Layout sin la tabla del heap
		mp.RootGrid.SetRows(-1, 1).
			SetColumns(44, 0).
			AddItem(mp.MenuGrid, 0, 0, 1, 1, 0, 0, false).
			AddItem(mp.Table.GetTable(), 0, 1, 1, 1, 0, 0, true).
			AddItem(mp.Footer, 1, 0, 1, 2, 0, 0, false)
	}
}

func (mp *MainPage) UpdateHeap(row int, value string, heapStart, heapEnd int) {

	if row >= heapStart && row <= heapEnd {

		mp.HeapTable.GetCell(row-heapStart, 1).SetText(value)

	}

}

func (mp *MainPage) CreateHeapTable(heapStart, heapEnd int) {

	mp.HeapTable = tview.NewTable().
		SetBorders(true).
		SetSelectable(false, false)

	mp.HeapTable.SetBorderColor(tcell.ColorRed)

	// Encabezados

	headers := []string{"  Memory Address  ", "  Data  "}

	for col, header := range headers {

		mp.HeapTable.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter))

	}

	// Rellenar la tabla

	for i := heapStart; i <= heapEnd; i++ {

		mp.HeapTable.SetCell(i-heapStart+1, 0, tview.NewTableCell(fmt.Sprintf("%03X", i)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		mp.HeapTable.SetCell(i-heapStart+1, 1, tview.NewTableCell("NOP").
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

	}

}

func (ui *UI) CreateHeapView(heapStart, heapEnd int) {

	ui.MainPage.HeapTable = tview.NewTable()

	ui.MainPage.HeapTable.SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(false, false).
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle("[ Heap ]").
		SetTitleAlign(tview.AlignCenter)

	ui.MainPage.HeapTable.SetBordersColor(tcell.ColorRed)

	// Configurar encabezados

	headers := []string{" Memory Address ", " Data "}

	for col, header := range headers {

		ui.MainPage.HeapTable.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter))

	}

	// Rellenar la tabla con datos iniciales (vacíos)

	for i := heapStart; i <= heapEnd; i++ {

		ui.MainPage.HeapTable.SetCell(i-heapStart+1, 0, tview.NewTableCell(fmt.Sprintf("%03X", i)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

		ui.MainPage.HeapTable.SetCell(i-heapStart+1, 1, tview.NewTableCell("0000").
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter))

	}

}
