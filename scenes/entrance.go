package scenes

import (
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/textbox"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"github.com/solarlune/resolv"
	"log/slog"
)

type Entrance struct {
	BaseScene
}

func NewEntrance(id string, mapScene *graphics.MapScene) *Entrance {
	return &Entrance{
		BaseScene{ID: id, MapScene: mapScene, ui: NewHUI()},
	}
}

func (e *Entrance) Update() error {
	if e.uiUpdate() {
		return nil
	}
	skip, err := e.BaseScene.Update()
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	action := e.checkActionExecuted()
	if action != nil {
		e.Action(action)
	}
	return nil
}

func (e *Entrance) Action(collision *resolv.Collision) {
	if c, ok := collision.Objects[0].Data.(*npc.NPC); ok {
		e.TalkToNPC(c)
	}
}

func (e *Entrance) TalkToNPC(npc *npc.NPC) {
	answerFunc := func(answer string) {
		if answer != textbox.NoResponse && answer != textbox.No {
			for _, v := range e.state.World["town1"].Houses {
				if v.ID == answer {
					npc.SetHouse(v, e.state.Day)
					v.Owner = &npc.ID
				}
			}
			npc.Move = &common.Position{X: -16, Y: npc.Y}
			e.state.World["people"].RemoveNPC(npc.ID)
			slog.With("id", npc.ID).Info("removing npc in entrance")
			e.state.World["town1"].AddNPC(npc)
			slog.With("id", npc.ID).Info("adding npc in town")
		}
	}
	var options []string
	for _, v := range e.state.World["town1"].Houses {
		if v.Owner != nil {
			continue
		}
		options = append(options, v.ID)
	}

	if npc.DayIn == e.state.Day {
		e.Text.TalkToNPC(npc.ID, npc.Talk(e.state.Day))
	} else {
		options = append(options, textbox.NoResponse)
		e.Text.ShowAndQuestionNPC(npc.ID, npc.Talk(e.state.Day), options, answerFunc)
	}
}

func (e *Entrance) Load(st State, sm stagehand.SceneController[State]) {
	e.BaseScene.Load(st, sm)
}

func (e *Entrance) PreTransition(destination stagehand.Scene[State]) State {
	return e.state
}

func (e *Entrance) PostTransition(st State, original stagehand.Scene[State]) {
	e.state.Player.X, e.state.Player.Y = e.TransitionPoints.Position.X, e.TransitionPoints.Position.Y
	e.state.Status = Playing
}
