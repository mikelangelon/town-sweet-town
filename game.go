package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

const (
	screenWidth  = 800
	screenHeight = 600
	scale        = 2
)

type Game struct {
	MapScene *graphics.MapScene
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.MapScene.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return screenWidth, screenHeight
}
