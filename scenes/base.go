package scenes

import (
	imageNine "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/textbox"
	"github.com/mikelangelon/town-sweet-town/world"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"time"
)

type BaseScene struct {
	bounds image.Rectangle
	// state & management
	state State
	sm    *stagehand.SceneManager[State]

	// World related
	ID       string
	MapScene *graphics.MapScene
	NPCs     npc.NPCs
	Objects  []world.Object

	// Between scenes
	TransitionPoints Transition

	//UI
	Text    textbox.TextBox
	ui      *hui
	rulesUI *ruleUI
	goalsUI *goalsUI
}

func (bs *BaseScene) Layout(w, h int) (int, int) {
	bs.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (bs *BaseScene) Update() (bool, error) {
	if bs.Text.Visible() {
		bs.Text.Update()
		return true, nil
	}
	if bs.ui != nil {
		bs.ui.Update(bs.state.Stats)
	}
	if bs.state.Status != Playing {
		return true, nil
	}

	// Deal with NPC moves
	for _, v := range bs.NPCs {
		v.Update()
	}

	// Deal with player moves
	pressed, err := bs.playerUpdate()
	if err != nil {
		return true, err
	}
	if !pressed {
		return true, nil
	}

	// Change scene
	v := bs.TransitionPoints
	x, y := bs.state.Player.X, bs.state.Player.Y
	if x >= v.Position.X && x < v.Position.X+16 &&
		y >= v.Position.Y && y < v.Position.Y+16 {
		bs.state.Status = Transitioning
		bs.sm.SwitchWithTransition(v.Scene, stagehand.NewTicksTimedSlideTransition[State](v.Direction, time.Second*time.Duration(1)))
		return false, nil
	}

	return false, nil
}

func (bs *BaseScene) Draw(screen *ebiten.Image) {
	bs.MapScene.Draw(screen)
	if bs.rulesUI != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(50, 50)
		bg := ebiten.NewImage(common.ScreenWidth-100, common.ScreenHeight-100)
		bg.Fill(color.RGBA{50, 50, 50, 150})
		bs.rulesUI.ui.Draw(bg)
		screen.DrawImage(bg, op)
		return
	}
	if bs.goalsUI != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(50, 50)
		bg := ebiten.NewImage(common.ScreenWidth-100, common.ScreenHeight-100)
		bg.Fill(color.RGBA{50, 50, 50, 150})
		bs.goalsUI.ui.Draw(bg)
		screen.DrawImage(bg, op)
		return
	}
	bs.state.Player.Draw(screen)
	for _, v := range bs.NPCs {
		v.Draw(screen)
	}
	for _, v := range bs.Objects {
		v.Draw(screen)
	}
	bs.Text.Draw(screen)
	if bs.ui != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(500, 0)
		bg := ebiten.NewImage(300, 60)
		bs.ui.ui.Draw(bg)
		screen.DrawImage(bg, op)
	}
}

func (bs *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	bs.state = st
	bs.sm = sm.(*stagehand.SceneManager[State])
	for _, v := range bs.state.World[bs.ID].Houses {
		if v.Owner != nil {
			npc := bs.state.World[bs.ID].NPCs.GetNPC(*v.Owner)
			npc.X = v.DoorPosition().X + 16
			npc.Y = v.DoorPosition().Y + 16
			npc.Move = nil
		}
	}
	bs.NPCs = bs.state.World[bs.ID].NPCs
	bs.Objects = bs.state.World[bs.ID].Objects
}

func (bs *BaseScene) Unload() State {
	return bs.state
}

func (bs *BaseScene) checkActionExecuted() *resolv.Collision {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return nil
	}
	// TODO it seems a bit inneficient to recreate this every time
	space := resolv.NewSpace(640, 480, 16, 16)
	player := resolv.NewObject(float64(bs.state.Player.X), float64(bs.state.Player.Y), 16, 16)
	space.Add(player)
	for _, v := range bs.NPCs {
		person := resolv.NewObject(float64(v.X), float64(v.Y), 16, 16)
		person.Data = v
		space.Add(person)
	}
	for _, v := range bs.Objects {
		person := resolv.NewObject(float64(v.Position().X), float64(v.Position().Y), 16, 16)
		person.Data = v
		space.Add(person)
	}

	if collision := player.Check(16, 0); collision != nil {
		return collision
	}
	if collision := player.Check(-16, 0); collision != nil {
		return collision
	}
	// If anything, show pause menu
	bs.state.Status = Pause
	bs.rulesUI = NewRulesUI(bs.state.GameLogic.GetRuler().Rules)

	return nil
}

// returns if something was pressed and the error
func (bs *BaseScene) playerUpdate() (bool, error) {
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
		pressed = true
	}
	if x > (common.ScreenWidth-16)/common.Scale || x < 0 ||
		y < 0 || y > (common.ScreenHeight-16)/common.Scale {
		return true, nil
	}

	if !pressed {
		return false, nil
	}
	for _, v := range bs.state.World[bs.ID].Objects {
		if v.Position().X == x && v.Position().Y == y {
			return true, nil
		}
	}
	if !bs.MapScene.AnyPropertyTileAs(int(x+16/2), int(y), "blocked", "true") {
		bs.state.Player.X, bs.state.Player.Y = x, y
	}

	return true, nil
}

func progress(current int) *widget.ProgressBar {
	return widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 10),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle: imageNine.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			&widget.ProgressBarImage{
				Idle: imageNine.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		widget.ProgressBarOpts.Values(0, 100, current),
	)
}
func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
