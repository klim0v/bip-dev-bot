package callback

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BuyCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *BuyCoinCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.EnterMinterAddress)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.SelectMinterAddressMarkup(callback.Localizer(), callback.Repository.MinterAddresses())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
