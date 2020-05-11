package main

import (
	"context"
	"log"

	"github.com/evgsolntsev/durnir_bot/bot"
	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/player"
	"github.com/evgsolntsev/durnir_bot/executables"
	"github.com/globalsign/mgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	CONFIGFILE  = "conf.json"
	SheetsCreds = "sheets_credentials.json"
	Config      executables.Configuration
	//Manager     GoogleSheetManager
)

//func initManager() {
//	if err := Manager.init(); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	if err := Config.Init(CONFIGFILE); err != nil {
		log.Fatal(err)
	}

	tgbot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		log.Fatal(err)
	}

	tgbot.Debug = Config.Debug
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	session, err := mgo.Dial(Config.MongoURL)
	if err != nil {
		panic(err)
	}
	playerDAO := player.NewDAO(ctx, session)
	playerManager := player.NewManager(ctx, playerDAO)
	fighterDAO := fighter.NewDAO(ctx, session)
	fighterManager := fighter.NewManager(ctx, fighterDAO)
	botManager := bot.NewManager(playerManager, fighterManager, tgbot)

	updates, err := tgbot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s (%v)] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		err := botManager.ProcessMessage(ctx, update)
		if err != nil {
			log.Printf("Failed to generate response: %v", err)
		}
	}
}
