package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
/*
const asciiTitle = 
" \n[cyan]░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░░▒▓███████▓▒░ "+
"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░        \n"+
"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░        \n"+
"░▒▓████████▓▒░▒▓███████▓▒░░▒▓████████▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░  \n"+
"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░ \n"+
"░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░ \n"+
"░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░ ░▒▓██████▓▒░░▒▓███████▓▒░  \n"+
"																				\n"+                                                                                
"																				\n"+                                                                                                                                                                
"[turquoise]░▒▓████████▓▒░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓████████▓▒░        \n"+
"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░               \n"+
"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░               \n"+
"░▒▓██████▓▒░ ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒▒▓███▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓██████▓▒░          \n"+
"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░               \n"+
"░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░               \n"+
"░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░        \n"
                          
*/
						  

const asciiTitle = "[turquoise]"+
"_____________                             \n"+
"___    |__  /_______ __________  _________\n"+
"__  /| |_  __ \\  __ `/  ___/  / / /_  ___/\n"+
"_  ___ |  /_/ / /_/ // /__ / /_/ /_(__  ) \n"+
"/_/  |_/_.___/\\__,_/ \\___/ \\__,_/ /____/  \n"+
"                                         \n"+
"__________              _____             \n"+
"___  ____/_____________ ___(_)___________ \n"+
"__  __/  __  __ \\_  __ `/_  /__  __ \\  _ \\\n"+
"_  /___  _  / / /  /_/ /_  / _  / / /  __/\n"+
"/_____/  /_/ /_/_\\__, / /_/  /_/ /_/\\___/ \n"+
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

	// Configurar encabezados
	headers := []string{" Memory address ", " OPCode/Data ", "Commentary"}
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

	// Agregar filas iniciales
	for row := 1; row <= m.rows; row++ {
		m.addRow(row)
	}

	// Configurar eventos de entrada
	m.table.SetInputCapture(m.handleInput)

	// Configurar cambio de selección
	m.table.SetSelectionChangedFunc(m.handleSelectionChange)

	// Estilo para la celda seleccionada
	m.table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorRed).
		Foreground(tcell.ColorWhite))
}

// addRow agrega una fila a la tabla con valores iniciales.
func (m *MemoryTable) addRow(row int) {
	m.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%03d", row-1)).
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

	// Ingreso de texto
	if event.Key() == tcell.KeyRune {
		if column == 1 {
			if cell.Text == "NOP" {
				cell.SetText("")
			}
			if len(cell.Text) < 4 { // Limitar a 4 caracteres
				cell.SetText(cell.Text + string(event.Rune()))
			}
		} else if column == 2 { // Sin límite para la columna de comentarios
			cell.SetText(cell.Text + string(event.Rune()))
		}
		m.table.SetCell(row, column, cell)
		return nil
	}

	// Borrar texto con Backspace
	if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
		if len(cell.Text) > 0 {
			cell.SetText(cell.Text[:len(cell.Text)-1])
		}
		m.table.SetCell(row, column, cell)
		return nil
	}

	// Borrar todo el contenido de la celda con la tecla "Supr"
	if event.Key() == tcell.KeyDelete {
		cell.SetText("")
		m.table.SetCell(row, column, cell)
		return nil
	}

	// Seleccionar la celda de abajo al presionar Enter
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
	// Verificar si la celda previa en OPCode/Data quedó vacía
	if m.prevCol == 1 {
		prevCell := m.table.GetCell(m.prevRow, m.prevCol)
		if prevCell.Text == "" {
			prevCell.SetText("NOP").SetTextColor(tcell.ColorGreen)
			m.table.SetCell(m.prevRow, m.prevCol, prevCell)
		}
	}

	// Actualizar las referencias de fila y columna seleccionadas
	m.prevRow = newRow
	m.prevCol = newCol
}

// GetTable devuelve el componente gráfico de la tabla.
func (m *MemoryTable) GetTable() *tview.Table {
	return m.table
}

// UI encapsula la estructura gráfica del programa.
type UI struct {
	Grid   *tview.Grid
	Table  *MemoryTable
	Info *tview.TextView
}

// NewUI crea una nueva interfaz de usuario con una tabla y un widget lateral.
func NewUI(rows int) *UI {
	ui := &UI{}

	// Crear tabla de memoria
	ui.Table = NewMemoryTable(rows)

// Crear widget lateral ajustado
ui.Info = tview.NewTextView().
	SetDynamicColors(true).
	SetWrap(false).
	SetTextAlign(tview.AlignCenter).
	SetText(asciiTitle)//.
	//SetRegions(true).
	//Highlight()

ui.Info.Box.SetBorder(true).
	SetTitle("Info").
	SetTitleAlign(tview.AlignCenter)

// Crear grid para organizar los elementos
ui.Grid = tview.NewGrid().
	SetRows(0).        // Una sola fila dinámica
	SetColumns(0, 48). // Ajustar el ancho de la segunda columna para contener el título
	AddItem(ui.Table.GetTable(), 0, 0, 1, 1, 0, 0, true).
	AddItem(ui.Info, 0, 1, 1, 1, 0, 0, false)



	return ui
}
