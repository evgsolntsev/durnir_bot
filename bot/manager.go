package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/player"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Manager struct {
	PlayerManager  player.Manager
	FighterManager fighter.Manager
	BotAPI         *tgbotapi.BotAPI
}

func NewManager(
	playerManager player.Manager,
	fighterManager fighter.Manager,
	botAPI *tgbotapi.BotAPI,
) *Manager {
	return &Manager{
		PlayerManager:  playerManager,
		FighterManager: fighterManager,
		BotAPI:         botAPI,
	}
}

func (m *Manager) ProcessMessage(ctx context.Context, u tgbotapi.Update) error {
	if u.Message == nil { // ignore any non-Message updates
		m.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Only text messages are supported."))
		return nil
	}

	player, err := m.PlayerManager.GetOneByTelegramId(ctx, u.Message.Chat.ID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	s := strings.Split(u.Message.Text, " ")
	command := s[0]
	args := s[1:]

	var response string
	if player != nil {
		response, err = m.processPlayerMessage(ctx, player, command, args)
	} else {
		response, err = m.processStrangerMessage(ctx, command, args)
	}
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, response)
	if _, err := m.BotAPI.Send(msg); err != nil {
		return err
	}
	return nil
}

func (m *Manager) processStrangerMessage(ctx context.Context, command string, args []string) (string, error) {
	var response string
	switch command {
	case "/start":
		response = "Привет! Я дурнирный бот."
	case "/me":
		response = "Я тебя не знаю!"
	default:
		response = "Извини, я тебя не понял. Попробуй ещё разок или пожалуйся @evgsol."
	}
	return response, nil
}

func (m *Manager) processPlayerMessage(
	ctx context.Context, player *player.Player, command string, args []string,
) (string, error) {
	var (
		response string
		fighter  *fighter.Fighter
		err      error
	)

	if player.FighterID != nil {
		fighter, err = m.FighterManager.GetOne(ctx, *player.FighterID)
	}
	if err != nil {
		log.Printf("Error getting fighter with ID `%s`: %s", *player.FighterID, err.Error())
	}

	switch command {
	case "/start":
		response = "Ты чего, мы же уже разговариваем."
	case "/me":
		response = description(player, fighter)
	case "/generate":
	default:
		response = "Извини, я тебя не понял. Попробуй ещё разок или пожалуйся @evgsol."
	}
	return response, nil
}

func description(player *player.Player, fighter *fighter.Fighter) string {
	fighterString := ""
	if fighter != nil {
		fighterString = fmt.Sprintf(
			"У тебя есть мoнстр %s.\nЖизней: %d\nМаны: %d",
			fighter.Name, fighter.Health, fighter.Mana,
		)
	}
	return fmt.Sprintf(
		"Ты %s.\nУ тебя %d золота.\n%s",
		player.Name, player.Gold, fighterString)

}
