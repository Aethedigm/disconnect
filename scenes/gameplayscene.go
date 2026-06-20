package scenes

import (
	"main/objects"
	"main/physics"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameplayScene struct {
	gameObjects       []objects.GameObject
	staticCollisions  []objects.Collisions
	dynamicCollisions []objects.DynamicCollisions
}

func NewGameplayScene() *GameplayScene {
	gScene := &GameplayScene{}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	gScene.addObject(objects.NewMecha(utils.Vector2{X: 100, Y: 100}))
	gScene.addObject(objects.NewPlayerMecha(utils.Vector2{X: 30, Y: 30}))
	gScene.addObject(objects.NewCursor())

	return gScene
}

func (g *GameplayScene) addObject(obj objects.GameObject) {
	g.gameObjects = append(g.gameObjects, obj)

	if sCollider, ok := obj.(objects.Collisions); ok {
		if dCollider, ok := obj.(objects.DynamicCollisions); ok {
			g.dynamicCollisions = append(g.dynamicCollisions, dCollider)
		} else {
			g.staticCollisions = append(g.staticCollisions, sCollider)
		}
	}
}

func (g *GameplayScene) Update(controller *SceneController) error {
	for i := range g.gameObjects {
		g.gameObjects[i].Update()
	}

	// Dynamic vs Static: Push Dynamic Away
	for _, dynamic := range g.dynamicCollisions {
		for _, static := range g.staticCollisions {
			push := physics.ResolveCircleOverlap(dynamic.Collider(), static.Collider())
			dynamic.Move(push)
		}
	}

	// Dynamic vs Dynamic: Move both back half
	for i := range g.dynamicCollisions {
		for j := i + 1; j < len(g.dynamicCollisions); j++ {
			push := physics.ResolveCircleOverlap(
				g.dynamicCollisions[i].Collider(),
				g.dynamicCollisions[j].Collider(),
			)

			push.MulScalar(0.5)
			g.dynamicCollisions[i].Move(push)

			push.MulScalar(-1)
			g.dynamicCollisions[j].Move(push)
		}
	}

	return nil
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	for _, gObj := range g.gameObjects {
		gObj.Draw(screen)
	}
}
