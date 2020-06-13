package skype

import (
	"fmt"
	"strings"
	"time"
)

type Message struct {
	Ackrequired         string `json:"ackrequired"`         // "https://client-s.gateway.messenger.live.com/v1/users/ME/conversations/ALL/messages/1451606400000/ack",
	Clientmessageid     string `json:"clientmessageid"`     // "1451606399999",
	Composetime         string `json:"composetime"`         // "2016-01-01T00:00:00.000Z",
	Content             string `json:"content"`             // "A message for the team.",
	Contenttype         string `json:"contenttype"`         // "text",
	ConversationLink    string `json:"conversationLink"`    // "https://client-s.gateway.messenger.live.com/v1/users/ME/conversations/19:a0b1c2...d3e4f5@thread.skype",
	From                string `json:"from"`                // "https://client-s.gateway.messenger.live.com/v1/users/ME/contacts/8:anna.7",
	Id                  string `json:"id"`                  // "1451606400000",
	Imdisplayname       string `json:"imdisplayname"`       // "Anna Cooper",
	Isactive            bool   `json:"isactive"`            // True,
	Messagetype         string `json:"messagetype"`         // "RichText",
	Originalarrivaltime string `json:"originalarrivaltime"` // "22016-01-01T00:00:00.000Z",
	Threadtopic         string `json:"threadtopic"`         // "Team chat",
	Type                string `json:"type"`                // "Message",
	Version             string `json:"version"`             // "1451606400000"
}

type Handler interface {
	HandleError(err error)
}

type JsonMessageHandler interface {
	Handler
	HandleJsonMessage(message string)
}

//messagetype: RichText
type TextMessageHandler interface {
	Handler
	HandleTextMessage(message Resource)
}

//messagetype:
type VideoMessageHandler interface {
	Handler
	HandleVideoMessage()
}

//messagetype:
type AudioMessageHandler interface {
	Handler
	HandleAudioMessage()
}

//A user connects to Skype with a new endpoint
type EndpointPresenceHandler interface {
	Handler
	HandleEndpointPresence()
}

//A user”s availability has changed
type UserPresenceHandler interface {
	Handler
	HandlePresence()
}

//A user”s availability has changed
type ConversationHandler interface {
	Handler
	HandleConversation()
}

/*
AddHandler adds an handler to the list of handler that receive dispatched messages.
The provided handler must at least implement the Handler interface. Additionally implemented
handlers(TextMessageHandler, ImageMessageHandler) are optional. At runtime it is checked if they are implemented
and they are called if so and needed.
*/
func (wac *Conn) AddHandler(handler Handler) {
	wac.handler = append(wac.handler, handler)
}

// RemoveHandler removes a handler from the list of handlers that receive dispatched messages.
func (wac *Conn) RemoveHandler(handler Handler) bool {
	i := -1
	for k, v := range wac.handler {
		if v == handler {
			i = k
			break
		}
	}
	if i > -1 {
		wac.handler = append(wac.handler[:i], wac.handler[i+1:]...)
		return true
	}
	return false
}

// RemoveHandlers empties the list of handlers that receive dispatched messages.
func (wac *Conn) RemoveHandlers() {
	wac.handler = make([]Handler, 0)
}

func (wac *Conn) handle(message Conversation) {
	wac.handleWithCustomHandlers(message, wac.handler)
}

func (wac *Conn) shouldCallSynchronously(handler Handler) bool {
	return false
}

type TextMessage struct {
	Resource
}

func (wac *Conn) handleWithCustomHandlers(message Conversation, handlers []Handler) {


	if message.Resource.MessageType == "RichText" {
		for _, h := range handlers {
			if x, ok := h.(TextMessageHandler); ok {
				ConversationLinkArr := strings.Split(message.Resource.ConversationLink, "/conversations/")
				t, _ := time.Parse(time.RFC3339,message.Resource.ComposeTime)
				message.Resource.Jid = ConversationLinkArr[1]
				message.Resource.Timestamp = t.Unix()
				if wac.shouldCallSynchronously(h) {
					x.HandleTextMessage(message.Resource)
				} else {
					go x.HandleTextMessage(message.Resource)
				}
			}
		}
	} else if message.Resource.MessageType == "Control/Typing" {

	} else if message.Resource.MessageType == "Control/ClearTyping" {

	} else {
		fmt.Println()
		fmt.Printf("unknown message type: %+v", message)
		fmt.Println()
	}

	//switch m := message.(type) {
	//case error:
	//	for _, h := range handlers {
	//		if wac.shouldCallSynchronously(h) {
	//			h.HandleError(m)
	//		} else {
	//			go h.HandleError(m)
	//		}
	//	}
	//case string:
	//	for _, h := range handlers {
	//		if x, ok := h.(JsonMessageHandler); ok {
	//			if wac.shouldCallSynchronously(h) {
	//				x.HandleJsonMessage(m)
	//			} else {
	//				go x.HandleJsonMessage(m)
	//			}
	//		}
	//	}
	//case TextMessage:
	//	for _, h := range handlers {
	//		if x, ok := h.(TextMessageHandler); ok {
	//			if wac.shouldCallSynchronously(h) {
	//				x.HandleTextMessage(m)
	//			} else {
	//				go x.HandleTextMessage(m)
	//			}
	//		}
	//	}
	//}
}



