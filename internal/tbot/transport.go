package tbot

import (
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"jane2/utils"
)

func NewTransport(endpoints Endpoints, b *tele.Bot, config utils.TBot) (*tele.Bot, error) {
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello man!")
	}, middleware.Whitelist(config.Channel))

	b.Handle(tele.OnText, endpoints.FreeMessage, middleware.Whitelist(config.Channel))

	b.Handle(tele.OnVoice, endpoints.OnVoice, middleware.Whitelist(config.Channel))

	b.Start()

	return b, nil
}
