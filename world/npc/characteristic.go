package npc

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	Love = "love"
	Hate = "hate"
	Meh  = "don't care about"
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
			phrases = append(phrases, fmt.Sprintf("I %s %s", v.Love(), strings.ToLower(v.Name)))
		}
	}
	return phrases
}

func (c Characteristic) Love() string {
	if c.Level > 3 {
		return Love
	} else if c.Level < 3 {
		return Hate
	}
	return Meh
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

func (c Chars) charMap() map[string]Characteristic {
	var m = make(map[string]Characteristic, len(c))
	for _, k := range c {
		m[k.Name] = k
	}
	return m
}

func WithRandom(amount int) Chars {
	var result Chars
	options := []string{
		Extrovert, Competitive, Optimistic, Adventurous, Sports, Workaholic, Reading, Cooking, Music, Animals,
	}
	for i := 0; i < amount; i++ {
		index := rand.Intn(len(options))
		result = append(result, Characteristic{
			Name:  options[index],
			Level: int64(rand.Intn(10)),
		})
		options = append(options[0:index], options[index+1:]...)
	}
	return result.WithRent(int64(rand.Intn(3 + 3)))
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
