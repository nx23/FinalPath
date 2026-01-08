package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/game"
)

func main() {
	g := game.NewGame()
	
	ebiten.SetWindowSize(config.Config.Width, config.Config.Height)
	ebiten.SetWindowTitle(config.Config.Title)
	
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
