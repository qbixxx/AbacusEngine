
package filemanager

import (
    "abacus_engine/internal/core"
    "encoding/csv"
    "fmt"
    "os"
    "path/filepath"
    "strings"
	"strconv"
	"unicode"
    "github.com/rivo/tview"
)

type FileManager struct {
    Path string
}

func NewFileManager(path string) *FileManager {
    return &FileManager{Path: path}
}

// LoadProgram carga el archivo CSV como un programa lógico.
func (fm *FileManager) LoadProgram() (core.Program, error) {
	file, err := os.Open(fm.Path)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV: %w", err)
	}

	var program core.Program
	var prevAddr int = -1

	for i, row := range rows {
		if len(row) < 1 || strings.TrimSpace(row[0]) == "" {
			continue // ignorar filas vacías
		}

		addrStr := strings.ToUpper(strings.TrimSpace(row[0]))
		if len(addrStr) != 3 || !isHex(addrStr) {
			return nil, fmt.Errorf("dirección inválida en línea %d: %s", i+1, addrStr)
		}

		addrVal, err := strconv.ParseInt(addrStr, 16, 64)
		if err != nil || addrVal < 0x000 || addrVal > 0xFFF {
			return nil, fmt.Errorf("dirección fuera de rango (000–FFF) en línea %d: %s", i+1, addrStr)
		}

		if prevAddr != -1 && addrVal != int64(prevAddr+1) {
			return nil, fmt.Errorf("direcciones no contiguas en línea %d: %s (esperado: %03X)", i+1, addrStr, prevAddr+1)
		}
		prevAddr = int(addrVal)

		var data string
		var comment string

		if len(row) > 1 {
			data = strings.ToUpper(strings.TrimSpace(row[1]))
			if data != "" && (!isHex(data) || len(data) != 4) {
				return nil, fmt.Errorf("dato inválido en línea %d: %s", i+1, data)
			}
		}
		if len(row) > 2 {
			comment = strings.TrimSpace(row[2])
		}

        		// 🛠 Si no hay dato, se asigna NOP
		if data == "" {
			data = "NOP"
		}

		cell := core.MemoryCell{
			Address: addrStr,
			Data:    data,
			Comment: comment,
		}
		program = append(program, cell)
	}

	if len(program) == 0 {
		return nil, fmt.Errorf("el archivo no contiene celdas válidas")
	}

	return program, nil
}


// SaveProgram guarda el programa lógico en un archivo CSV.
func (fm *FileManager) SaveProgram(program core.Program) error {
    file, err := os.Create(fm.Path)
    if err != nil {
        return fmt.Errorf("error creando archivo: %w", err)
    }
    defer file.Close()

    w := csv.NewWriter(file)
    w.Comma = ';'

    for _, cell := range program {
        record := []string{cell.Address, cell.Data, cell.Comment}
        if err := w.Write(record); err != nil {
            return fmt.Errorf("error escribiendo CSV: %w", err)
        }
    }

    w.Flush()
    return w.Error()
}

// BuildCSVTree genera un árbol de archivos CSV desde un directorio raíz.
func (fm *FileManager) BuildCSVTree(root string, onSelect func(string)) *tview.TreeNode {
    rootNode := tview.NewTreeNode(root).SetReference(root).SetExpanded(true)
    addCSVChildren(rootNode, root, onSelect)
    return rootNode
}

func addCSVChildren(node *tview.TreeNode, path string, onSelect func(string)) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return
    }

    for _, entry := range entries {
        fullPath := filepath.Join(path, entry.Name())
        if entry.IsDir() {
            dirNode := tview.NewTreeNode(entry.Name() + "/").
                SetReference(fullPath).
                SetExpanded(false).
                SetColor(tview.Styles.SecondaryTextColor)
            dirNode.SetSelectedFunc(func() {
                if len(dirNode.GetChildren()) == 0 {
                    addCSVChildren(dirNode, fullPath, onSelect)
                }
                dirNode.SetExpanded(!dirNode.IsExpanded())
            })
            node.AddChild(dirNode)
        } else if strings.HasSuffix(entry.Name(), ".csv") {
            fileNode := tview.NewTreeNode(entry.Name()).
                SetReference(fullPath).
                SetSelectedFunc(func() {
                    onSelect(fullPath)
                })
            node.AddChild(fileNode)
        }
    }
}

func isHex(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && !(r >= 'A' && r <= 'F') && !(r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}