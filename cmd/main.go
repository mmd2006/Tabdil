package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"Tabdil/internal/service"
)

var userState = make(map[int64]string)

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64, text string) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("محاسبه طول"),
			tgbotapi.NewKeyboardButton("محاسبه ارز"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func main() {

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

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// start
		if text == "/start" {
			userState[chatID] = ""
			sendMainMenu(bot, chatID, "سلام 👋\nبه بات تبدیل خوش آمدید\nیک گزینه انتخاب کن:")
			continue
		}

		// دکمه محاسبه طول
		if text == "محاسبه طول" {
			lengthKeyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("KM ➜ Mile"),
					tgbotapi.NewKeyboardButton("KG ➜ Pound"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("CM ➜ Foot"),
					tgbotapi.NewKeyboardButton("megabit ➜ megabyte"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("بازگشت به منوی اصلی"),
				),
			)

			msg := tgbotapi.NewMessage(chatID, "یک نوع تبدیل انتخاب کن:")
			msg.ReplyMarkup = lengthKeyboard
			bot.Send(msg)
			continue
		}

		// بازگشت
		if text == "بازگشت به منوی اصلی" {
			userState[chatID] = ""
			sendMainMenu(bot, chatID, "منوی اصلی:")
			continue
		}

		// دکمه محاسبه ارز
		if text == "محاسبه ارز" {

			rate, err := service.FetchDollarRate()
			if err != nil {
				send(bot, chatID, "خطا در دریافت نرخ دلار 😕")
				continue
			}

			dollarInToman := rate / 10
			userState[chatID] = "currency"

			send(bot, chatID,
				"💵 قیمت روز دلار: "+
					strconv.FormatFloat(dollarInToman, 'f', 0, 64)+
					" تومان\n\nمقدار تومان را وارد کنید:")

			continue
		}

		// انتخاب نوع تبدیل طول
		switch text {
		case "KM ➜ Mile":
			userState[chatID] = "km_to_mile"
			send(bot, chatID, "مقدار کیلومتر را وارد کن:")
			continue

		case "KG ➜ Pound":
			userState[chatID] = "kg_to_pound"
			send(bot, chatID, "مقدار کیلوگرم را وارد کن:")
			continue

		case "CM ➜ Foot":
			userState[chatID] = "cm_to_foot"
			send(bot, chatID, "مقدار سانتی‌متر را وارد کن:")
			continue

		case "megabit ➜ megabyte":
			userState[chatID] = "mbit_to_mbyte"
			send(bot, chatID, "مقدار مگابیت را وارد کن:")
			continue
		}

		// اگر کاربر عدد وارد کرده
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

			userState[chatID] = ""
			sendMainMenu(bot, chatID,
				"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" Mile\n\nیک گزینه دیگر انتخاب کن:")

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

			userState[chatID] = ""
			sendMainMenu(bot, chatID,
				"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" Pound\n\nیک گزینه دیگر انتخاب کن:")

		case "cm_to_foot":
			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				send(bot, chatID, "عدد معتبر وارد کن")
				continue
			}

			result, err := service.CentimeterToFoot(value)
			if err != nil {
				send(bot, chatID, err.Error())
				continue
			}

			userState[chatID] = ""
			sendMainMenu(bot, chatID,
				"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" Foot\n\nیک گزینه دیگر انتخاب کن:")

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

			// تبدیل تومان به ریال
			value = value * 10

			result, err := service.ConvertTomanToDollarWithAPI(value)
			if err != nil {
				send(bot, chatID, "خطا در دریافت نرخ دلار: "+err.Error())
				continue
			}

			userState[chatID] = ""
			sendMainMenu(bot, chatID,
				"نتیجه: "+strconv.FormatFloat(result, 'f', 4, 64)+" USD\n\nیک گزینه دیگر انتخاب کن:")
		}
	}
}
