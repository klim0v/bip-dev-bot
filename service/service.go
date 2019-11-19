package service

import (
	"bip-dev/service/message"
	"bip-dev/service/message/callback"
	"bip-dev/service/message/command"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
	"strings"
)

type Application struct {
	Rds            *redis.Client
	Pgql           *sqlx.DB
	languageBundle *i18n.Bundle
	logger         *log.Logger
}

func (s *Application) Localizer(lang ...string) *i18n.Localizer {
	return i18n.NewLocalizer(s.languageBundle, lang...)
}

func (s *Application) SaveLanguage(chatID int64, lang string) {
	if err := s.saveLanguage(chatID, lang); err != nil {
		s.logger.Println(err)
	}
}

func (s *Application) saveLanguage(chatID int64, lang string) error {
	if err := s.Rds.Set(fmt.Sprintf("%d:lang", chatID), lang, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Application) SaveReply(chatID int64, lang string) {
	if err := s.saveReply(chatID, lang); err != nil {
		s.logger.Println(err)
	}
}

func (s *Application) saveReply(chatID int64, lang string) error {
	if err := s.Rds.Set(fmt.Sprintf("%d:last", chatID), lang, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Application) language(chatID int64) (string, error) {
	lang, err := s.Rds.Get(fmt.Sprintf("%d:lang", chatID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return lang, nil
}

func (s *Application) Language(chatID int64) string {
	lang, err := s.language(chatID)
	if err != nil {
		s.logger.Println(err)
		return ""
	}
	return lang
}

func (s *Application) lastCommand(chatID int64) (string, error) {
	lang, err := s.Rds.Get(fmt.Sprintf("%d:last", chatID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return lang, nil
}

func (s *Application) LastCommand(chatID int64) string {
	lang, err := s.lastCommand(chatID)
	if err != nil {
		s.logger.Println(err)
		return ""
	}
	return lang
}

func (s *Application) Log(err error) {
	s.logger.Println(err)
}

func NewApplication(rds *redis.Client, pgql *sqlx.DB, languageBundle *i18n.Bundle, logger *log.Logger) *Application {
	return &Application{
		Rds:            rds,
		Pgql:           pgql,
		languageBundle: languageBundle,
		logger:         logger,
	}
}

func (s *Application) NewFactory(update tgbotapi.Update) *message.AbstractFactory {
	if update.CallbackQuery != nil {
		fields := strings.Fields(update.CallbackQuery.Data)
		var args string
		if len(fields) == 2 {
			args = fields[1]
		}

		var concreteFactory message.Factory
		var msg message.Message
		msg.SetMessageLang(update.CallbackQuery.Message.From.LanguageCode)
		msg.SetChatID(update.CallbackQuery.Message.Chat.ID)
		callbackFactory := callback.CallbackFactory{
			Message:         msg,
			MessageUpdateID: update.CallbackQuery.Message.MessageID,
			Command:         fields[0],
			Args:            args,
			Repository:      message.NewRepository(s.Rds, s.Pgql),
		}

		switch callbackFactory.Command {
		case message.CheckSendDeposit:
			concreteFactory = &callback.CheckSendDepositCallbackFactory{CallbackFactory: callbackFactory, QueryID: update.CallbackQuery.ID}
		case message.SellCoin:
			concreteFactory = &callback.SellCoinCallbackFactory{CallbackFactory: callbackFactory}
		case message.BuyCoin:
			concreteFactory = &callback.BuyCoinCallbackFactory{CallbackFactory: callbackFactory}
		case message.UseEmailAddress:
			concreteFactory = &callback.UseEmailAddressCallbackFactory{CallbackFactory: callbackFactory}
		case message.UseMinterAddress:
			concreteFactory = &callback.UseMinterAddressCallbackFactory{CallbackFactory: callbackFactory}
		case message.CheckSell:
			concreteFactory = &callback.CheckSellCallbackFactory{CallbackFactory: callbackFactory, QueryID: update.CallbackQuery.ID}
		case message.UseBitcoinAddress:
			concreteFactory = &callback.UseBitcoinAddressCallbackFactory{CallbackFactory: callbackFactory}
		default:
			concreteFactory = &callback.HelpCallbackFactory{CallbackFactory: callbackFactory}
		}

		return &message.AbstractFactory{
			Factory:  concreteFactory,
			Resource: s,
		}
	}

	cmd := update.Message.Command()
	commandArguments := update.Message.CommandArguments()
	if !update.Message.IsCommand() {
		cmd = s.LastCommand(update.Message.Chat.ID)
		commandArguments = update.Message.Text
	}

	var concreteFactory message.Factory
	var msg message.Message
	msg.SetMessageLang(update.Message.From.LanguageCode)
	msg.SetChatID(update.Message.Chat.ID)
	commandFactory := command.CommandFactory{
		Message:    msg,
		Command:    cmd,
		Args:       commandArguments,
		Repository: message.NewRepository(s.Rds, s.Pgql),
	}

	switch commandFactory.Command {
	case message.EnterCoinName:
		concreteFactory = &command.EnterCoinNameCommandFactory{CommandFactory: commandFactory}
	case message.EnterPriceCoin:
		concreteFactory = &command.EnterPriceCoinCommandFactory{CommandFactory: commandFactory}
	case message.EnterBitcoinAddress:
		concreteFactory = &command.EnterBitcoinAddressCommandFactory{CommandFactory: commandFactory}
	case message.EnterEmailAddress:
		concreteFactory = &command.EnterEmailAddressCommandFactory{CommandFactory: commandFactory}
	case message.EnterMinterAddress:
		concreteFactory = &command.EnterMinterAddressCommandFactory{CommandFactory: commandFactory}
	default:
		concreteFactory = &command.HelpCommandFactory{CommandFactory: commandFactory}
	}

	return &message.AbstractFactory{
		Factory:  concreteFactory,
		Resource: s,
	}
}
