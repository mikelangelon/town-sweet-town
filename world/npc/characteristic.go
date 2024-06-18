package npc

import (
	"fmt"
	"strings"
)

type Characteristic struct {
	Name  string
	Level int64
}

type Chars []Characteristic

type CharMap map[string]int64

func (c Chars) AsPhrases() []string {
	var phrases []string
	for _, v := range c {
		switch v.Name {
		case Extrovert:
			if v.Level > 3 {
				phrases = append(phrases, "I'm an extrovert.")
				continue
			}
			phrases = append(phrases, "I'm an introvert.")
		case Competitive:
			if v.Level > 3 {
				phrases = append(phrases, "I'm a competitive person.")
				continue
			}
			phrases = append(phrases, "I'm a cooperative.")
		case Optimistic:
			if v.Level > 3 {
				phrases = append(phrases, "I'm optimist.")
				continue
			}
			phrases = append(phrases, "I'm pessimist.")
		default:
			var phrase = "I"
			if v.Level > 3 {
				phrase += " love "
			} else {
				phrase += " hate "
			}
			phrases = append(phrases, fmt.Sprintf("%s%s", phrase, strings.ToLower(v.Name)))
		}
	}
	return phrases
}

func AddHappyCharacteristics(values ...string) Chars {
	var chs []Characteristic
	for _, v := range values {
		chs = append(chs, Characteristic{
			Name:  v,
			Level: 7,
		})
	}
	return chs
}

func (c Chars) charLevelMap() map[string]int {
	var m = make(map[string]int, len(c))
	for _, k := range c {
		m[k.Name] = int(k.Level)
	}
	return m
}

func (c Chars) Stats() Stats {
	return Stats{
		Security: c.Security(),
		Food:     c.Food(),
		Cultural: c.Cultural(),
		Health:   c.Healthy(),
	}
}

func (c Chars) Security() int {
	var total int
	m := c.charLevelMap()
	if v1, ok1 := m[Stuff]; ok1 {
		if v2, ok1 := m[Optimistic]; ok1 {
			return v1 * v2
		}
		if v2, ok1 := m[Adventurous]; ok1 {
			return v1 * v2
		}
	}
	return total
}

func (c Chars) Food() int {
	var total int
	m := c.charLevelMap()
	if v1, ok1 := m[Cooking]; ok1 {
		total += v1
	}
	if v1, ok1 := m[Eating]; ok1 {
		total -= v1
	}
	return total
}

func (c Chars) Cultural() int {
	var total int
	m := c.charLevelMap()
	if v1, ok1 := m[Reading]; ok1 {
		total += v1
	}
	if v1, ok1 := m[Music]; ok1 {
		total += v1
	}
	return total
}

func (c Chars) Healthy() int {
	var total int
	m := c.charLevelMap()
	if v1, ok1 := m[Workaholic]; ok1 && v1 > 7 {
		total -= v1
	}
	if v1, ok1 := m[Eating]; ok1 && v1 > 7 {
		total -= v1
	}
	if v1, ok1 := m[Sports]; ok1 {
		total += v1
	}
	return total
}

func CheckHappiness(c []Chars) int {
	var total int
	var animals []int
	var competitive []int
	for _, v := range c {
		m := v.charLevelMap()
		if v1, ok1 := m[Animals]; ok1 {
			animals = append(animals, v1)
		}
	}
	for _, v := range c {
		m := v.charLevelMap()
		if v1, ok1 := m[Competitive]; ok1 {
			competitive = append(competitive, v1)
		}
	}
	if len(animals) > 1 {
		total += 14
	}
	if len(competitive) > 1 {
		total -= 10
	}
	return total
}

const (
	Extrovert   = "Extrovert"
	Competitive = "Competitive"
	Optimistic  = "Optimistic"
	Adventurous = "Adventurous"

	Workaholic = "to Work"
	Sports     = "Sports"
	Reading    = "Reading"
	Cooking    = "Cooking"
	Music      = "Music"
	Animals    = "Animals"
	Stuff      = "Stuff"
	Eating     = "Food"

	Quite = "Quite"
	Big   = "Big"
	Humid = "Humid"
	Fancy = "Fancy"
)
