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
	"golang.org/x/image/font"
	"image/color"
	"time"
)

type MenuScene struct {
	state State
	sm    *stagehand.SceneManager[State]
	img   *ebiten.Image

	ui         *ebitenui.UI
	Buttons    []*widget.Button
	FocusIndex int
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

	buttonImage, _ := loadButtonImage()

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
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top:    350,
				Left:   100,
				Right:  100,
				Bottom: 0,
			}),
		)),
	)
	easyButton := menu.createButton(buttonImage, face, widget.AnchorLayoutPositionStart, "Easy")
	normalButton := menu.createButton(buttonImage, face, widget.AnchorLayoutPositionCenter, "Normal")
	hardButton := menu.createButton(buttonImage, face, widget.AnchorLayoutPositionEnd, "Hard")
	rootContainer.AddChild(easyButton)
	rootContainer.AddChild(normalButton)
	rootContainer.AddChild(hardButton)
	menu.Buttons = []*widget.Button{easyButton, normalButton, hardButton}
	easyButton.Focus(true)
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	menu.ui = &ui
	return menu

}

func (bs *MenuScene) Layout(w, h int) (int, int) {
	return w, h
}

func (m *MenuScene) Update() error {
	m.ui.Update()
	if time.Since(m.TimeInit) < 500*time.Millisecond {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		m.Buttons[m.FocusIndex].Click()
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		m.FocusIndex--
		if m.FocusIndex < 0 {
			m.FocusIndex = 2
		}
		m.Buttons[m.FocusIndex].Focus(true)
		m.TimeInit = time.Now().Add(450 * time.Millisecond)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		m.FocusIndex++
		if m.FocusIndex > 2 {
			m.FocusIndex = 0
		}
		m.Buttons[m.FocusIndex].Focus(true)
		m.TimeInit = time.Now().Add(450 * time.Millisecond)
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

func (m *MenuScene) StartGame(difficulty string) {
	m.state.Difficulty = difficulty
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

func (m *MenuScene) createButton(buttonImage *widget.ButtonImage, face font.Face, position widget.AnchorLayoutPosition, difficulty string) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: position,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(difficulty, face, &widget.ButtonTextColor{
			Idle:  color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			Hover: color.NRGBA{0, 255, 128, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.StartGame(difficulty)
		}),
	)
}
