package world

import "fmt"

type Goal struct {
	Day       int
	Stat      string
	Value     int
	GiftStat  string
	GiftValue int
	Mandatory bool
}

func (g Goal) String() string {
	mandatory := "mandatory"
	if !g.Mandatory {
		mandatory = "optional"
	}

	return fmt.Sprintf("Day %d --> Get to %s %d [%s]", g.Day, g.Stat, g.Value, mandatory)
}

func (g Goal) Price() string {
	price := "HAPPY END"
	if g.GiftStat != "" {
		price = fmt.Sprintf("Price %s +%d", g.GiftStat, g.GiftValue)
	}
	return price
}
