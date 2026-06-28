package scenes

import (
	"main/data"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CustomizeMechaScene struct {
	MechaUppers []*ebiten.Image
	MechaLowers []*ebiten.Image

	selectRow int

	upperSelection int
	lowerSelection int

	Selector *ebiten.Image
}

func NewCustomizeMechaScene() *CustomizeMechaScene {
	scene := &CustomizeMechaScene{
		MechaUppers: []*ebiten.Image{
			utils.ImageDecode(data.MechaTop),
			utils.ImageDecode(data.MechaTopTwo),
			utils.ImageDecode(data.TankTopOne),
			utils.ImageDecode(data.MechaTopCommander),
		},

		MechaLowers: []*ebiten.Image{
			utils.ImageDecode(data.MechaBottomLegs),
			utils.ImageDecode(data.TankBottomOne),
		},

		Selector: utils.ImageDecode(data.SelectorOutline),
	}

	return scene
}

func (c *CustomizeMechaScene) Update(controller *SceneController) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		controller.ChangeScene(NewMainMenuScene())
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		controller.ChangeScene(NewGameplayScene(Loadout{
			Upper: c.upperSelection,
			Lower: c.lowerSelection,
		}))

		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		c.selectRow++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		c.selectRow--
	}

	if c.selectRow < 0 {
		c.selectRow = 0
	} else if c.selectRow > 1 {
		c.selectRow = 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if c.selectRow == 0 { // Upper
			c.upperSelection--
		} else {
			c.lowerSelection--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if c.selectRow == 0 {
			c.upperSelection++
		} else {
			c.lowerSelection++
		}
	}

	if c.upperSelection < 0 {
		c.upperSelection = len(c.MechaUppers) - 1
	} else if c.upperSelection > len(c.MechaUppers)-1 {
		c.upperSelection = 0
	}

	if c.lowerSelection < 0 {
		c.lowerSelection = len(c.MechaLowers) - 1
	} else if c.lowerSelection > len(c.MechaLowers)-1 {
		c.lowerSelection = 0
	}

	return nil
}

func (c *CustomizeMechaScene) DrawParts(screen *ebiten.Image, arrData []*ebiten.Image, offset float64) {
	for i := range arrData {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i*30), offset)
		screen.DrawImage(arrData[i], op)
	}
}

func (c *CustomizeMechaScene) Draw(screen *ebiten.Image) {
	c.DrawParts(screen, c.MechaUppers, 10)
	c.DrawParts(screen, c.MechaLowers, 45)

	// Draw selectors
	upOp := &ebiten.DrawImageOptions{}
	upOp.GeoM.Translate(float64(c.upperSelection)*30, 10)
	screen.DrawImage(c.Selector, upOp)

	downOp := &ebiten.DrawImageOptions{}
	downOp.GeoM.Translate(float64(c.lowerSelection)*30, 45)
	screen.DrawImage(c.Selector, downOp)

	// Draw preview
	mOp := &ebiten.DrawImageOptions{}
	w, h := ebiten.WindowSize()
	mOp.GeoM.Translate(float64(w/2), float64(h/2))
	screen.DrawImage(c.MechaLowers[c.lowerSelection], mOp)
	screen.DrawImage(c.MechaUppers[c.upperSelection], mOp)
}
