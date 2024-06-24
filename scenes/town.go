package scenes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	// For now, only for initial presentation
	Ack *time.Time
}

func NewTown(id string, mapScene *graphics.MapScene) *Town {
	return &Town{
		BaseScene: BaseScene{ID: id, MapScene: mapScene, ui: NewHUI()},
	}
}

func (t *Town) Update() error {
	if t.rulesUI != nil {
		t.rulesUI.ui.Update()

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			t.state.Status = Playing
			t.rulesUI = nil
		}

		return nil
	}
	if t.state.Status == InitialState {
		_, err := t.playerUpdate()
		if err != nil {
			return err
		}
	}
	if t.state.Status == InitialState && t.Ack != nil && time.Since(*t.Ack) > 5*time.Second {
		t.state.Status = InitExplanation
		t.Text.Show([]string{
			"Hello there! \nLet me introduce you a bit on the rules of this world. \nYou can move with the arrows\nPress ENTER to do an action.\n(As passing this dialog)",
			"You are the architect of this town,\nyour goal is to build a happy community in 2 weeks.",
			"To build a house, press ENTER next to a signal post.\nYou will need money for it.\n(See top values)",
			"To recruit villagers for your town,\ngo to the east.",
			"Every 2 days new people will come there.",
			"To end the day,\npress ENTER next to the fire.",
			"Villagers have their own characteristics, \ntheir combinations would make improve or decrease \nyour town stats.",
			"There are some active rules that you will need to be aware.\nPress ENTER to see them.",
			"Every day new rules\nwill come up!",
			"You have 2 weeks! Don't fail me!",
		})
	}
	if t.state.Status == InitExplanation && !t.Text.Visible() {
		time.Sleep(500 * time.Millisecond)
		t.Text.Show([]string{"And put some clothes on!", "I have faith in you!"})
		t.state.Status = NoClothes
	}
	if t.state.Status == NoClothes && !t.Text.Visible() {
		t.state.Status = Playing
		t.state = t.state.GameLogic.ChangePlayer(t.state)
	}
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
	if t.rulesUI != nil {
		return nil
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
				t.endOfDay = createShowEndOfDay(t.state.GameLogic.GetRuler(), t.NPCs, t.state.Day, t.state.Stats)
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
			if answer == textbox.NoResponse {
				return
			}
			info := house.MapHouseBulding[answer]
			if t.state.Stats["money"] < info.Cost {
				t.Text.ShowAndQuestion([]string{"", "Not enough money. Sorry"}, nil, func(s string) {})
				return
			}
			t.state.Stats["money"] -= info.Cost
			newHouse := t.state.GameLogic.CreateHouse(fmt.Sprintf("%s %s", info.Name, signal.ID), info.Type)
			newHouse.House.Offset = signal.HousePlace
			t.MapScene.Child = append(t.MapScene.Child, &newHouse.House)
			t.state.World["town1"].Houses = append(t.state.World["town1"].Houses, &newHouse)
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

	t.Text.ShowAndQuestionNPC(npc.ID,
		append(npc.Sentences(), "How can I help you?"),
		options, answerFunc)
}

func (t *Town) PreTransition(destination stagehand.Scene[State]) State {
	return t.state
}

func (t *Town) Load(st State, sm stagehand.SceneController[State]) {
	t.BaseScene.Load(st, sm)

	for _, v := range t.state.World["town1"].Houses {
		t.MapScene.Child = append(t.MapScene.Child, &v.House)
	}

	if t.state.Status == InitialState {
		t.state.Stats = make(map[string]int)
		t.state.Stats[npc.Money] = 20
		t.state.Stats[npc.Happiness] = 10
		t.state.Stats[npc.Security] = 15
		t.state.Stats[npc.Food] = 10
		t.state.Stats[npc.Health] = 30
		t.state.Day = 1

		now := time.Now()
		t.Ack = &now
		return
	}
}

func (t *Town) PostTransition(st State, original stagehand.Scene[State]) {
	t.state.Player.X, t.state.Player.Y = t.TransitionPoints.Position.X, t.TransitionPoints.Position.Y
	t.state.Status = Playing
}
