package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	CONFIGFILE  = "conf.json"
	SheetsCreds = "sheets_credentials.json"
	Config      Configuration
	Manager     ItemsManager
	unknownText = "Извини, я тебя не понял. Попробуй ещё разок или пожалуйся @evgsol."
)

func initManager() {
	if err := Manager.init(); err != nil {
		log.Fatal(err)
	}
}

func process(id int, input string) (string, error) {
	s := strings.Split(input, " ")
	command := s[0]
	args := s[1:]

	switch command {
	case "/start":
		return "Привет! Я дурнирный бот.", nil
	case "/list":
		initManager()
		if len(args) != 0 {
			return unknownText, nil
		}
		return Manager.list(id), nil
	case "/give":
		initManager()
		if len(args) != 3 || args[1] != "to" {
			return unknownText, nil
		}
		return fmt.Sprintf("(args: %v) Not implemented yet.", args), nil
	default:
		return unknownText, nil
	}
}

func main() {
	if err := Config.init(CONFIGFILE); err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = Config.Debug

	log.Printf("Authorized on %s", bot.Self.UserName)

	//_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("WEBHOOK")))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//info, err := bot.GetWebhookInfo()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if info.LastErrorDate != 0 {
	//	log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	//}
	//
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//go http.ListenAndServeTLS("0.0.0.0:8443")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Only text messages are supported."))
			continue
		}

		log.Printf("[%s (%v)] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		response, err := process(update.Message.From.ID, update.Message.Text)
		if err != nil {
			log.Printf("Failed to generate response: %v", err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
	}
}
