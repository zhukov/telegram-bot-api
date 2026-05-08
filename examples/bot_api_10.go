package main

import (
	"log"

	api "github.com/OvyFlash/telegram-bot-api"
)

func polling_guest_bot() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60
	updateConfig.AllowedUpdates = []string{api.UpdateTypeGuestMessage}

	updatesChannel := bot.GetUpdatesChan(updateConfig)
	for update := range updatesChannel {
		if update.GuestMessage == nil || update.GuestMessage.GuestQueryID == "" {
			continue
		}

		result := api.NewInlineQueryResultArticle("guest-reply", "Reply", "Hello from guest mode")
		if _, err := bot.AnswerGuestQuery(api.NewAnswerGuestQuery(update.GuestMessage.GuestQueryID, result)); err != nil {
			log.Println(err)
		}
	}
}

func send_live_photo() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	msg := api.NewLivePhoto(123456, api.FilePath("live-photo.mp4"), api.FilePath("live-photo.jpg"))
	msg.Caption = "Live photo"

	if _, err := bot.SendLivePhoto(msg); err != nil {
		log.Println(err)
	}
}

func send_poll_with_media() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	optionMedia := api.NewInputMediaSticker(api.FileID("sticker_file_id"))
	poll := api.NewPoll(
		123456,
		"Where should we meet?",
		api.NewPollOptionWithMedia("At the cafe", &optionMedia),
		api.NewPollOption("At the office"),
	)
	location := api.NewInputMediaLocation(40.7128, -74.0060)
	poll.Media = &location
	poll.MembersOnly = true
	poll.CountryCodes = []string{"US"}

	if _, err := bot.Send(poll); err != nil {
		log.Println(err)
	}
}

func remove_message_reactions() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	reaction := api.NewDeleteMessageReaction(123456, 10)
	reaction.UserID = 777000
	if _, err := bot.DeleteMessageReaction(reaction); err != nil {
		log.Println(err)
	}

	clearRecent := api.NewDeleteAllMessageReactions(123456)
	clearRecent.UserID = 777000
	if _, err := bot.DeleteAllMessageReactions(clearRecent); err != nil {
		log.Println(err)
	}
}
