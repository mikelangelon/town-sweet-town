package npc

import "fmt"

type NPCs []*NPC

func (n NPCs) GetNPC(id string) *NPC {
	for _, v := range n {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (n NPCs) OldWishes(day int) CompleteWishes {
	var result []CompleteWish
	for _, v := range n {
		if v.House == nil {
			continue
		}
		for _, j := range v.Wishes {
			if j.DayEnd == day {
				result = append(result, CompleteWish{
					Wish: j,
					NPC:  v,
				})
			}
		}
	}
	return result
}

func (n NPCs) NewWishesString(day int, stats map[string]int) CompleteWishes {
	var result []CompleteWish
	for _, v := range n {
		if v.House == nil {
			continue
		}
		for i, j := range v.Wishes {
			if j.DayStart == day {
				if stats[j.Stat] >= j.Value {
					v.Wishes = append(v.Wishes[0:i], v.Wishes[i+1:]...)
					break
				}
				result = append(result, CompleteWish{
					Wish: j,
					NPC:  v,
				})
			}
		}
	}
	return result
}

func (n NPCs) ApplyWishesString(day int) []string {
	var result []string
	for _, v := range n {
		if v.House == nil {
			continue
		}
		for _, j := range v.Wishes {
			if j.DayStart == day {
				result = append(result, fmt.Sprintf("%s wishes the town should have more %s", v.ID, j.Stat))
			}
		}
	}
	return result
}
func addSteps(steps []StatStep, v int, charID *string, name, text string) []StatStep {
	return append(steps, StatStep{
		Name:   name,
		CharID: charID,
		Value:  v,
		Text:   text,
	})
}
