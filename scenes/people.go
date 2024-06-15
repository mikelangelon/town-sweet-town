package scenes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type PeopleScene struct {
	BaseScene
}

func (s *PeopleScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.state.count++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {

	}
	return nil
}

func (s *PeopleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 255, 255}) // Fill Blue
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.state.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}
