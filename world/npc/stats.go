package npc

import (
	"fmt"
	"strings"
)

const (
	Happiness = "happiness"
	Cultural  = "cultural"
	Health    = "health"
	Security  = "security"
	Food      = "food"
	Money     = "money"
)

type Stats struct {
	Food      int
	Money     int
	Happiness int
	Security  int
	Cultural  int
	Health    int
}

func (s1 Stats) Merge(s2 Stats) Stats {
	return Stats{
		Happiness: s1.Happiness + s2.Happiness,
		Security:  s1.Security + s2.Security,
		Food:      s1.Food + s2.Food,
		Cultural:  s1.Cultural + s2.Cultural,
		Health:    s1.Health + s2.Health,
	}
}

type Stat struct {
	name  string
	steps []StatStep
}

type StatStep struct {
	Name   string
	CharID *string
	Value  int
	Text   string
}

func (s StatStep) FormatText() string {
	return fmt.Sprintf("%s%s--> %s: %d", strings.ToUpper(s.Name), s.FormatCharID(), s.Text, s.Value)
}

func (s StatStep) FormatCharID() string {
	if s.CharID == nil {
		return ""
	}
	return fmt.Sprintf(" %s ", *s.CharID)
}
