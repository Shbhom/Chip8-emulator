package main

import (
	"flag"
	"log"
	"time"

	e "github.com/Shbhom/chip8-emu/emulator"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	filname := flag.String("rom", "", "name of the rom file to be loaded")
	scale := flag.Int("scale", 10, "64x32 is too small for modern monitors so put the scale value, default is 10")
	cycleDelay := flag.Int("cd", 2, "Cycle delay for the emulator, defaults to 10")
	flag.Parse()

	if *filname == "" {
		log.Fatal("no rom file provided")
	}
	cd := time.Duration(*cycleDelay) * time.Millisecond
	g := e.NewGame(*filname, cd, *scale)
	ebiten.SetWindowTitle("Chip8 Emulator")
	ebiten.SetWindowSize(
		int(e.VIDEO_WIDTH)*(*scale),
		int(e.VIDEO_HEIGHT)*(*scale),
	)
	ebiten.RunGame(g)
}
