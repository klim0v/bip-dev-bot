package callback

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type UseMinterAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseMinterAddressCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {

	minterID, err := strconv.Atoi(callback.Args)
	if err != nil {
		return err
	}

	if err := callback.Repository.SaveBuyMinterAddress(callback.ChatID(), minterID); err != nil {
		return err
	}

	callback.Message.SetReply(message.EnterEmailAddress)

	addresses := callback.Repository.EmailAddresses()

	var msg tgbotapi.EditMessageTextConfig
	if len(addresses) == 0 {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate("new_email"))
	} else {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	}

	markup := message.SelectEmailAddressMarkup(callback.Localizer(), addresses)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
