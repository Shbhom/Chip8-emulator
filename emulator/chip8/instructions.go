package chip8

// 00E0: CLS
// By setting 0 we are turning off all the pixels in the display
func (c8 *Chip8) OP_00E0() {
	for i := range c8.Video {
		c8.Video[i] = 0
	}
}

// 00EE: RET
// decreases the Stack pointer and updates the PC value
func (c8 *Chip8) OP_00EE() {
	c8.SP--
	c8.PC = c8.Stack[c8.SP]
}

// 1nnn: JP addr
// Jump to location nnn
func (c8 *Chip8) OP_1nnn() {
	address := c8.Opcode & 0x0FFF
	c8.PC = address
}

// 2nnn - CALL addr
// Call subroutine at nnn.
func (c8 *Chip8) OP_2nnn() {
	address := c8.Opcode & 0x0FFF
	c8.Stack[c8.SP] = c8.PC
	c8.SP++
	c8.PC = address
}

// Instruction is split into 3 sections OPcode Register Imediate Value

// 3xkk - SE Vx, byte
// skip if equal instruction
func (c8 *Chip8) OP_3xkk() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	value := uint8(c8.Opcode & 0x00FF)

	if c8.Registers[vx] == value {
		c8.PC += 2
	}
}

// 4xkk - SNE Vx, byte
// Skip to next instrucation if vx != value
func (c8 *Chip8) OP_4xkk() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	value := uint8(c8.Opcode & 0x00FF)

	if c8.Registers[vx] != value {
		c8.PC += 2
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
func (c8 *Chip8) OP_5xy0() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	vy := uint8((c8.Opcode & 0x00F0) >> 4)

	if c8.Registers[vx] == c8.Registers[vy] {
		c8.PC += 2
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
func (c8 *Chip8) OP_6xkk() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	value := uint8(c8.Opcode & 0x00FF)

	c8.Registers[vx] = value
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
func (c8 *Chip8) OP_7xkk() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	value := uint8(c8.Opcode & 0x00FF)

	c8.Registers[vx] += value
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
func (c8 *Chip8) OP_8xy0() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	vy := uint8((c8.Opcode & 0x00F0) >> 4)

	c8.Registers[vx] = c8.Registers[vy]
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
func (c8 *Chip8) OP_8xy1() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	vy := uint8((c8.Opcode & 0x00F0) >> 4)

	c8.Registers[vx] |= c8.Registers[vy]
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
func (c8 *Chip8) OP_8xy2() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	vy := uint8((c8.Opcode & 0x00F0) >> 4)

	c8.Registers[vx] &= c8.Registers[vy]
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy
func (c8 *Chip8) OP_8xy3() {
	vx := uint8((c8.Opcode & 0x0F00) >> 8)
	vy := uint8((c8.Opcode & 0x00F0) >> 4)

	c8.Registers[vx] ^= c8.Registers[vy]
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
func (c8 *Chip8) OP_8xy4() {
	vx := (c8.Opcode & 0x0F00) >> 8
	vy := (c8.Opcode & 0x00F0) >> 4

	sum := uint16(c8.Registers[vx]) + uint16(c8.Registers[vy])

	if sum > 255 {
		c8.Registers[0xF] = 1
	} else {
		c8.Registers[0xF] = 0
	}

	c8.Registers[vx] = uint8(sum) & 0xFF
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
func (c8 *Chip8) OP_8xy5() {
	vx := (c8.Opcode & 0x0F00) >> 8
	vy := (c8.Opcode & 0x00F0) >> 4

	if c8.Registers[vx] > c8.Registers[vy] {
		c8.Registers[0xF] = 1
	} else {
		c8.Registers[0xF] = 0
	}
	c8.Registers[vx] -= c8.Registers[vy]
}

// 8xy6 - SHR Vx
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (c8 *Chip8) OP_8xy6() {
	vx := (c8.Opcode & 0x0F00) >> 8

	//put the least signficant bit in vf register
	c8.Registers[0xF] = c8.Registers[vx] & 0x1

	//right shift by 1, means division by 2
	c8.Registers[vx] >>= 1
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func (c8 *Chip8) OP_8xy7() {
	vx := (c8.Opcode & 0x0F00) >> 8
	vy := (c8.Opcode & 0x00F0) >> 4

	if c8.Registers[vy] > c8.Registers[vx] {
		c8.Registers[0xF] = 1
	} else {
		c8.Registers[0xF] = 0
	}
	c8.Registers[vx] = c8.Registers[vy] - c8.Registers[vx]
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
func (c8 *Chip8) OP_8xyE() {
	vx := (c8.Opcode & 0x0F00) >> 8

	//put the Most signficant bit in vf register
	c8.Registers[0xF] = (c8.Registers[vx] & 0x80) >> 7

	//right shift by 1, means division by 2
	c8.Registers[vx] <<= 1
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
func (c8 *Chip8) OP_9xy0() {
	vx := (c8.Opcode & 0x0F00) >> 8
	vy := (c8.Opcode & 0x00F0) >> 4

	if c8.Registers[vx] != c8.Registers[vy] {
		c8.PC += 2
	}
}

// Annn - LD I, addr
// Set I = nnn.
func (c8 *Chip8) OP_Annn() {
	address := c8.Opcode & 0x0FFF
	c8.Index = address
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
func (c8 *Chip8) OP_Bnnn() {
	address := c8.Opcode & 0x0FFF
	c8.PC = uint16(c8.Registers[0]) + address
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
func (c8 *Chip8) OP_Cxkk() {
	vx := (c8.Opcode & 0x0F00) >> 8
	value := uint8(c8.Opcode & 0x00FF)

	c8.Registers[vx] = uint8(c8.RNG.IntN(256)) & value
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func (c8 *Chip8) OP_Dxyn() {
	vx := (c8.Opcode & 0x0F00) >> 8
	vy := (c8.Opcode & 0x00F0) >> 4
	height := c8.Opcode & 0x000F

	// Wrap if going beyond screen boundaries
	xpos := uint16(c8.Registers[vx])
	ypos := uint16(c8.Registers[vy])

	c8.Registers[0xF] = 0

	for row := uint16(0); row < height; row++ {
		spriteByte := c8.Memory[c8.Index+row]
		for col := uint16(0); col < 8; col++ {

			spritePixel := spriteByte & (0x80 >> col)
			if spritePixel == 0 {
				continue
			}

			screenX := (xpos + col) % uint16(VIDEO_WIDTH)
			screenY := (ypos + row) % uint16(VIDEO_HEIGHT)
			idx := screenY*uint16(VIDEO_WIDTH) + screenX

			screenPixel := &c8.Video[idx]
			//check if sprite pixel is on
			if spritePixel != 0 {
				if *screenPixel == 0xFFFFFFFF {
					//collision
					c8.Registers[0xF] = 1
				}
				*screenPixel ^= 0xFFFFFFFF
			}
		}
	}
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
func (c8 *Chip8) OP_Ex9E() {
	vx := (c8.Opcode & 0x0F00) >> 8
	key := c8.Registers[vx]
	if c8.Keypad[key] != 0 {
		c8.PC += 2
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
func (c8 *Chip8) OP_ExA1() {
	vx := (c8.Opcode & 0x0F00) >> 8
	key := c8.Registers[vx]
	if c8.Keypad[key] == 0 {
		c8.PC += 2
	}
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
func (c8 *Chip8) OP_Fx07() {
	vx := (c8.Opcode & 0x0F00) >> 8
	c8.Registers[vx] = c8.DelayTimer
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
func (c8 *Chip8) OP_Fx0A() {
	vx := (c8.Opcode & 0x0F00) >> 8

	if c8.Keypad[0] != 0 {
		c8.Registers[vx] = 0
	} else if c8.Keypad[1] != 0 {
		c8.Registers[vx] = 1
	} else if c8.Keypad[2] != 0 {
		c8.Registers[vx] = 2
	} else if c8.Keypad[3] != 0 {
		c8.Registers[vx] = 3
	} else if c8.Keypad[4] != 0 {
		c8.Registers[vx] = 4
	} else if c8.Keypad[5] != 0 {
		c8.Registers[vx] = 5
	} else if c8.Keypad[6] != 0 {
		c8.Registers[vx] = 6
	} else if c8.Keypad[7] != 0 {
		c8.Registers[vx] = 7
	} else if c8.Keypad[8] != 0 {
		c8.Registers[vx] = 8
	} else if c8.Keypad[9] != 0 {
		c8.Registers[vx] = 9
	} else if c8.Keypad[10] != 0 {
		c8.Registers[vx] = 10
	} else if c8.Keypad[11] != 0 {
		c8.Registers[vx] = 11
	} else if c8.Keypad[12] != 0 {
		c8.Registers[vx] = 12
	} else if c8.Keypad[13] != 0 {
		c8.Registers[vx] = 13
	} else if c8.Keypad[14] != 0 {
		c8.Registers[vx] = 14
	} else if c8.Keypad[15] != 0 {
		c8.Registers[vx] = 15
	} else {
		c8.PC -= 2
	}
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
func (c8 *Chip8) OP_Fx15() {
	vx := (c8.Opcode & 0x0F00) >> 8
	c8.DelayTimer = c8.Registers[vx]
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
func (c8 *Chip8) OP_Fx18() {
	vx := (c8.Opcode & 0x0F00) >> 8
	c8.SoundTimer = c8.Registers[vx]
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
func (c8 *Chip8) OP_Fx1E() {
	vx := (c8.Opcode & 0x0F00) >> 8
	c8.Index += uint16(c8.Registers[vx])
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
func (c8 *Chip8) OP_Fx29() {
	vx := (c8.Opcode & 0x0F00) >> 8
	digit := c8.Registers[vx]
	c8.Index = uint16(FONTSET_START_ADDRESS + (5 * digit))
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (c8 *Chip8) OP_Fx33() {
	vx := (c8.Opcode & 0x0F00) >> 8
	value := c8.Registers[vx]

	c8.Memory[c8.Index+2] = value % 10
	value /= 10

	c8.Memory[c8.Index+1] = value % 10
	value /= 10

	c8.Memory[c8.Index] = value % 10
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
func (c8 *Chip8) OP_Fx55() {
	vx := (c8.Opcode & 0x0F00) >> 8
	var i uint16
	for i = 0; i <= vx; i++ {
		c8.Memory[c8.Index+i] = c8.Registers[i]
	}
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
func (c8 *Chip8) OP_Fx65() {
	vx := (c8.Opcode & 0x0F00) >> 8
	var i uint16
	for i = 0; i <= vx; i++ {
		c8.Registers[i] = c8.Memory[c8.Index+i]
	}
}
