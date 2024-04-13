// y0_AgAAAAABp1fqAATuwQAAAADu9iZJBqtW5P-DTF-JWQ4AOMYF8IG3DFc
package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	tele "gopkg.in/telebot.v3"
	openapi2 "jane2/internal/openapi"
	"jane2/internal/tbot"
	"jane2/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"

	l, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	logger := l.Sugar()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	pref := tele.Settings{
		Token:  cfg.TBot.Token,
		Poller: &tele.LongPoller{Timeout: cfg.TBot.Timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatalf("config loading error %v", err.Error())
	}

	openApiClient := openapi2.NewClient(cfg.App.OpenApiToken)
	//tts := tts.NewClient(cfg.App.YandexKitToken)
	tBotEndpoints := tbot.MakeEndpoints(ctx, openApiClient)

	go func() {
		log.Println("Start telegram channel")
		_, err := tbot.NewTransport(tBotEndpoints, b, cfg.TBot)
		if err != nil {
			log.Println("Error: Telegram bot wasn`t running.", err)
			errs <- err
		}
	}()

	log.Println("exit", <-errs)
}
