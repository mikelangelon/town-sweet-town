package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(ow, oh int) (int, int) {
	return screenWidth, screenHeight
}
