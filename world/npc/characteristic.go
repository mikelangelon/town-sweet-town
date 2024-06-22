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
		case Adventurous:
			if v.Level > 3 {
				phrases = append(phrases, "I'm adventurous.")
				continue
			}
			phrases = append(phrases, "I'm a coward.")
		case Rent:
			phrases = append(phrases, fmt.Sprintf("I can pay %d coins as rent", v.Level))
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

func (c Chars) WithRent(rent int64) Chars {
	return append(c, Characteristic{
		Name:  Rent,
		Level: rent,
	})
}
func WithCharacteristic(values ...string) Chars {
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

const (
	Rent        = "Rent"
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
