package telegramBot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"


var StartKeyboardRegistered = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(LOCATION),
		tgbotapi.NewKeyboardButton(CONTACT),
		tgbotapi.NewKeyboardButton(ABOUT_US),
	),

)