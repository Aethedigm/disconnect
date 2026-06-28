package scenes

import (
	"bytes"
	"log"
	"main/data"
	"main/utils"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MainMenuScene struct {
	textSource    *text.GoTextFaceSource
	selectedIndex int
	Sprite        *ebiten.Image
}

func NewMainMenuScene() *MainMenuScene {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	return &MainMenuScene{
		textSource: s,
		Sprite:     utils.ImageDecode(data.SelectorArrow),
	}
}

func (m *MainMenuScene) Update(controller *SceneController) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.selectedIndex++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.selectedIndex--
	}

	if m.selectedIndex < 0 {
		m.selectedIndex = len(data.MenuOptions) - 1
	}

	if m.selectedIndex > len(data.MenuOptions)-1 {
		m.selectedIndex = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch data.MenuOptions[m.selectedIndex] {
		case "Play":
			// controller.ChangeScene(NewGameplayScene())
			controller.ChangeScene(NewCustomizeMechaScene())
		case "Options":
			controller.ChangeScene(NewOptionsScene())
		case "Quit":
			os.Exit(0)
		default:
			log.Println("Invalid selection")
		}

	}

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

	arrowSelectOp := &ebiten.DrawImageOptions{}
	arrowSelectOp.GeoM.Translate(0, float64(35*m.selectedIndex)+10)

	screen.DrawImage(m.Sprite, arrowSelectOp)
}
