package chip8

import (
	"log"

	"github.com/nsf/termbox-go"
)

type Pixel struct {
	X, Y  int
	Color int
}

type Display struct {
	Cpu    *CPU
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
			termbox.SetCell(pixel.X, pixel.Y, ' ', termbox.ColorBlack, termbox.ColorWhite)
		} else {
			termbox.SetCell(pixel.X, pixel.Y, ' ', termbox.ColorBlack, termbox.ColorBlack)
		}
	}

	termbox.Flush()
}

func (d *Display) Draw(x, y, h uint16) {
	d.Cpu.V[0xF] = 0

	registerX := uint16(d.Cpu.V[x])
	registerY := uint16(d.Cpu.V[y])

	var spr uint16

	for y := uint16(0); y < h; y++ {
		spr = uint16(d.Cpu.Memory[d.Cpu.I+y])

		for x := uint16(0); x < 8; x++ {
			if (spr & 0x80) > 0 {
				if d.GetPixel(int(registerX+x), int(registerY+y)).Color == 1 {
					d.Cpu.V[0xF] = 1
				}

				d.SetPixel(int(registerX+x), int(registerY+y), d.GetPixel(int(registerX+x), int(registerY+y)).Color^1)
			}

			spr <<= 1
		}
	}
}
