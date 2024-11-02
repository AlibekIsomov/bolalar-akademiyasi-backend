package telegramBot

import (
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	TelegramBot struct {
		Token         string `yaml:"token"`
		OwnerUsername string `yaml:"owner_username"`
	} `yaml:"telegram_bot"`
}

func loadConfig() (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		return nil, err
	}

	// Parse the YAML file
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

var admins = make(map[string]bool)

func isAdmin(username string) bool {
	return admins[username]
}

var (
	userState = make(map[int64]string)
	userData  = make(map[int64]map[string]string)
)

func Telegrambot() {
	// Load the config
	cfg, err := loadConfig()
	if err != nil {
		log.Panic("Failed to load config:", err)
	}

	token := cfg.TelegramBot.Token
	if token == "" {
		log.Panic("Telegram bot token not set in config.yml")
	}

	ownerUsername := cfg.TelegramBot.OwnerUsername
	if ownerUsername == "" {
		log.Panic("Owner username not set in config.yml")
	}

	// Add the owner as an admin automatically
	admins[ownerUsername] = true

	// Create the bot using the token from the config file
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 240

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")

		switch userState[chatID] {
		case ENTER_NAME:
			handleName(update.Message, &msg)
		case ENTER_PHONE_NUMBER:
			handlePhoneNumber(update.Message, &msg)
		case ENTER_AGE:
			handleAge(update.Message, &msg)
		default:
			if update.Message.IsCommand() {
				if !handleCommand(update.Message, &msg) {
					msg.Text = "There's no such command."
				}
			} else {
				if !handleMessage(update.Message, &msg, bot) {
					msg.Text = "There's no such thing."
				}
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func handleCommand(message *tgbotapi.Message, msg *tgbotapi.MessageConfig) bool {
	switch message.Command() {
	case "start":
		chatID := message.Chat.ID
		if isRegistered(chatID) {
			msg.Text = "Qaytib kelganingizdan xursandmiz!"
			msg.ReplyMarkup = StartKeyboardRegistered
		} else {
			msg.Text = "Xush kelibsiz! Iltimos, roʻyxatdan oʻting."
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Ro'yxatdan o'tish"),
				),
			)
		}
		return true
	}
	return false
}

func getUserPhoneNumber(chatID int64) string {
	var client models.Client
	result := database.DB.First(&client, "chat_id = ?", chatID)
	if result.Error != nil {
		return ""
	}
	return client.PhoneNumber
}

func handleMessage(message *tgbotapi.Message, msg *tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) bool {
	chatID := message.Chat.ID

	switch userState[chatID] {
	case ENTER_NAME:
		handleName(message, msg)
	case ENTER_PHONE_NUMBER:
		handlePhoneNumber(message, msg)
	case ENTER_AGE:
		handleAge(message, msg)
	default:
		switch message.Text {
		case "Ro'yxatdan o'tish":
			if !isRegistered(chatID) {
				userState[chatID] = ENTER_NAME
				msg.Text = "Iltimos, Ismingizni kiritng."
			} else {
				msg.Text = "Siz oldindan ro'yxatdan o'tgansiz."
				msg.ReplyMarkup = StartKeyboardRegistered
			}
		case LOCATION:
			handleLocation(msg)
		case CONTACT:
			handleContact(msg)
		case ABOUT_US:
			handleAboutUs(msg)
		default:
			return false
		}
	}
	return msg.Text != ""
}

func handleLocation(msg *tgbotapi.MessageConfig) {
	msg.Text = "Bizning akademiyas shu yerda joylashgan: BLAH BLAH BLAH"
	// Optionally, you can send a location using tgbotapi.NewLocation()
}

func handleContact(msg *tgbotapi.MessageConfig) {
	msg.Text = "Biz bilan bog'laning!:\nTelefon: +1234567890\nEmail: info@example.com"
}

func handleAboutUs(msg *tgbotapi.MessageConfig) {
	msg.Text = "We are a leading academy dedicated to children's education and development."
}

func handleName(message *tgbotapi.Message, msg *tgbotapi.MessageConfig) {
	if len(message.Text) < 3 {
		msg.Text = "Iltimos, ismingizni to'g'ri kiriting!"
		return
	}
	userData[message.Chat.ID] = make(map[string]string)
	userData[message.Chat.ID]["name"] = message.Text
	userState[message.Chat.ID] = ENTER_PHONE_NUMBER
	msg.Text = "Telefon raqamingizni kiriting:"
}
func handlePhoneNumber(message *tgbotapi.Message, msg *tgbotapi.MessageConfig) {
	if isRegistered(message.Chat.ID) {
		msg.Text = "Siz oldindan ro'yxatdan o'tgansiz"
		return
	}
	if !isValidPhoneNumber(message.Text) {
		msg.Text = "Iltimos raqamingizni +998 bilan boshlang!"
		return
	}
	userData[message.Chat.ID]["phone_number"] = message.Text
	userState[message.Chat.ID] = ENTER_AGE
	msg.Text = "Iltimos, Bolangiz yoshini kiriting:"
}

func handleAge(message *tgbotapi.Message, msg *tgbotapi.MessageConfig) {
	age, err := strconv.Atoi(message.Text)
	if err != nil || age <= 2 || age >= 18 {
		msg.Text = "Bunaqa yoshadagi bolalar biz uchun mos emas!"
		return
	}

	client := models.Client{
		Name:        userData[message.Chat.ID]["name"],
		PhoneNumber: userData[message.Chat.ID]["phone_number"],
		Age:         age,
		ChatID:      message.Chat.ID, // Save the chat ID
		Status:      models.Active,          // Set the status (e.g., "Active")
		Source:      models.Telegram,
	}

	result := database.DB.Create(&client)
	if result.Error != nil {
		msg.Text = "Ro'yxatdan o'tishda muammolik, iltimos qaytadan urinib ko'ring."
	} else {
		msg.Text = "Siz muvaffaqiyatli ro'yxatdan o'tdingiz!"
		msg.ReplyMarkup = StartKeyboardRegistered
		delete(userState, message.Chat.ID)
		delete(userData, message.Chat.ID)
	}
}

func isRegistered(chatID int64) bool {
	var client models.Client
	result := database.DB.First(&client, "chat_id = ?", chatID)
	return result.Error == nil
}

func isValidPhoneNumber(phone string) bool {
	return len(phone) >= 12 && strings.HasPrefix(phone, "+998")
}
