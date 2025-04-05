package ui

import (
	"abacus_engine/internal/core"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type MemoryTableAdapter struct {
	Table *tview.Table
}

func NewMemoryTableAdapter(table *tview.Table) *MemoryTableAdapter {
	return &MemoryTableAdapter{Table: table}
}

func (a *MemoryTableAdapter) LoadProgram(program core.Program) {
	for _, cell := range program {
		// Convertir la dirección de memoria desde string hexadecimal a entero
		addr := strings.ToUpper(cell.Address)
		row, err := strconv.ParseInt(addr, 16, 0)
		if err != nil || row < 0 {
			continue
		}
		rowIndex := int(row) + 1 // +1 porque la fila 0 es el encabezado

		a.Table.SetCell(rowIndex, 0, MakeCell(cell.Address))
		a.Table.SetCell(rowIndex, 1, MakeCell(cell.Data))
		a.Table.SetCell(rowIndex, 2, MakeCell(cell.Comment))
	}
}

// ExtractProgram extrae un core.Program desde la tabla visual.
func (a *MemoryTableAdapter) ExtractProgram() core.Program {
	var program core.Program
	rows := a.Table.GetRowCount()

	for i := 0; i < rows; i++ {
		addr := a.Table.GetCell(i, 0)
		data := a.Table.GetCell(i, 1)
		comm := a.Table.GetCell(i, 2)

		program = append(program, core.MemoryCell{
			Address: getCellText(addr),
			Data:    getCellText(data),
			Comment: getCellText(comm),
		})
	}

	return program
}

func getCellText(cell *tview.TableCell) string {
	if cell == nil {
		return ""
	}
	return cell.Text
}
