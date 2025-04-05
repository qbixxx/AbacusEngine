package ui

import (
	"abacus_engine/internal/filemanager"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewFileManagerPage crea y retorna la vista del file manager con lógica integrada
func NewFileManagerPage(adapter *MemoryTableAdapter, switchToMain func(), showModal func(string)) tview.Primitive {
	onSelect := func(path string) {
		fm := filemanager.NewFileManager(path)
		program, err := fm.LoadProgram()
		if err != nil {
			showModal("Error al cargar CSV:\n" + err.Error())
			return
		}
		adapter.LoadProgram(program)
		switchToMain()
	}

	fm := &filemanager.FileManager{}
	root := fm.BuildCSVTree("/", onSelect)
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
