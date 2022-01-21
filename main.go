package main

import (
	"flag"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/handler"
	"github.com/um4aru-ch4n/healthBot/service"
)

func main() {
	var pathToConfig string

	flag.StringVar(&pathToConfig, "config", "./config.yml", "Specify a path to config file")
	flag.Parse()

	config, err := config.NewConfig(pathToConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(u)

	service := service.NewService(config)

	routerHandler := handler.NewRouter(service, bot, config)

	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
