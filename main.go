package main

import "github.com/BehzadE/telbox/pkg/app"


func main() {
    bot := app.New()
    bot.ListenAndServe()
}
