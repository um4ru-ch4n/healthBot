package main

import (
	"flag"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/handler"
	"github.com/um4aru-ch4n/healthBot/repository"
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

	psqlPool, err := repository.NewPostgreConnPool(repository.PostgreConfig{
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		Username: config.Postgres.Username,
		DBName:   config.Postgres.DBName,
		SSLMode:  config.Postgres.SSLMode,
		Password: config.Postgres.Password,
	})
	if err != nil {
		fmt.Println("postgre connection failed: ", err)
	}

	repo := repository.NewRepository(psqlPool)

	service := service.NewService(config, repo)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	routerHandler := handler.NewRouter(service, bot, config)

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
