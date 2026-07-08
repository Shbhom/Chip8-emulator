# CHIP-8 Emulator

A CHIP-8 emulator written in Go using the [Ebitengine](https://ebiten.org/) game library.

## Demo

![CHIP-8 Emulator Demo](asset/chip8-demo.gif)

## Features
- Fully functional CHIP-8 CPU instruction set
- Audio and visual rendering via Ebitengine
- Keyboard input mapping
- Scalable window and adjustable CPU execution speed

## Requirements
- Go 1.26 or later

## Running the Emulator

Run the emulator by passing the path to a CHIP-8 ROM:

```bash
go run main.go -rom <path-to-rom>
```

### Optional Flags
- `-scale <int>`: Window scale factor (default: `10`)
- `-cpu <int>`: CPU execution speed/frequency (default: `700`)

For example, to run space invaders with a larger window:
```bash
go run main.go -rom space_invaders.ch8 -scale 15 -cpu 800
```

## Blog Post
For a detailed write-up on how this emulator was built, check out the [blog article](https://example.com/placeholder-blog-link).
