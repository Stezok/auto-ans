package main

import (
	"log"
	"time"

	"github.com/Stezok/auto-ans/internal/bot/telegram"
	"github.com/Stezok/auto-ans/internal/timemanager"
)

func main() {
	timeManager := timemanager.NewManager()
	timeManager.Push(
		timemanager.NewTimeSegment(time.Hour*0, time.Hour*10, "Здравствуйте, спасибо за обращение!\nОператор ответит на все ваши вопросы в рабочее время с 10:00 до 20:00 (GMT +2:00), с понедельника по пятницу."),
		timemanager.NewTimeSegment(time.Hour*10+time.Second, time.Hour*20, "Здравствуйте, спасибо за обращение! Оператор свяжется с вами в ближайшее время. Пока задайте свой вопрос."),
		timemanager.NewTimeSegment(time.Hour*20, time.Hour*24, "Здравствуйте, спасибо за обращение!\nОператор ответит на все ваши вопросы в рабочее время с 10:00 до 20:00 (GMT +2:00), с понедельника по пятницу."),
	)

	bot, err := telegram.NewTelegramBot("1797876808:AAFUeiOkIVr93hZvQ_0FSUjUXrS1rbJOwEQ", timeManager, 496823111)
	// bot, err := telegram.NewTelegramBot("1057474366:AAGloprVxM9DP13UvgE75A5cBVyyJAf1M4s", timeManager, 768413750)
	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
