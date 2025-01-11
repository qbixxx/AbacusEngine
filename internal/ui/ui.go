package ui

import (
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

//const asciiTitle = "[lightgreen]" +
//"                         \n"+
//" █████╗ ███╗   ██╗ █████╗ \n"+
//"██╔══██╗████╗  ██║██╔══██╗\n"+
//"███████║██╔██╗ ██║███████║\n"+
//"██╔══██║██║╚██╗██║██╔══██║\n"+
//"██║  ██║██║ ╚████║██║  ██║\n"+
//"╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝\n"+
//"Abacus's Not Assembly :)"

//const asciiTitle = "[lightgreen]\n" +
//
//"                                                 \n"+
//" █████╗ ██████╗  █████╗  ██████╗██╗   ██╗███████╗\n"+
//"██╔══██╗██╔══██╗██╔══██╗██╔════╝██║   ██║██╔════╝\n"+
//"██████║██████╔╝███████║██║     ██║   ██║███████╗\n"+
//"██╔══██║██╔══██╗██╔══██║██║     ██║   ██║╚════██║\n"+
//"██║  ██║██████╔╝██║  ██║╚██████╗╚██████╔╝███████║\n"+
//"╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝\n"+
//"                                                 \n"+
//" █████╗  ██████╗  ██████╗  ██████╗               \n"+
//"██╔══██╗██╔═████╗██╔═████╗██╔═████╗              \n"+
//"╚██████║██║██╔██║██║██╔██║██║██╔██║              \n"+
//" ╚═══██║████╔╝██║████╔╝██║████╔╝██║              \n"+
//" █████╔╝╚██████╔╝╚██████╔╝╚██████╔╝              \n"+
//" ╚════╝  ╚═════╝  ╚═════╝  ╚═════╝               \n"
//
//

// MemoryTable encapsula la lógica y el comportamiento de la tabla de memoria.
type MemoryTable struct {
	table        *tview.Table // Componente gráfico de la tabla
	rows         int          // Cantidad de filas
	prevRow      int          // Fila previamente seleccionada
	prevCol      int          // Columna previamente seleccionada
	OnHeapUpdate func(row int, value string) // Callback para actualizaciones del heap
}

// NewMemoryTable crea y configura una nueva instancia de MemoryTable.
func NewMemoryTable(rows int) *MemoryTable {
	memoryTable := &MemoryTable{
		table:   tview.NewTable(), //.SetOffset(8,0),
		rows:    rows,
		prevRow: 1,
		prevCol: 1,
	}
	memoryTable.initTable()
	return memoryTable
}

func (m *MemoryTable) ScrollToCurrentRow(row int) {
	// Calcular el nuevo desplazamiento
	visibleRows := 10 // hacer constante global
	if row >= visibleRows {
		m.table.SetOffset(row-visibleRows+1, 0)
	} else {
		m.table.SetOffset(0, 0)
	}
}

func (m *MemoryTable) GetSize() int {
	return m.rows
}

// initTable configura la tabla con encabezados, filas iniciales y comportamiento.
func (m *MemoryTable) initTable() {
	m.table.SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(true, true)

	headers := []string{" Memory address ", " Instruction/Data ", "Commentary"}
	for col, header := range headers {
		expansion := 0
		if col == 2 {
			expansion = 1
		}
		m.table.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(expansion))
	}

	for row := 1; row <= m.rows; row++ {
		m.addRow(row)
	}

	m.table.SetInputCapture(m.handleInput)
	m.table.SetSelectionChangedFunc(m.handleSelectionChange)
	m.table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorRed).
		Foreground(tcell.ColorWhite))
}

// addRow agrega una fila a la tabla con valores iniciales.
func (m *MemoryTable) addRow(row int) {
	m.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%03X", row-1)).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetSelectable(false).
		SetAlign(tview.AlignCenter).
		SetExpansion(0))

	m.table.SetCell(row, 1, tview.NewTableCell("NOP").
		SetTextColor(tcell.ColorGreen).
		SetSelectable(true).
		SetAlign(tview.AlignCenter).
		SetMaxWidth(4).
		SetExpansion(0))

	m.table.SetCell(row, 2, tview.NewTableCell("").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetSelectable(true).
		SetAlign(tview.AlignLeft).
		SetExpansion(18))
}

func (m *MemoryTable) ColorInitAddr(addr int) {
	cell := m.table.GetCell(addr+1, 0)
	if addr < 0 {
		cell.SetTextColor(tcell.ColorWhite)
	} else {
		cell.SetTextColor(tcell.ColorRed)
	}

}

// handleInput maneja los eventos de teclado para la tabla.
func (m *MemoryTable) handleInput(event *tcell.EventKey) *tcell.EventKey {
	row, column := m.table.GetSelection()
	cell := m.table.GetCell(row, column)

	if event.Key() == tcell.KeyRune {
		if column == 1 {
			if cell.Text == "NOP" {
				cell.SetText("")
			}
			if len(cell.Text) < 4 {
				cell.SetText(cell.Text + string(event.Rune()))
			}

			m.updateCellColor(cell)

		} else if column == 2 {
			cell.SetText(cell.Text + string(event.Rune()))

		}
		m.table.SetCell(row, column, cell)
		return nil
	}

	if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
		if len(cell.Text) > 0 {
			cell.SetText(cell.Text[:len(cell.Text)-1])
		}
		m.table.SetCell(row, column, cell)
		return nil
	}

	if event.Key() == tcell.KeyDelete {
		cell.SetText("")
		m.table.SetCell(row, column, cell)
		return nil
	}

	if event.Key() == tcell.KeyEnter {
		if row < m.rows-1 {
			m.table.Select(row+1, column)
		}
		return nil
	}

	return event
}

func (m *MemoryTable) GetCell(row int) *tview.TableCell {
	return m.table.GetCell(row, 1)
}

func (m *MemoryTable) updateCellColor(cell *tview.TableCell) {
	switch cell.Text[0] {
	case '0':
		cell.SetTextColor(tcell.ColorYellow)
	case '1':
		cell.SetTextColor(tcell.ColorOrange)
	case '2':
		cell.SetTextColor(tcell.ColorRed)
	case '3':
		cell.SetTextColor(tcell.ColorMediumVioletRed)
	case '4':
		cell.SetTextColor(tcell.ColorHotPink)
	case '7':
		cell.SetTextColor(tcell.ColorLime)
	case '8':
		cell.SetTextColor(tcell.ColorTurquoise)
	case '9':
		cell.SetTextColor(tcell.ColorViolet)
	case 'F':
		cell.SetTextColor(tcell.ColorGrey)
	default:
		cell.SetTextColor(tcell.ColorWhite)
	}
	if cell.Text == "NOP" {
		cell.SetTextColor(tcell.ColorGreen)
	}
}
func (m *MemoryTable) Goto(row int, column int) {
	m.table.Select(row+1, column)
}

func (m *MemoryTable) GetInstruction(row int) string {
	cell := m.table.GetCell(row+1, 1) // +1 para ignorar el encabezado
	if cell == nil {
		return "NOP"
	}
	return cell.Text
}
func (m *MemoryTable) ResetTable() {
	for row := 1; row <= m.rows; row++ {
		m.table.GetCell(row, 1).SetText("NOP")
		m.table.GetCell(row, 2).SetText("")
		m.updateCellColor(m.table.GetCell(row, 1))
	}
}

func (m *MemoryTable) WriteCell(row int, accumulator int) {
	cell := m.table.GetCell(row+1, 1) // +1 porque la fila 0 es para encabezados

	// Formatear el acumulador como string con ceros completados
	formattedData := fmt.Sprintf("%04d", accumulator)
	cell.SetText(formattedData)

	m.table.SetCell(row+1, 1, cell) // +1 porque la fila 0 es el encabezado
	m.updateCellColor(cell)

	// Llamar al callback si está configurado
	if m.OnHeapUpdate != nil {
		m.OnHeapUpdate(row, formattedData)
	}
}
// handleSelectionChange maneja el evento de cambio de selección.
func (m *MemoryTable) handleSelectionChange(newRow, newCol int) {
	if m.prevCol == 1 {
		prevCell := m.table.GetCell(m.prevRow, m.prevCol)
		if prevCell.Text == "" {
			prevCell.SetText("NOP").SetTextColor(tcell.ColorGreen)
			m.table.SetCell(m.prevRow, m.prevCol, prevCell)
		}
	}
	m.prevRow = newRow
	m.prevCol = newCol
}

// GetTable devuelve el componente gráfico de la tabla.
func (m *MemoryTable) GetTable() *tview.Table {
	return m.table
}

// UI encapsula la estructura gráfica del programa.
type UI struct {
	Pages         *tview.Pages
	MainPage      MainPage
	onInitAddress func(string) // Callback para manejar el initAddress
	//inputField *tview.TextView
}

// SetInitAddressCallback configura la función callback.
func (ui *UI) SetInitAddressCallback(callback func(string)) {
	ui.onInitAddress = callback
}

type MainPage struct {
	RootGrid        *tview.Grid
	Table           *MemoryTable
	HeapTable       *tview.Table    // Tabla del heap
	MenuGrid        *tview.Grid
	Title           *tview.TextView
	InfoState       *tview.TextView
	InfoInterpreter *tview.TextView
	Footer          *tview.TextView
	IsHeapVisible   bool            // Indica si la tabla del heap está visible
}

func (ui *UI) CreateHeapView(heapStart, heapEnd int) {
    ui.MainPage.HeapTable = tview.NewTable()
    ui.MainPage.HeapTable.SetBorders(true).
        SetSelectable(false, false).
        SetBorder(true).
        SetBorderColor(tcell.ColorRed).
        SetTitle("- Heap -").
        SetTitleAlign(tview.AlignCenter)

	
    // Configurar encabezados
    headers := []string{"  Memory Address  ", "  Data "}
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
        ui.MainPage.HeapTable.SetCell(i-heapStart+1, 1, tview.NewTableCell("NOP").
            SetTextColor(tview.Styles.PrimaryTextColor).
            SetAlign(tview.AlignCenter))
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

func (mp *MainPage) UpdateHeap(row int, value string, heapStart, heapEnd int) {
	if row >= heapStart && row <= heapEnd {
		mp.HeapTable.GetCell(row-heapStart+1, 1).SetText(value)
	}
}

func (mp *MainPage) UpdateLayout() {
	mp.RootGrid.Clear() // Limpia la disposición actual

	if mp.IsHeapVisible {
		// Layout con la tabla del heap visible
		mp.RootGrid.SetRows(-1, 1).
			SetColumns(44, -1, 30). // Se añade una columna para el heap
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


func NewUI(rows int) *UI {

	pages := tview.NewPages()

	ui := &UI{
		Pages: pages,
	}

	// Inicializar MainPage
	mainPage := &MainPage{
		Table:           NewMemoryTable(rows),
		Title:           tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText(asciiTitle),
		InfoState:       tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter),
		InfoInterpreter: tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter),
		Footer:          tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("[black:green]^E[white:black] Edit	[black:green]^D[white:black] Debug	[black:green]^R[white:black] Run	[black:green]^I[white:black] Set Init Address	[black:green]^K[white:black] Reset"),
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

	// Asignar MainPage a UI
	ui.MainPage = *mainPage

	var inputField *tview.InputField
	var initAddress string

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
				pages.SwitchToPage("main")

			}
		})

	inputField.SetFieldBackgroundColor(tcell.ColorLightGreen).SetFieldTextColor(tcell.ColorBlack)
	inputField.Box.SetBorder(true)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}
	ui.Pages.AddPage("main", mainPage.RootGrid, true, true)
	ui.Pages.AddPage("input", modal(inputField, 42, 3), true, true)
	ui.Pages.SwitchToPage("main")
	return ui
}

func (ui *UI) ShowInputField() {
	ui.Pages.ShowPage("input")

}

func (ui *UI) ToggleHeapView() {
    if ui.MainPage.IsHeapVisible {
        ui.MainPage.RootGrid.RemoveItem(ui.MainPage.HeapTable)
      	ui.MainPage.IsHeapVisible = false
    } else {
        ui.MainPage.RootGrid.AddItem(ui.MainPage.HeapTable, 0, 2, 1, 1, 0, 0, false)
        ui.MainPage.IsHeapVisible = true
    }
}


// UpdateStateTitle actualiza la sección State con el estado actual.
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
		ui.MainPage.InfoInterpreter.SetText(fmt.Sprintf("RIP: [red]undefined[white]\nAccumulator: %d\nInit Address: "+colorEnable+"%03X\n"+"[white]Enabled: "+colorEnable+"%v", acc, initAdr, runnable))
	} else {
		ui.MainPage.InfoInterpreter.SetText(fmt.Sprintf("RIP: %03X\nAccumulator: %d\nInit Address: "+colorEnable+"%03X\n"+"[white]Enabled: "+colorEnable+"%v", rip, acc, initAdr, runnable))
	}

}
