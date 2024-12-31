
package interpreter

import (
	"abacus_engine/internal/memory"
	"fmt"
)

type Interpreter struct {
	Memory      *memory.Memory
	Accumulator int
	LoadAddress int
	Running     bool
}

func NewInterpreter(memory *memory.Memory) *Interpreter {
	return &Interpreter{
		Memory:      memory,
		Accumulator: 0,
		LoadAddress: -1,
		Running:     false,
	}
}

func (i *Interpreter) SetLoadAddress(address int) error {
	if address < 0 || address >= i.Memory.Size {
		return fmt.Errorf("Dirección inválida")
	}
	i.LoadAddress = address
	return nil
}

func (i *Interpreter) Step() error {
	if i.LoadAddress == -1 {
		return fmt.Errorf("Dirección de carga no inicializada")
	}
	opcode, data := i.Memory.ReadInstruction(i.LoadAddress)
	switch opcode {
	case 0x0: // Cargar inmediato
		i.Accumulator = data
	case 0x1: // Cargar de memoria
		i.Accumulator = i.Memory.Read(data)
	case 0x2: // Almacenar en memoria
		i.Memory.Write(data, i.Accumulator)
	case 0x3: // Sumar al acumulador
		i.Accumulator += data
	case 0x5: // Saltar si AC == 0
		if i.Accumulator == 0 {
			i.LoadAddress = data
			return nil
		}
	case 0xF: // Fin del programa
		i.Running = false
		return nil
	default:
		return fmt.Errorf("Opcode desconocido: %x", opcode)
	}
	i.LoadAddress++
	return nil
}

func (i *Interpreter) Run() {
	i.Running = true
	for i.Running {
		if err := i.Step(); err != nil {
			fmt.Println("Error: ", err)
			i.Running = false
		}
	}
}
