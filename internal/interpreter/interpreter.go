package interpreter

import "abacus_engine/internal/ui"

// Interpreter representa al intérprete que lee directamente desde la tabla de memoria.
type Interpreter struct {
	memoryTable        *ui.MemoryTable
	initAddress        int
	instructionPointer int
	accumulator        int
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

// Step ejecuta la instrucción en la posición actual del puntero.
func (i *Interpreter) Step() {
	instruction := i.memoryTable.GetInstruction(i.instructionPointer)
	switch instruction {
	case "NOP":
		// No hacer nada
	case "INC":
		i.accumulator++
	case "DEC":
		i.accumulator--
		// Más instrucciones según sea necesario
	}
	i.instructionPointer++
}

func (i *Interpreter) SetForDebug() {
	if i.IsRunnable(){
		i.instructionPointer = i.initAddress
		i.memoryTable.Goto(i.instructionPointer, 1)
		i.accumulator = 0
	}
}

// Reset reinicia el estado del intérprete.
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
