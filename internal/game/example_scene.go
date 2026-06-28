package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ExampleScene struct{}

func NewExampleScene() *ExampleScene {
	return nil
}

func (s *ExampleScene) Update() (Scene, error) {
	return nil, nil
}

func (s *ExampleScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, I am ExampleScene.")
}
