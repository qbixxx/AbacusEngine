// internal/csvloader/csvloader.go
package csvloader

import (
	"encoding/csv"
	//"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"abacus_engine/internal/memory"
	"github.com/rivo/tview"
)

// LoadCSVToTable valida y carga un archivo CSV simple con formato: direccion, dato, comentario
func LoadCSVToTable(path string, memTable *tview.Table) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("error leyendo CSV: %w", err)
	}

	type rowData struct {
		addr string
		data string
		comm string
	}
	var dataRows []rowData
	addrValues := []int{}
	addrSeen := map[int]bool{}

	for i, row := range rows {
		if len(row) < 2 {
			continue
		}

		rawAddr := strings.TrimSpace(row[0])
		rawData := strings.TrimSpace(row[1])
		comment := ""
		if len(row) >= 3 {
			comment = strings.TrimSpace(row[2])
		}

		if rawAddr == "" || rawData == "" {
			continue
		}

		addrInt, err := strconv.ParseInt(rawAddr, 16, 0)
		if err != nil || addrInt < 0x000 || addrInt > 0xFFF {
			return fmt.Errorf("dirección inválida en fila %d: %s", i+1, rawAddr)
		}
		if len(rawData) != 4 {
			return fmt.Errorf("dato inválido en fila %d (debe tener 4 dígitos hex): %s", i+1, rawData)
		}
		if _, err := strconv.ParseUint(rawData, 16, 16); err != nil {
			return fmt.Errorf("dato inválido en fila %d: %s", i+1, rawData)
		}

		addr := int(addrInt)
		if addrSeen[addr] {
			return fmt.Errorf("dirección duplicada en fila %d: %s", i+1, rawAddr)
		}
		addrSeen[addr] = true
		addrValues = append(addrValues, addr)
		dataRows = append(dataRows, rowData{
			addr: strings.ToUpper(rawAddr),
			data: strings.ToUpper(rawData),
			comm: comment,
		})
	}

	if len(addrValues) == 0 {
		return fmt.Errorf("no se encontraron datos válidos en el archivo")
	}

	// Verificar continuidad de direcciones
	min, max := addrValues[0], addrValues[0]
	for _, v := range addrValues {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	if len(addrValues) != (max - min + 1) {
		return fmt.Errorf("las direcciones de memoria no son contiguas entre %03X y %03X", min, max)
	}

	// Cargar en la tabla
	for _, row := range dataRows {
		addrInt, _ := strconv.ParseInt(row.addr, 16, 0)
		rowIndex := int(addrInt) + 1
		memTable.SetCell(rowIndex, 0, memory.MakeCell(row.addr))
		memTable.SetCell(rowIndex, 1, memory.MakeCell(row.data))
		memTable.SetCell(rowIndex, 2, memory.MakeCell(row.comm))
	}

	return nil
}
