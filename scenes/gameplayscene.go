package scenes

import (
	"main/objects"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameplayScene struct {
	gameObjects []objects.GameObject
}

func NewGameplayScene() *GameplayScene {
	gScene := &GameplayScene{}

	gScene.gameObjects = append(gScene.gameObjects, objects.NewMecha(utils.Vector2{30, 30}))

	return gScene
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
