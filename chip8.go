package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	BaseAddress = 0x200
	MemorySize  = 4096
	l           = 64
	L           = 32
	NbrOpcode = 35
)

type JMP struct {
	mask [NbrOpcode]uint16
	ids [NbrOpcode]uint16
}

func (jp *JMP) Init() {
	jp.mask[0]= 0x0000; jp.ids[0]=0x0FFF          /* 0NNN */
	jp.mask[1]= 0xFFFF; jp.ids[1]=0x00E0          /* 00E0 */
	jp.mask[2]= 0xFFFF; jp.ids[2]=0x00EE          /* 00EE */
	jp.mask[3]= 0xF000; jp.ids[3]=0x1000          /* 1NNN */
	jp.mask[4]= 0xF000; jp.ids[4]=0x2000          /* 2NNN */
	jp.mask[5]= 0xF000; jp.ids[5]=0x3000          /* 3XNN */
	jp.mask[6]= 0xF000; jp.ids[6]=0x4000          /* 4XNN */
	jp.mask[7]= 0xF00F; jp.ids[7]=0x5000          /* 5XY0 */
	jp.mask[8]= 0xF000; jp.ids[8]=0x6000          /* 6XNN */
	jp.mask[9]= 0xF000; jp.ids[9]=0x7000          /* 7XNN */
	jp.mask[10]= 0xF00F; jp.ids[10]=0x8000        /* 8XY0 */
	jp.mask[11]= 0xF00F; jp.ids[11]=0x8001        /* 8XY1 */
	jp.mask[12]= 0xF00F; jp.ids[12]=0x8002        /* 8XY2 */
	jp.mask[13]= 0xF00F; jp.ids[13]=0x8003        /* BXY3 */
	jp.mask[14]= 0xF00F; jp.ids[14]=0x8004        /* 8XY4 */
	jp.mask[15]= 0xF00F; jp.ids[15]=0x8005        /* 8XY5 */
	jp.mask[16]= 0xF00F; jp.ids[16]=0x8006        /* 8XY6 */
	jp.mask[17]= 0xF00F; jp.ids[17]=0x8007        /* 8XY7 */
	jp.mask[18]= 0xF00F; jp.ids[18]=0x800E        /* 8XYE */
	jp.mask[19]= 0xF00F; jp.ids[19]=0x9000        /* 9XY0 */
	jp.mask[20]= 0xF000; jp.ids[20]=0xA000        /* ANNN */
	jp.mask[21]= 0xF000; jp.ids[21]=0xB000        /* BNNN */
	jp.mask[22]= 0xF000; jp.ids[22]=0xC000        /* CXNN */
	jp.mask[23]= 0xF000; jp.ids[23]=0xD000        /* DXYN */
	jp.mask[24]= 0xF0FF; jp.ids[24]=0xE09E        /* EX9E */
	jp.mask[25]= 0xF0FF; jp.ids[25]=0xE0A1        /* EXA1 */
	jp.mask[26]= 0xF0FF; jp.ids[26]=0xF007        /* FX07 */
	jp.mask[27]= 0xF0FF; jp.ids[27]=0xF00A        /* FX0A */
	jp.mask[28]= 0xF0FF; jp.ids[28]=0xF015        /* FX15 */
	jp.mask[29]= 0xF0FF; jp.ids[29]=0xF018        /* FX18 */
	jp.mask[30]= 0xF0FF; jp.ids[30]=0xF01E        /* FX1E */
	jp.mask[31]= 0xF0FF; jp.ids[31]=0xF029        /* FX29 */
	jp.mask[32]= 0xF0FF; jp.ids[32]=0xF033        /* FX33 */
	jp.mask[33]= 0xF0FF; jp.ids[33]=0xF055        /* FX55 */
	jp.mask[34]= 0xF0FF; jp.ids[34]=0xF065        /* FX65 */
}

type CPU struct {
	memory [4096]uint8
	v [16]uint8
	i uint16
	jmp [16]uint16
	jmpNbr uint8
	gameClock uint8
	soundClock uint8
	pc uint16 // program counter
}

func (cpu *CPU) Init() {
	for i := 0; i < MemorySize; i++ {
		cpu.memory[i] = 0x0
	}

	for i := 0; i < 16; i++ {
		cpu.v[i] = 0x0
		cpu.jmp[i] = 0x0
	}

	cpu.pc = BaseAddress
	cpu.jmpNbr = 0x0
	cpu.gameClock = 0x0
	cpu.soundClock = 0x0
	cpu.i = 0x0

	cpu.LoadFont()
}

func (cpu *CPU) LoadFont() {
		cpu.memory[0] = 0xF0; cpu.memory[1] = 0x90; cpu.memory[2] = 0x90; cpu.memory[3] = 0x90; cpu.memory[4] = 0xF0 // O
		cpu.memory[5] = 0x20; cpu.memory[6] = 0x60; cpu.memory[7] = 0x20; cpu.memory[8] = 0x20; cpu.memory[9] = 0x70 // 1
		cpu.memory[10] = 0xF0; cpu.memory[11] = 0x10; cpu.memory[12] = 0xF0; cpu.memory[13] = 0x80; cpu.memory[14] = 0xF0 // 2
		cpu.memory[15] = 0xF0; cpu.memory[16] = 0x10; cpu.memory[17] = 0xF0; cpu.memory[18] = 0x10; cpu.memory[19] = 0xF0 // 3
		cpu.memory[20] = 0x90; cpu.memory[21] = 0x90; cpu.memory[22] = 0xF0; cpu.memory[23] = 0x10; cpu.memory[24] = 0x10 // 4
		cpu.memory[25] = 0xF0; cpu.memory[26] = 0x80; cpu.memory[27] = 0xF0; cpu.memory[28] = 0x10; cpu.memory[29] = 0xF0 // 5
		cpu.memory[30] = 0xF0; cpu.memory[31] = 0x80; cpu.memory[32] = 0xF0; cpu.memory[33] = 0x90; cpu.memory[34] = 0xF0 // 6
		cpu.memory[35] = 0xF0; cpu.memory[36] = 0x10; cpu.memory[37] = 0x20; cpu.memory[38] = 0x40; cpu.memory[39] = 0x40 // 7
		cpu.memory[40] = 0xF0; cpu.memory[41] = 0x90; cpu.memory[42] = 0xF0; cpu.memory[43] = 0x90; cpu.memory[44] = 0xF0 // 8
		cpu.memory[45] = 0xF0; cpu.memory[46] = 0x90; cpu.memory[47] = 0xF0; cpu.memory[48] = 0x10; cpu.memory[49] = 0xF0 // 9
		cpu.memory[50] = 0xF0; cpu.memory[51] = 0x90; cpu.memory[52] = 0xF0; cpu.memory[53] = 0x90; cpu.memory[54] = 0x90 // A
		cpu.memory[55] = 0xE0; cpu.memory[56] = 0x90; cpu.memory[57] = 0xE0; cpu.memory[58] = 0x90; cpu.memory[59] = 0xE0 // B
		cpu.memory[60] = 0xF0; cpu.memory[61] = 0x80; cpu.memory[62] = 0x80; cpu.memory[63] = 0x80; cpu.memory[64] = 0xF0 // C
		cpu.memory[65] = 0xE0; cpu.memory[66] = 0x90; cpu.memory[67] = 0x90; cpu.memory[68] = 0x90; cpu.memory[69] = 0xE0 // D
		cpu.memory[70] = 0xF0; cpu.memory[71] = 0x80; cpu.memory[72] = 0xF0; cpu.memory[73] = 0x80; cpu.memory[74] = 0xF0 // E
		cpu.memory[75] = 0xF0; cpu.memory[76] = 0x80; cpu.memory[77] = 0xF0; cpu.memory[78] = 0x80; cpu.memory[79] = 0x80 // F
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

	if size > MemorySize - BaseAddress {
		panic(errors.New("corrupted program"))
	}

	bytes := make([]byte, size)
	reader := bufio.NewReader(file)
	_, err = reader.Read(bytes)
	if err != nil {
		panic(err)
	}

	copy(cpu.memory[BaseAddress:], bytes)
}

func (cpu *CPU) ReadOpcode() uint16 {
	return (uint16(cpu.memory[cpu.pc]) << 8) | uint16(cpu.memory[cpu.pc + 1])
}

func (cpu *CPU) ReadAction(opcode uint16) uint8 {
	var action uint8
	var result uint16

	for action = 0; action < NbrOpcode; action++ {
		result = Jmp.mask[action]&opcode

		if result == Jmp.ids[action] {
			break
		}
	}

	return action
}

func (cpu *CPU) Exec(opcode uint16) {
	b4 := cpu.ReadAction(opcode)
	b3 := (opcode & (0x0F00)) >> 8 // X
	b2 := (opcode & (0x00F0)) >> 4 // Y
	b1 := opcode & (0x000F) // N

	switch b4 {
	case 0:
		// 0nnn
		break

	case 1:
		// 00E0
		Screen.Clear()
		break

	case 2:
		// 00EE
		if cpu.jmpNbr > 0 {
			cpu.jmpNbr--
			cpu.pc = cpu.jmp[cpu.jmpNbr]
		}
		break

	case 3:
		// 1nnn
		cpu.pc = (b3 << 8) + (b2 << 4) + b1
		cpu.pc -= 2
		break

	case 4:
		// 2nnn
		cpu.jmp[cpu.jmpNbr] = cpu.pc

		if cpu.jmpNbr < 15 {
			cpu.jmpNbr++
		}

		cpu.pc = (b3 << 8) + (b2 << 4) + b1
		cpu.pc -= 2
		break

	case 5:
		// 3xkk
		if uint16(cpu.v[b3]) == ((b2 << 4) + b1) {
			cpu.pc += 2
		}
		break

	case 6:
		// 4xkk
		if uint16(cpu.v[b3]) != ((b2 << 4) + b1) {
			cpu.pc += 2
		}
		break

	case 7:
		// 5xy0
		if cpu.v[b3] == cpu.v[b2] {
			cpu.pc += 2
		}
		break

	case 8:
		// 6xkk
		cpu.v[b3] = uint8((b2 << 4) + b1)
		break

	case 9:
		// 7xkk
		cpu.v[b3] += uint8((b2 << 4) + b1)
		break

	case 10:
		// 8xy0
		cpu.v[b3] = cpu.v[b2]
		break

	case 11:
		// 8xy1
		cpu.v[b3] = cpu.v[b3] | cpu.v[b2]
		break

	case 12:
		// 8xy2
		cpu.v[b3] = cpu.v[b3] & cpu.v[b2]

		break

	case 13:
		// 8xy3
		cpu.v[b3] = cpu.v[b3] ^ cpu.v[b2]
		break

	case 14:
		// 8xy4
		if (cpu.v[b3] + cpu.v[b2]) > 0xFF {
			cpu.v[0xF] = 1
		} else {
			cpu.v[0xF] = 0
		}

		cpu.v[b3] += cpu.v[b2]
		break

	case 15:
		// 8xy5
		if cpu.v[b3] < cpu.v[b2] {
			cpu.v[0xF] = 0
		} else {
			cpu.v[0xF] = 1
		}

		cpu.v[b3] -= cpu.v[b2]
		break

	case 16:
		// 8xy6
		cpu.v[0xF] = cpu.v[b3] & (0x01)
		cpu.v[b3] = cpu.v[b3] >> 1
		break

	case 17:
		// 8xy7
		if cpu.v[b2] < cpu.v[b3] {
			cpu.v[0xF] = 0
		} else {
			cpu.v[0xF] = 1
		}

		cpu.v[b3] = cpu.v[b2] - cpu.v[b3]
		break

	case 18:
		// 8xyE
		cpu.v[0xF] = cpu.v[b3] >> 7
		cpu.v[b3] = cpu.v[b3] << 1
		break

	case 19:
		// 9xy0
		if cpu.v[b3] != cpu.v[b2] {
			cpu.pc += 2
		}
		break

	case 20:
		// Annn
		cpu.i = (b3 << 8) + (b2 << 4) + b1
		break

	case 21:
		// Bnnn
		cpu.pc = (b3 << 8) + (b2 << 4) + b1 + uint16(cpu.v[0])
		cpu.pc -= 2
		break

	case 22:
		//Cxkk
		cpu.v[b3] = uint8(uint16(rand.Intn(255)) % ((b2 << 4) + b1 + 1))

		break

	case 23:
		// Dxyn
		Screen.Draw(b1, b2, b3)
		break

	case 24:
		// Ex9E
		break

	case 25:
		// ExA1
		break

	case 26:
		// Fx07
		cpu.v[b3] = cpu.gameClock
		break

	case 27:
		// Fx0A
		break

	case 28:
		// Fx15
		cpu.gameClock = cpu.v[b3]

		break

	case 29:
		// Fx18
		cpu.soundClock = cpu.v[b3]
		break

	case 30:
		// Fx1E
		if cpu.i + uint16(cpu.v[b3]) > 4095 {
			cpu.v[0xF] = 1
		} else {
			cpu.v[0xF] = 0
		}

		cpu.i += uint16(cpu.v[b3])
		break

	case 31:
		// Fx29
		cpu.i = uint16(cpu.v[b3] * 5)
		break

	case 32:
		// Fx33
		cpu.memory[cpu.i] = (cpu.v[b3] - cpu.v[b3] % 100) / 100
		cpu.memory[cpu.i + 1] = ((cpu.v[b3] - cpu.v[b3] % 10) / 10) % 10
		cpu.memory[cpu.i + 2] = cpu.v[b3] - cpu.memory[cpu.i] * 100 - cpu.memory[cpu.i + 1] * 10
		break

	case 33:
		// Fx55
		for i := 0; i <= int(b3); i++{
			cpu.memory[int(cpu.i) + i] = cpu.v[i]
		}
		break

	case 34:
		// Fx65
		for i := 0; i <= int(b3); i++{
			cpu.v[i] = cpu.memory[int(cpu.i) + i]
		}
		break
	}

	cpu.pc += 2
}

func (cpu *CPU) Run() {
	if cpu.gameClock > 0 {
		cpu.gameClock--
	}

	if cpu.soundClock > 0 {
		cpu.soundClock--
	}
}

type Pixel struct {
	X, Y int
	Color int
}

type Display struct {
	pixels []*Pixel
}

func (d *Display) Init() {
	err := termbox.Init()

	if err != nil {
		log.Fatal("Cannot initialize display")
	}

	for x := 0; x < l; x++ {
		for y := 0; y < L; y++ {
			d.SetPixel(x, y, 0)
		}
	}

	d.Update()
}

func (d *Display) SetPixel(x, y int, color int) {
	for _, p := range d.pixels {
		if p.X == x && p.Y == y {
			*p = Pixel{X: x, Y: y, Color: color}
			return
		}
	}

	d.pixels = append(d.pixels, &Pixel{X: x, Y: y, Color: color})
}

func (d *Display) GetPixel(x, y int) *Pixel {
	for _, p := range d.pixels {
		if p.X == x && p.Y == y {
			return p
		}
	}

	d.SetPixel(x, y, 0)

	return d.GetPixel(x, y)
}

func (d *Display) Clear() {
	d.pixels = []*Pixel{}
}

func (d *Display) Update() {
	for _, pixel := range d.pixels {
		if pixel.Color == 1 {
			termbox.SetCell(pixel.X, pixel.Y, ' ', termbox.ColorWhite, termbox.ColorWhite)
		} else {
			termbox.SetCell(pixel.X, pixel.Y, ' ', termbox.ColorBlack, termbox.ColorBlack)
		}
	}

	termbox.Flush()
}

func (d *Display) Draw(b1, b2, b3 uint16) {
	// ToDo
}

var (
	Cpu *CPU
	Screen *Display
	Jmp *JMP
)

func main() {
	var path string

	fmt.Println("Enter path to your ROM")
	fmt.Print("> ")

	fmt.Scan(&path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("ROM not found!")
		return
	}

	Cpu = &CPU{}
	Screen = &Display{pixels: []*Pixel{}}
	Jmp = &JMP{}

	Cpu.Init()
	Screen.Init()
	Jmp.Init()

	Cpu.SetRom(path)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
					return
				}
			}

		default:
			for i := 0; i < 4; i++ {
				Cpu.Exec(Cpu.ReadOpcode())
			}

			Screen.Update()
			Cpu.Run()
			time.Sleep(16 * time.Millisecond)
		}
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
