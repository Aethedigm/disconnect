package scenes

import (
	"image/color"
	"main/camera"
	"main/objects"
	"main/physics"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameplayScene struct {
	gameObjects       []objects.GameObject
	staticCollisions  []objects.Collisions
	dynamicCollisions []objects.DynamicCollisions
	radioCollisions   []objects.RadioNode
	capturers         []objects.Capturer
	towers            []*objects.Tower
	projectiles       []objects.Projectile

	Cam         *camera.Camera
	PlayerMecha *objects.Mecha

	captureAmounts map[objects.Team]int

	isPaused bool
}

func NewGameplayScene() *GameplayScene {
	gScene := &GameplayScene{}

	gScene.Cam = camera.GetCamera()

	// TODO: Load levels
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	gScene.addObject(objects.NewFriendlyMecha(utils.Vector2{X: 100, Y: 100}))
	gScene.addObject(objects.NewFriendlyMecha(utils.Vector2{X: 000, Y: 000}))
	gScene.addObject(objects.NewFriendlyMecha(utils.Vector2{X: 100, Y: 000}))
	gScene.addObject(objects.NewFriendlyMecha(utils.Vector2{X: 000, Y: 100}))

	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 1100, Y: 1100}))
	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 1200, Y: 1100}))
	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 1100, Y: 1200}))
	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 1200, Y: 1200}))
	gScene.addObject(objects.NewEnemyMecha(utils.Vector2{X: 1300, Y: 1300}))

	pMecha := objects.NewPlayerMecha(utils.Vector2{X: 30, Y: 30})
	gScene.addObject(pMecha)
	gScene.PlayerMecha = pMecha

	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 200, Y: 200}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 1000, Y: 1000}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 200, Y: 1000}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 1000, Y: 200}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 650, Y: 200}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 200, Y: 650}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 1000, Y: 650}))
	gScene.addObject(objects.NewNeutralTower(utils.Vector2{X: 650, Y: 1000}))
	gScene.addObject(objects.NewCursor())

	gScene.captureAmounts = make(map[objects.Team]int)

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

	if radioCollider, ok := obj.(objects.RadioNode); ok && radioCollider.HasRadio() {
		g.radioCollisions = append(g.radioCollisions, radioCollider)
	}
}

func (g *GameplayScene) radioNetworkFor(mecha *objects.Mecha) map[objects.RadioNode]bool {
	connected := map[objects.RadioNode]bool{}
	var insideNodes []objects.RadioNode

	// Loop all radio nodes, collect same Team and mecha is inside of (likely only max 2)
	for _, radioNode := range g.radioCollisions {
		if radioNode.TeamOwned() != mecha.TeamOwned() {
			continue
		}

		if physics.CircleCollidersCollided(radioNode.RadioCollider(), mecha.Collider()) {
			insideNodes = append(insideNodes, radioNode)
		}
	}

	// Loop and walk the touching same team radio nodes
	for len(insideNodes) > 0 {
		node := insideNodes[len(insideNodes)-1]
		insideNodes = insideNodes[:len(insideNodes)-1]

		// if we're already connected, hop out of this node
		if connected[node] {
			continue
		}
		connected[node] = true

		for _, next := range g.radioCollisions {
			if connected[next] || next.TeamOwned() != mecha.TeamOwned() {
				continue
			}

			if physics.CircleCollidersCollided(node.RadioCollider(), next.RadioCollider()) {
				insideNodes = append(insideNodes, next)
			}
		}
	}

	return connected
}

func (g *GameplayScene) radioNetworkContainsMecha(network map[objects.RadioNode]bool, target *objects.Mecha) bool {
	for node := range network {
		if physics.CircleCollidersCollided(node.RadioCollider(), target.Collider()) {
			return true
		}
	}

	return false
}

func (g *GameplayScene) buildWorldContext(mecha *objects.Mecha) objects.WorldContext {
	var wc objects.WorldContext

	network := g.radioNetworkFor(mecha)

	for _, dynamic := range g.dynamicCollisions {
		// Only Mechas
		obj, ok := dynamic.(*objects.Mecha)
		if !ok || obj == mecha {
			continue
		}

		delta := obj.Position.Subbed(mecha.Position)
		delLength := delta.Length()
		if delLength < 350 || g.radioNetworkContainsMecha(network, obj) {
			wc.NearbyMecha = append(wc.NearbyMecha, objects.ObjectInfo{
				Position: obj.Position,
				Team:     obj.TeamOwned(),
				Distance: delLength,
			})
		}
	}

	// Loop Towers
	for _, tower := range g.towers {
		delta := tower.Position.Subbed(mecha.Position)
		wc.NearbyTowers = append(wc.NearbyTowers, objects.ObjectInfo{
			Position: tower.Position,
			Team:     tower.TeamOwned(),
			Distance: delta.Length(),
			Progress: tower.CaptureProgress,
		})
	}

	return wc
}

func (g *GameplayScene) CheckForGameEnded() bool {
	fTeam := 0
	eTeam := 0

	for _, tower := range g.towers {
		if tower.Team == objects.TeamNone {
			return false
		}

		if tower.Team == objects.TeamEnemy {
			eTeam++
		} else {
			fTeam++
		}

		if eTeam > 0 && fTeam > 0 {
			return false
		}
	}

	return true
}

func (g *GameplayScene) Update(controller *SceneController) error {
	var spawned []objects.GameObject
	var destroyed []objects.GameObject

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.isPaused = !g.isPaused
	}

	if g.isPaused {
		return nil
	}

	// Check for End State
	if g.CheckForGameEnded() {
		// Game is complete, go to Win/Loss screen

	}

	for i := range g.gameObjects {
		if explosion, ok := g.gameObjects[i].(*objects.Explosion); ok {
			for _, dynamics := range g.dynamicCollisions {
				if mecha, ok := dynamics.(*objects.Mecha); ok {
					if mecha.Team == explosion.Team {
						continue
					}

					val := physics.ResolveCircleOverlap(explosion.Explosion(), mecha.Collider())
					if !val.Equals(utils.Vector2Zero()) {
						mecha.ApplyDamage(explosion.SplashDamage)
					}
				}
			}

			destroyed = append(destroyed, explosion)
		}

		if mecha, ok := g.gameObjects[i].(*objects.Mecha); ok {
			mecha.SetWorldContext(
				g.buildWorldContext(mecha),
			)
		}

		res := g.gameObjects[i].Update()
		spawned = append(spawned, res.Spawn...)

		if res.Destroy {
			destroyed = append(destroyed, g.gameObjects[i])
		}
	}
	g.Cam.Move(g.PlayerMecha.Position)

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

	// Projectiles vs Dynamic+Static
	for i := range g.projectiles {
		if g.projectiles[i].IsDestroyed() {
			continue
		}

		// Static just destroys the projectile
		for _, static := range g.staticCollisions {
			val := physics.ResolveCircleOverlap(g.projectiles[i].Collider(), static.Collider())
			if !val.Equals(utils.Vector2Zero()) {
				g.projectiles[i].Destroy()

				if rocket, ok := g.projectiles[i].(*objects.Rocket); ok {
					spawned = append(spawned, rocket.Explode().Spawn...)
				}

				break
			}
		}

		for j := range g.dynamicCollisions {
			if targetTeam, ok := g.dynamicCollisions[j].(objects.HasTeam); ok {
				if targetTeam.TeamOwned() == g.projectiles[i].TeamOwned() {
					continue
				}
			}

			if g.projectiles[i].IsDestroyed() {
				continue
			}

			val := physics.ResolveCircleOverlap(g.projectiles[i].Collider(), g.dynamicCollisions[j].Collider())
			if !val.Equals(utils.Vector2Zero()) {
				g.projectiles[i].Destroy()

				if damageTaker, ok := g.dynamicCollisions[j].(objects.DamageTarget); ok {
					damageTaker.ApplyDamage(g.projectiles[i].ProjectileDamage())

					if rocket, ok := g.projectiles[i].(*objects.Rocket); ok {
						spawned = append(spawned, rocket.Explode().Spawn...)
					}
				}
				break
			}
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
				g.captureAmounts[tower.TeamOwned()]++
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
