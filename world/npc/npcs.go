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
