package command

import (
	"bip-dev/service/message"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"strconv"
	"strings"
)

var matchCoinName = regexp.MustCompile("^[0-9-A-Z-a-z]{3,10}$")
var matchEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var matchBitcoin = regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$")

func isValidPriceCoin(name string, value string) bool {
	price, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	if name == "BIP" {
		return 0.01 <= price && price <= 0.32
	}

	return 0.01 <= price && price <= 1000
}
func isValidMinterAddress(address string) bool {
	address = strings.TrimSpace(address)

	if address == "Mx00000000000000000000000000000000000001" {
		return false
	}

	return len(address) == 42 && address[:2] != "Mx"
}

func isValidEmailAddress(email string) bool {
	if !matchEmail.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

func isValidCoinName(coinName string) bool {
	if !matchCoinName.MatchString(coinName) {
		return false
	}
	return true
}

func isValidBitcoinAddress(address string) bool {
	return matchBitcoin.MatchString(address)
}

type CommandFactory struct {
	message.Message
	Command    string
	Args       string
	Repository *message.Repository
}

type HelpCommandFactory struct {
	CommandFactory
}

func (command *HelpCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	command.Message.SetReply(message.Help)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Translate(command.Message.Reply()),
	)
	msg.ReplyMarkup = message.HelpMarkup(command.Localizer())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type EnterMinterAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterMinterAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !isValidMinterAddress(command.Args) {
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

type EnterEmailAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterEmailAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !isValidEmailAddress(command.Args) {
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

type EnterCoinNameCommandFactory struct {
	CommandFactory
}

func (command *EnterCoinNameCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !isValidCoinName(command.Args) {
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

	if err := command.Repository.SaveSellCoinName(command.ChatID(), command.Args); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterPriceCoin)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply())),
	)
	//msg.ReplyMarkup = todo
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type EnterPriceCoinCommandFactory struct {
	CommandFactory
}

func (command *EnterPriceCoinCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
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

	if !isValidPriceCoin(coinName, command.Args) {
		command.Message.SetReply(message.EnterPriceCoin)
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

	if err := command.Repository.SaveSellPrice(command.ChatID(), command.Args); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterBitcoinAddress)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply())),
	)
	msg.ReplyMarkup = message.SelectBitcoinMarkup(command.Localizer(), command.Repository.BtcAddresses())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type EnterBitcoinAddressCommandFactory struct {
	CommandFactory
}

func (command *EnterBitcoinAddressCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !isValidBitcoinAddress(command.Args) {
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
