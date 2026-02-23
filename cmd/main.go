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

// تابع کمکی برای ارسال منوی اصلی
func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64, text string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("محاسبه طول", "length"),
			tgbotapi.NewInlineKeyboardButtonData("محاسبه ارز", "currency"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// تابع کمکی برای ارسال پیام ساده
func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func main() {

	// خواندن فایل .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set")
	}

	apiURL := "https://api.example.com/dollar"

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

			// دستور /start
			if text == "/start" {
				sendMainMenu(bot, chatID, "سلام \nبه بات تبدیل خوش آمدید\nیک گزینه انتخاب کن:")
				continue
			}

			// اگر کاربر در حالت وارد کردن عدد است
			switch userState[chatID] {

			case "km_to_mile":
				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن")
					continue
				}

				result, err := service.KilometerToMile(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				resultText := strconv.FormatFloat(result, 'f', 4, 64) + " Mile"
				userState[chatID] = ""
				sendMainMenu(bot, chatID, "نتیجه: "+resultText+"\n\nیک گزینه دیگر انتخاب کن:")

			case "kg_to_pound":
				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن")
					continue
				}

				result, err := service.KilogramToPound(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				resultText := strconv.FormatFloat(result, 'f', 4, 64) + " Pound"
				userState[chatID] = ""
				sendMainMenu(bot, chatID, "نتیجه: "+resultText+"\n\nیک گزینه دیگر انتخاب کن:")

			case "cm_to_foot":
				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				result, err := service.CentimeterToFoot(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				userState[chatID] = ""
				sendMainMenu(bot, chatID, "نتیجه :"+strconv.FormatFloat(result, 'f', 4, 64)+"Foot\n\n یک گزینه دیگر انتخاب کن :")

			case "mbit_to_mbyte":
				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن")
					continue
				}

				result, err := service.MegabitToMegabyte(value)
				if err != nil {
					send(bot, chatID, err.Error())
					continue

				}

				userState[chatID] = ""
				sendMainMenu(bot, chatID,
					"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" MB\n\nیک گزینه دیگر انتخاب کن:")

			case "currency":
				value, err := strconv.ParseFloat(text, 64)
				if err != nil {
					send(bot, chatID, "عدد معتبر وارد کن")
					continue
				}

				rate, err := service.FetchDollarRate(apiURL, apiKey)
				if err != nil {
					send(bot, chatID, "خطا در دریافت نرخ دلار")
					continue
				}

				result, err := service.TomanToDollar(value, rate)
				if err != nil {
					send(bot, chatID, err.Error())
					continue
				}

				userState[chatID] = ""
				sendMainMenu(bot, chatID,
					"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" USD\n\nیک گزینه دیگر انتخاب کن:")
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
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("CM ➜ Foot", "cm_to_foot"),
						tgbotapi.NewInlineKeyboardButtonData("megabit ➜ megabyte", "mbit_to_mbyte"),
					),
				)

				msg := tgbotapi.NewMessage(chatID, "یک نوع تبدیل انتخاب کن:")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "km_to_mile":
				userState[chatID] = "km_to_mile"
				send(bot, chatID, "مقدار کیلومتر را وارد کن:")

			case "kg_to_pound":
				userState[chatID] = "kg_to_pound"
				send(bot, chatID, "مقدار کیلوگرم را وارد کن:")
			case "cm_to_foot":
				userState[chatID] = "cm_to_foot"
				send(bot, chatID, "مقدار سانتی‌متر را وارد کن:")

			case "mbit_to_mbyte":
				userState[chatID] = "mbit_to_mbyte"
				send(bot, chatID, "مقدار مگابیت را وارد کن:")

			case "currency":
				userState[chatID] = "currency"
				send(bot, chatID, "مقدار تومان را وارد کن:")
			}

			// پاسخ به Callback برای جلوگیری از حالت loading
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.Request(callback)
		}
	}
}
