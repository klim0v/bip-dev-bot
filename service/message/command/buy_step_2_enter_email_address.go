package command

import (
	"bip-dev/service/message"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EnterEmailAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterEmailAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !IsValidEmailAddress(command.Args) {
		command.Message.SetReply(message.EnterEmailAddress)
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Translate(command.Message.Reply()+"_invalid"),
		)
		msg.ParseMode = "markdown"

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}

	emailID, err := command.Repository.AddEmailAddress(command.ChatID(), command.Args)
	if err != nil {
		return err
	}

	if err := command.Repository.SaveEmailAddressForBuy(command.ChatID(), emailID); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterBitcoinAddress)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply()), 0.0184, -24.28, 516841, 4.00, command.Repository.BtcAddresses()),
	)
	msg.ReplyMarkup = message.SendBTCAddressMarkup(command.Localizer())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
