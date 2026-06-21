package main

import (
	"main/camera"
	"main/scenes"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneController scenes.SceneController
}

func NewGame() *Game {
	return &Game{
		// sceneController: *scenes.NewSceneController(scenes.NewMainMenuScene()),
		sceneController: *scenes.NewSceneController(scenes.NewGameplayScene()),
	}
}

func (g *Game) Update() error {
	return g.sceneController.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneController.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	h := float64(outsideHeight)

	if h > 480 {
		h = 480
	}

	asp := float64(outsideWidth) / float64(outsideHeight)

	w := math.Min(asp*h, h*4)

	camera.GetCamera().Bounds = utils.Vector2{
		X: w,
		Y: h,
	}

	return int(w), int(h)
}
