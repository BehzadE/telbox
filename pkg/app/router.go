package app

import (
	"sync"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Router struct {
	routes      map[string]handler
	userURLS    *BiMap[int, string]
	inputStates map[int]int
	botName     string
}

type BiMap[T1 comparable, T2 comparable] struct {
	ab   map[T1]T2
	ba   map[T2]T1
	lock *sync.RWMutex
}

func NewBiMap[T1 comparable, T2 comparable]() *BiMap[T1, T2] {
	ab := map[T1]T2{}
	ba := map[T2]T1{}
	lock := sync.RWMutex{}
	return &BiMap[T1, T2]{ab: ab, ba: ba, lock: &lock}
}

func (b *BiMap[T1, T2]) GetByKey(key T1) (T2, bool) {
	b.lock.RLock()
	result, ok := b.ab[key]
	b.lock.RUnlock()
	return result, ok
}

func (b *BiMap[T1, T2]) GetByValue(value T2) (T1, bool) {
	b.lock.RLock()
	result, ok := b.ba[value]
	b.lock.RUnlock()
	return result, ok
}

func (b *BiMap[T1, T2]) Insert(key T1, value T2) {
	b.lock.Lock()
	b.ab[key] = value
	b.ba[value] = key
	b.lock.Unlock()
}

func (b *BiMap[T1, T2]) Delete(key T1) {
	b.lock.Lock()
	value, ok := b.ab[key]
	if !ok {
		b.lock.Unlock()
		return
	}
	delete(b.ab, key)
	delete(b.ba, value)
	b.lock.Unlock()
}

type Response struct {
	text        string
	fileID      string
	chatID      int
	messageType MessageType
}
type MessageType uint

const (
	TEXT MessageType = iota
	STICKER
	PHOTO
	VIDEO
	AUDIO
	VOICE
	DOCUMENT
	ANIMATION
)

type handler func(request *telegram.Message, response *Response)

func NewRouter(botName string) *Router {
	var router Router
	router.userURLS = NewBiMap[int, string]()
	router.inputStates = map[int]int{}
	router.botName = botName
	{
		r := map[string]handler{}
		r["start"] = router.startHandler
		r["init"] = router.initHandler
		r["stopreceive"] = router.stopReceiveHandler
		r["stopsend"] = router.stopSendHandler
		router.routes = r
	}
	return &router
}

func (r *Router) RouteMessage(request *telegram.Message) *Response {
	command := request.Command()

	for path, handler := range r.routes {
		if path == command {
			response := Response{}
			handler(request, &response)
			return &response
		}
	}

	if target, ok := r.inputStates[request.From.ID]; ok {
		return r.HandleInput(request, target)
	}

	msg := Response{chatID: request.From.ID, text: "Use a valid command."}
	return &msg
}

func (r *Router) HandleInput(request *telegram.Message, target int) *Response {
	switch true {
	case request.Text != "":
		resp := Response{chatID: target, text: request.Text}
		return &resp
	case request.Sticker != nil:
		resp := Response{chatID: target, fileID: request.Sticker.FileID, messageType: STICKER}
		return &resp
	case request.Document != nil:
		resp := Response{chatID: target, fileID: request.Document.FileID, messageType: DOCUMENT}
		return &resp
	case request.Audio != nil:
		resp := Response{chatID: target, fileID: request.Audio.FileID, messageType: AUDIO}
		return &resp
	case request.Video != nil:
		resp := Response{chatID: target, fileID: request.Video.FileID, messageType: VIDEO}
		return &resp
	case request.Animation != nil:
		resp := Response{chatID: target, fileID: request.Animation.FileID, messageType: ANIMATION}
		return &resp
    case request.Photo != nil:
        resp := Response{chatID: target, fileID: (*request.Photo)[0].FileID, messageType: PHOTO}
        return &resp
	default:
		resp := Response{chatID: request.From.ID, text: "Unsupported file type"}
		return &resp
	}
}
