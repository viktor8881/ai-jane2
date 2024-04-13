package tbot

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"jane2/internal/openapi"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Service interface {
	AutoTrading(ctx context.Context) error
}

type Endpoints struct {
	FreeMessage telebot.HandlerFunc
	OnVoice     telebot.HandlerFunc
}

func MakeEndpoints(ctx context.Context, openapiClient *openapi.Client) Endpoints {
	return Endpoints{
		FreeMessage: makeFreeMessageEndpoint(ctx, openapiClient),
		OnVoice:     makeOnVoiceEndpoint(ctx, openapiClient),
	}
}

func makeFreeMessageEndpoint(ctx context.Context, openapiClient *openapi.Client) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		text, err := openapiClient.Chat(ctx, c.Message().Text)
		if err != nil {
			_ = c.Send("some thing went wrong: " + err.Error())
			return err
		}
		_ = c.Send(text)

		resp, err := openapiClient.CreateVoice(ctx, text)
		if err != nil {
			_ = c.Send("create voice error: " + err.Error())
			return err
		}

		return c.Send(&telebot.Voice{File: telebot.FromReader(resp)})
	}
}

func makeOnVoiceEndpoint(ctx context.Context, openapiClient *openapi.Client) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		fnameMp3 := "./data/request.mp3"

		err := downloadAndConvertFileToMp3(c, fnameMp3)
		if err != nil {
			_ = c.Send("conver file error: " + err.Error())
			return err
		}

		// wisper-3
		respText, err := openapiClient.CreateText(ctx, fnameMp3)
		if err != nil {
			_ = c.Send("whisper error: " + err.Error())
			return err
		}

		lowerStr := strings.ToLower(respText)
		_ = c.Send(lowerStr)
		if lowerStr == "новый разговор" || lowerStr == "новый диалог" || lowerStr == "новый разговор." || lowerStr == "новый диалог." {
			openapiClient.ClearChatContext()
			return c.Send("Начинаем новый диалог.")
		}

		text, err := openapiClient.Chat(ctx, respText)
		if err != nil {
			_ = c.Send("openApi error: " + err.Error())
			return err
		}
		_ = c.Send(text)

		resp, err := openapiClient.CreateVoice(ctx, text)
		if err != nil {
			_ = c.Send("create voice error: " + err.Error())
			return err
		}

		return c.Send(&telebot.Voice{File: telebot.FromReader(resp)})
	}
}

func downloadAndConvertFileToMp3(c telebot.Context, mp3path string) error {
	os.Remove(mp3path)
	//return nil
	tmpFile := fmt.Sprintf("%s%d", mp3path, time.Now().UnixNano()/int64(time.Millisecond))

	b := c.Bot()
	err := b.Download(&c.Message().Voice.File, tmpFile)
	if err != nil {
		return err
	}

	// Конвертирование файла в mp3
	cmd := exec.Command("ffmpeg", "-i", tmpFile, "-q:a", "0", mp3path)
	err = cmd.Run()
	if err != nil {
		return err
	}

	os.Remove(tmpFile)

	return nil
}
