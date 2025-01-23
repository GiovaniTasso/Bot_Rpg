package spells

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manobra struct {
	Name      string `json:"nome"`
	Descricao string `json:"descricao"`
}

var ManobraLista []Manobra

func LoadManeuvers(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ManobraLista)
	if err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}
	return nil
}
