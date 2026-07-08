package main

import (
	"flag"
	"log"

	e "github.com/Shbhom/chip8-emu/emulator"
	"github.com/Shbhom/chip8-emu/emulator/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	filname := flag.String("rom", "", "name of the rom file to be loaded")
	scale := flag.Int("scale", 10, "64x32 is too small for modern monitors so put the scale value, default is 10")
	cpuhz := flag.Int("cpu", 700, "Cycle delay for the emulator, defaults to 10")
	flag.Parse()

	if *filname == "" {
		log.Fatal("no rom file provided")
	}
	g := e.NewGame(*filname, *cpuhz, *scale)
	ebiten.SetWindowTitle("Chip8 Emulator")
	ebiten.SetWindowSize(
		int(chip8.VIDEO_WIDTH)*(*scale),
		int(chip8.VIDEO_HEIGHT)*(*scale),
	)
	ebiten.RunGame(g)

}
