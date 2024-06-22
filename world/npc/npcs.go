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
	steps = append(steps, n.Money()...)
	steps = append(steps, n.Food()...)
	steps = append(steps, n.Cultural()...)
	steps = append(steps, n.Health()...)
	steps = append(steps, n.Security()...)
	steps = append(steps, n.Happiness()...)
	return steps
}

func (n NPCs) Money() []StatStep {
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok := m[Rent]; ok {
			steps = addSteps(steps, v1, &v.ID, Money, "Rent")
		}
	}
	return steps
}
func (n NPCs) Food() []StatStep {
	var stat = Food
	var steps []StatStep
	steps = addSteps(steps, -2*(len(n)+1), nil, stat, fmt.Sprintf("%d Villagers + yourself", len(n)))
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Cooking]; ok1 {
			steps = addSteps(steps, v1, &v.ID, stat, "Cooking")
		}
		if v1, ok1 := m[Eating]; ok1 {
			steps = addSteps(steps, v1, &v.ID, stat, "Eating")
		}
	}
	return steps
}

func (n NPCs) Cultural() []StatStep {
	var stat = Cultural
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Reading]; ok1 {
			steps = addSteps(steps, v1, &v.ID, stat, "Reading")
		}
		if v1, ok1 := m[Music]; ok1 {
			steps = addSteps(steps, v1, &v.ID, stat, "Music")
		}
	}
	return steps
}

func (n NPCs) Health() []StatStep {
	var stat = Health
	var steps []StatStep
	for _, v := range n {
		m := v.Chars.charLevelMap()
		if v1, ok1 := m[Workaholic]; ok1 && v1 > 7 {
			steps = addSteps(steps, v1, &v.ID, stat, "Workalcoholic too much")
		}
		if v1, ok1 := m[Eating]; ok1 && v1 > 7 {
			steps = addSteps(steps, v1, &v.ID, stat, "Eating too much")
		}
		if v1, ok1 := m[Sports]; ok1 {
			steps = addSteps(steps, v1, &v.ID, stat, "Sport")
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
				steps = addSteps(steps, -(v1 * v2), &v.ID, Security, "Optimistic & Stuff")
			}
			if v2, ok1 := m[Adventurous]; ok1 {
				steps = addSteps(steps, -(v1 * v2), &v.ID, Security, "Adventurous & Stuff")
			}
		}

	}
	return steps
}

func (n NPCs) Happiness() []StatStep {
	var stat = Happiness
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
		steps = addSteps(steps, 14, nil, stat, "Multiple animals")
	}
	if len(competitive) > 1 {
		steps = addSteps(steps, -10, nil, stat, "Multiple competitive")
	}
	return steps
}

func addSteps(steps []StatStep, v int, charID *string, name, text string) []StatStep {
	return append(steps, StatStep{
		Name:   name,
		CharID: charID,
		Value:  v,
		Text:   text,
	})
}
