package scenes

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/world"
	"image/color"
)

type goalsUI struct {
	ui ebitenui.UI
}

func NewGoals(goals []world.Goal) *goalsUI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)
	label := widget.NewText(
		widget.TextOpts.Text("Goals", common.BigFont, color.White),
	)
	rootContainer.AddChild(label)
	for i, v := range goals {
		title := widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("Goal %d", i), common.NormalFont, color.White),
			widget.TextOpts.Insets(widget.NewInsetsSimple(20)),
		)
		goal := widget.NewText(
			widget.TextOpts.Text(v.String(), common.TinyFont, color.NRGBA{200, 250, 250, 255}),
			widget.TextOpts.Insets(widget.NewInsetsSimple(20)),
		)
		price := widget.NewText(
			widget.TextOpts.Text(v.Price(), common.TinyFont, color.NRGBA{200, 250, 250, 255}),
			widget.TextOpts.Insets(widget.NewInsetsSimple(20)),
		)
		rootContainer.AddChild(title)
		rootContainer.AddChild(goal)
		rootContainer.AddChild(price)
	}

	return &goalsUI{
		ui: ebitenui.UI{
			Container: rootContainer,
		},
	}
}
