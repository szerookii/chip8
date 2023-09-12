package chip8

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
)

type CPU struct {
	Memory  [4096]uint8
	Display *Display
	V       [16]uint8
	I       uint16
	stack   [16]uint16
	SP      uint8
	DT      uint8
	ST      uint8
	PC      uint16 // program counter
}

func NewCPU() *CPU {
	cpu := new(CPU)
	display := &Display{Cpu: cpu}

	display.Init()

	cpu.Display = display

	return cpu
}

func (cpu *CPU) Init() {
	for i := 0; i < MemorySize; i++ {
		cpu.Memory[i] = 0x0
	}

	for i := 0; i < 16; i++ {
		cpu.V[i] = 0x0
		cpu.stack[i] = 0x0
	}

	cpu.PC = BaseAddress
	cpu.SP = 0x0
	cpu.DT = 0x0
	cpu.ST = 0x0
	cpu.I = 0x0

	cpu.LoadFont()
}

func (cpu *CPU) LoadFont() {
	cpu.Memory[0] = 0xF0
	cpu.Memory[1] = 0x90
	cpu.Memory[2] = 0x90
	cpu.Memory[3] = 0x90
	cpu.Memory[4] = 0xF0 // O
	cpu.Memory[5] = 0x20
	cpu.Memory[6] = 0x60
	cpu.Memory[7] = 0x20
	cpu.Memory[8] = 0x20
	cpu.Memory[9] = 0x70 // 1
	cpu.Memory[10] = 0xF0
	cpu.Memory[11] = 0x10
	cpu.Memory[12] = 0xF0
	cpu.Memory[13] = 0x80
	cpu.Memory[14] = 0xF0 // 2
	cpu.Memory[15] = 0xF0
	cpu.Memory[16] = 0x10
	cpu.Memory[17] = 0xF0
	cpu.Memory[18] = 0x10
	cpu.Memory[19] = 0xF0 // 3
	cpu.Memory[20] = 0x90
	cpu.Memory[21] = 0x90
	cpu.Memory[22] = 0xF0
	cpu.Memory[23] = 0x10
	cpu.Memory[24] = 0x10 // 4
	cpu.Memory[25] = 0xF0
	cpu.Memory[26] = 0x80
	cpu.Memory[27] = 0xF0
	cpu.Memory[28] = 0x10
	cpu.Memory[29] = 0xF0 // 5
	cpu.Memory[30] = 0xF0
	cpu.Memory[31] = 0x80
	cpu.Memory[32] = 0xF0
	cpu.Memory[33] = 0x90
	cpu.Memory[34] = 0xF0 // 6
	cpu.Memory[35] = 0xF0
	cpu.Memory[36] = 0x10
	cpu.Memory[37] = 0x20
	cpu.Memory[38] = 0x40
	cpu.Memory[39] = 0x40 // 7
	cpu.Memory[40] = 0xF0
	cpu.Memory[41] = 0x90
	cpu.Memory[42] = 0xF0
	cpu.Memory[43] = 0x90
	cpu.Memory[44] = 0xF0 // 8
	cpu.Memory[45] = 0xF0
	cpu.Memory[46] = 0x90
	cpu.Memory[47] = 0xF0
	cpu.Memory[48] = 0x10
	cpu.Memory[49] = 0xF0 // 9
	cpu.Memory[50] = 0xF0
	cpu.Memory[51] = 0x90
	cpu.Memory[52] = 0xF0
	cpu.Memory[53] = 0x90
	cpu.Memory[54] = 0x90 // A
	cpu.Memory[55] = 0xE0
	cpu.Memory[56] = 0x90
	cpu.Memory[57] = 0xE0
	cpu.Memory[58] = 0x90
	cpu.Memory[59] = 0xE0 // B
	cpu.Memory[60] = 0xF0
	cpu.Memory[61] = 0x80
	cpu.Memory[62] = 0x80
	cpu.Memory[63] = 0x80
	cpu.Memory[64] = 0xF0 // C
	cpu.Memory[65] = 0xE0
	cpu.Memory[66] = 0x90
	cpu.Memory[67] = 0x90
	cpu.Memory[68] = 0x90
	cpu.Memory[69] = 0xE0 // D
	cpu.Memory[70] = 0xF0
	cpu.Memory[71] = 0x80
	cpu.Memory[72] = 0xF0
	cpu.Memory[73] = 0x80
	cpu.Memory[74] = 0xF0 // E
	cpu.Memory[75] = 0xF0
	cpu.Memory[76] = 0x80
	cpu.Memory[77] = 0xF0
	cpu.Memory[78] = 0x80
	cpu.Memory[79] = 0x80 // F
}

func (cpu *CPU) SetRom(path string) {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	stats, err := file.Stat()

	if err != nil {
		panic(err)
	}

	size := stats.Size()

	if size > MemorySize-BaseAddress {
		panic(errors.New("corrupted program"))
	}

	bytes := make([]byte, size)
	reader := bufio.NewReader(file)
	_, err = reader.Read(bytes)
	if err != nil {
		panic(err)
	}

	copy(cpu.Memory[BaseAddress:], bytes)
}

func (cpu *CPU) ReadOpcode() uint16 {
	return (uint16(cpu.Memory[cpu.PC]) << 8) | uint16(cpu.Memory[cpu.PC+1])
}

func (cpu *CPU) Exec(opcode uint16) {
	inst, args := Disassemble(opcode)

	//fmt.Println(inst)

	switch inst {
	case "CLS":
		cpu.Display.Clear()
		cpu.PC += 2
		break

	case "RET":
		cpu.PC = cpu.stack[cpu.SP]
		cpu.SP--
		break

	case "JP_ADDR":
		cpu.PC = args[0]
		break

	case "CALL_ADDR":
		cpu.SP++
		cpu.stack[cpu.SP] = cpu.PC + 2
		cpu.PC = args[0]
		break

	case "SE_VX_NN":
		if uint16(cpu.V[args[0]]) == args[1] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
		break

	case "SNE_VX_NN":
		if uint16(cpu.V[args[0]]) != args[1] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
		break

	case "SE_VX_VY":
		if cpu.V[args[0]] != cpu.V[args[1]] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
		break

	case "LD_VX_NN":
		cpu.V[args[0]] = uint8(args[1])
		cpu.PC += 2
		break

	case "ADD_VX_NN":
		v := cpu.V[args[0]] + uint8(args[1])

		if v > 255 {
			v -= 255
		}

		cpu.V[args[0]] = v
		cpu.PC += 2

		break

	case "LD_VX_VY":
		cpu.V[args[0]] = cpu.V[args[1]]
		cpu.PC += 2

		break

	case "OR_VX_VY":
		cpu.V[args[0]] |= cpu.V[args[1]]
		cpu.PC += 2
		break

	case "AND_VX_VY":
		cpu.V[args[0]] &= cpu.V[args[1]]
		cpu.PC += 2
		break

	case "XOR_VX_VY":
		cpu.V[args[0]] ^= cpu.V[args[1]]
		cpu.PC += 2
		break

	case "ADD_VX_VY":
		if cpu.V[args[0]]+cpu.V[args[1]] > 0xFF {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}

		cpu.V[args[0]] += cpu.V[args[1]]
		cpu.PC += 2
		break

	case "SUB_VX_VY":
		if cpu.V[args[0]] > cpu.V[args[1]] {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}

		cpu.V[args[0]] -= cpu.V[args[1]]
		cpu.PC += 2
		break

	case "SUBN_VX_VY":
		cpu.V[0xF] = cpu.V[args[0]] >> 7
		cpu.V[args[0]] <<= 1
		cpu.PC += 2
		break

	case "SNE_VX_VY":
		if cpu.V[args[0]] != cpu.V[args[1]] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
		break

	case "LD_I_ADDR":
		cpu.I = args[1]
		cpu.PC += 2
		break

	case "JP_V0_ADDR":
		cpu.PC = uint16(cpu.V[args[0]]) + args[1]
		break

	case "RND_VX_NN":
		cpu.V[args[0]] = uint8(rand.Intn(255)) & uint8(args[1])
		cpu.PC += 2
		break

	case "DRW_VX_VY_N":
		cpu.Display.Draw(args[0], args[1], args[2])
		cpu.PC += 2
		break

	case "SKP_VX":
		isKeyPressed := cpu.V[args[0]] == 1

		if isKeyPressed {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}

		break

	case "SKNP_VX":
		isKeyPressed := cpu.V[args[0]] == 1

		if !isKeyPressed {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}

		break

	case "LD_VX_DT":
		cpu.V[args[0]] = cpu.DT
		cpu.PC += 2
		break

	case "LD_VX_N":
		// ToDo : Wait for key press
		break

	case "LD_DT_VX":
		cpu.DT = cpu.V[args[1]]
		cpu.PC += 2
		break

	case "LD_ST_VX":
		cpu.ST = cpu.V[args[1]]

		if cpu.ST > 0 {
			// ToDo : Enable sound
		}

		cpu.PC += 2
		break

	case "ADD_I_VX":
		cpu.I += uint16(cpu.V[args[1]])
		cpu.PC += 2
		break

	case "LD_F_VX":
		cpu.I = uint16(cpu.V[args[1]] * 5)
		cpu.PC += 2
		break

	case "LD_B_VX":
		cpu.Memory[cpu.I] = (cpu.V[args[1]] - cpu.V[args[1]]%100) / 100
		cpu.Memory[cpu.I+1] = ((cpu.V[args[1]] - cpu.V[args[1]]%10) / 10) % 10
		cpu.Memory[cpu.I+2] = cpu.V[args[1]] - cpu.Memory[cpu.I]*100 - 10*cpu.Memory[cpu.I+1]
		cpu.PC += 2
		break

	case "LD_I_VX":
		for i := uint16(0); i < args[1]; i++ {
			cpu.Memory[cpu.I+i] = cpu.V[i]
		}

		cpu.PC += 2
		break

	case "LD_VX_I":
		for i := uint16(0); i < args[1]; i++ {
			cpu.V[i] = cpu.Memory[cpu.I+i]
		}

		cpu.PC += 2
		break
	}
}

func (cpu *CPU) Run() {
	if cpu.DT > 0 {
		cpu.DT--
	}

	if cpu.ST > 0 {
		cpu.ST--
	}
}
