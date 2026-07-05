package emulator

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Chip8         *Chip8
	CycleDelay    time.Duration
	LastCycle     time.Time
	LastTimerTick time.Time
	keyMap        [16]ebiten.Key
	Image         *ebiten.Image
	Pixels        []byte
	ScaleFactor   int
}

func NewGame(fileName string, cd time.Duration, scale int) *Game {
	c8 := NewChip8()
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

	img := ebiten.NewImage(int(VIDEO_WIDTH), int(VIDEO_HEIGHT))

	return &Game{
		Chip8:         c8,
		CycleDelay:    cd,
		LastCycle:     time.Now(),
		LastTimerTick: time.Now(),
		keyMap:        keyMap,
		Image:         img,
		Pixels:        make([]byte, VIDEO_WIDTH*VIDEO_HEIGHT*4), // as each pixel requires 4 Bytes each for R,G,B, A
		ScaleFactor:   scale,
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

	if time.Since(g.LastCycle) >= g.CycleDelay {
		g.LastCycle = time.Now()
		for i := 0; i < 10; i++ {
			g.Chip8.Cycle()
		}
	}

	const timerFrequency = time.Second / 60

	if time.Since(g.LastTimerTick) >= timerFrequency {
		g.LastTimerTick = time.Now()
		g.Chip8.TickTimers()
	}
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
	return int(VIDEO_WIDTH), int(VIDEO_HEIGHT)
}
