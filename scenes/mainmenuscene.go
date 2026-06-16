package scenes

import (
	"bytes"
	"log"
	"main/data"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MainMenuScene struct {
	textSource *text.GoTextFaceSource
}

func NewMainMenuScene() *MainMenuScene {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	return &MainMenuScene{
		textSource: s,
	}
}

func (m *MainMenuScene) Update(controller *SceneController) error {
	return nil
}

func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(30, 10)

	for _, opt := range data.MenuOptions {
		text.Draw(screen, opt, &text.GoTextFace{
			Source: m.textSource,
			Size:   24,
		}, op)
		op.GeoM.Translate(0, 35)
	}
}
