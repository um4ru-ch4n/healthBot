package service

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/um4aru-ch4n/healthBot/config"
	"github.com/um4aru-ch4n/healthBot/domain"
)

const CheckTimeSleep = 10 * time.Second
const RegisterNewUser = "register_new_user"

type Service struct {
	chatInfo map[int64]*domain.ChatInfo
	pollChat map[string]int64
	repo     repository
}

type repository interface {
}

func NewService(cfg *config.Config, repo repository) *Service {
	newService := &Service{
		chatInfo: make(map[int64]*domain.ChatInfo),
		pollChat: make(map[string]int64),
		repo:     repo,
	}

	newService.chatInfo[-727028014] = &domain.ChatInfo{
		IsWorking: false,
		Done:      make(chan struct{}),
		PollInfo: &domain.PollInfo{
			Results: &domain.PollResults{
				Results: []domain.Result{
					{
						Title: "Здоров",
						Users: make(map[int64]bool),
					},
					{
						Title: "Болен",
						Users: make(map[int64]bool),
					},
					{
						Title: "Сдал",
						Users: make(map[int64]bool),
					},
					{
						Title: "Отрицательный",
						Users: make(map[int64]bool),
					},
					{
						Title: "Положительный",
						Users: make(map[int64]bool),
					},
				},
				All: make(map[int64]bool),
			},
		},
		HeadPerson: &domain.User{
			ID:        371947069,
			Username:  "oooMRXooo",
			Firstname: "Alexander",
			Lastname:  "Oleynikov",
			ChatID:    371947069,
		},
		Users: make(map[int64]domain.User),
	}

	return newService
}

func (srv *Service) HelpGroup(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "About group bot...")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) HelpPrivate(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "About private bot...")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) UnknownCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, but I don't know this command( Please, type /help command")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) Start(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Welcome!")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) StartRoutine(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if srv.chatInfo[msg.Chat.ID].IsWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already started")
		_, err := bot.Send(newMsg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	times, err := parseTimeSliceFromString(msg.CommandArguments())
	if err != nil {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		newMsg.ParseMode = "MarkdownV2"
		_, err := bot.Send(newMsg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	srv.chatInfo[msg.Chat.ID].PollInfo.Times = times

	srv.chatInfo[msg.Chat.ID].Done = make(chan struct{}, 1)
	srv.chatInfo[msg.Chat.ID].IsWorking = true

	go srv.createPolls(bot, msg.Chat.ID)

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot started...")
	_, err = bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func parseTimeSliceFromString(argTime string) ([]domain.MentionTime, error) {
	times := strings.Split(argTime, " ")

	if len(times) < 2 {
		return nil, fmt.Errorf("you must enter 2 times\\! `/start_routine [poll creation time] [first reminder time] [second reminder time] ...`")
	}

	parsedTimes := make([]domain.MentionTime, len(times))

	for i, t := range times {
		parsedTime, err := time.Parse("15:04:05", t)
		if err != nil {
			return nil, fmt.Errorf("error in %d time format", i+1)
		}
		parsedTimes[i] = domain.MentionTime{
			MenTime: parsedTime,
			Done:    false,
		}
	}

	return parsedTimes, nil
}

func (srv *Service) createPolls(bot *tgbotapi.BotAPI, chatID int64) {
	for {
		select {
		case <-srv.chatInfo[chatID].Done:
			return
		default:
			time.Sleep(10 * time.Second)
		}

		times := srv.chatInfo[chatID].PollInfo.Times
		timeNowRaw := time.Now()
		timeNow := time.Date(0000, 01, 01, timeNowRaw.Hour(), timeNowRaw.Minute(), timeNowRaw.Second(), 0, time.UTC)

		if timeNow.Sub(srv.chatInfo[chatID].PollInfo.CreationDate).Hours() >= 24 {
			for i := range srv.chatInfo[chatID].PollInfo.Times {
				srv.chatInfo[chatID].PollInfo.Times[i].Done = false
			}
		}

		// create poll - times[0]
		if !times[0].Done && timeNow.After(times[0].MenTime) && timeNow.Before(times[1].MenTime) {
			newPoll := tgbotapi.NewPoll(
				chatID,
				fmt.Sprintf("Здоровье %s", timeNowRaw.Format("02.01")),
				[]string{
					"Здоров",
					"Болен",
					"Сдал",
					"Отрицательный",
					"Положительный",
				}...,
			)

			newPoll.IsAnonymous = false

			poll, err := bot.Send(newPoll)
			if err != nil {
				fmt.Println(err)
				return
			}

			srv.pollChat[poll.Poll.ID] = poll.Chat.ID

			srv.chatInfo[chatID].PollInfo.ID = poll.Poll.ID
			srv.chatInfo[chatID].PollInfo.CreationDate = timeNowRaw
			srv.chatInfo[chatID].PollInfo.Times[0].Done = true
		}

		for i := 1; i < len(times)-2; i++ {
			if !times[i].Done && timeNow.After(times[i].MenTime) && timeNow.Before(times[i+1].MenTime) {
				var mentionUsers string

				for key := range srv.chatInfo[chatID].Users {
					if _, ok := srv.chatInfo[chatID].PollInfo.Results.All[key]; !ok {
						tmpUser := srv.chatInfo[chatID].Users[key]
						firstLetter, _ := utf8.DecodeRuneInString(tmpUser.Firstname)
						mentionUsers += fmt.Sprintf("[%s %s\\.](tg://user?id=%d)\n", tmpUser.Lastname, string(firstLetter), key)
					}
				}

				newMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Lets take a survey:\n%s", mentionUsers))
				newMsg.ParseMode = "MarkdownV2"

				_, err := bot.Send(newMsg)
				if err != nil {
					fmt.Println(err)
				}

				srv.chatInfo[chatID].PollInfo.Times[i].Done = true
			}
		}

		// mention head - times[len(times)-2]
		if !times[len(times)-2].Done && timeNow.After(times[len(times)-2].MenTime) && timeNow.Before(times[len(times)-1].MenTime) {
			var mentionUsers string

			for key := range srv.chatInfo[chatID].Users {
				if _, ok := srv.chatInfo[chatID].PollInfo.Results.All[key]; !ok {
					u := srv.chatInfo[chatID].Users[key]
					firstLetter, _ := utf8.DecodeRuneInString(u.Firstname)
					mentionUsers += fmt.Sprintf("%s %s\\.\n", u.Lastname, string(firstLetter))
				}
			}

			headPerson := srv.chatInfo[chatID].HeadPerson
			firstLetter, _ := utf8.DecodeRuneInString(headPerson.Firstname)

			newMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
				"[%s %s\\.](tg://user?id=%d),\nthese users haven't passed the survey\\:\n%s",
				headPerson.Lastname,
				string(firstLetter),
				headPerson.ID,
				mentionUsers,
			))
			newMsg.ParseMode = "MarkdownV2"

			_, err := bot.Send(newMsg)
			if err != nil {
				fmt.Println(err)
			}

			srv.chatInfo[chatID].PollInfo.Times[len(times)-2].Done = true
		}

		// send poll - times[len(times)-1]
		if !times[len(times)-1].Done &&
			timeNow.After(times[len(times)-1].MenTime) &&
			timeNow.Before(srv.chatInfo[chatID].PollInfo.CreationDate.Add(24*time.Hour)) {

			srv.SendResultsTimeout(bot, chatID)

			srv.chatInfo[chatID].PollInfo.Times[len(times)-1].Done = true
		}
	}
}

func (srv *Service) StopRoutine(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	var newMsg tgbotapi.MessageConfig

	if !srv.chatInfo[msg.Chat.ID].IsWorking {
		newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot is already stopped")
		_, err := bot.Send(newMsg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	srv.chatInfo[msg.Chat.ID].Done <- struct{}{}
	srv.chatInfo[msg.Chat.ID].IsWorking = false

	newMsg = tgbotapi.NewMessage(msg.Chat.ID, "Bot stopped...")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) AddNewChat(bot *tgbotapi.BotAPI, chatID int64) {
	srv.chatInfo[chatID] = &domain.ChatInfo{
		IsWorking: false,
		Done:      make(chan struct{}),
		PollInfo: &domain.PollInfo{
			Results: &domain.PollResults{
				Results: []domain.Result{
					{
						Title: "Здоров",
						Users: make(map[int64]bool),
					},
					{
						Title: "Болен",
						Users: make(map[int64]bool),
					},
					{
						Title: "Сдал",
						Users: make(map[int64]bool),
					},
					{
						Title: "Отрицательный",
						Users: make(map[int64]bool),
					},
					{
						Title: "Положительный",
						Users: make(map[int64]bool),
					},
				},
				All: make(map[int64]bool),
			},
		},
		HeadPerson: &domain.User{
			ID:        371947069,
			Username:  "oooMRXooo",
			Firstname: "Alexander",
			Lastname:  "Oleynikov",
			ChatID:    371947069,
		},
		Users: make(map[int64]domain.User),
	}

	keyboard := tgbotapi.NewInlineKeyboardButtonData("Hi!", RegisterNewUser)
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{keyboard})

	newMsg := tgbotapi.NewMessage(chatID, "Hello everyone!")
	newMsg.ReplyMarkup = keyboardMarkup

	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Added new chat")
}

func (srv *Service) RemoveChat(bot *tgbotapi.BotAPI, chatID int64) {
	delete(srv.chatInfo, chatID)

	newMsg := tgbotapi.NewMessage(chatID, "Goodbuy everyone!")
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Removed chat")
}

func (srv *Service) RegisterNewUser(bot *tgbotapi.BotAPI, chatID int64, user *domain.User) {
	var newMsg tgbotapi.MessageConfig

	oldUser, ok := srv.chatInfo[chatID].Users[user.ID]
	firstLetter, _ := utf8.DecodeRuneInString(oldUser.Firstname)
	if ok {
		newMsg = tgbotapi.NewMessage(
			chatID,
			fmt.Sprintf("And hello again %s %s., but we have already greeted each other)",
				oldUser.Lastname,
				string(firstLetter),
			),
		)
		_, err := bot.Send(newMsg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// add to DB
	srv.chatInfo[chatID].Users[user.ID] = *user

	firstLetter, _ = utf8.DecodeRuneInString(user.Firstname)

	newMsg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Hello %s %s.!", user.Lastname, string(firstLetter)))
	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) UpdatePollResults(bot *tgbotapi.BotAPI, pollID string, userID int64, optionIDs []int) {
	chatID := srv.pollChat[pollID]

	srv.chatInfo[chatID].PollInfo.Results.All[userID] = true
	srv.chatInfo[chatID].PollInfo.Results.Results[int64(optionIDs[0])].Users[userID] = true

	if len(srv.chatInfo[chatID].PollInfo.Results.All) == len(srv.chatInfo[chatID].Users) {
		srv.SendResults(bot, chatID)
	}
}

func (srv *Service) SendResults(bot *tgbotapi.BotAPI, chatID int64) {
	var (
		resultsMsg string
		sickCount  int32
	)

	for _, res := range srv.chatInfo[chatID].PollInfo.Results.Results {
		resultsMsg += fmt.Sprintf("%s:\n", res.Title)

		for userID := range res.Users {
			if res.Title == "Болен" || res.Title == "Положительный" {
				sickCount++
			}

			tmpUser := srv.chatInfo[chatID].Users[userID]
			firstLetter, _ := utf8.DecodeRuneInString(tmpUser.Firstname)

			resultsMsg += fmt.Sprintf("%s %s.\n", tmpUser.Lastname, string(firstLetter))
		}
		resultsMsg += "\n"
	}

	newMsg := tgbotapi.NewMessage(srv.chatInfo[chatID].HeadPerson.ID, resultsMsg)

	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}

	if sickCount == 0 {
		newMsg = tgbotapi.NewMessage(srv.chatInfo[chatID].HeadPerson.ID, "Сегодня заболевших нет!")
		_, err := bot.Send(newMsg)
		if err != nil {
			fmt.Println(err)
		}

		newSticker := tgbotapi.NewSticker(
			srv.chatInfo[chatID].HeadPerson.ID,
			tgbotapi.FileID("CAACAgIAAxkBAAIBwmHsamDfSid1MDODOgWbAAG-FbfHCgACegEAAhAabSKcnbdITu9y3iME"),
		)

		_, err = bot.Send(newSticker)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (srv *Service) SendResultsTimeout(bot *tgbotapi.BotAPI, chatID int64) {
	var (
		resultsMsg string
		sickCount  int32
	)

	for _, res := range srv.chatInfo[chatID].PollInfo.Results.Results {
		resultsMsg += fmt.Sprintf("%s:\n", res.Title)

		for userID := range res.Users {
			if res.Title == "Болен" || res.Title == "Положительный" {
				sickCount++
			}

			tmpUser := srv.chatInfo[chatID].Users[userID]
			firstLetter, _ := utf8.DecodeRuneInString(tmpUser.Firstname)

			resultsMsg += fmt.Sprintf("%s %s.\n", tmpUser.Lastname, string(firstLetter))
		}
		resultsMsg += "\n"
	}

	newMsg := tgbotapi.NewMessage(srv.chatInfo[chatID].HeadPerson.ID, resultsMsg)

	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}

	var mentionUsers string

	for key := range srv.chatInfo[chatID].Users {
		if _, ok := srv.chatInfo[chatID].PollInfo.Results.All[key]; !ok {
			tmpUser := srv.chatInfo[chatID].Users[key]
			firstLetter, _ := utf8.DecodeRuneInString(tmpUser.Firstname)
			mentionUsers += fmt.Sprintf("%s %s.\n", tmpUser.Lastname, string(firstLetter))
		}
	}

	newMsg = tgbotapi.NewMessage(srv.chatInfo[chatID].HeadPerson.ID, fmt.Sprintf("This people hasn't passed the survey:\n%s", mentionUsers))

	_, err = bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Service) SetHeadPerson(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	user := msg.From

	srv.chatInfo[msg.Chat.ID].HeadPerson = &domain.User{
		ID:        user.ID,
		Username:  user.UserName,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		ChatID:    user.ID,
	}

	newMsg := tgbotapi.NewMessage(user.ID, "Now you are head person, congratulations!)")

	_, err := bot.Send(newMsg)
	if err != nil {
		fmt.Println(err)
	}
}
