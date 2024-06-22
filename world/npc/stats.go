package npc

// TODO This should not be in npc package
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
