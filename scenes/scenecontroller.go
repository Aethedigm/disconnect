package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneController struct {
	current Scene
}

var (
	GlobalSceneController *SceneController
)

func NewSceneController(initial Scene) *SceneController {
	return &SceneController{
		current: initial,
	}
}

func (s *SceneController) Update() error {
	return s.current.Update(s)
}

func (s *SceneController) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
}

func (s *SceneController) ChangeScene(next Scene) {
	s.current = next
}
