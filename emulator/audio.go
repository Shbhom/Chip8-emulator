package emulator

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate   = 44100
	ToneFreq     = 440.0        // A4 note
	Amplitude    = int16(12000) // Max int16 is 32767
	ToneDuration = time.Second  // Generate 1 second of audio
)

func generateSquareWave() []byte {
	numSamples := int(float64(SampleRate) * ToneDuration.Seconds())

	// 2 bytes per sample (int16)
	buf := make([]byte, numSamples*2)

	period := float64(SampleRate) / ToneFreq

	for i := 0; i < numSamples; i++ {
		var sample int16

		if math.Mod(float64(i), period) < period/2 {
			sample = Amplitude
		} else {
			sample = -Amplitude
		}

		binary.LittleEndian.PutUint16(
			buf[i*2:],
			uint16(sample),
		)
	}

	return buf
}

type AudioManager struct {
	ctx     *audio.Context
	player  *audio.Player
	playing bool
}

func NewAudioManager() (*AudioManager, error) {
	ctx := audio.NewContext(SampleRate)

	beep := generateSquareWave()

	reader := bytes.NewReader(beep)

	loop := audio.NewInfiniteLoop(
		reader,
		int64(len(beep)),
	)

	player, err := ctx.NewPlayer(loop)
	if err != nil {
		return nil, err
	}
	return &AudioManager{
		ctx:    ctx,
		player: player,
	}, nil
}

func (a *AudioManager) Update(enabled bool) {
	switch {
	case enabled && !a.playing:
		a.player.Play()
		a.playing = true

	case !enabled && a.playing:
		a.player.Pause()
		a.playing = false
	}
}

func (a *AudioManager) Close() error {
	return a.player.Close()
}
