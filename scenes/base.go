package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/npc"
	"github.com/mikelangelon/town-sweet-town/textbox"
	"github.com/solarlune/resolv"
	"image"
	"time"
)

type BaseScene struct {
	bounds image.Rectangle
	state  State
	sm     *stagehand.SceneManager[State]

	ID               string
	MapScene         *graphics.MapScene
	NPCs             []*npc.NPC
	TransitionPoints Transition
	Text             textbox.TextBox
}

func (bs *BaseScene) Layout(w, h int) (int, int) {
	bs.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (bs *BaseScene) Unload() State {
	return bs.state
}

func (bs *BaseScene) Draw(screen *ebiten.Image) {
	bs.MapScene.Draw(screen)
	bs.state.Player.Draw(screen)
	for _, v := range bs.NPCs {
		v.Draw(screen)
	}
	bs.Text.Draw(screen)
}

func (bs *BaseScene) Update() error {
	if bs.Text.Visible() {
		bs.Text.Update()
		return nil
	}

	for _, v := range bs.NPCs {
		v.Update()
	}
	var speed int64 = common.TileSize
	var pressed = false

	x, y := bs.state.Player.X, bs.state.Player.Y
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// TODO it seems a bit inneficient to recreate this every time
		space := resolv.NewSpace(640, 480, 16, 16)
		player := resolv.NewObject(float64(bs.state.Player.X), float64(bs.state.Player.Y), 16, 16)
		space.Add(player)
		for _, v := range bs.NPCs {
			npc := resolv.NewObject(float64(v.X), float64(v.Y), 16, 16)
			npc.Data = v
			space.Add(npc)
		}

		if collision := player.Check(16, 0); collision != nil {
			if c, ok := collision.Objects[0].Data.(*npc.NPC); ok {
				bs.Text.ShowAndQuestion(c, []string{"House 1", "House 2", textbox.NoResponse})
			}
		}
		if collision := player.Check(-16, 0); collision != nil {
			if c, ok := collision.Objects[0].Data.(*npc.NPC); ok {
				bs.Text.ShowAndQuestion(c, []string{"House 1", "House 2", textbox.NoResponse})
			}
		}
	}

	if x > (common.ScreenWidth-16)/common.Scale || x < 0 ||
		y < 0 || y > (common.ScreenHeight-16)/common.Scale {
		return nil
	}

	if !pressed {
		return nil
	}

	t := bs.MapScene.TileForPos(int(x+16/2), int(y)) // to consider as position the middle-bottom pixel
	if !t.Properties.HasPropertyAs("blocked", "true") {
		bs.state.Player.X, bs.state.Player.Y = x, y
	}
	v := bs.TransitionPoints
	if x >= v.Position.X && x < v.Position.X+16 &&
		y >= v.Position.Y && y < v.Position.Y+16 {
		bs.sm.SwitchWithTransition(v.Scene, stagehand.NewTicksTimedSlideTransition[State](v.Direction, time.Second*time.Duration(1)))
		return nil
	}

	return nil
}

func (bs *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	if st.Status == InitialState {
		st.Status = Playing
		bs.state = st
		bs.sm = sm.(*stagehand.SceneManager[State])
		return
	}
	bs.state = st
	bs.sm = sm.(*stagehand.SceneManager[State])
	timer := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer.C
		bs.state.Player.X, bs.state.Player.Y = bs.TransitionPoints.Position.X, bs.TransitionPoints.Position.Y
	}()
}
