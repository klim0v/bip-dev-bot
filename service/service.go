package service

import (
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
	"strings"
)

type Application struct {
	rds            *redis.Client
	pgql           *sqlx.DB
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
	if err := s.rds.Set(fmt.Sprintf("%d:lang", chatID), lang, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Application) SaveReply(chatID int64, lang string) {
	if err := s.saveLanguage(chatID, lang); err != nil {
		s.logger.Println(err)
	}
}

func (s *Application) saveReply(chatID int64, lang string) error {
	if err := s.rds.Set(fmt.Sprintf("%d:last", chatID), lang, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Application) language(chatID int64) (string, error) {
	lang, err := s.rds.Get(fmt.Sprintf("%d:lang", chatID)).Result()
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
	lang, err := s.rds.Get(fmt.Sprintf("%d:last", chatID)).Result()
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
		rds:            rds,
		pgql:           pgql,
		languageBundle: languageBundle,
		logger:         logger,
	}
}

type Resource interface {
	Localizer(lang ...string) *i18n.Localizer
	SaveLanguage(chatID int64, lang string)
	SaveReply(chatID int64, reply string)
	Language(chatID int64) string
	LastCommand(chatID int64) string
	Log(error)
}

type Factory interface {
	ChatID() int64
	SetLocalizer(*i18n.Localizer)
	Localizer() *i18n.Localizer
	MessageLang() string
	Reply() string
	CreateMessage() tgbotapi.Chattable
	SaveArgs() error
}

type AbstractFactory struct {
	factory  Factory
	resource Resource
}

func (a *AbstractFactory) CreateMessage() tgbotapi.Chattable {
	return a.factory.CreateMessage()
}

func (a *AbstractFactory) SaveArgs() error {
	return a.factory.SaveArgs()
}

func (a *AbstractFactory) SaveLanguage(lang string) {
	a.factory.SetLocalizer(a.resource.Localizer(lang, a.factory.MessageLang()))
	a.resource.SaveLanguage(a.factory.ChatID(), lang)
}

func (a *AbstractFactory) SaveReply() {
	a.resource.SaveReply(a.factory.ChatID(), a.factory.Reply())
}

func (a *AbstractFactory) SetLocalizer() {
	a.factory.SetLocalizer(a.resource.Localizer(a.resource.Language(a.factory.ChatID()), a.factory.MessageLang()))
}

func (a *AbstractFactory) Log(err error) {
	a.resource.Log(err)
}

func (s *Application) NewFactory(update tgbotapi.Update) *AbstractFactory {
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" && update.CallbackQuery.Message != nil {
		fields := strings.Fields(update.CallbackQuery.Data)
		var args string
		if len(fields) == 2 {
			args = fields[1]
		}
		return &AbstractFactory{
			factory: &CallbackFactory{
				Message: Message{
					chatID:      update.Message.Chat.ID,
					messageLang: update.Message.From.LanguageCode,
					localizer:   nil,
					reply:       "",
				},
				MessageUpdateID: update.CallbackQuery.Message.MessageID,
				Command:         fields[0],
				Args:            args,
				Repository:      NewRepository(s.pgql),
			},
			resource: s,
		}
	}

	if update.Message != nil {
		var lastCommand string
		if !update.Message.IsCommand() {
			lastCommand = s.LastCommand(update.Message.Chat.ID)
		}
		return &AbstractFactory{
			factory: &CommandFactory{
				Message: Message{
					chatID:      update.Message.Chat.ID,
					messageLang: update.Message.From.LanguageCode,
					localizer:   nil,
					reply:       "",
				},
				LastMessage: lastCommand,
				Command:     update.Message.Command(),
				Args:        update.Message.CommandArguments(),
				Repository:  NewRepository(s.pgql),
			},
			resource: s,
		}
	}

	return nil
}
