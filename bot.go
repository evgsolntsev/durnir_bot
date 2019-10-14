package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	//bot.Debug = true

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
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot is still under development.")
		bot.Send(msg)
	}
}
