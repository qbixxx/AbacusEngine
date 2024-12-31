package memory


// Memory gestiona los datos en la memoria simulada.
type Memory struct {
	cells []int
	Size  int
}

// NewMemory crea una nueva instancia de memoria.
func NewMemory(size int) *Memory {
	return &Memory{
		cells: make([]int, size),
		Size:  size,
	}
}

// Read lee el valor almacenado en una dirección de memoria.
func (m *Memory) Read(address int) int {
	if address < 0 || address >= m.Size {
		panic("Dirección de memoria fuera de rango")
	}
	return m.cells[address]
}

// Write escribe un valor en una dirección de memoria.
func (m *Memory) Write(address int, value int) {
	if address < 0 || address >= m.Size {
		panic("Dirección de memoria fuera de rango")
	}
	m.cells[address] = value
}

// ReadInstruction devuelve el opcode y los datos de una dirección específica.
func (m *Memory) ReadInstruction(address int) (int, int) {
	value := m.Read(address)
	opcode := (value >> 8) & 0xFF
	data := value & 0xFF
	return opcode, data
}

