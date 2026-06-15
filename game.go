package main

import (
	"main/objects"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	GameObjects []objects.GameObject
}

func (g *Game) Update() error {
	for _, gOjbect := range g.GameObjects {
		gOjbect.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, gObject := range g.GameObjects {
		gObject.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return outsideWidth, outsideHeight
}
