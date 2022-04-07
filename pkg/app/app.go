package app

import (
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (a *App) ListenAndServe() {
	log.Printf(a.Router.botName + " started")
	u := telegram.NewUpdate(0)
	u.Timeout = 60
	updates, err := a.TClient.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message != nil {
			response := a.Router.RouteMessage(update.Message)
			var msg telegram.Chattable
			switch response.messageType {
			case TEXT:
				msg = telegram.NewMessage(int64(response.chatID), response.text)
			case STICKER:
				msg = telegram.NewStickerShare(int64(response.chatID), response.fileID)
			case DOCUMENT:
				msg = telegram.NewDocumentShare(int64(response.chatID), response.fileID)
			case AUDIO:
				msg = telegram.NewAudioShare(int64(response.chatID), response.fileID)
			case VIDEO:
				msg = telegram.NewVideoShare(int64(response.chatID), response.fileID)
			case ANIMATION:
				msg = telegram.NewAnimationShare(int64(response.chatID), response.fileID)
			case PHOTO:
				msg = telegram.NewPhotoShare(int64(response.chatID), response.fileID)
			}
			_, err := a.TClient.Send(msg)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
