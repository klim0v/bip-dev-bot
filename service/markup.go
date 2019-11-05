package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func helpMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "buy_coin"}), "buy_coin"),
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "sell_coin"}), "sell_coin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "my_orders"}), "my_orders"),
		),
	)
}

func sendBTCAddressMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), "check_sell"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "buy_coin"),
		),
	)
}

func selectBitcoinMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("use_bitcoin_address %d", address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "send_price_coin"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func selectEmailAddressMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("use_email_address %d", address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "buy_coin"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func selectMinterAddressMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("use_minter_address %d", address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "help"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func selectCoinNameMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	var keyboards [][]tgbotapi.InlineKeyboardButton
	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "help"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(keyboards...)
}
