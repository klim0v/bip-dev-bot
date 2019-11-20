package command

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EnterMinterAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterMinterAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !IsValidMinterAddress(command.Args) {
		command.Message.SetReply(message.EnterMinterAddress)
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

	emailID, err := command.Repository.AddMinterAddress(command.ChatID(), command.Args)
	if err != nil {
		return err
	}

	if err := command.Repository.SaveEmailAddressForBuy(command.ChatID(), emailID); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterEmailAddress)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Translate(command.Message.Reply()),
	)
	msg.ReplyMarkup = message.SelectEmailAddressMarkup(command.Localizer(), command.Repository.EmailAddresses())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
