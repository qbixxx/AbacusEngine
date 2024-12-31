package ui

import (
	"fmt"
	//"abacus_engine/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const asciiTitle = "[turquoise]" +
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

// MemoryTable encapsula la lógica y el comportamiento de la tabla de memoria.
type MemoryTable struct {
	table   *tview.Table // Componente gráfico de la tabla
	rows    int          // Cantidad de filas
	prevRow int          // Fila previamente seleccionada
	prevCol int          // Columna previamente seleccionada
}

// NewMemoryTable crea y configura una nueva instancia de MemoryTable.
func NewMemoryTable(rows int) *MemoryTable {
	memoryTable := &MemoryTable{
		table:   tview.NewTable(),
		rows:    rows,
		prevRow: 1,
		prevCol: 1,
	}
	memoryTable.initTable()
	return memoryTable
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
	RootGrid  		*tview.Grid
	Table 			*MemoryTable
	MenuGrid		*tview.Grid
	Title	  		*tview.TextView
	InfoState		*tview.TextView
	InfoInterpreter	*tview.TextView
}

// NewUI crea una nueva interfaz de usuario con una tabla y un widget lateral.
func NewUI(rows int) *UI {
	
	ui := &UI{}
	
	ui.Table = NewMemoryTable(rows)
	ui.Title = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(asciiTitle)

	ui.Title.//.SetBorder(true).
		//SetTitle("Title").
		SetTitleAlign(tview.AlignCenter)
	
		ui.InfoState = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	ui.InfoInterpreter = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)


	ui.MenuGrid = tview.NewGrid().
			SetRows(16,1, -1).
			//SetColumns(0).
			AddItem(ui.Title, 0, 0, 1, 1, 0, 0, false).
			AddItem(ui.InfoState, 1, 0, 1, 1, 0, 0, true).
			AddItem(ui.InfoInterpreter, 2, 0, 1, 1, 0, 0, true)

	ui.UpdateInterpreterInfo("null")

	ui.RootGrid = tview.NewGrid().
		SetRows(0).
		SetColumns(0, 48).
		AddItem(ui.Table.GetTable(), 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.MenuGrid, 0, 1, 1, 1, 0, 0, false)

	return ui
}

// UpdateStateTitle actualiza la sección Title con el estado actual.
func (ui *UI) UpdateInterpreterInfo(info string) {
	ui.InfoInterpreter.SetText(fmt.Sprintf("[red]Loading point address: %s", info))
}

// UpdateStateTitle actualiza la sección Title con el estado actual.
func (ui *UI) UpdateStateInfo(state string) {
	ui.InfoState.SetText(fmt.Sprintf("[green]Mode: %s", state))
}
