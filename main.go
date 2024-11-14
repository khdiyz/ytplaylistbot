package main

import (
	"log"
	"ytplaylistbot/telegram"
	"ytplaylistbot/youtube"
)

const (
	telegramToken = "6099513704:AAFD9QKrpwv7KTfVExdnGEvzv_wqc3Nb_lo"
	youtubeAPIKey = "AIzaSyCIzRYNuvvDuE40kWe012ljq2gPkVKgBHs"
)

func main() {
	ytClient := youtube.NewYouTubeClient(youtubeAPIKey)
	tgBot := telegram.NewTelegramBot(telegramToken, ytClient)

	log.Println("Bot is running...")
	tgBot.Start()
}
