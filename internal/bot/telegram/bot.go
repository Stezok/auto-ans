package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Stezok/auto-ans/internal/timemanager"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Appeal struct {
	id       int
	username string
	text     string
}

type TelegramBot struct {
	bot         *tgbotapi.BotAPI
	TimeManager *timemanager.Manager
	banList     []int
	admin       int64

	closeChan chan struct{}
}

func (tgbot *TelegramBot) IsBanned(id int) bool {
	for _, bannedID := range tgbot.banList {
		if id == bannedID {
			return true
		}
	}
	return false
}

func (tgbot *TelegramBot) NotifyAdmin(appeal Appeal) error {
	text := fmt.Sprintf("Новое обращение от %s.\n\n%s", appeal.username, appeal.text)
	mes := tgbotapi.NewMessage(tgbot.admin, text)

	data := fmt.Sprintf("k%d", appeal.id)
	mes.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принято", data),
		),
	)

	_, err := tgbot.bot.Send(mes)
	return err
}

func (tgbot *TelegramBot) handleCallback(update tgbotapi.Update) {
	defer recover()

	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil || update.CallbackQuery.Message.Chat == nil {
		return
	}

	command := update.CallbackQuery.Data[0]
	id := update.CallbackQuery.Data[1:]
	idInt, _ := strconv.Atoi(id)
	switch command {
	case 'k':
		tgbot.banList = append(tgbot.banList, idInt)

		messageID := update.CallbackQuery.Message.MessageID
		chatID := update.CallbackQuery.Message.Chat.ID
		data := fmt.Sprintf("d%s", id)
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Рассмотрено", data),
			),
		))
		_, err := tgbot.bot.Send(edit)
		if err != nil {
			log.Print(err)
		}

		message := tgbotapi.NewMessage(int64(idInt), "Оператор присоединился к чату.")
		_, err = tgbot.bot.Send(message)
		if err != nil {
			log.Print(err)
		}
	case 'd':
		for i, bannedID := range tgbot.banList {
			if bannedID == idInt {
				tgbot.banList = append(tgbot.banList[:i], tgbot.banList[i+1:]...)
			}
		}

		messageID := update.CallbackQuery.Message.MessageID
		chatID := update.CallbackQuery.Message.Chat.ID
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		_, err := tgbot.bot.Send(del)
		if err != nil {
			log.Print(err)
		}

		message := tgbotapi.NewMessage(int64(idInt), "Оператор покинул чат.")
		_, err = tgbot.bot.Send(message)
		if err != nil {
			log.Print(err)
		}
	}
}

func (tgbot *TelegramBot) handle(update tgbotapi.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	if tgbot.IsBanned(update.Message.From.ID) {
		return
	}

	messages := tgbot.TimeManager.Messages()
	if len(messages) == 0 {
		log.Printf("No handlers to time near %v", time.Now())
		return
	}

	text := messages[0]
	mes := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	_, err := tgbot.bot.Send(mes)
	if err != nil {
		log.Print(err)
	}

	err = tgbot.NotifyAdmin(Appeal{
		id:       update.Message.From.ID,
		username: update.Message.From.UserName,
		text:     update.Message.Text,
	})

	if err != nil {
		log.Print(err)
	}
}

func (tgbot *TelegramBot) Handle(updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case update := <-updates:
			tgbot.handle(update)
			tgbot.handleCallback(update)
		case <-tgbot.closeChan:
			return
		}
	}
}

func (tgbot *TelegramBot) Run() error {
	u := tgbotapi.NewUpdate(0)
	updates, err := tgbot.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	tgbot.Handle(updates)

	return nil
}

func (tgbot *TelegramBot) Close() {
	for i := 0; i < 1; i++ {
		tgbot.closeChan <- struct{}{}
	}
}

func NewTelegramBot(token string, timeManager *timemanager.Manager, admin int64) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	me, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:         bot,
		TimeManager: timeManager,
		admin:       admin,
		closeChan:   make(chan struct{}, 1),
		banList:     []int{me.ID},
	}, nil
}
