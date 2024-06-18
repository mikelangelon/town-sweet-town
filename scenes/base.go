package scenes

import (
	"github.com/ebitenui/ebitenui"
	imageNine "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/textbox"
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
	NPCs     []*npc.NPC
	Objects  []*graphics.Char

	// Between scenes
	TransitionPoints Transition
	TransitionSleep  uint8

	//UI
	Text     textbox.TextBox
	ui       *ebitenui.UI
	endOfDay *ebitenui.UI
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
		} else {
			if bs.endOfDay == nil {
				bs.ShowEndOfDay()
			}
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		bg := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
		bg.Fill(colorGoal)
		screen.DrawImage(bg, op)
	}
	if bs.ui != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(500, 0)
		bg := ebiten.NewImage(300, 50)
		bs.ui.Draw(bg)
		screen.DrawImage(bg, op)
	}
	if bs.endOfDay != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(50, 150)
		bg := ebiten.NewImage(500, 300)
		bs.endOfDay.Draw(bg)
		screen.DrawImage(bg, op)
	}
}

func (bs *BaseScene) Update() error {
	if bs.Text.Visible() {
		bs.Text.Update()
		return nil
	}
	if bs.ui != nil {
		bs.ui.Update()
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

func (bs *BaseScene) SetupUI() {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	// Construct a container to hold the progress bars.
	progressBarsContainer := widget.NewContainer(
		// The container will use a vertical row layout to lay out the progress
		// bars in a vertical row.
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
		// Set the required anchor layout data to determine where in the root
		// container to place the progress bars.
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// Construct a horizontal progress bar.
	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: imageNine.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Disabled).
			&widget.ProgressBarImage{
				Idle: imageNine.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)

	/*
		To update a progress bar programmatically you can use
		hProgressbar.SetCurrent(value)
		hProgressbar.GetCurrent()
		hProgressbar.Min = 5
		hProgressbar.Max = 10
	*/
	// Add the progress bars as a child of their container.
	progressBarsContainer.AddChild(hProgressbar)
	rootContainer.AddChild(progressBarsContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	bs.ui = &ui
}

func (bs *BaseScene) ShowEndOfDay() {

	var total npc.Stats
	var allChars []npc.Chars
	for _, v := range bs.NPCs {
		total = total.Merge(v.Chars.Stats())
		allChars = append(allChars, v.Chars)
	}
	happiness := npc.CheckHappiness(allChars)
	total.Happiness += happiness

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(imageNine.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 100})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)

	secondaryContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	leftContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	rightContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	face, _ := loadFont(20)
	label1 := widget.NewText(
		widget.TextOpts.Text("Day 1 - Results", face, color.White),
	)
	lRent := widget.NewText(
		widget.TextOpts.Text("Rent: 10 euros", face, color.White),
	)
	lFood := widget.NewText(
		widget.TextOpts.Text("Food: 10", face, color.White),
	)
	lHappiness := widget.NewText(
		widget.TextOpts.Text("Happiness", face, color.White),
	)
	labelHealth := widget.NewText(
		widget.TextOpts.Text("Health", face, color.White),
	)
	lSecurity := widget.NewText(
		widget.TextOpts.Text("Security", face, color.White),
	)
	lCultural := widget.NewText(
		widget.TextOpts.Text("Cultural", face, color.White),
	)
	currentStuff := widget.NewText(
		widget.TextOpts.Text("House 1 - Nice upgrade", face, color.White),
	)

	rootContainer.AddChild(label1)
	rootContainer.AddChild(secondaryContainer)
	secondaryContainer.AddChild(leftContainer)
	secondaryContainer.AddChild(rightContainer)
	leftContainer.AddChild(lRent)
	leftContainer.AddChild(lSecurity)
	leftContainer.AddChild(progress(total.Security))
	leftContainer.AddChild(lCultural)
	leftContainer.AddChild(progress(total.Cultural))
	rightContainer.AddChild(lFood)
	rightContainer.AddChild(lHappiness)
	rightContainer.AddChild(progress(total.Happiness))
	rightContainer.AddChild(labelHealth)
	rightContainer.AddChild(progress(total.Health))
	rootContainer.AddChild(currentStuff)

	bs.endOfDay = &ui
}

func (bs *BaseScene) calculateDay() {
	// TODO Calculate town security
	//

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
