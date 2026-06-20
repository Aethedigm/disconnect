package scenes

import (
	"main/objects"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameplayScene struct {
	gameObjects []objects.GameObject
	collisions  []objects.Collisions
}

func NewGameplayScene() *GameplayScene {
	gScene := &GameplayScene{}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	gScene.addObject(objects.NewPlayerMecha(utils.Vector2{X: 30, Y: 30}))
	gScene.addObject(objects.NewCursor())

	return gScene
}

func (g *GameplayScene) addObject(obj objects.GameObject) {
	g.gameObjects = append(g.gameObjects, obj)

	if collider, ok := obj.(objects.Collisions); ok {
		g.collisions = append(g.collisions, collider)
	}
}

func (g *GameplayScene) Update(controller *SceneController) error {
	for i := range g.gameObjects {
		g.gameObjects[i].Update()
	}

	return nil
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	for _, gObj := range g.gameObjects {
		gObj.Draw(screen)
	}
}
