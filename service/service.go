package service

import (
	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
)

type Service struct {
	bot            *tgbotapi.BotAPI
	languageBundle *i18n.Bundle
}

func NewService(bot *tgbotapi.BotAPI) *Service {
	bot.Debug = true
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.en.toml"))
	bundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.ru.toml"))
	return &Service{bot: bot, languageBundle: bundle}
}

func (s *Service) UpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return s.bot.GetUpdatesChan(u)

}

func (s *Service) Handle(update tgbotapi.Update) error {
	err := s.replyMessage(update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) executeCommand(update tgbotapi.Update, localizer *i18n.Localizer) (msg tgbotapi.MessageConfig, err error) {
	switch update.Message.Command() {
	case "help":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "помощи не будет, давай уж как нибудь сам :/")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(
					localizer.MustLocalize(
						&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{
								ID: "ByCoin",
							},
						},
					),
				),
				tgbotapi.NewKeyboardButton("Продать"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Мои заявки"),
			),
		)
	}
	return msg, nil
}

type Chat struct {
	LastBotMessage string
	Lang           string
}

func (s *Service) replyMessage(update tgbotapi.Update) (err error) {
	// todo get lang and last bot message by update.Message.Chat.ID from redis
	//if err != nil {
	//	return msg, err
	//}
	chat := Chat{}

	localizer := i18n.NewLocalizer(s.languageBundle, chat.Lang)

	if update.Message.IsCommand() {
		msg, err := s.executeCommand(update, localizer)
		if err != nil {

		}
		_, err = s.bot.Send(msg)
		if err != nil {

		}
		return nil
	}

	//message, err := localizer.LocalizeMessage()
	//if err != nil {
	//	return msg, err
	//}
	//msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
	//msg.ReplyToMessageID = update.Message.MessageID

	//to do save msg to radis

	return nil
}
