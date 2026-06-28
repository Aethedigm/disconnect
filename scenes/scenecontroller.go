package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneController struct {
	current Scene
}

var (
	GlobalSceneController *SceneController
)

func NewSceneController(initial Scene) *SceneController {
	if GlobalSceneController == nil {

		GlobalSceneController = &SceneController{
			current: initial,
		}
	} else {
		log.Fatal("Scene Controller Already Exists")
	}

	return GlobalSceneController
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
