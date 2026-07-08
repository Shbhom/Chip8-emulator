package chip8

import "math/rand/v2"

const START_ADDR = 0x200
const FONTSET_SIZE = 80
const FONTSET_START_ADDRESS = 0x50
const VIDEO_HEIGHT uint32 = 32
const VIDEO_WIDTH uint32 = 64
const debug = false

type INSTRUCTIONS func()

var (
	FILE_NAME         = "Soccer.ch8"
	fontset   []uint8 = []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
)

type Chip8 struct {
	Registers  [16]uint8
	Memory     [4096]uint8
	Index      uint16
	PC         uint16
	Stack      [16]uint16
	SP         uint8
	DelayTimer uint8
	SoundTimer uint8
	Keypad     [16]uint8
	Video      [VIDEO_WIDTH * VIDEO_HEIGHT]uint32 //flattened buffer to store display pixel
	Opcode     uint16
	RNG        *rand.Rand //random number generator
	table      [16]INSTRUCTIONS
	table0     [16]INSTRUCTIONS
	table8     [16]INSTRUCTIONS
	tableE     [16]INSTRUCTIONS
	tableF     [128]INSTRUCTIONS
}

func NewChip8() *Chip8 {
	cp := &Chip8{}
	cp.PC = uint16(START_ADDR)
	for i := 0; i < FONTSET_SIZE; i++ {
		cp.Memory[FONTSET_START_ADDRESS+i] = fontset[i]
	}
	cp.RNG = rand.New(rand.NewPCG(
		rand.Uint64(),
		rand.Uint64(),
	))

	cp.table[0x0] = cp.Table0
	cp.table[0x1] = cp.OP_1nnn
	cp.table[0x2] = cp.OP_2nnn
	cp.table[0x3] = cp.OP_3xkk
	cp.table[0x4] = cp.OP_4xkk
	cp.table[0x5] = cp.OP_5xy0
	cp.table[0x6] = cp.OP_6xkk
	cp.table[0x7] = cp.OP_7xkk
	cp.table[0x8] = cp.Table8
	cp.table[0x9] = cp.OP_9xy0
	cp.table[0xA] = cp.OP_Annn
	cp.table[0xB] = cp.OP_Bnnn
	cp.table[0xC] = cp.OP_Cxkk
	cp.table[0xD] = cp.OP_Dxyn
	cp.table[0xE] = cp.TableE
	cp.table[0xF] = cp.TableF

	for i := 0; i <= 0xE; i++ {
		cp.table0[i] = cp.OP_NULL
		cp.table8[i] = cp.OP_NULL
		cp.tableE[i] = cp.OP_NULL
	}

	cp.table0[0x0] = cp.OP_00E0
	cp.table0[0xE] = cp.OP_00EE

	cp.table8[0x0] = cp.OP_8xy0
	cp.table8[0x1] = cp.OP_8xy1
	cp.table8[0x2] = cp.OP_8xy2
	cp.table8[0x3] = cp.OP_8xy3
	cp.table8[0x4] = cp.OP_8xy4
	cp.table8[0x5] = cp.OP_8xy5
	cp.table8[0x6] = cp.OP_8xy6
	cp.table8[0x7] = cp.OP_8xy7
	cp.table8[0xE] = cp.OP_8xyE

	cp.tableE[0x1] = cp.OP_ExA1
	cp.tableE[0xE] = cp.OP_Ex9E

	for i := 0; i <= 0x65; i++ {
		cp.tableF[i] = cp.OP_NULL
	}

	cp.tableF[0x07] = cp.OP_Fx07
	cp.tableF[0x0A] = cp.OP_Fx0A
	cp.tableF[0x15] = cp.OP_Fx15
	cp.tableF[0x18] = cp.OP_Fx18
	cp.tableF[0x1E] = cp.OP_Fx1E
	cp.tableF[0x29] = cp.OP_Fx29
	cp.tableF[0x33] = cp.OP_Fx33
	cp.tableF[0x55] = cp.OP_Fx55
	cp.tableF[0x65] = cp.OP_Fx65

	return cp
}
