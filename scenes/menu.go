package scenes

import (
	"bytes"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/assets"
	"image/color"
	"time"
)

type MenuScene struct {
	state State
	sm    *stagehand.SceneManager[State]
	img   *ebiten.Image

	ui         *ebitenui.UI
	btn        *widget.Button
	StartScene Transition

	gameLogic Brainer
	TimeInit  time.Time
}

func NewMenu(scene stagehand.Scene[State], logic Brainer) *MenuScene {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Title))
	if err != nil {
		return nil
	}

	ebitenImage := ebiten.NewImageFromImage(img)

	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(30)

	menu := &MenuScene{
		img: ebitenImage,
		StartScene: Transition{
			Scene:     scene,
			Direction: stagehand.BottomToTop,
		},
		gameLogic: logic,
		TimeInit:  time.Now(),
	}
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//Set how much padding before displaying content
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top:    350,
				Left:   0,
				Right:  0,
				Bottom: 0,
			}),
		)),
	)
	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text("Start", face, &widget.ButtonTextColor{
			Idle:  color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			Hover: color.NRGBA{0, 255, 128, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			menu.StartGame()
		}),

		// Indicate that this button should not be submitted when enter or space are pressed
		// widget.ButtonOpts.DisableDefaultKeys(),

	)

	rootContainer.AddChild(button)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	menu.ui = &ui
	menu.btn = button
	return menu

}

func (bs *MenuScene) Layout(w, h int) (int, int) {
	return w, h
}

func (m *MenuScene) Update() error {
	m.ui.Update()
	if time.Since(m.TimeInit) > 500*time.Millisecond && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		m.btn.Click()
	}

	return nil
}

func (m *MenuScene) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, 1)
	screen.DrawImage(m.img, op)
	m.ui.Draw(screen)
}

func (m *MenuScene) Load(st State, sm stagehand.SceneController[State]) {
	m.TimeInit = time.Now()
	m.state = st
	m.sm = sm.(*stagehand.SceneManager[State])
	if !m.state.MenuSong.IsPlaying() {
		m.state.MenuSong.Play()
	}
	if m.state.TownSillySong.IsPlaying() {
		m.state.TownSillySong.Pause()
	}

}

func (m *MenuScene) Unload() State {
	m.state.MenuSong.Pause()
	if !m.state.TownSillySong.IsPlaying() {
		m.state.TownSillySong.Play()
	}
	return m.state
}

func (m *MenuScene) StartGame() {
	m.state = m.gameLogic.NextDay(m.state)
	m.sm.SwitchWithTransition(m.StartScene.Scene, stagehand.NewTicksTimedSlideTransition[State](m.StartScene.Direction, time.Second*time.Duration(1)))
}
func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
