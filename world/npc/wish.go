package npc

import "fmt"

type Wish struct {
	DayStart  int
	DayEnd    int
	Stat      string
	Value     int
	Happiness int
}

type CompleteWish struct {
	Wish Wish
	NPC  *NPC
}

type CompleteWishes []CompleteWish

func (c CompleteWish) String() string {
	return fmt.Sprintf("%s wishes the village would have more %s", c.NPC.ID, c.Wish.Stat)
}

func (w Wish) IamText() string {
	return fmt.Sprintf("I wish the village would have more %s", w.Stat)
}

func (c CompleteWishes) String() []string {
	var result []string
	for _, v := range c {
		result = append(result, v.String())
	}
	return result
}
