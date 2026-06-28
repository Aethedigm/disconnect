package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type OptionsScene struct{}

func NewOptionsScene() *OptionsScene {
	return &OptionsScene{}
}

func (o *OptionsScene) Update(controller *SceneController) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		controller.ChangeScene(NewMainMenuScene())
		return nil
	}

	return nil
}

func (o *OptionsScene) Draw(screen *ebiten.Image) {}
