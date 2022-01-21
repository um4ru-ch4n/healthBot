package service

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/domain"
)

type Service struct {
	chatInfo map[int64]*domain.ChatInfo
}

func NewService(cfg *config.Config) *Service {
	newService := &Service{
		chatInfo: make(map[int64]*domain.ChatInfo),
	}

	newService.chatInfo[-727028014] = &domain.ChatInfo{
		IsWorking: false,
		Done:      make(chan struct{}),
		PollInfo:  &domain.PollInfo{},
		HeadPerson: &domain.HeadPerson{
			Username: "oooMRXooo",
			ChatID:   371947069,
		},
	}

	return newService
}

func (srv *Service) HelpGroup(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "About group bot...")
	bot.Send(newMsg)
}

func (srv *Service) HelpPrivate(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "About private bot...")
	bot.Send(newMsg)
}

func (srv *Service) UnknownCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, but I don't know this command( Please, type /help command")
	bot.Send(newMsg)
}

func (srv *Service) Start(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Welcome!")
	bot.Send(newMsg)
}

func (srv *Service) StartRoutine(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if srv.chatInfo[msg.Chat.ID].IsWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already started")
		bot.Send(newMsg)
		return
	}

	times, err := parseTimeSliceFromString(msg.CommandArguments())
	if err != nil {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ParseMode = "MarkdownV2"
		bot.Send(newMsg)
		return
	}

	srv.chatInfo[msg.Chat.ID].PollInfo.Times = times

	srv.chatInfo[msg.Chat.ID].Done = make(chan struct{}, 1)
	srv.chatInfo[msg.Chat.ID].IsWorking = true

	go srv.createPolls(bot, msg.Chat.ID)

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot started...")
	bot.Send(newMsg)
}

func parseTimeSliceFromString(argTime string) ([]time.Time, error) {
	times := strings.Split(argTime, " ")

	if len(times) < 2 {
		return nil, fmt.Errorf("you must enter 2 times\\! `/start_routine [poll creation time] [first reminder time] [second reminder time] ...`")
	}

	parsedTimes := make([]time.Time, len(times))

	for i, t := range times {
		parsedTime, err := time.Parse("15:04:05", t)
		if err != nil {
			return nil, fmt.Errorf("error in %d time format", i+1)
		}
		parsedTimes = append(parsedTimes, parsedTime)
	}

	return parsedTimes, nil
}

func (srv *Service) createPolls(bot *tgbotapi.BotAPI, chatID int64) {
	for {
		select {
		case <-srv.chatInfo[chatID].Done:
			return
		default:
		}

		fmt.Println("qwer")

		time.Sleep(10 * time.Second)
	}
}

func (srv *Service) StopRoutine(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if !srv.chatInfo[msg.Chat.ID].IsWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already stopped")
		bot.Send(newMsg)
		return
	}
	srv.chatInfo[msg.Chat.ID].Done <- struct{}{}
	srv.chatInfo[msg.Chat.ID].IsWorking = false

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot stopped...")
	bot.Send(newMsg)
}

func (srv *Service) AddNewChat(bot *tgbotapi.BotAPI, chatID int64) {
	srv.chatInfo[chatID] = &domain.ChatInfo{
		IsWorking: false,
		Done:      make(chan struct{}),
		PollInfo:  &domain.PollInfo{},
	}

	newMsg := tgbotapi.NewMessage(chatID, "Hello everyone!")
	bot.Send(newMsg)
	fmt.Println("Added new chat")
}

func (srv *Service) RemoveChat(bot *tgbotapi.BotAPI, chatID int64) {
	delete(srv.chatInfo, chatID)

	newMsg := tgbotapi.NewMessage(chatID, "Goodbuy everyone!")
	bot.Send(newMsg)
	fmt.Println("Removed chat")
}
