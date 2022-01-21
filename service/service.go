package service

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
)

type Service struct {
	chatInfo map[int64]*ChatInfo
}

type ChatInfo struct {
	isWorking bool
	done      chan struct{}
	pollID    int64
}

func NewService(cfg *config.Config) *Service {
	newService := &Service{
		chatInfo: make(map[int64]*ChatInfo),
	}

	newService.chatInfo[-727028014] = &ChatInfo{
		isWorking: false,
		done:      make(chan struct{}),
		pollID:    0,
	}

	return newService
}

func (srv *Service) Help(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "About bot...")
	bot.Send(newMsg)
}

func (srv *Service) Start(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if srv.chatInfo[msg.Chat.ID].isWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already started")
		bot.Send(newMsg)
		return
	}
	srv.chatInfo[msg.Chat.ID].done = make(chan struct{}, 1)
	srv.chatInfo[msg.Chat.ID].isWorking = true

	go srv.createPolls(bot, msg.Chat.ID)

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot started...")
	bot.Send(newMsg)
}

func (srv *Service) createPolls(bot *tgbotapi.BotAPI, chatID int64) {
	for {
		select {
		case <-srv.chatInfo[chatID].done:
			return
		default:
		}

		fmt.Println("qwer")

		time.Sleep(3 * time.Second)
	}
}

func (srv *Service) Stop(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if !srv.chatInfo[msg.Chat.ID].isWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already stopped")
		bot.Send(newMsg)
		return
	}
	srv.chatInfo[msg.Chat.ID].done <- struct{}{}
	srv.chatInfo[msg.Chat.ID].isWorking = false

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot stopped...")
	bot.Send(newMsg)
}

func (srv *Service) AddNewChat(bot *tgbotapi.BotAPI, chatID int64) {
	srv.chatInfo[chatID] = &ChatInfo{
		isWorking: false,
		done:      make(chan struct{}),
		pollID:    0,
	}

	newMsg := tgbotapi.NewMessage(chatID, "Hello everyone!")
	bot.Send(newMsg)
	fmt.Println("Added new chat")
	spew.Dump(srv.chatInfo)
}

func (srv *Service) RemoveChat(bot *tgbotapi.BotAPI, chatID int64) {
	delete(srv.chatInfo, chatID)

	newMsg := tgbotapi.NewMessage(chatID, "Goodbuy everyone!")
	bot.Send(newMsg)
	fmt.Println("Removed chat")
	spew.Dump(srv.chatInfo)
}
