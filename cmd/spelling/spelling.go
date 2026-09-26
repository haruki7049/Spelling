package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/haruki7049/spelling/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.WindowWidth, game.WindowHeight)
	ebiten.SetWindowTitle("Spelling")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game.NewGame(game.NewExampleScene())); err != nil {
		log.Fatal(err)
	}
}
