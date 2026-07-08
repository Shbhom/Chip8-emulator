package chip8

import (
	"fmt"
	"io"
	"os"
)

func (cp *Chip8) Cycle() {

	if debug {
		fmt.Printf(
			"PC=%03X Opcode=%04X SP=%d\n",
			cp.PC,
			cp.Opcode,
			cp.SP,
		)
	}

	switch cp.Opcode & 0xF0FF {
	case 0xE09E:
		fmt.Printf("SKP  PC=%03X\n", cp.PC-2)

	case 0xE0A1:
		fmt.Printf("SKNP PC=%03X\n", cp.PC-2)

	case 0xF00A:
		fmt.Printf("WAITKEY PC=%03X\n", cp.PC-2)
	}

	high := uint16(cp.Memory[cp.PC])
	low := uint16(cp.Memory[cp.PC+1])
	cp.Opcode = high<<8 | low
	cp.PC += 2
	cp.table[(cp.Opcode&0xF000)>>12]()
}

func (cp *Chip8) TickTimers() {
	if cp.DelayTimer > 0 {
		cp.DelayTimer--
	}
	if cp.SoundTimer > 0 {
		cp.SoundTimer--
	}
}

func (cp *Chip8) Table0() {
	cp.table0[cp.Opcode&0x000F]()
}

func (cp *Chip8) Table8() {
	cp.table8[cp.Opcode&0x000F]()
}

func (cp *Chip8) TableE() {
	cp.tableE[cp.Opcode&0x000F]()
}

func (cp *Chip8) TableF() {
	cp.tableF[cp.Opcode&0x00FF]()
}

func (c8 *Chip8) ReadRom(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error while fetching stats, %w", err)
	}
	romSize := stat.Size()
	buffer := make([]byte, romSize)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error while seek op, %w", err)
	}
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("error while read operation, %w", err)
	}

	var i uint = 0

	for i < uint(romSize) {
		c8.Memory[START_ADDR+i] = buffer[i]
		i++
	}

	fmt.Printf("%02X %02X %02X %02X\n",
		c8.Memory[0x200],
		c8.Memory[0x201],
		c8.Memory[0x202],
		c8.Memory[0x203],
	)
	return nil
}

func (c *Chip8) IsBuzzing() bool {
	return c.SoundTimer > 0
}

// NULL function for funciton pointer table
func (c8 *Chip8) OP_NULL() {}
