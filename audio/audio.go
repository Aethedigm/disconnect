package audio

import (
	"bytes"
	"io"
	"log"
	"main/data"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type SFX int

const (
	ShootSFX SFX = iota
	ExplosionSFX
)

type AudioManager struct {
	ctx     *audio.Context
	players map[SFX][]*audio.Player
	muted   bool
}

func NewAudioManager() *AudioManager {
	playerMap := make(map[SFX][]*audio.Player)
	a := &AudioManager{
		ctx:     audio.NewContext(44100),
		muted:   false,
		players: playerMap,
	}

	a.loadSFX(ShootSFX, data.SFX_Shoot, 8)
	a.loadSFX(ExplosionSFX, data.SFX_Explode, 8)
	return a
}

func (a *AudioManager) loadSFX(sfx SFX, raw []byte, channels int) {
	stream, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(raw))
	if err != nil {
		log.Println(err)
		return
	}

	decoded, err := io.ReadAll(stream)
	if err != nil {
		log.Println(err)
		return
	}

	for range channels {
		a.players[sfx] = append(a.players[sfx], a.ctx.NewPlayerFromBytes(decoded))
	}
}

func (a *AudioManager) Play(sfx SFX) {
	if a == nil || a.muted {
		return
	}

	for _, player := range a.players[sfx] {
		if !player.IsPlaying() {
			_ = player.Rewind()
			player.Play()

			return
		}
	}
}

func (a *AudioManager) ToggleMute() {
	a.muted = !a.muted
}

func (a *AudioManager) IsMuted() bool {
	return a.muted
}
