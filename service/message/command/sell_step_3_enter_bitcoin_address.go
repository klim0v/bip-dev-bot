package command

import (
	"bip-dev/service/message"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EnterBitcoinAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterBitcoinAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !IsValidBitcoinAddress(command.Args) {
		command.Message.SetReply(message.EnterBitcoinAddress)
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

	id, err := command.Repository.AddBitcoinAddress(command.ChatID(), command.Args)
	if err != nil {
		return err
	}

	err = command.Repository.SaveSellBitcoinAddress(command.ChatID(), id)
	if err != nil {
		return err
	}

	coinName, err := command.Repository.SellCoinName(command.ChatID())
	if err != nil {
		if err != redis.Nil {
			//todo: logging
		}

		command.Message.SetReply(message.EnterCoinName)
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

	link := "www.example.com"

	command.Message.SetReply(message.SendYourCoins)
	msg1 := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply()), coinName, link),
	)
	msg1.ReplyMarkup = message.ShareMarkup(command.Localizer(), link)
	msg1.ParseMode = "markdown"

	if _, err := bot.Send(msg1); err != nil {
		return err
	}

	msg2 := tgbotapi.NewMessage(
		command.ChatID(),
		"`Mx233750d042b2098409242d9fdfeee8aa51137738`",
	)
	msg2.ReplyMarkup = message.SendYourCoinsMarkup(command.Localizer())
	msg2.ParseMode = "markdown"

	if _, err := bot.Send(msg2); err != nil {
		return err
	}

	return nil
}
