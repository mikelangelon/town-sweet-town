package scenes

import (
	"fmt"
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
	"log/slog"
	"math/rand"
	"time"
)

type Town struct {
	BaseScene

	endOfDay        *endOfDay
	TransitionSleep uint8
	MenuScene       Transition
	// For now, only for initial presentation
	Ack *time.Time
}

func NewTown(id string, mapScene *graphics.MapScene) *Town {
	return &Town{
		BaseScene: BaseScene{ID: id, MapScene: mapScene, ui: NewHUI()},
	}
}

func (t *Town) Update() error {
	if t.uiUpdate() {
		return nil
	}
	if !t.Text.Visible() && t.state.Status == GoingToMenu {
		t.sm.SwitchWithTransition(t.MenuScene.Scene, stagehand.NewTicksTimedSlideTransition[State](t.MenuScene.Direction, time.Millisecond*time.Duration(200)))
		t.state.Day = 0
		t.state.Status = Menu
	}
	if !t.Text.Visible() && t.state.Status == HappyEnd {
		t.Text.ShowAndQuestion([]string{"Today is the last day of the 2 weeks period", "You beat the game! Congratulations!"}, nil, func(s string) {})
		t.state.Status = GoingToMenu
	}
	if !t.Text.Visible() && t.state.Status == GameOver {
		t.Text.ShowAndQuestion([]string{"Today is the last day of the 2 weeks period", "Sorry... you didn't make it"}, nil, func(s string) {})
		t.state.Status = GoingToMenu
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
			"There are some active rules that you will need to be aware \nto improve your town stats.\nPress ENTER to see them.",
			"Most of them are based on villager stats. \nCombine them to improve your village!",
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
		t.goalsUI = NewGoals(t.state.Goals)
	}
	if t.endOfDay != nil {
		t.endOfDay.Update()
		if t.endOfDay.done {
			for _, v := range t.NPCs {
				random := rand.Intn(100)
				fmt.Printf("NPC %s --> random of %d, with  nitpicky of %d, adapted as %d\n", v.ID, random, v.NitPicky, v.AdaptNitpicky())
				if random < v.AdaptNitpicky() {
					values := []string{npc.Cultural, npc.Health, npc.Security, npc.Happiness}
					v.Wishes = append(v.Wishes, npc.Wish{
						DayStart:  t.state.Day + 1,
						DayEnd:    t.state.Day + 3,
						Stat:      values[rand.Intn(4)],
						Value:     8 + t.state.Day*v.NitPickyLevel,
						Happiness: 3,
					})
					fmt.Printf("Adding wish stat %s value %d\n", v.Wishes[len(v.Wishes)-1].Stat, v.Wishes[len(v.Wishes)-1].Value)
				}
			}
			step, end := t.state.GameLogic.GetRuler().CheckGoals(t.state.Goals, t.state.Day, t.state.Stats)
			switch end {
			case 0:
				if step != nil {
					t.Text.ShowAndQuestion([]string{step[0].Text}, nil, func(s string) {})
					t.state.Stats[step[0].Name] = step[0].Value
				}
			case 1:
				t.state.Status = HappyEnd
				t.endOfDay = nil
				return nil
			case -1:
				t.state.Status = GameOver
				t.endOfDay = nil
				return nil
			}
			var byebyeNPCsMessages []string
			removeNPC := func(id string, message string) {
				for _, h := range t.state.World["town1"].Houses {
					if h.Owner != nil && *h.Owner == id {
						h.Owner = nil
					}
				}
				for i, j := range t.NPCs {
					if j.ID == id {
						byebyeNPCsMessages = append(byebyeNPCsMessages, fmt.Sprintf("%s %s", id, message))
						t.NPCs = append(t.NPCs[0:i], t.NPCs[i+1:]...)
						break
					}
				}
				t.state.World["town1"].RemoveNPC(id)
			}
			// Dying NPCs due to food

			var dyingNPC string
			if t.state.Stats[npc.Food] < 0 {
				dyingNPC = t.NPCs[rand.Intn(len(t.NPCs))].ID
				removeNPC(dyingNPC, "died due to the lack of food")
				byebyeNPCsMessages = append(byebyeNPCsMessages, "Good news! Now you have additional food!")
				t.state.Stats[npc.Food] += 5
			}
			// Leaving NPCs due to wishes
			for _, v := range t.endOfDay.leavingNPCs {
				removeNPC(v.ID, "was sad and left the village")
			}
			if len(byebyeNPCsMessages) > 0 {
				t.Text.ShowAndQuestion(byebyeNPCsMessages, nil, func(s string) {})
			}

			//t.leavingNPCs = t.endOfDay.leavingNPCs
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
				t.NPCs = deduplicate(t.NPCs)
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
			t.state.Status = CheckWishes
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		bg := ebiten.NewImage(common.ScreenWidth, common.ScreenHeight)
		bg.Fill(colorGoal)
		screen.DrawImage(bg, op)
	}
	if t.state.Status == CheckWishes {
		wishes := t.NPCs.NewWishesString(t.state.Day, t.state.Stats)
		if len(wishes) > 0 {
			t.Text.ShowAndQuestion(wishes.String(), nil, func(s string) {})
		}
		t.state.Status = Playing
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
	if c, ok := collision.Objects[0].Data.(*house.Signal); ok {
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

func (t *Town) SignalAction(signal *house.Signal) {
	if signal.HouseOptions == nil {
		options := house.MapHouseBulding.GiveMeThree()
		signal.HouseOptions = &options
	}
	var question = "Which house do you want to build?"
	var existingHouse *house.House
	for _, v := range t.state.World["town1"].Houses {
		if v.House.Offset == signal.HousePlace {
			existingHouse = v
			question = "Which upgrade do you want?"
		}
	}

	t.Text.ShowAndQuestion(
		[]string{question},
		append(*signal.HouseOptions,
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
			signal.HouseOptions = house.MapHouseBulding.ReplaceThree(*signal.HouseOptions, answer)
			t.state.Stats["money"] -= info.Cost
			newHouse := t.state.GameLogic.CreateHouse(fmt.Sprintf("%s %s", info.Name, signal.ID), info.Type)
			newHouse.House.Offset = signal.HousePlace
			newHouse.Type = info.Type
			if existingHouse == nil {
				t.MapScene.Child = append(t.MapScene.Child, &newHouse.House)
				t.state.World["town1"].Houses = append(t.state.World["town1"].Houses, &newHouse)
				return
			}
			existingHouse.House = newHouse.House
			existingHouse.Type = newHouse.Type

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
			slog.With("id", npc.ID).Info("adding npc in entrance")
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
		now := time.Now()
		t.Ack = &now
		return
	}
}

func (t *Town) PostTransition(st State, original stagehand.Scene[State]) {
	if t.state.Status == InitialState {
		return
	}
	t.state.Player.X, t.state.Player.Y = t.TransitionPoints.Position.X, t.TransitionPoints.Position.Y
	t.state.Status = Playing
}

func deduplicate(n npc.NPCs) npc.NPCs {
	var unique = make(map[string]*npc.NPC)
	for _, v := range n {
		if _, ok := unique[v.ID]; !ok {
			unique[v.ID] = v
		}
	}
	var result npc.NPCs
	for _, v := range unique {
		result = append(result, v)
	}
	return result
}
