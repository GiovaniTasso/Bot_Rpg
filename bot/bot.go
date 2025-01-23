package bot

import (
	"BotRpg/spells"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

var AwaitingSpellName map[int64]bool

func InitializeBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI("7820372433:AAF5b_TlxZx1TtWdRv7QqlBrIFUxI0sz7i0")
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func ShowMainMenu(chatID int64, bot *tgbotapi.BotAPI) {
	menuText := "🔮Bem-vindo ao bot de magias de D&D!🔮\n\nEscolha uma opção abaixo:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔍 Buscar Magia", "buscar_magia")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📜 Listar Magias por Classe", "listar_classes")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⚔️ Manobras", "listar_manobras")),
	)
	msg := tgbotapi.NewMessage(chatID, menuText)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleButtonPress(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	switch data {
	case "buscar_magia":
		AwaitingSpellName[chatID] = true
		bot.Send(tgbotapi.NewMessage(chatID, "Digite o nome da magia que deseja buscar:"))
	case "listar_classes":
		SendClassSelection(chatID, bot)
	case "listar_manobras":
		ShowManobrasList(chatID, bot)
	default:
		if strings.HasPrefix(data, "spell_") {
			spellName := strings.TrimPrefix(data, "spell_")
			SendSpellDetails(spellName, chatID, bot)
		}
		if strings.HasPrefix(data, "manobra_") {
			maneuverName := strings.TrimPrefix(data, "manobra_")
			ShowManeuverDetails(chatID, maneuverName, bot)
		}
	}
	bot.Request(tgbotapi.NewCallback(callbackQuery.ID, "Processando..."))
}

func SendClassSelection(chatID int64, bot *tgbotapi.BotAPI) {
	classes := []string{"Mago", "Bruxo", "Clérigo", "Bardo", "Druida", "Feiticeiro", "Paladino", "Ranger"}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, class := range classes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(class, "classe_"+strings.ToLower(class)),
		))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, "Escolha uma classe para listar magias:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func SendSpellDetails(spellName string, chatID int64, bot *tgbotapi.BotAPI) {
	matchingSpells := spells.SearchSpellsByName(spellName)

	if len(matchingSpells) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Nenhuma magia encontrada com esse nome."))
		return
	}

	for _, spell := range matchingSpells {
		response := fmt.Sprintf("✨ **%s**\n📜 Nível: %d\n🏫 Escola: %s\n 🧙‍♂️ Conjuradores: %s\n⏳ Tempo: %s\n🎯 Alcance: %s\n🔄 Duração: %s\n📖 Descrição:\n%s\n",
			spell.Name, spell.Level, spell.School, strings.Join(spell.Classes, ", "), spell.CastingTime, spell.Range, spell.Duration, strings.Join(spell.Description, "\n"))
		bot.Send(tgbotapi.NewMessage(chatID, response))
	}
}

func HandleClassSelection(class string, chatID int64, bot *tgbotapi.BotAPI) {
	matchingSpells := spells.ListSpellsByClass(class)

	if len(matchingSpells) > 0 {
		response := fmt.Sprintf("📜 **Magias da Classe %s**:\n", strings.Title(class))
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, spell := range matchingSpells {
			schoolEmoji := getMagicSchoolEmoji(spell.School)
			buttonText := fmt.Sprintf("%s (Nível %d) %s", spell.Name, spell.Level, schoolEmoji)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, "spell_"+spell.Name),
			))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, "Nenhuma magia encontrada para esta classe."))
	}
}

func getMagicSchoolEmoji(school string) string {
	school = strings.ToLower(school)
	switch school {
	case "abjuração":
		return "🛡️"
	case "adivinhação":
		return "🔮"
	case "conjuração":
		return "🌀"
	case "encantamento":
		return "✨"
	case "evocação":
		return "🔥"
	case "ilusão":
		return "🎭"
	case "necromancia":
		return "💀"
	case "transmutação":
		return "🔄"
	default:
		return "📖"
	}
}

func HandleUserInput(chatID int64, text string, botAPI *tgbotapi.BotAPI) {
	if AwaitingSpellName[chatID] {
		AwaitingSpellName[chatID] = false

		matchingSpells := spells.SearchSpellsByName(text)

		if len(matchingSpells) > 0 {
			response := "🔍 **Magias encontradas:**\n"

			var rows [][]tgbotapi.InlineKeyboardButton
			for _, spell := range matchingSpells {
				schoolEmoji := getMagicSchoolEmoji(spell.School)

				buttonText := fmt.Sprintf("%s (Nível %d) %s", spell.Name, spell.Level, schoolEmoji)
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(buttonText, "spell_"+spell.Name),
				))
			}

			keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
			msg := tgbotapi.NewMessage(chatID, response)
			msg.ReplyMarkup = keyboard
			botAPI.Send(msg)
		} else {
			botAPI.Send(tgbotapi.NewMessage(chatID, "Nenhuma magia encontrada com esse nome."))
		}
	} else {
		ShowMainMenu(chatID, botAPI)
	}
}

func ShowManobrasList(chatID int64, bot *tgbotapi.BotAPI) {
	response := "⚔️ **Lista de Manobras:**\n"

	if len(spells.ManobraLista) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Não há manobras disponíveis no momento."))
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, manobras := range spells.ManobraLista {
		buttonText := fmt.Sprintf(manobras.Name)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, "manobra_"+manobras.Name),
		))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, response)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func ShowManeuverDetails(chatID int64, maneuverName string, bot *tgbotapi.BotAPI) {
	// Procurar a manobra pelo nome
	var selectedManeuver *spells.Manobra
	for _, maneuver := range spells.ManobraLista {
		if maneuver.Name == maneuverName {
			selectedManeuver = &maneuver
			break
		}
	}

	if selectedManeuver == nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Manobra não encontrada."))
		return
	}

	// Criar a resposta com os detalhes da manobra
	response := fmt.Sprintf(
		"⚔️ **Manobra: %s**\n\n"+
			"📜 **Descrição**:\n%s",
		selectedManeuver.Name,
		selectedManeuver.Descricao,
	)

	// Enviar a mensagem para o chat
	msg := tgbotapi.NewMessage(chatID, response)
	bot.Send(msg)
}
