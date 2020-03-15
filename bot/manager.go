package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/evgsolntsev/durnir_bot/player"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Manager struct {
	PlayerManager player.Manager
	BotAPI        *tgbotapi.BotAPI
}

func NewManager(
	playerManager player.Manager,
	botAPI *tgbotapi.BotAPI,
) *Manager {
	return &Manager{
		PlayerManager: playerManager,
		BotAPI:        botAPI,
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
	var response string
	switch command {
	case "/start":
		response = "Ты чего, мы же уже разговариваем."
	case "/me":
		response = fmt.Sprintf(
			"Ты **%s**.\nУ тебя **%d** золота.",
			player.Name, player.Gold,
		)
	default:
		response = "Извини, я тебя не понял. Попробуй ещё разок или пожалуйся @evgsol."
	}
	return response, nil
}
