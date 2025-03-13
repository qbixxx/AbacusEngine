package interpreter

import (
	"abacus_engine/internal/memory"
	"fmt"
	"strconv"
)

// Interpreter representa al intérprete que lee directamente desde la tabla de memoria.
type Interpreter struct {
	memoryTable        *memory.MemoryTable
	initAddress        int
	instructionPointer int
	accumulator        int
}

func (i *Interpreter) GetRIP() int {
	return i.instructionPointer
}

// NewInterpreter crea una nueva instancia del intérprete conectado a la MemoryTable.
func NewInterpreter(memoryTable *memory.MemoryTable) *Interpreter {
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
	if !i.IsRunnable() {
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

	case '1': // Carga desde memoria
		d := i.memoryTable.GetInstruction(numericData)  // Obtener el valor de la celda
		numericValue, err := strconv.ParseInt(d, 16, 0) // Interpretar como hexadecimal
		if err != nil {
			fmt.Printf("Error al convertir el valor %s a número: %v\n", d, err)
			i.accumulator = 0 // Opcional: Establecer acumulador a 0 en caso de error
		} else {
			i.accumulator = int(numericValue) // Convertir a entero y asignar al acumulador
		}
		i.instructionPointer++

	case '2': // Guardado - escritura
		i.memoryTable.WriteCell(numericData+1, i.accumulator)
		i.instructionPointer++

	case '3': // Suma
		d := i.memoryTable.GetInstruction(numericData)
		numericValue, _ := strconv.ParseInt(d, 16, 0) // Interpretar como hexadecimal
		i.accumulator += int(numericValue)
		i.instructionPointer++

	case '4': // NOT
		i.accumulator = i.accumulator * -1
		i.instructionPointer++

	case '5': // No definido
		i.instructionPointer++

	case '6': // No definido
		i.instructionPointer++

	case '7': // Bifurcacion si AC == 0
		if i.accumulator == 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 1)
		} else {
			i.instructionPointer++
		}
	case '8': // Bifurcacion si AC > 0
		if i.accumulator > 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 1)
		} else {
			i.instructionPointer++
		}

	case '9': // Bifurcacion si AC < 0
		if i.accumulator < 0 {
			i.instructionPointer = numericData
			i.memoryTable.Goto(i.instructionPointer-1, 1)
		} else {
			i.instructionPointer++
		}

	case 'F': // Fin del programa
		i.instructionPointer = -1
		i.initAddress = -1

	default:
		i.instructionPointer++
	}

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
}

// GetState devuelve el estado actual del intérprete.
func (i *Interpreter) GetState() (int, int, int, bool) {
	return i.instructionPointer, i.accumulator, i.initAddress, i.IsRunnable()
}
