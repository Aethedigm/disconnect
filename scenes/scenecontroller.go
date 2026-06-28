package scenes

import (
	"main/audio"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneController struct {
	current      Scene
	AudioManager *audio.AudioManager
}

var (
	GlobalSceneController *SceneController
)

func NewSceneController(initial Scene) *SceneController {
	sceneController := &SceneController{
		current:      initial,
		AudioManager: audio.NewAudioManager(),
	}
	return sceneController
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
