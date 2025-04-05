package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MakeCell crea una celda con formato personalizado según el texto
func MakeCell(text string) *tview.TableCell {
	cell := tview.NewTableCell(text).
		SetAlign(tview.AlignCenter).
		SetSelectable(true)

	// Estilos por columna (por heurística)
	if len(text) == 3 && isHex(text) {
		// dirección de memoria
		cell.SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetSelectable(false).
			SetExpansion(0)
	} else if len(text) == 4 && isHex(text) {
		// instrucción / dato
		switch text[0] {
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
		cell.SetMaxWidth(4).SetExpansion(0)
	} else {
		// comentario o celda vacía
		cell.SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft).
			SetExpansion(18)
	}

	return cell
}

func isHex(s string) bool {
	for _, c := range s {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'F') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
