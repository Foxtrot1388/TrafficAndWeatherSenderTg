package telegram

import (
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	apiEndpoint = "https://api.telegram.org/bot%s/%s"
)

type telegramBot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func New(token string, chatID int64) (*telegramBot, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(token, apiEndpoint, &http.Client{})
	if err != nil {
		return nil, err
	}
	return &telegramBot{bot: bot, chatID: chatID}, nil
}

func (tg *telegramBot) SendPhoto(name, text string, bytes []byte) error {

	b := tgbotapi.FileBytes{Name: name, Bytes: bytes}

	msg := tgbotapi.NewPhoto(tg.chatID, b)
	msg.Caption = text
	_, err := tg.bot.Send(msg)
	return err

}

func (tg *telegramBot) SendText(text string) error {

	msg := tgbotapi.NewMessage(tg.chatID, text)
	_, err := tg.bot.Send(msg)
	return err

}
