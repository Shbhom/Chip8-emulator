package emulator

import (
	"fmt"
	"log"
	"time"

	"github.com/Shbhom/chip8-emu/emulator/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Chip8       *chip8.Chip8
	Audio       *AudioManager
	keyMap      [16]ebiten.Key
	Image       *ebiten.Image
	Pixels      []byte
	ScaleFactor int
	LastUpdate  time.Time

	CPUAccumulator   time.Duration
	TimerAccumulator time.Duration

	CPUStep   time.Duration
	TimerStep time.Duration
}

func NewGame(fileName string, cpuhz, scale int) *Game {
	c8 := chip8.NewChip8()
	audioMan, err := NewAudioManager()
	if err != nil {
		log.Fatal(fmt.Errorf("Error while generating audio manager: %w", err))
	}
	if err := c8.ReadRom(fileName); err != nil {
		log.Fatal(fmt.Errorf("Error while loading game rom: %w", err))
	}

	var keyMap [16]ebiten.Key

	keyMap[0] = ebiten.KeyX
	keyMap[1] = ebiten.Key1
	keyMap[2] = ebiten.Key2
	keyMap[3] = ebiten.Key3
	keyMap[0xC] = ebiten.Key4
	keyMap[4] = ebiten.KeyQ
	keyMap[5] = ebiten.KeyW
	keyMap[6] = ebiten.KeyE
	keyMap[0xD] = ebiten.KeyR
	keyMap[7] = ebiten.KeyA
	keyMap[8] = ebiten.KeyS
	keyMap[9] = ebiten.KeyD
	keyMap[0xE] = ebiten.Key0
	keyMap[0xA] = ebiten.KeyZ
	keyMap[0xB] = ebiten.KeyC
	keyMap[0xF] = ebiten.KeyV

	img := ebiten.NewImage(int(chip8.VIDEO_WIDTH), int(chip8.VIDEO_HEIGHT))

	if cpuhz <= 0 {
		log.Fatal("CPU frequency must be greater than zero")
	}

	return &Game{
		Chip8:       c8,
		Audio:       audioMan,
		LastUpdate:  time.Now(),
		CPUStep:     time.Second / time.Duration(cpuhz),
		keyMap:      keyMap,
		Image:       img,
		Pixels:      make([]byte, chip8.VIDEO_WIDTH*chip8.VIDEO_HEIGHT*4), // as each pixel requires 4 Bytes each for R,G,B, A
		ScaleFactor: scale,
		TimerStep:   time.Second / 60,
	}
}

func (g *Game) Update() error {
	// Read key input
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	for i, key := range g.keyMap {
		if ebiten.IsKeyPressed(key) {
			g.Chip8.Keypad[i] = 1
		} else {
			g.Chip8.Keypad[i] = 0
		}
	}

	now := time.Now()
	elapsed := now.Sub(g.LastUpdate)
	if elapsed > 100*time.Millisecond {
		elapsed = 100 * time.Millisecond
	}
	g.LastUpdate = now

	g.CPUAccumulator += elapsed
	g.TimerAccumulator += elapsed

	for g.CPUAccumulator >= g.CPUStep {
		g.Chip8.Cycle()
		g.CPUAccumulator -= g.CPUStep
	}

	for g.TimerAccumulator >= g.TimerStep {
		g.Chip8.TickTimers()
		g.TimerAccumulator -= g.TimerStep
	}

	g.Audio.Update(g.Chip8.IsBuzzing())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < len(g.Chip8.Video); i++ {
		j := i * 4 //index for pixels buffer, as for each element in video pixels has 4 elements
		if g.Chip8.Video[i] == 0 {
			g.Pixels[j] = 0
			g.Pixels[j+1] = 0
			g.Pixels[j+2] = 0
			g.Pixels[j+3] = 0
		} else {
			g.Pixels[j] = 255
			g.Pixels[j+1] = 255
			g.Pixels[j+2] = 255
			g.Pixels[j+3] = 255
		}
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.ScaleFactor), float64(g.ScaleFactor))
	g.Image.WritePixels(g.Pixels)
	screen.DrawImage(g.Image, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(chip8.VIDEO_WIDTH), int(chip8.VIDEO_HEIGHT)
}
