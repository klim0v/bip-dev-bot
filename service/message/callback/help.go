package callback

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HelpCallbackFactory struct {
	CallbackFactory
}

func (callback *HelpCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.Help)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.HelpMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
