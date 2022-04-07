package app

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)


type App struct {
    TClient *telegram.BotAPI
    Router *Router
}

func New() *App {
    var app App
    var err error
	app.TClient, err = telegram.NewBotAPI(os.Getenv("BOT_SECRET"))
    if err!=nil {
        log.Fatal(err)
    }

    app.Router = NewRouter(os.Getenv("BOT_NAME"))
    return &app
}
