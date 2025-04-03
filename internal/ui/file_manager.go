// internal/ui/filemanager.go
package ui

import (
	"github.com/rivo/tview"
	"abacus_engine/internal/fsutil"
	"abacus_engine/internal/csvloader"
	"github.com/gdamore/tcell/v2"

)

// NewFileManagerPage crea y retorna la vista del file manager con lógica integrada
func NewFileManagerPage(table *tview.Table, switchToMain func(), showModal func(string)) tview.Primitive {
	onSelect := func(path string) {
		err := csvloader.LoadCSVToTable(path, table)
		if err != nil {
			showModal("Error al cargar CSV:\n" + err.Error())
		} else {
			switchToMain()
		}
	}

	root := fsutil.BuildCSVTree("/", onSelect)
	root.SetExpanded(true)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle(" File Manager (CSV) ")

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			switchToMain()
			return nil
		}
		return event
	})

	return tree
}
