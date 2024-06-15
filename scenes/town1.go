package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"time"
)

type Town1Scene struct {
	BaseScene
	MapScene         *graphics.MapScene
	TransitionPoints map[common.Position]stagehand.Scene[State]
}

func (ts *Town1Scene) Draw(screen *ebiten.Image) {
	ts.MapScene.Draw(screen)
	ts.state.Player.Draw(screen)
}

func (ts *Town1Scene) Update() error {
	var speed int64 = common.TileSize
	var pressed = false

	x, y := ts.state.Player.X, ts.state.Player.Y
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
	if x > (common.ScreenWidth-16)/common.Scale || x < 0 ||
		y < 0 || y > (common.ScreenHeight-16)/common.Scale {
		return nil
	}

	if !pressed {
		return nil
	}

	t := ts.MapScene.TileForPos(int(x+16/2), int(y)) // to consider as position the middle-bottom pixel
	if !t.Properties.HasPropertyAs("blocked", "true") {
		ts.state.Player.X, ts.state.Player.Y = x, y
	}
	for k, v := range ts.TransitionPoints {
		if x >= k.X && x < k.X+16 &&
			y >= k.Y && y < k.Y+16 {
			ts.sm.SwitchWithTransition(v, stagehand.NewTicksTimedSlideTransition[State](stagehand.RightToLeft, time.Second*time.Duration(1)))
			return nil
		}
	}
	return nil
}
