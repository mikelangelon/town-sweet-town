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
		// the container will use an grid layout to layout its ScrollableContainer and Slider
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
	//Create the new ScrollContainer object
	scrollContainer := widget.NewScrollContainer(
		//Set the content that will be scrolled
		widget.ScrollContainerOpts.Content(content),
		//Tell the container to stretch the content width to match available space
		widget.ScrollContainerOpts.StretchContentWidth(),
		//Set the background images for the scrollable container
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0}),
			Mask: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)
	//Add the scrollable container to the left side of the window
	secondaryContainer.AddChild(scrollContainer)

	//Create a function to return the page size used by the slider
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}
	//Create a vertical Slider bar to control the ScrollableContainer
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		//On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
	)
	//Set the slider's position if the scrollContainer is scrolled by other means than the slider
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	//Add the slider to the second slot in the root container
	secondaryContainer.AddChild(vSlider)
	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Active rules", common.BigFont, color.White),
	))
	rootContainer.AddChild(secondaryContainer)
	return &ruleUI{
		ui: ebitenui.UI{
			Container: rootContainer,
		},
	}
}
