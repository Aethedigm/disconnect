package scenes

import (
	"image/color"
	"main/objects"
	"main/physics"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameplayScene struct {
	gameObjects       []objects.GameObject
	staticCollisions  []objects.Collisions
	dynamicCollisions []objects.DynamicCollisions
	radioCollisions   []objects.RadioNode
	capturers         []objects.Capturer
	towers            []*objects.Tower
	projectiles       []objects.Projectile
}

func NewGameplayScene() *GameplayScene {
	gScene := &GameplayScene{}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	gScene.addObject(objects.NewFriendlyMecha(utils.Vector2{X: 100, Y: 100}))
	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 500, Y: 300}))
	gScene.addObject(objects.NewPlayerMecha(utils.Vector2{X: 30, Y: 30}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 800, Y: 350}))
	gScene.addObject(objects.NewCursor())

	return gScene
}

func (g *GameplayScene) addObject(obj objects.GameObject) {
	g.gameObjects = append(g.gameObjects, obj)

	if projectile, ok := obj.(objects.Projectile); ok {
		g.projectiles = append(g.projectiles, projectile)
	} else if sCollider, ok := obj.(objects.Collisions); ok {
		if dCollider, ok := obj.(objects.DynamicCollisions); ok {
			g.dynamicCollisions = append(g.dynamicCollisions, dCollider)
		} else {
			g.staticCollisions = append(g.staticCollisions, sCollider)
		}
	}

	if capture, ok := obj.(objects.Capturer); ok {
		g.capturers = append(g.capturers, capture)
	}

	if tower, ok := obj.(*objects.Tower); ok {
		g.towers = append(g.towers, tower)
	}

	if radioCollider, ok := obj.(objects.RadioNode); ok {
		g.radioCollisions = append(g.radioCollisions, radioCollider)
	}
}

func (g *GameplayScene) Update(controller *SceneController) error {
	var spawned []objects.GameObject
	var destroyed []objects.GameObject

	for i := range g.gameObjects {
		res := g.gameObjects[i].Update()
		spawned = append(spawned, res.Spawn...)

		if res.Destroy {
			destroyed = append(destroyed, g.gameObjects[i])
		}
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

	// Tower Capture
	for _, tower := range g.towers {
		g.updateTowerCapture(tower)
	}

	for _, obj := range destroyed {
		g.removeObject(obj)
	}

	for _, obj := range spawned {
		g.addObject(obj)
	}

	return nil
}

func (g *GameplayScene) updateTowerCapture(tower *objects.Tower) {
	// if owning team is not null, we set owning team to first target
	// if object in capture radius, and same team, we increase capture by 2
	// if object in capture radius, and different team, we decrease capture by 2
	// once captured, towers should not be null team again
	for _, capture := range g.capturers {
		// Detect if in radius
		val := physics.ResolveCircleOverlap(tower.CaptureCollider(), capture.Collider())
		if !val.Equals(utils.Vector2Zero()) { // Collision of some kind
			if tower.TeamOwned() == objects.TeamNone { // Should only happen once per tower
				tower.Team = capture.TeamOwned()
			}

			if tower.TeamOwned() == capture.TeamOwned() {
				tower.CaptureProgress += 2
			} else {
				tower.CaptureProgress -= 2
			}

			if tower.CaptureProgress < 1 {
				tower.Team = capture.TeamOwned()
				tower.CaptureProgress = 10
			}
		}
	}
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{162, 169, 71, 0})
	for _, gObj := range g.gameObjects {
		gObj.Draw(screen)
	}
}

func removeComparable[T comparable](items []T, target T) []T {
	for i := range items {
		if items[i] == target {
			return append(items[:i], items[i+1:]...)
		}
	}

	return items
}

func (g *GameplayScene) removeObject(obj objects.GameObject) {
	g.gameObjects = removeComparable(g.gameObjects, obj)

	if projectile, ok := obj.(objects.Projectile); ok {
		g.projectiles = removeComparable(g.projectiles, projectile)
	}

	if dynamic, ok := obj.(objects.DynamicCollisions); ok {
		g.dynamicCollisions = removeComparable(g.dynamicCollisions, dynamic)
	}

	if capture, ok := obj.(objects.Capturer); ok {
		g.capturers = removeComparable(g.capturers, capture)
	}
}
