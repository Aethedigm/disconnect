package scenes

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameOverScreen struct {
	GOS        GameOverStats
	textSource *text.GoTextFaceSource
}

type GameOverStats struct {
	PlayerWon   bool
	FriendlyCap int
	EnemyCap    int
}

func NewGameOverScreen(gos GameOverStats) *GameOverScreen {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	return &GameOverScreen{
		GOS:        gos,
		textSource: s,
	}
}

func (g *GameOverScreen) Update(controller *SceneController) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		controller.ChangeScene(NewMainMenuScene())
		return nil
	}

	return nil
}

func (g *GameOverScreen) Draw(screen *ebiten.Image) {
	// Draw final stats + win/lose
	var message string
	if g.GOS.PlayerWon == true {
		message = "You Win!"
	} else {
		message = "You Lost!"
	}

	w := screen.Bounds().Dx()

	face := &text.GoTextFace{
		Source: g.textSource,
		Size:   24,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(w/4), 15)
	text.Draw(screen, message, face, op)

	message = fmt.Sprintf("Player Capture Score: %d", g.GOS.FriendlyCap)
	op.GeoM.Translate(0, 50)
	text.Draw(screen, message, face, op)

	message = fmt.Sprintf("Enemy Capture Score: %d", g.GOS.EnemyCap)
	op.GeoM.Translate(0, 50)
	text.Draw(screen, message, face, op)

	message = "Press ESC or RETURN (Enter) to continue..."
	op.GeoM.Translate(0, 50)
	text.Draw(screen, message, face, op)
}
