package scenes

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	imageNine "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"image/color"
	"time"
)

type endOfDay struct {
	ui            *ebitenui.UI
	endOfDaySteps []npc.StatStep
	endOfDayIndex int
	endOfDayTimer *time.Time
	mapStats      map[string]*widget.Text
	mapProgress   map[string]*widget.ProgressBar
	currentStuff  *widget.Text

	stats map[string]int
	done  bool
}

func (e *endOfDay) Update() {
	if e == nil {
		return
	}
	if e.endOfDayTimer != nil && time.Since(*e.endOfDayTimer) > 1*time.Second {
		if e.endOfDayIndex < len(e.endOfDaySteps) {
			t := time.Now()
			e.endOfDayTimer = &t
			fmt.Println(e.endOfDaySteps[e.endOfDayIndex].FormatText())
			e.currentStuff.Label = e.endOfDaySteps[e.endOfDayIndex].FormatText()
			name := e.endOfDaySteps[e.endOfDayIndex].Name
			e.stats[name] += e.endOfDaySteps[e.endOfDayIndex].Value
			switch name {
			case npc.Security, npc.Health, npc.Cultural, npc.Happiness:
				e.mapProgress[name].SetCurrent(e.stats[name])
			case npc.Food:
				e.mapStats[name].Label = fmt.Sprintf("Food: %d", e.stats[name])
			case npc.Money:
				e.mapStats[name].Label = fmt.Sprintf("Money: %d", e.stats[name])
			}
			e.endOfDayIndex++
		} else {
			e.done = true
		}
	}
}

func createShowEndOfDay(npcs npc.NPCs, day int, stats map[string]int) *endOfDay {

	var e endOfDay
	e.endOfDaySteps = npcs.AllSteps()
	t := time.Now()
	e.endOfDayTimer = &t

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
		widget.TextOpts.Text(fmt.Sprintf("Day %d - Results", day), face, color.White),
	)
	e.stats = stats
	e.mapStats = make(map[string]*widget.Text)
	e.mapProgress = make(map[string]*widget.ProgressBar)
	e.mapStats[npc.Money] = widget.NewText(
		widget.TextOpts.Text(moneyText(stats[npc.Money]), face, color.White),
	)
	e.mapStats[npc.Food] = widget.NewText(
		widget.TextOpts.Text(foodText(stats[npc.Food]), face, color.White),
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
	e.currentStuff = widget.NewText(
		widget.TextOpts.Text("...", face, color.White),
	)

	e.mapProgress[npc.Security] = progress(stats[npc.Security])
	e.mapProgress[npc.Happiness] = progress(stats[npc.Happiness])
	e.mapProgress[npc.Health] = progress(stats[npc.Health])
	e.mapProgress[npc.Cultural] = progress(stats[npc.Cultural])
	rootContainer.AddChild(label1)
	rootContainer.AddChild(secondaryContainer)
	secondaryContainer.AddChild(leftContainer)
	secondaryContainer.AddChild(rightContainer)
	leftContainer.AddChild(e.mapStats[npc.Money])
	leftContainer.AddChild(lSecurity)
	leftContainer.AddChild(e.mapProgress[npc.Security])
	leftContainer.AddChild(lCultural)
	leftContainer.AddChild(e.mapProgress[npc.Cultural])
	rightContainer.AddChild(e.mapStats[npc.Food])
	rightContainer.AddChild(lHappiness)
	rightContainer.AddChild(e.mapProgress[npc.Happiness])
	rightContainer.AddChild(labelHealth)
	rightContainer.AddChild(e.mapProgress[npc.Health])
	rootContainer.AddChild(e.currentStuff)

	e.ui = &ui
	return &e
}

func moneyText(money int) string {
	return fmt.Sprintf("Money: %d", money)
}
func foodText(food int) string {
	return fmt.Sprintf("Food: %d", food)
}
