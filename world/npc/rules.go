package npc

import "fmt"

type RuleApplier struct {
	Rules []Rule
}

func (r RuleApplier) ApplyRules(n NPCs) []StatStep {
	var steps []StatStep
	for _, v := range r.Rules {
		steps = v.Func(n, steps)
	}
	return steps
}

type Rule struct {
	Name        string
	Description string
	Func        RuleFunc
}
type RuleFunc func(n NPCs, steps []StatStep) []StatStep

var (
	EatingRule = Rule{
		Name:        "Being alive",
		Description: "-2 units of Food per villager & yourself",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			return addSteps(steps, -2*(len(n)+1), nil, Food, fmt.Sprintf("%d Villagers + yourself", len(n)))
		},
	}

	NoFood = Rule{
		Name:        "No food",
		Description: "-20 Health if there is no food",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			return addSteps(steps, -2*(len(n)+1), nil, Food, fmt.Sprintf("%d Villagers + yourself", len(n)))
		},
	}

	RentRule = Rule{
		Name:        "Rent",
		Description: "Every villager pays their rent (+X money)",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charLevelMap()
				if v1, ok := m[Rent]; ok {
					steps = addSteps(steps, v1, &v.ID, Money, "Rent")
				}
			}
			return steps
		},
	}

	CookingBonus = Rule{
		Name:        "Cooking",
		Description: "+8 food for every villager that likes cooking",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Cooking]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 8, &v.ID, Food, "Cooking")
				}
			}
			return steps
		},
	}

	EatingTooMuch = Rule{
		Name:        "Eating too much",
		Description: "-4 food for every villager that likes eatings",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Eating]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, -4, &v.ID, Food, "Eating too much")
				}
			}
			return steps
		},
	}

	GoodCulture = Rule{
		Name:        "Good culture",
		Description: "+4 culture for every villager that likes Reading or Music",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Reading]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 4, &v.ID, Cultural, "Reading")
				}
				if v1, ok1 := m[Music]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 4, &v.ID, Cultural, "Music")
				}
			}
			return steps
		},
	}

	BadCulture = Rule{
		Name:        "Bad culture",
		Description: "-4 culture for every villager that hates Reading or Music",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Reading]; ok1 && v1.Love() == Hate {
					steps = addSteps(steps, -4, &v.ID, Cultural, "Reading")
				}
				if v1, ok1 := m[Music]; ok1 && v1.Love() == Hate {
					steps = addSteps(steps, -4, &v.ID, Cultural, "Music")
				}
			}
			return steps
		},
	}

	WorkTooMuch = Rule{
		Name:        "Work too much",
		Description: "-8 Health if a villager likes to work",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Workaholic]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, -8, &v.ID, Health, "Works too much")
				}
			}
			return steps
		},
	}

	HealthyGuy = Rule{
		Name:        "Healthy Sport",
		Description: "+4 Health if a villager likes to sport",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Sports]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 4, &v.ID, Health, "Healthy sport")
				}
			}
			return steps
		},
	}

	OptimisticThief = Rule{
		Name:        "Optimistic & Stuff",
		Description: "-10 Security if a villager likes Stuff & is optimist",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				v1, ok1 := m[Stuff]
				v2, ok2 := m[Optimistic]
				if ok1 && ok2 && v1.Love() == Love && v2.Love() == Love {
					steps = addSteps(steps, -10, &v.ID, Security, "Optimistic & Stuff")
				}
			}
			return steps
		},
	}

	AdventurousThief = Rule{
		Name:        "Adventurous & Stuff",
		Description: "-10 Security if a villager likes Stuff & is Adventurous",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				v1, ok1 := m[Stuff]
				v2, ok2 := m[Adventurous]
				if ok1 && ok2 && v1.Love() == Love && v2.Love() == Love {
					steps = addSteps(steps, -10, &v.ID, Security, "Adventurous & Stuff")
				}
			}
			return steps
		},
	}

	AnimalLovers = Rule{
		Name:        "Animal Lovers",
		Description: "+10 Happiness if more than 1 villager likes animals",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var animals []bool
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Animals]; ok1 && v1.Love() == Love {
					animals = append(animals, true)
				}
			}
			if len(animals) > 1 {
				steps = addSteps(steps, 10, nil, Happiness, "Animal lovers")
			}
			return steps
		},
	}

	CompetitionTooMuch = Rule{
		Name:        "Too Competitive",
		Description: "-12 Happiness if more than 1 villager is competitive",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var result []bool
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Competitive]; ok1 && v1.Love() == Love {
					result = append(result, true)
				}
			}
			if len(result) > 1 {
				steps = addSteps(steps, -12, nil, Happiness, "Too Competitive")
			}
			return steps
		},
	}

	AnimalConflict = Rule{
		Name:        "Animal conflict",
		Description: "-20 Happiness if villagers like & hate animals",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var animals bool
			var animalHaters bool
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Animals]; ok1 {
					switch v1.Love() {
					case Love:
						animals = true
					case Hate:
						animalHaters = true
					}
				}
			}
			if animals && animalHaters {
				steps = addSteps(steps, -20, nil, Happiness, "Animal conflict")
			}
			return steps
		},
	}
)
