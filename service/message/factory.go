package message

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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
	Answer(*tgbotapi.BotAPI) error
}

type AbstractFactory struct {
	Factory  Factory
	Resource Resource
}

func (a *AbstractFactory) Answer(bot *tgbotapi.BotAPI) error {
	return a.Factory.Answer(bot)
}

func (a *AbstractFactory) SaveLanguage(lang string) {
	a.Factory.SetLocalizer(a.Resource.Localizer(lang, a.Factory.MessageLang()))
	a.Resource.SaveLanguage(a.Factory.ChatID(), lang)
}

func (a *AbstractFactory) SaveReply() {
	a.Resource.SaveReply(a.Factory.ChatID(), a.Factory.Reply())
}

func (a *AbstractFactory) SetLocalizer() {
	a.Factory.SetLocalizer(a.Resource.Localizer(a.Resource.Language(a.Factory.ChatID()), a.Factory.MessageLang()))
}

func (a *AbstractFactory) Log(err error) {
	a.Resource.Log(err)
}
