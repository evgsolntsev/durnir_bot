package main

import (
	"context"
	"log"
	"time"

	"github.com/evgsolntsev/durnir_bot/bot"
	"github.com/evgsolntsev/durnir_bot/executables"
	"github.com/evgsolntsev/durnir_bot/fight"
	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/evgsolntsev/durnir_bot/player"
	"github.com/globalsign/mgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	CONFIGFILE = "conf.json"
	Config     executables.Configuration
)

func main() {
	if err := Config.Init(CONFIGFILE); err != nil {
		log.Fatal(err)
	}

	tgbot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		log.Fatal(err)
	}

	tgbot.Debug = Config.Debug

	ctx := context.Background()
	session, err := mgo.Dial(Config.MongoURL)
	if err != nil {
		panic(err)
	}
	fighterDAO := fighter.NewDAO(ctx, session)
	fighterManager := fighter.NewManager(ctx, fighterDAO)
	playerDAO := player.NewDAO(ctx, session)
	playerManager := player.NewManager(ctx, playerDAO, fighterManager)
	botManager := bot.NewManager(playerManager, fighterManager, tgbot)
	fightDAO := fight.NewDAO(ctx, session)
	fightManager := fight.NewManager(ctx, playerManager, fighterManager, fightDAO, botManager)

	for {
		err := fightManager.Step(ctx, idtype.StartHex)
		result := "OK"
		if err != nil {
			result = err.Error()
		}

		log.Printf(result)
		time.Sleep(time.Second)
	}
}
