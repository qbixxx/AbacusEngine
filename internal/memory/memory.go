// internal/memory/memory.go


package memory


import (
	"abacus_engine/internal/styles"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"fmt"

)

// MemoryTable encapsula la lógica y el comportamiento de la tabla de memoria.
type MemoryTable struct {
	table        *tview.Table                // Componente gráfico de la tabla
	rows         int                         // Cantidad de filas
	prevRow      int                         // Fila previamente seleccionada
	prevCol      int                         // Columna previamente seleccionada
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

			newCell := styles.MakeCell(cell.Text)
			m.table.SetCell(row, column, newCell)
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


func (m *MemoryTable) WriteCell(row int, accumulator int) {
	cell := m.table.GetCell(row, 1) // +1 porque la fila 0 es para encabezados

	// Formatear el acumulador como string con ceros completados
	formattedData := fmt.Sprintf("%04X", accumulator)
	cell.SetText(formattedData)

	m.table.SetCell(row, 1, cell) // +1 porque la fila 0 es el encabezado
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



