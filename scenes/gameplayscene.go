package scenes

import (
	"main/objects"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameplayScene struct {
	gameObjects []objects.GameObject
}

func NewGameplayScene() *GameplayScene {
	return &GameplayScene{}
}

func (g *GameplayScene) Update(controller *SceneController) error {
	for _, gObj := range g.gameObjects {
		gObj.Update()
	}

	return nil
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	for _, gObj := range g.gameObjects {
		gObj.Draw(screen)
	}
}
