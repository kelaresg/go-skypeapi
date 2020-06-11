package skype

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

//messagetype: RichText
type TextMessageHandler interface {
	Handler
	HandleTextMessage()
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



