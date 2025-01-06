package interpreter

import (
	"abacus_engine/internal/ui"
	"fmt"
	"strconv"
)

// Interpreter representa al intérprete que lee directamente desde la tabla de memoria.
type Interpreter struct {
	memoryTable        *ui.MemoryTable
	initAddress        int
	instructionPointer int
	accumulator        int
}

func (i *Interpreter) GetRIP() int{
	return i.instructionPointer
}
// NewInterpreter crea una nueva instancia del intérprete conectado a la MemoryTable.
func NewInterpreter(memoryTable *ui.MemoryTable) *Interpreter {
	return &Interpreter{
		memoryTable:        memoryTable,
		initAddress:        -1,
		instructionPointer: -1,
		accumulator:        0,
	}
}

func (i *Interpreter) IsRunnable() bool {
	return i.memoryTable.GetSize() > i.initAddress && i.initAddress >= 0
}

func (i *Interpreter) SetInitAddress(address int) {
	i.initAddress = address
}

func (i *Interpreter) Step() {
	if !i.IsRunnable(){
		return
	}
	// Obtener la instrucción completa desde la memoria
	instruction := i.memoryTable.GetInstruction(i.instructionPointer)

	// Separar el opcode y los datos
	opcode := instruction[0] // Primer carácter (opcode)
	data := instruction[1:]  // Últimos tres caracteres (data)

	numericData64, _ := strconv.ParseInt(data, 16, 0)
	numericData := int(numericData64)
	// Ejecutar según el opcode
	switch opcode {
	case 'N':
		i.instructionPointer++

	case '0': // Carga inmediata
		i.accumulator = numericData
		i.instructionPointer++

	case '1': // Carga
		d := i.memoryTable.GetInstruction(numericData)
		var n int
		fmt.Sscanf(d, "%d", &n)
		i.accumulator = n
		i.instructionPointer++

	case '2': // Guardado - escritura
		i.memoryTable.WriteCell(numericData+1, i.accumulator)
		i.instructionPointer++

	case '3': // Suma
		d := i.memoryTable.GetInstruction(numericData)
		var n int
		fmt.Sscanf(d, "%d", &n)
		i.accumulator += n
		i.instructionPointer++
	
	case '4': // NOT
		i.accumulator = i.accumulator * -1
		i.instructionPointer++
	
	case '5':// No definido
	case '6':// No definido
	
	case '7': // Bifurcacion si AC == 0
		if i.accumulator == 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 0)
		} else {
			i.instructionPointer++
		}
	case '8': // Bifurcacion si AC > 0
		if i.accumulator > 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 0)
		} else {
			i.instructionPointer++
		}

	case '9':  // Bifurcacion si AC < 0
		if i.accumulator < 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 0)
		} else {
			i.instructionPointer++
		}

	case 'F': // Fin del programa
		i.instructionPointer = -1
		i.initAddress = -1
		return
	default:
		i.instructionPointer++
	}

	// Avanzar el puntero de instrucción
	//i.instructionPointer++

	// Desplazar la tabla a la fila actual
	i.memoryTable.ScrollToCurrentRow(i.instructionPointer)
}

func (i *Interpreter) SetForDebug() {
	if i.IsRunnable() {

		i.instructionPointer = i.initAddress
		i.memoryTable.Goto(i.instructionPointer, 1)
		i.accumulator = 0

	}
}

// Clean limpia el interprete sin afectar su habilitacion
func (i *Interpreter) Clean() {
	i.instructionPointer = -1
	i.accumulator = 0
}

// Reset reinicia el estado del intérprete y la tabla
func (i *Interpreter) Reset() {
	i.instructionPointer = -1
	i.accumulator = 0
	i.initAddress = -1
	i.memoryTable.ResetTable()
	//i.memoryTable.ColorInitAddr(i.initAddress)
}

// GetState devuelve el estado actual del intérprete.
func (i *Interpreter) GetState() (int, int, int, bool) {
	return i.instructionPointer, i.accumulator, i.initAddress, i.IsRunnable()
}
