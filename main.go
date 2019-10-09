package main

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) == 1{
		log.Fatal("Error: Please pass the bot token as the first argument to the bot.")
	}
	//Load bot
	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	if err != nil {
		panic("Cannot initialize the bot: " + err.Error())
	}
	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Text == ""{ // ignore any non-Message or non string Updates
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
			case "help":
				msg.Text = "Welcome to Paste Ubuntu bot! Here you can paste your text to create a paste on paste.ubuntu.com\nTo get started, you can send the bot the text. The bot will send you the share link after.\nThe expiry date is never, the language is text and the poster name is your telegram first name and last name (not your ID)\n/about"
			case "about":
				msg.Text = "Created by Hirbod Behnam\nhttps://github.com/HirbodBehnam/"
			default:
				msg.Text = "Command not recognized! Try /help"
			}
			_ , _= bot.Send(msg)
			continue
		}
		go func(m tgbotapi.Message) {
			msg := tgbotapi.NewMessage(m.Chat.ID, "")
			params := url.Values{}
			params.Add("poster", m.From.FirstName + " " + m.From.LastName)
			params.Add("syntax", "text")
			params.Add("content", m.Text)

			resp, err := http.Post("https://paste.ubuntu.com", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(params.Encode())))
			if err != nil {
				msg.Text = err.Error()
			}else {
				msg.Text = resp.Request.URL.Host + resp.Request.URL.Path
			}
			_ , _= bot.Send(msg)
		}(*update.Message)
	}
}