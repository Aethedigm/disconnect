package scenes

import "github.com/hajimehoshi/ebiten/v2"

type OptionsScene struct{}

func NewOptionsScene() *OptionsScene {
	return &OptionsScene{}
}

func (o *OptionsScene) Update(controller *SceneController) error {
	return nil
}

func (o *OptionsScene) Draw(screen *ebiten.Image) {}
