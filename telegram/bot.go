package telegram

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"ytplaylistbot/youtube"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot     *tgbotapi.BotAPI
	YouTube *youtube.YouTubeClient
}

func NewTelegramBot(token string, ytClient *youtube.YouTubeClient) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}
	return &TelegramBot{Bot: bot, YouTube: ytClient}
}

func (t *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				t.HandleStart(update.Message)
			case "playlist":
				t.HandlePlaylist(update.Message)
			}
		}

		// CallbackQueryni qayta ishlash
		if update.CallbackQuery != nil {
			t.HandleCallbackQuery(update.CallbackQuery) // Callback queryni qayta ishlash
		}
	}
}

func (t *TelegramBot) HandleStart(msg *tgbotapi.Message) {
	text := "Send a YouTube playlist URL using /playlist <url>."
	t.SendMessage(msg.Chat.ID, text)
}

func (t *TelegramBot) HandlePlaylist(msg *tgbotapi.Message) {
	playlistURL := msg.CommandArguments()
	playlistID := extractPlaylistID(playlistURL) // URL'dan IDni oling

	videos, err := t.YouTube.GetPlaylistVideos(playlistID)
	if err != nil {
		t.SendMessage(msg.Chat.ID, "Failed to fetch playlist videos.")
		return
	}

	// Inline buttonlarni yaratish
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, video := range videos {
		button := tgbotapi.InlineKeyboardButton{
			Text:         video.Snippet.Title,
			CallbackData: &video.Snippet.ResourceId.VideoId,
		}
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	message := tgbotapi.NewMessage(msg.Chat.ID, "Here are the videos from the playlist:")
	message.ReplyMarkup = replyMarkup
	t.Bot.Send(message)
}

func (t *TelegramBot) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	// Callback data orqali video IDni olish
	videoID := callback.Data
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	fmt.Println("kirdi")

	// Video yuklab olish
	err := t.DownloadAndSendVideo(callback.Message.Chat.ID, videoURL)
	if err != nil {
		t.SendMessage(callback.Message.Chat.ID, "Failed to download video.")
		return
	}

	// Callbackni yo'q qilish va javob qaytarish
	callbackResponse := tgbotapi.NewCallback(callback.ID, "Video is downloading...")
	t.Bot.Send(callbackResponse) // Send callback response
}

func (t *TelegramBot) DownloadAndSendVideo(chatID int64, videoURL string) error {
	// Temporary fayl nomini yaratish
	outputPath := filepath.Join("/tmp", fmt.Sprintf("%s.mp4", strings.Split(videoURL, "=")[1]))
	fmt.Println("here")
	fmt.Println(videoURL, outputPath)
	err := downloadVideo(videoURL, outputPath)
	if err != nil {
		return err
	}

	// Faylni yuborish
	file := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(outputPath))
	_, err = t.Bot.Send(file)
	return err
}

func downloadVideo(videoURL, outputPath string) error {
	// yt-dlp yordamida video yuklab olish
	cmd := exec.Command("yt-dlp", "-o", outputPath, videoURL)
	fmt.Println(outputPath, videoURL)
	return cmd.Run()
}

func (t *TelegramBot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	t.Bot.Send(msg)
}

func extractPlaylistID(url string) string {
	// YouTube playlist URL uchun regex pattern
	re := regexp.MustCompile(`(?:list=)([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1] // IDni qaytaradi
	}
	log.Println("Invalid YouTube playlist URL")
	return ""
}
