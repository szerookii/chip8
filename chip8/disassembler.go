package chip8

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func Disassemble(opcode uint16) (string, []uint16) {
	instructions := GetInstructions()

	var instruction *Instruction

	for i := 0; i < NbrOpcode; i++ {
		if (instructions[i].Mask & opcode) == instructions[i].Pattern {
			instruction = instructions[i]
			break
		}
	}

	if instruction != nil {
		var args []uint16

		for _, arg := range instruction.Arguments {
			args = append(args, (opcode&arg.Mask)>>arg.Shift)
		}

		return instruction.Id, args
	}

	return "", nil
}

type Argument struct {
	Mask  uint16 `json:"mask"`
	Shift int    `json:"shift"`
}

type Instruction struct {
	Id        string      `json:"id"`
	Mask      uint16      `json:"mask"`
	Pattern   uint16      `json:"pattern"`
	Arguments []*Argument `json:"arguments"`
}

func GetInstructions() []*Instruction {
	bytes, err := ioutil.ReadFile(filepath.Join("data", "instructions.json"))

	if err != nil {
		panic(err)
	}

	var instructions []*Instruction

	json.Unmarshal(bytes, &instructions)

	return instructions
}
