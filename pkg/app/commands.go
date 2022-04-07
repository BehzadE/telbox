package app

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

func generateUserURL(botName, userUUID string) string {
	return "https://t.me/" + botName + "?start=" + userUUID
}

func (r *Router) initHandler(request *telegram.Message, response *Response) {
	if userUUID, ok := r.userURLS.GetByKey(request.From.ID); ok {
        response.chatID = request.From.ID
		response.text = generateUserURL(r.botName, userUUID)
		return
	}
	userUUID := uuid.New()
	r.userURLS.Insert(request.From.ID, userUUID.String())
    response.chatID = request.From.ID
	response.text = generateUserURL(r.botName, userUUID.String())
}

func (r *Router) startHandler(request *telegram.Message, response *Response) {
    args := request.CommandArguments()
    if args == "" {
        response.chatID = request.From.ID
        response.text = "use the /init command to receive your own url"
        return
    }
	userID, ok := r.userURLS.GetByValue(request.CommandArguments())
	if !ok {
        response.chatID = request.From.ID
		response.text = "user not found"
		return
	}
    response.chatID = request.From.ID
	response.text = "enter your messages or files. use the /stopsend when done."
	r.inputStates[request.From.ID] = userID
}

func (r *Router) stopReceiveHandler(request *telegram.Message, response *Response) {
    r.userURLS.Delete(request.From.ID)
    response.chatID = request.From.ID
    response.text = "done"
}

func (r *Router) stopSendHandler(request *telegram.Message, response *Response) {
    delete(r.inputStates, request.From.ID)
    response.chatID = request.From.ID
    response.text = "stopped sending messages"
}
