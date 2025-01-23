package spells

import (
	"encoding/json"
	"os"
	"strings"
)

type Spell struct {
	Name        string   `json:"nome"`
	Level       int      `json:"nivel"`
	School      string   `json:"escola"`
	Description []string `json:"descricao"`
	Classes     []string `json:"conjuradores"`
	CastingTime string   `json:"tempo_conjuracao"`
	Range       string   `json:"alcance"`
	Components  []string `json:"componentes"`
	Duration    string   `json:"duracao"`
	Material    string   `json:"descricao_componentes"`
}

var Spells []Spell

func LoadSpells(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&Spells)
}

func SearchSpellsByName(name string) []Spell {
	var result []Spell
	for _, spell := range Spells {
		if strings.Contains(strings.ToLower(spell.Name), strings.ToLower(name)) {
			result = append(result, spell)
		}
	}
	return result
}

func ListSpellsByClass(class string) []Spell {
	var result []Spell
	for _, spell := range Spells {
		for _, c := range spell.Classes {
			if strings.EqualFold(c, class) {
				result = append(result, spell)
			}
		}
	}
	return result
}
