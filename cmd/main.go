package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"Tabdil/internal/service"
)

// نگهداری وضعیت هر کاربر
var userState = make(map[int64]string)

func main() {

	// خواندن فایل env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		// ===============================
		// هندل پیام‌های متنی
		// ===============================
		if update.Message != nil {

			chatID := update.Message.Chat.ID
			text := update.Message.Text

			// دستور start
			if text == "/start" {

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("محاسبه طول", "length"),
						tgbotapi.NewInlineKeyboardButtonData("محاسبه ارز", "currency"),
					),
				)

				msg := tgbotapi.NewMessage(chatID,
					"سلام 👋\nبه بات تبدیل خوش آمدید\nیک گزینه انتخاب کن:")
				msg.ReplyMarkup = keyboard

				bot.Send(msg)
				continue
			}

			// اگر کاربر در حالت وارد کردن عدد است
			switch userState[chatID] {

			case "km_to_mile":

				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن ❌")
					continue
				}

				result, err := service.KilometerToMile(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				send(bot, chatID,
					strconv.FormatFloat(result, 'f', 4, 64)+" Mile")

				userState[chatID] = ""

			case "kg_to_pound":

				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن ❌")
					continue
				}

				result, err := service.KilogramToPound(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				send(bot, chatID,
					strconv.FormatFloat(result, 'f', 4, 64)+" Pound")

				userState[chatID] = ""
			}
		}

		// ===============================
		// هندل دکمه‌ها
		// ===============================
		if update.CallbackQuery != nil {

			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			switch data {

			case "length":

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("KM ➜ Mile", "km_to_mile"),
						tgbotapi.NewInlineKeyboardButtonData("KG ➜ Pound", "kg_to_pound"),
					),
				)

				msg := tgbotapi.NewMessage(chatID,
					"یک نوع تبدیل انتخاب کن:")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "km_to_mile":
				userState[chatID] = "km_to_mile"
				send(bot, chatID, "مقدار کیلومتر رو وارد کن:")

			case "kg_to_pound":
				userState[chatID] = "kg_to_pound"
				send(bot, chatID, "مقدار کیلوگرم رو وارد کن:")

			case "currency":
				send(bot, chatID,
					"بخش ارز هنوز پیاده‌سازی نشده 😉")
			}

			// پاسخ به Callback برای جلوگیری از لودینگ
			callback := tgbotapi.NewCallback(
				update.CallbackQuery.ID, "")
			bot.Request(callback)
		}
	}
}

// تابع کمکی برای ارسال پیام
func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
