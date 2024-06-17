package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/textbox"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"github.com/solarlune/resolv"
	"image"
	"image/color"
	"time"
)

type BaseScene struct {
	bounds image.Rectangle
	state  State
	sm     *stagehand.SceneManager[State]

	ID               string
	MapScene         *graphics.MapScene
	NPCs             []*npc.NPC
	Objects          []*graphics.Char
	TransitionPoints Transition
	Text             textbox.TextBox

	TransitionSleep uint8
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
	for _, v := range bs.Objects {
		v.Draw(screen)
	}
	bs.Text.Draw(screen)
	if bs.state.Status == DayEnding {
		colorGoal := color.RGBA{10, 10, 10, bs.TransitionSleep}
		if bs.TransitionSleep < 200 {
			bs.TransitionSleep++
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		bg := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
		bg.Fill(colorGoal)
		screen.DrawImage(bg, op)
	}
}

func (bs *BaseScene) Update() error {
	if bs.Text.Visible() {
		bs.Text.Update()
		return nil
	}
	if bs.state.Status != Playing {
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
		for _, v := range bs.Objects {
			npc := resolv.NewObject(float64(v.X), float64(v.Y), 16, 16)
			npc.Data = v
			space.Add(npc)
		}

		if collision := player.Check(16, 0); collision != nil {
			bs.Action(collision)
		}
		if collision := player.Check(-16, 0); collision != nil {
			bs.Action(collision)
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

func (bs *BaseScene) Action(collision *resolv.Collision) {
	if c, ok := collision.Objects[0].Data.(*npc.NPC); ok {
		bs.TalkToNPC(c)
	}
	if c, ok := collision.Objects[0].Data.(*graphics.Char); ok {
		bs.ActionToObject(c)
	}
}

func (bs *BaseScene) TalkToNPC(npc *npc.NPC) {
	answerFunc := func(answer string) {
		if answer != textbox.NoResponse && answer != textbox.No {
			for _, v := range bs.state.World["town1"].Houses {
				if v.ID == answer {
					npc.House = v
					v.Owner = &npc.ID
				}
			}
			npc.Move = &common.Position{X: -16, Y: npc.Y}
			bs.state.World["people"].RemoveNPC(npc.ID)
			bs.state.World["town1"].AddNPC(npc)
		}
	}
	var options []string
	for _, v := range bs.state.World["town1"].Houses {
		if v.Owner != nil {
			continue
		}
		options = append(options, v.ID)
	}
	options = append(options, textbox.NoResponse)
	bs.Text.ShowAndQuestion(npc.Talk(), options, answerFunc)
}

func (bs *BaseScene) ActionToObject(object *graphics.Char) {
	bs.Text.ShowAndQuestion(
		[]string{"Go to the next day?"},
		[]string{"Yes", textbox.No},
		func(answer string) {
			bs.state.Status = DayEnding
		},
	)
}

func (bs *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	bs.state = st
	bs.sm = sm.(*stagehand.SceneManager[State])
	for _, v := range bs.state.World[bs.ID].Houses {
		if v.Owner != nil {
			npc := bs.state.World[bs.ID].NPCs.GetNPC(*v.Owner)
			npc.X = v.Position.X + 16
			npc.Y = v.Position.Y + 16
			npc.Move = nil
		}
	}
	bs.NPCs = bs.state.World[bs.ID].NPCs
	bs.Objects = bs.state.World[bs.ID].Objects

	if bs.state.Status == InitialState {
		bs.state.Status = Playing
		return
	}

	timer := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer.C
		bs.state.Player.X, bs.state.Player.Y = bs.TransitionPoints.Position.X, bs.TransitionPoints.Position.Y
	}()
}
