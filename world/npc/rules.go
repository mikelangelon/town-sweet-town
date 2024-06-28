package npc

import (
	"fmt"
	"github.com/mikelangelon/town-sweet-town/world"
	"math/rand"
)

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

func (r RuleApplier) CheckGoals(goals []world.Goal, currentDay int, stats map[string]int) ([]StatStep, int) {
	var endStatus = 0
	for i, v := range goals {
		if v.Day != currentDay {
			continue
		}
		if stats[v.Stat] >= v.Value {
			if v.Mandatory {
				endStatus = 1
			}
			return []StatStep{
				{
					Name:  v.GiftStat,
					Value: v.Value,
					Text:  fmt.Sprintf("You made Goal %d! (Get to %s %d)", i+1, v.Stat, v.Value),
				},
			}, endStatus
		}
		if v.Mandatory {
			endStatus = -1
		}
		return []StatStep{
			{
				Name:  v.GiftStat,
				Value: 0,
				Text:  fmt.Sprintf("You failed Goal %d! (Get to %s %d)", i+1, v.Stat, v.Value),
			},
		}, endStatus
	}
	return nil, endStatus
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

	RoofRule = Rule{
		Name:        "Having a roof",
		Description: "+1 happiness for every villager",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			return addSteps(steps, len(n), nil, Happiness, fmt.Sprintf("%d Villagers having a house", len(n)))
		},
	}

	CompleteTownRule = Rule{
		Name:        "Town completed",
		Description: "+3 happiness if you have 4 villagers",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			if len(n) > 4 {
				return addSteps(steps, 3, nil, Happiness, "Town completed")
			}
			return steps
		},
	}

	TendsRule = Rule{
		Name:        "Camping life",
		Description: "+5 happiness if there are at least 2 villagers living in a tend",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var count = 0
			for _, v := range n {
				if v.House.Type == 3 {
					count++
				}
			}
			if count > 1 {
				return addSteps(steps, 5, nil, Happiness, "Camping life")
			}
			return steps
		},
	}

	FancyHappinessRule = Rule{
		Name:        "Fancy house",
		Description: "+4 happiness for every villager living in a fancy house",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				if v.House.Type == 4 {
					steps = addSteps(steps, 4, &v.ID, Happiness, "Living in a fancy house")
				}
			}
			return steps
		},
	}

	SpaceLover = Rule{
		Name:        "Space house",
		Description: "+5 happiness for every Stuff lover living in a Big House",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok := m[Stuff]; ok && v1.Love() == Love {
					steps = addSteps(steps, 4, &v.ID, Happiness, "Stufflovers dream")
				}
			}
			return steps
		},
	}

	ThemeRule = Rule{
		Name:        "Theme town",
		Description: "+3 happiness if all houses are the same style",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var style *int
			for _, v := range n {
				if style == nil {
					style = &v.House.Type
					continue
				}
				if *style != v.House.Type {
					return steps
				}
			}
			return addSteps(steps, 3, nil, Happiness, "Theme town")
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
		Description: "-4 food for every villager that likes eating",
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

	HappyCooperation = Rule{
		Name:        "Happy cooperation",
		Description: "+10 happiness if there are more than 2 villager cooperative",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var cooperation int
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Competitive]; ok1 {
					switch v1.Love() {
					case Hate:
						cooperation++
					}
				}
			}
			if cooperation > 1 {
				steps = addSteps(steps, 10, nil, Happiness, "Happy cooperation")
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

	WorkingPower = Rule{
		Name:        "WorkingPower",
		Description: "+3 coins for every villager loving to work",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Workaholic]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 5, &v.ID, Money, "Generates extra money")
				}
			}
			return steps
		},
	}

	ExtrovertPower = Rule{
		Name:        "Extrovert Power",
		Description: "+3 Happiness for extrovert person",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Happiness]; ok1 && v1.Love() == Love {
					steps = addSteps(steps, 5, &v.ID, Happiness, "Extrovert Power")
				}
			}
			return steps
		},
	}

	IntrovertCulture = Rule{
		Name:        "Introvert Culture",
		Description: "+4 Culture for every introvert that likes Reading or Music",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				v1, ok1 := m[Extrovert]
				v2, ok2 := m[Reading]
				v3, ok3 := m[Music]
				if ok1 && v1.Love() == Hate && ((ok2 && v2.Love() == Love) || (ok3 && v3.Love() == Love)) {
					steps = addSteps(steps, 4, &v.ID, Cultural, "Introvert culture")

				}
			}
			return steps
		},
	}

	CultureLeak = Rule{
		Name:        "Culture Leak",
		Description: "-3 Culture for every person hating Reading or Music",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Cultural]; ok1 && v1.Love() != Love {
					steps = addSteps(steps, -3, &v.ID, Cultural, "Culture Leak")
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
			if isConflict(n, Animals) {
				steps = addSteps(steps, -20, nil, Happiness, "Animal conflict")
			}
			return steps
		},
	}

	BravenessConflict = Rule{
		Name:        "Braveness Conflict",
		Description: "-10 Happiness when there are adventorous and coward villagers",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			if isConflict(n, Adventurous) {
				steps = addSteps(steps, -10, nil, Happiness, "Braveness Conflict")
			}
			return steps
		},
	}

	HappyCapitalism = Rule{
		Name:        "Happy capitalism",
		Description: "+12 happiness for each villager loving stuff",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var result []bool
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Stuff]; ok1 && v1.Love() == Love {
					result = append(result, true)
				}
			}
			if len(result) > 1 {
				steps = addSteps(steps, 12, nil, Happiness, "Happy capitalism")
			}
			return steps
		},
	}

	UnethicalFood = Rule{
		Name:        "Unethical food",
		Description: "+10 food for every villager that hates animals",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Animals]; ok1 && v1.Love() == Hate {
					steps = addSteps(steps, 8, &v.ID, Food, "Unethical food")
				}
			}

			return steps
		},
	}

	PessimistFood = Rule{
		Name:        "Saving food",
		Description: "+6 food for every villager that is pessimist",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Optimistic]; ok1 && v1.Love() == Hate {
					steps = addSteps(steps, 8, &v.ID, Food, "Saving food")
				}
			}

			return steps
		},
	}

	OptimistTension = Rule{
		Name:        "Optimist vs Pessimist",
		Description: "-10 Health if there are optimist and pessimist together",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			if isConflict(n, Optimistic) {
				steps = addSteps(steps, -10, nil, Health, "Optimist Conflict")
			}
			return steps
		},
	}

	OptimistBoost = Rule{
		Name:        "Optimist Boost",
		Description: "+10 Happiness if there is a pessimist and at least 2 Optimists",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var optimists []bool
			var pessimist []bool
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Optimistic]; ok1 && v1.Love() == Love {
					optimists = append(optimists, true)
				}
				if v1, ok1 := m[Optimistic]; ok1 && v1.Love() == Hate {
					pessimist = append(pessimist, true)
				}
			}
			if len(optimists) > 1 && len(pessimist) == 1 {
				steps = addSteps(steps, 10, nil, Happiness, "Optimist Boost")
			}
			return steps
		},
	}

	CowardFancy = Rule{
		Name:        "Coward Fancy House",
		Description: "+8 Happiness if coward villager lives in a Fancy house",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Adventurous]; ok1 && v1.Love() == Hate && v.House.Type == 4 {
					steps = addSteps(steps, 8, &v.ID, Happiness, "Coward Fancy House")
				}
			}
			return steps
		},
	}

	AdventurousTend = Rule{
		Name:        "Adventurous tend",
		Description: "+8 Happiness if an adventurous villager lives in a tend",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Adventurous]; ok1 && v1.Love() == Love && v.House.Type == 3 {
					steps = addSteps(steps, 8, &v.ID, Happiness, "Adventurous Tend")
				}
			}
			return steps
		},
	}

	ExtraRent = Rule{
		Name:        "Extra Rent",
		Description: "+5 Rent for each villager living in a Big or Fancy House",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				if v.House.Type == 4 || v.House.Type == 2 {
					steps = addSteps(steps, 5, &v.ID, Money, "Extra rent")
				}
			}
			return steps
		},
	}

	ExtraSecurity = Rule{
		Name:        "House Protection",
		Description: "+3 Security for each villager leaving in Fancy House",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				if v.House.Type == 4 {
					steps = addSteps(steps, 5, &v.ID, Money, "House protection")
				}
			}
			return steps
		},
	}

	TendMisery = Rule{
		Name:        "Misery tend",
		Description: "-3 Happiness for each villager leaving in a Tend",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			for _, v := range n {
				if v.House.Type == 3 {
					steps = addSteps(steps, -3, &v.ID, Happiness, "Misery tend")
				}
			}
			return steps
		},
	}

	EatingTogether = Rule{
		Name:        "Eating together",
		Description: "+12 happiness if there is at least 1 villager loving cooking and another loving eating",
		Func: func(n NPCs, steps []StatStep) []StatStep {
			var eating []string
			var cooking []string
			for _, v := range n {
				m := v.Chars.charMap()
				if v1, ok1 := m[Cooking]; ok1 && v1.Love() == Love {
					cooking = append(cooking, v.ID)
				}
				if v1, ok1 := m[Eating]; ok1 && v1.Love() == Love {
					eating = append(eating, v.ID)
				}
			}
			if len(eating) == 0 || len(cooking) == 0 {
				return steps
			}
			for _, v := range eating {
				for _, j := range cooking {
					if v != j {
						steps = addSteps(steps, 12, nil, Happiness, "Eating together")
					}
				}
			}
			return steps
		},
	}
)

func isConflict(n NPCs, characteristic string) bool {
	var side bool
	var opposite bool
	for _, v := range n {
		m := v.Chars.charMap()
		if v1, ok1 := m[characteristic]; ok1 {
			switch v1.Love() {
			case Love:
				side = true
			case Hate:
				opposite = true
			}
		}
	}
	return side && opposite
}

var AllAvailableRules = []Rule{
	CookingBonus,
	AdventurousThief,
	CompetitionTooMuch,
	EatingTooMuch,
	GoodCulture,
	BadCulture,
	WorkTooMuch,
	WorkingPower,
	HealthyGuy,
	OptimisticThief,
	AnimalLovers,
	AnimalConflict,
	HappyCooperation,
	UnethicalFood,
	HappyCapitalism,
	BravenessConflict,
	CultureLeak,
	IntrovertCulture,
	ExtrovertPower,
	TendsRule,
	ThemeRule,
	FancyHappinessRule,
	OptimistTension,
	PessimistFood,
	OptimistBoost,
	CowardFancy,
	AdventurousTend,
	ExtraRent,
	ExtraSecurity,
	TendMisery,
	SpaceLover,
	EatingTogether,
}

func RandomRule() Rule {
	index := rand.Intn(len(AllAvailableRules))
	rule := AllAvailableRules[index]
	AllAvailableRules = append(AllAvailableRules[0:index], AllAvailableRules[index+1:]...)
	return rule
}
