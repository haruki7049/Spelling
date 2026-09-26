package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	WindowWidth  = 1280
	WindowHeight = 720
)

// Game is the root ebiten.Game implementation. It only owns the currently
// active Scene and delegates Update/Draw to it.
type Game struct {
	scene Scene
}

// NewGame returns a Game that starts on the given Scene.
func NewGame(initial Scene) *Game {
	return &Game{scene: initial}
}

func (g *Game) Update() error {
	next, err := g.scene.Update()
	if err != nil {
		return err
	}
	if next != nil {
		g.scene = next
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}
