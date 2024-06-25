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
	return fmt.Sprintf("%s wishes the town should have more %s", c.NPC.ID, c.Wish.Stat)
}

func (c CompleteWishes) String() []string {
	var result []string
	for _, v := range c {
		result = append(result, v.String())
	}
	return result
}
