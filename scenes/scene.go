package scenes

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update(controller *SceneController) error
	Draw(screen *ebiten.Image)
}
