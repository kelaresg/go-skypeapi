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

