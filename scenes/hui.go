package scenes

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	imageNine "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"image/color"
)

type hui struct {
	lMoney, lFood, lHappiness, lHealth, lSecurity, lCultural *widget.Text
	ui                                                       ebitenui.UI
}

func NewHUI() *hui {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(imageNine.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 100})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.GridLayoutOpts.Spacing(5, 5),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true}, []bool{true, true, true}),
		)),
	)
	face, _ := loadFont(12)
	lMoney := widget.NewText(
		widget.TextOpts.Text("Money: 000", face, color.White),
	)
	lFood := widget.NewText(
		widget.TextOpts.Text("Food: 000", face, color.White),
	)
	lHappiness := widget.NewText(
		widget.TextOpts.Text("Happiness: 000", face, color.White),
	)
	lHealth := widget.NewText(
		widget.TextOpts.Text("Health: 000", face, color.White),
	)
	lSecurity := widget.NewText(
		widget.TextOpts.Text("Security: 000", face, color.White),
	)
	lCultural := widget.NewText(
		widget.TextOpts.Text("Cultural: 000", face, color.White),
	)
	rootContainer.AddChild(lMoney)
	rootContainer.AddChild(lHappiness)
	rootContainer.AddChild(lHealth)
	rootContainer.AddChild(lFood)
	rootContainer.AddChild(lSecurity)
	rootContainer.AddChild(lCultural)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	return &hui{
		lMoney:     lMoney,
		lFood:      lFood,
		lHappiness: lHappiness,
		lHealth:    lHealth,
		lSecurity:  lSecurity,
		lCultural:  lCultural,
		ui:         ui,
	}
}

func (h *hui) Update(stats map[string]int) {
	h.lMoney.Label = fmt.Sprintf("Money: %03d", stats[npc.Money])
	h.lFood.Label = fmt.Sprintf("Food: %03d", stats[npc.Food])
	h.lHappiness.Label = fmt.Sprintf("Happiness: %03d", stats[npc.Happiness])
	h.lHealth.Label = fmt.Sprintf("Health: %03d", stats[npc.Health])
	h.lSecurity.Label = fmt.Sprintf("Secutiry: %03d", stats[npc.Security])
	h.lCultural.Label = fmt.Sprintf("Culture: %03d", stats[npc.Cultural])
}
