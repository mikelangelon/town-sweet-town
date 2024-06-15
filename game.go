package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

const (
	screenWidth  = 800
	screenHeight = 600
	scale        = 2
	tileSize     = 16
)

type Game struct {
	MapScene *graphics.MapScene
	Player   *graphics.Char
}

func (g *Game) Update() error {
	var speed int64 = tileSize
	var pressed = false

	x, y := g.Player.X, g.Player.Y
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		pressed = true
		y -= speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		pressed = true
		y += speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		pressed = true
		x -= speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		pressed = true
		x += speed
	}
	if !pressed {
		return nil
	}
	t := g.MapScene.TileForPos(int(x+16/2), int(y)) // to consider as position the middle-bottom pixel
	if !t.Properties.HasPropertyAs("blocked", "true") {
		g.Player.X, g.Player.Y = x, y
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.MapScene.Draw(screen)
	g.Player.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return screenWidth, screenHeight
}
