package main

import (
	"fmt"
	"github.com/Seyz123/chip8/chip8"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

var Cpu *chip8.CPU

func main() {
	var path string

	fmt.Println("Enter path to your ROM")
	fmt.Print("> ")

	fmt.Scan(&path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("ROM not found!")
		return
	}

	Cpu = chip8.NewCPU()
	Cpu.Init()

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

			Cpu.Display.Update()
			Cpu.Run()
			time.Sleep(16 * time.Millisecond)
		}
	}
}
