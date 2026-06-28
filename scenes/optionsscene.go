package scenes

import (
	"bytes"
	"log"
	"main/data"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type OptionsScene struct {
	textSource *text.GoTextFaceSource

	selector bool

	selectorSprite *ebiten.Image
}

func NewOptionsScene() *OptionsScene {
	optScene := &OptionsScene{}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	optScene.textSource = s

	optScene.selectorSprite = utils.ImageDecode(data.SelectorArrow)

	return optScene
}

func (o *OptionsScene) Update(controller *SceneController) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		controller.ChangeScene(NewMainMenuScene())
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		o.selector = !o.selector
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		shouldMute := o.selector
		if controller.AudioManager.IsMuted() != shouldMute {
			controller.AudioManager.ToggleMute()
		}
	}

	return nil
}

func (o *OptionsScene) Draw(screen *ebiten.Image) {
	face := &text.GoTextFace{
		Source: o.textSource,
		Size:   24,
	}

	// Draw text
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10)
	text.Draw(screen, "Options", face, textOp)

	// Draw Audio Mute
	textOp.GeoM.Translate(30, 30)
	text.Draw(screen, "Mute Audio", face, textOp)

	textOp.GeoM.Translate(0, 30)
	text.Draw(screen, "Mute | Unmute", face, textOp)

	selecOp := &ebiten.DrawImageOptions{}
	selecOp.GeoM.Rotate(-math.Pi / 2)
	if o.selector {
		selecOp.GeoM.Translate(50, 130)
	} else {
		selecOp.GeoM.Translate(160, 130)
	}
	screen.DrawImage(o.selectorSprite, selecOp)

	textOp.GeoM.Translate(-30, 100)
	text.Draw(screen, "Press Enter to accept, ESC to exit", face, textOp)
}
