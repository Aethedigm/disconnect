package objects

import (
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Position utils.Vector2
}

func (p *Player) Update() {
}

func (p *Player) Draw(screen *ebiten.Image) {

}
