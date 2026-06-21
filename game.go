package main

import (
	"main/scenes"

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
	return outsideWidth, outsideHeight
}
