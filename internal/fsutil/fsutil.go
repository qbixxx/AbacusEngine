// internal/fsutil/fsutil.go
package fsutil

import (
	"os"
	"path/filepath"
	"github.com/rivo/tview"
)

// BuildCSVTree genera un TreeNode con archivos CSV a partir de root
func BuildCSVTree(root string, onFileSelected func(string)) *tview.TreeNode {
	rootNode := tview.NewTreeNode(root).SetReference(root).SetExpanded(true)
	addCSVChildren(rootNode, root, onFileSelected)
	return rootNode
}

func addCSVChildren(node *tview.TreeNode, path string, onFileSelected func(string)) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			dirNode := tview.NewTreeNode(entry.Name() + "/").SetReference(fullPath).SetExpanded(false)
			dirNode.SetColor(tview.Styles.SecondaryTextColor)
			dirNode.SetSelectedFunc(func() {
				if len(dirNode.GetChildren()) == 0 {
					addCSVChildren(dirNode, fullPath, onFileSelected)
				}
				dirNode.SetExpanded(!dirNode.IsExpanded())
			})
			node.AddChild(dirNode)
		} else if filepath.Ext(entry.Name()) == ".csv" {
			fileNode := tview.NewTreeNode(entry.Name()).SetReference(fullPath)
			fileNode.SetSelectedFunc(func() {
				onFileSelected(fullPath)
			})
			node.AddChild(fileNode)
		}
	}
}