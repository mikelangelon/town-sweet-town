package scenes

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"image/color"
	"math"
)

type ruleUI struct {
	ui ebitenui.UI
}

func NewRulesUI(rules []npc.Rule) *ruleUI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)
	secondaryContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 0),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:   true,
				Position:  widget.RowLayoutPositionCenter,
				MaxWidth:  600,
				MaxHeight: 450,
			}),
		),
	)

	//Create the container with the content that should be scrolled
	content := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(10),
	)))

	for _, v := range rules {
		rule := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
		label := widget.NewText(
			widget.TextOpts.Text(v.Name, common.NormalFont, color.White),
		)
		label2 := widget.NewText(
			widget.TextOpts.Text(v.Description, common.MegaTinyFont, color.NRGBA{200, 250, 250, 255}),
			widget.TextOpts.Insets(widget.NewInsetsSimple(20)))
		rule.AddChild(label)
		rule.AddChild(label2)
		content.AddChild(rule)
	}

	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0}),
			Mask: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)

	secondaryContainer.AddChild(scrollContainer)

	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(

			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},

			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
	)
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})
	secondaryContainer.AddChild(vSlider)
	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Active rules", common.NormalFont, color.White),
	))
	rootContainer.AddChild(secondaryContainer)
	return &ruleUI{
		ui: ebitenui.UI{
			Container: rootContainer,
		},
	}
}
