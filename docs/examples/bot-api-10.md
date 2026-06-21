# Bot API 10.0

## Guest Bots

Guest Mode lets a user invoke a bot in any chat without adding it as a member.
Request the `guest_message` update and answer with `answerGuestQuery`.

```go
updateConfig := tgbotapi.NewUpdate(0)
updateConfig.AllowedUpdates = []string{tgbotapi.UpdateTypeGuestMessage}

for update := range bot.GetUpdatesChan(updateConfig) {
	if update.GuestMessage == nil || update.GuestMessage.GuestQueryID == "" {
		continue
	}

	result := tgbotapi.NewInlineQueryResultArticle("guest-reply", "Reply", "Hello")
	_, err := bot.AnswerGuestQuery(tgbotapi.NewAnswerGuestQuery(update.GuestMessage.GuestQueryID, result))
	if err != nil {
		log.Println(err)
	}
}
```

`Message.GuestBotCallerUser` and `Message.GuestBotCallerChat` contain caller context when Telegram provides it.

## Live Photos

Live photos use two file values: the video part and the static photo.

```go
msg := tgbotapi.NewLivePhoto(chatID, tgbotapi.FilePath("live.mp4"), tgbotapi.FilePath("live.jpg"))
msg.Caption = "Live photo"
_, err := bot.SendLivePhoto(msg)
```

For paid media, wrap `InputMediaLivePhoto` with `NewInputPaidMediaLivePhoto`.

## Poll Media

Polls can include media on the poll itself, quiz explanations, and options.

```go
sticker := tgbotapi.NewInputMediaSticker(tgbotapi.FileID("sticker_file_id"))
poll := tgbotapi.NewPoll(chatID, "Choose one", tgbotapi.NewPollOptionWithMedia("A", &sticker), tgbotapi.NewPollOption("B"))

location := tgbotapi.NewInputMediaLocation(40.7128, -74.0060)
poll.Media = &location
poll.MembersOnly = true
poll.CountryCodes = []string{"US"}

_, err := bot.Send(poll)
```

## Reaction Moderation

Admin bots can remove a reaction from one message, or clear recent reactions from a user or actor chat.

```go
reaction := tgbotapi.NewDeleteMessageReaction(chatID, messageID)
reaction.UserID = userID
_, err := bot.DeleteMessageReaction(reaction)

clearRecent := tgbotapi.NewDeleteAllMessageReactions(chatID)
clearRecent.UserID = userID
_, err = bot.DeleteAllMessageReactions(clearRecent)
```
