package tbot

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type client struct {
	bot *tgbotapi.BotAPI
}
