package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/domain"
	"github.com/um4aru-ch4n/healthBot/service"
)

type ChatMemberStatus string

const (
	StatusLeft ChatMemberStatus = "left"
	StatusJoin ChatMemberStatus = "member"
)

type Router struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
	cfg     *config.Config
}

func NewRouter(srv *service.Service, bot *tgbotapi.BotAPI, cfg *config.Config) *Router {
	return &Router{
		bot:     bot,
		service: srv,
		cfg:     cfg,
	}
}

func (r *Router) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			fmt.Printf("recovered from panic: %v", panicValue)
		}
	}()

	switch {
	case update.Message != nil:
		r.HandleMessage(update.Message)
	case update.PollAnswer != nil:
		r.HandlePoll(update.PollAnswer)
	case update.MyChatMember != nil:
		status := update.MyChatMember.NewChatMember.Status

		if status == string(StatusJoin) {
			r.HandleUpdateMember(StatusJoin, update.MyChatMember.Chat.ID)
			return
		}

		r.HandleUpdateMember(StatusLeft, update.MyChatMember.Chat.ID)
	case update.CallbackQuery != nil:
		r.HandleCallBackQuery(update.CallbackQuery)
	default:
		// spew.Dump(update)
	}
}

func (r *Router) HandleMessage(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	if msg.Chat.IsGroup() {
		switch msg.Command() {
		case "help":
			r.service.HelpGroup(r.bot, msg)
		case "start_routine":
			r.service.StartRoutine(r.bot, msg)
		case "stop_routine":
			r.service.StopRoutine(r.bot, msg)
		default:
			r.service.UnknownCommand(r.bot, msg)
		}
		return
	}

	if msg.Chat.IsPrivate() {
		switch msg.Command() {
		case "help":
			r.service.HelpPrivate(r.bot, msg)
		case "start":
			r.service.Start(r.bot, msg)

		default:
			r.service.UnknownCommand(r.bot, msg)
		}
	}
}

func (r *Router) HandlePoll(poll *tgbotapi.PollAnswer) {
	r.service.UpdatePollResults(r.bot, poll.PollID, poll.User.ID, poll.OptionIDs)
}

func (r *Router) HandleUpdateMember(status ChatMemberStatus, chatID int64) {
	if status == StatusJoin {
		r.service.AddNewChat(r.bot, chatID)
		return
	}

	r.service.RemoveChat(r.bot, chatID)
}

func (r *Router) HandleCallBackQuery(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {
	case service.RegisterNewUser:
		r.service.RegisterNewUser(
			r.bot,
			callback.Message.Chat.ID,
			&domain.User{
				ID:        callback.From.ID,
				Username:  callback.From.UserName,
				Firstname: callback.From.FirstName,
				Lastname:  callback.From.LastName,
				ChatID:    0,
			},
		)
	default:
		fmt.Println("Unknown callback data: ", callback.Data)
	}
}
