package npc

type Stats struct {
	Happiness int
	Security  int
	Food      int
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
	Name  string
	Value int
	Text  string
}
