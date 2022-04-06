package app

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)


type App struct {
    TClient *telegram.Bot
    Router *Router
}

func New() *App {
    var app App
    var err error
	app.TClient, err = telegram.NewBotAPI(os.Getenv("TELEGRAM_BOT"))
    if err!=nil {
        log.Fatal(err)
    }

    app.Router = NewRouter()
    return &app

}
