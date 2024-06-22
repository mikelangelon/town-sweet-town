package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/textbox"
	"github.com/mikelangelon/town-sweet-town/world"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"github.com/solarlune/resolv"
	"image/color"
	"time"
)

type Town struct {
	BaseScene

	endOfDay        *endOfDay
	TransitionSleep uint8
}

func NewTown(id string, mapScene *graphics.MapScene) *Town {
	return &Town{
		BaseScene: BaseScene{ID: id, MapScene: mapScene},
	}
}

func (t *Town) Update() error {
	if t.endOfDay != nil {
		t.endOfDay.Update()
		if t.endOfDay.done {
			t.endOfDay = nil
			t.state.Status = DayStarting
			t.state.GameLogic.NextDay(t.state)
		}
	}
	skip, err := t.BaseScene.Update()
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	action := t.checkActionExecuted()
	if action != nil {
		t.Action(action)
	}
	return nil
}

func (t *Town) Draw(screen *ebiten.Image) {
	t.BaseScene.Draw(screen)
	if t.state.Status == DayEnding {
		colorGoal := color.RGBA{10, 10, 10, t.TransitionSleep}
		if t.TransitionSleep < 200 {
			t.TransitionSleep++
		} else {
			if t.endOfDay == nil {
				t.endOfDay = createShowEndOfDay(t.NPCs, t.state.Day, t.state.Stats)
			}
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		bg := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
		bg.Fill(colorGoal)
		screen.DrawImage(bg, op)
	}
	if t.state.Status == DayStarting {
		colorGoal := color.RGBA{0, 0, 0, t.TransitionSleep}
		if t.TransitionSleep > 1 {
			t.TransitionSleep--
		} else {
			t.state.Day++
			t.state.Status = Playing
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		bg := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
		bg.Fill(colorGoal)
		screen.DrawImage(bg, op)
	}
	if t.endOfDay != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(50, 150)
		bg := ebiten.NewImage(500, 300)
		t.endOfDay.ui.Draw(bg)
		screen.DrawImage(bg, op)
	}
}

func (t *Town) Action(collision *resolv.Collision) {
	if c, ok := collision.Objects[0].Data.(*npc.NPC); ok {
		t.KickOutHouse(c)
	}
	if _, ok := collision.Objects[0].Data.(world.Fire); ok {
		t.FireAction()
	}
	if c, ok := collision.Objects[0].Data.(house.Signal); ok {
		t.SignalAction(c)
	}

}

func (t *Town) FireAction() {
	t.Text.ShowAndQuestion(
		[]string{"Go to the next day?"},
		[]string{"Yes", textbox.No},
		func(answer string) {
			if answer == "Yes" {
				t.state.Status = DayEnding
			} else {
				t.state.Status = Playing
			}
		},
	)
}

func (t *Town) SignalAction(signal house.Signal) {
	options := house.MapHouseBulding.GiveMeThree()
	t.Text.ShowAndQuestion(
		[]string{"Which house do you want to build?"},
		append(options,
			textbox.NoResponse),
		func(answer string) {
			info := house.MapHouseBulding[answer]
			if t.state.Stats["money"] < info.Cost {
				t.Text.ShowAndQuestion([]string{"", "Not enough money. Sorry"}, nil, func(s string) {})
				return
			}
			t.state.Stats["money"] -= info.Cost
			house := t.state.GameLogic.CreateHouse(signal.ID, info.Type)
			house.House.Offset = signal.HousePlace
			t.MapScene.Child = append(t.MapScene.Child, &house.House)
		},
	)
}

func (t *Town) KickOutHouse(npc *npc.NPC) {
	options := []string{"Sorry, leave the house", textbox.NoResponse}
	answerFunc := func(answer string) {

		if answer == "Sorry, leave the house" {
			for _, v := range t.state.World["town1"].Houses {
				if v.ID == npc.House.ID {
					v.Owner = nil
					npc.SetHouse(nil, 0)
					break
				}
			}
			t.state.World["town1"].RemoveNPC(npc.ID)
			newNpc := *npc
			newNpc.X = 16 * 20
			newNpc.Y = 16 * 6
			newNpc.Move = &common.Position{X: 16 * 6, Y: 16 * 6}
			t.state.World["people"].AddNPC(&newNpc)
			npc.Move = &common.Position{X: common.ScreenWidth + 16, Y: npc.Y}

		}
	}

	t.Text.ShowAndQuestion([]string{"How can I help you?"}, options, answerFunc)
}

func (t *Town) Load(st State, sm stagehand.SceneController[State]) {
	t.BaseScene.Load(st, sm)

	for _, v := range t.state.World["town1"].Houses {
		t.MapScene.Child = append(t.MapScene.Child, &v.House)
	}

	if t.state.Status == InitialState {
		t.state.Status = Playing
		t.state.Stats = make(map[string]int)
		t.state.Stats[npc.Money] = 13
		t.state.Stats[npc.Happiness] = 10
		t.state.Stats[npc.Security] = 15
		t.state.Stats[npc.Food] = 10
		t.state.Stats[npc.Health] = 30
		t.state.Day = 1
		return
	}

	timer := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer.C
		t.state.Player.X, t.state.Player.Y = t.TransitionPoints.Position.X, t.TransitionPoints.Position.Y
	}()
}
