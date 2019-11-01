package service

import "github.com/nicksnyder/go-i18n/v2/i18n"

type Message struct {
	chatID      int64
	messageLang string
	localizer   *i18n.Localizer
	reply       string
}

func (message *Message) Reply() string {
	return message.reply
}

func (message *Message) SetReply(reply string) {
	message.reply = reply
}

func (message *Message) ChatID() int64 {
	return message.chatID
}

func (message *Message) SetLocalizer(localizer *i18n.Localizer) {
	message.localizer = localizer
}

func (message *Message) Localizer() *i18n.Localizer {
	return message.localizer
}

func (message *Message) SetMessageLang(messageLang string) {
	message.messageLang = messageLang
}

func (message *Message) MessageLang() string {
	return message.messageLang
}

func (message *Message) translateReply() string {
	return message.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: message.reply})
}
