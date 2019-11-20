package callback

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SellCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *SellCoinCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.EnterCoinName)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.SelectCoinNameMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
