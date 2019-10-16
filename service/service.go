package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Service struct {
	bot *tgbotapi.BotAPI
}

func NewService(bot *tgbotapi.BotAPI) *Service {
	bot.Debug = true
	return &Service{bot: bot}
}

func (s *Service) UpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return s.bot.GetUpdatesChan(u)

}

func (s *Service) Handle(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := s.bot.Send(msg)
	return err
}
