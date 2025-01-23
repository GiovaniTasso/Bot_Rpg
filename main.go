package main

import (
	"BotRpg/bot"
	"BotRpg/spells"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func main() {
	if err := spells.LoadSpells("spell.json"); err != nil {
		log.Fatalf("Erro ao carregar magias: %v", err)
	}

	err := spells.LoadManeuvers("manobras.json")
	if err != nil {
		log.Fatalf("Erro ao carregar manobras: %v", err)
	}

	botAPI, err := bot.InitializeBot()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot iniciado: %s", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			bot.HandleUserInput(chatID, update.Message.Text, botAPI)
		} else if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			if strings.HasPrefix(callbackData, "classe_") {
				class := strings.TrimPrefix(callbackData, "classe_")
				bot.HandleClassSelection(class, update.CallbackQuery.Message.Chat.ID, botAPI)
			} else {
				bot.HandleButtonPress(update.CallbackQuery, botAPI)
			}
		}
	}
}
