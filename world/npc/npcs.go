package npc

type NPCs []*NPC

func (n NPCs) GetNPC(id string) *NPC {
	for _, v := range n {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func addSteps(steps []StatStep, v int, charID *string, name, text string) []StatStep {
	return append(steps, StatStep{
		Name:   name,
		CharID: charID,
		Value:  v,
		Text:   text,
	})
}
