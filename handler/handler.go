package handler

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/service"
)

type ChatMemberStatus string

const (
	StatusLeft ChatMemberStatus = "left"
	StatusJoin ChatMemberStatus = "member"
)

type ChatType string

const (
	ChatPrivate ChatType = "private"
	ChatGroup   ChatType = "group"
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
	case update.Poll != nil:
		r.HandlePoll(update.Poll, update.FromChat().ID)
	case update.MyChatMember != nil:
		status := update.MyChatMember.NewChatMember.Status

		if status == string(StatusJoin) {
			r.HandleUpdateMember(StatusJoin, update.MyChatMember.Chat.ID)
			return
		}

		r.HandleUpdateMember(StatusLeft, update.MyChatMember.Chat.ID)

	default:
		spew.Dump(update)
	}
}

func (r *Router) HandleMessage(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	if msg.Chat.Type == string(ChatGroup) {
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

	if msg.Chat.Type == string(ChatPrivate) {
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

func (r *Router) HandlePoll(msg *tgbotapi.Poll, chatID int64) {

}

func (r *Router) HandleUpdateMember(status ChatMemberStatus, chatID int64) {
	if status == StatusJoin {
		r.service.AddNewChat(r.bot, chatID)
		return
	}

	r.service.RemoveChat(r.bot, chatID)
}
