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

func (n NPCs) AllSteps() []StatStep {
	var steps []StatStep
	steps = append(steps, n.Food()...)
	steps = append(steps, n.Cultural()...)
	steps = append(steps, n.Health()...)
	steps = append(steps, n.Security()...)
	steps = append(steps, n.Happiness()...)
	return steps
}
func (n NPCs) Food() []StatStep {
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Cooking]; ok1 {
			steps = addSteps(steps, v1, "food", fmt.Sprintf("FOOD %s -> Cooking: %d", v.ID, v1))
		}
		if v1, ok1 := m[Eating]; ok1 {
			steps = addSteps(steps, v1, "food", fmt.Sprintf("FOOD %s -> Eating: -%d", v.ID, -v1))
		}
	}
	return steps
}

func (n NPCs) Cultural() []StatStep {
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Reading]; ok1 {
			steps = addSteps(steps, v1, "cultural", fmt.Sprintf("CULTURAL %s -> Reading: %d", v.ID, v1))
		}
		if v1, ok1 := m[Reading]; ok1 {
			steps = addSteps(steps, v1, "cultural", fmt.Sprintf("CULTURAL %s -> Reading: %d", v.ID, v1))
		}
	}
	return steps
}

func (n NPCs) Health() []StatStep {
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Workaholic]; ok1 && v1 > 7 {
			steps = addSteps(steps, v1, "health", fmt.Sprintf("HEALTH %s -> Workalcoholic too much: -%d", v.ID, -v1))
		}
		if v1, ok1 := m[Eating]; ok1 && v1 > 7 {
			steps = addSteps(steps, v1, "health", fmt.Sprintf("HEALTH %s -> Eating too much: -%d", v.ID, -v1))
		}
		if v1, ok1 := m[Sports]; ok1 {
			steps = addSteps(steps, v1, "health", fmt.Sprintf("HEALTH %s -> Sport: %d", v.ID, v1))
		}
	}
	return steps
}
func (n NPCs) Security() []StatStep {
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Stuff]; ok1 {
			if v2, ok1 := m[Optimistic]; ok1 {
				steps = addSteps(steps, -(v1 * v2), "security", fmt.Sprintf("SECURITY %s -> Optimistic & Stuff: %d", v.ID, -(v1*v2)))
			}
			if v2, ok1 := m[Adventurous]; ok1 {
				steps = addSteps(steps, -(v1 * v2), "security", fmt.Sprintf("SECURITY %s -> Adventurous & Stuff: %d", v.ID, -(v1*v2)))
			}
		}

	}
	return steps
}

func (n NPCs) Happiness() []StatStep {
	var steps []StatStep

	var animals []int
	var competitive []int
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Animals]; ok1 {
			animals = append(animals, v1)
		}
	}
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Competitive]; ok1 {
			competitive = append(competitive, v1)
		}
	}
	if len(animals) > 1 {
		steps = addSteps(steps, 14, "happiness", fmt.Sprintf("HAPPINESS -> Multiple animals: %d", 14))
	}
	if len(competitive) > 1 {
		steps = addSteps(steps, 10, "happiness", fmt.Sprintf("HAPPINESS -> Multiple competitive: -%d", -10))
	}
	return steps
}

func addSteps(steps []StatStep, v int, name, text string) []StatStep {
	return append(steps, StatStep{
		Name:  name,
		Value: v,
		Text:  text,
	})
}
