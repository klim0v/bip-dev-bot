package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func helpMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: buyCoin}), buyCoin),
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: sellCoin}), sellCoin),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: myOrders}), myOrders),
		),
	)
}

func sendBTCAddressMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), checkSell),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), buyCoin),
		),
	)
}

func sendYourCoinsMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), checkSell),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), enterPriceCoin),
		),
	)
}

func shareMarkup(localizer *i18n.Localizer, link string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "share"}), link),
		),
	)
}

func selectBitcoinMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("%s %d", useBitcoinAddress, address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), enterPriceCoin),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func selectEmailAddressMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("%s %d", useEmailAddress, address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), buyCoin),
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
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), help),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func selectCoinNameMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	var keyboards [][]tgbotapi.InlineKeyboardButton
	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), help),
	))
	return tgbotapi.NewInlineKeyboardMarkup(keyboards...)
}
