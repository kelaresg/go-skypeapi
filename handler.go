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
	HandleJsonMessage(message Resource)
}

//messagetype: RichText
type TextMessageHandler interface {
	Handler
	HandleTextMessage(message Resource)
}

//messagetype: RichText/UriObject
type ImageMessageHandler interface {
	Handler
	HandleImageMessage(message Resource)
}

//messagetype: RichText/UriObject
type ContactMessageHandler interface {
	Handler
	HandleContactMessage(message Resource)
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
	HandlePresence(message Resource)
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
func (c *Conn) AddHandler(handler Handler) {
	c.handler = append(c.handler, handler)
}

// RemoveHandler removes a handler from the list of handlers that receive dispatched messages.
func (c *Conn) RemoveHandler(handler Handler) bool {
	i := -1
	for k, v := range c.handler {
		if v == handler {
			i = k
			break
		}
	}
	if i > -1 {
		c.handler = append(c.handler[:i], c.handler[i+1:]...)
		return true
	}
	return false
}

// RemoveHandlers empties the list of handlers that receive dispatched messages.
func (c *Conn) RemoveHandlers() {
	c.handler = make([]Handler, 0)
}

func (c *Conn) handle(message Conversation) {
	c.handleWithCustomHandlers(message, c.handler)
}

func (c *Conn) shouldCallSynchronously(handler Handler) bool {
	return false
}

type TextMessage struct {
	Resource
}

type ChatUpdateHandler interface {
	Handler
	HandleChatUpdate(message Resource)
}

func (c *Conn) handleWithCustomHandlers(message Conversation, handlers []Handler) {
	//switch m := message.(type) {
	if message.ResourceType == "NewMessage" {
		//resource := Resource{}
		//resource, ok := (message.Resource).(Resource)
		//if !ok {
		//	fmt.Println("handleWithCustomHandlers: not resource type")
		//	return
		//}
		//_ = json.Unmarshal([]byte(message.Resource), &resource)
		ConversationLinkArr := strings.Split(message.Resource.ConversationLink, "/conversations/")
		t, _ := time.Parse(time.RFC3339,message.Resource.ComposeTime)
		message.Resource.Jid = ConversationLinkArr[1]
		message.Resource.Timestamp = t.Unix()
		if message.Resource.MessageType == "RichText" || message.Resource.MessageType == "Text" {
			for _, h := range handlers {
				if x, ok := h.(TextMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleTextMessage(message.Resource)
					} else {
						go x.HandleTextMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/UriObject" {
			for _, h := range handlers {
				if x, ok := h.(ImageMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleImageMessage(message.Resource)
					} else {
						go x.HandleImageMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/Contacts" {
			for _, h := range handlers {
				if x, ok := h.(ContactMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleContactMessage(message.Resource)
					} else {
						go x.HandleContactMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/Media_GenericFile" {
			for _, h := range handlers {
				if x, ok := h.(ImageMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleImageMessage(message.Resource)
					} else {
						go x.HandleImageMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/Media_Album" {
			for _, h := range handlers {
				if x, ok := h.(ChatUpdateHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleChatUpdate(message.Resource)
					} else {
						go x.HandleChatUpdate(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/Media_Video" {
			for _, h := range handlers {
				if x, ok := h.(ImageMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleImageMessage(message.Resource)
					} else {
						go x.HandleImageMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "RichText/Media_AudioMsg" {
			for _, h := range handlers {
				if x, ok := h.(ImageMessageHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleImageMessage(message.Resource)
					} else {
						go x.HandleImageMessage(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "ThreadActivity/TopicUpdate" {
			for _, h := range handlers {
				if x, ok := h.(ChatUpdateHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleChatUpdate(message.Resource)
					} else {
						go x.HandleChatUpdate(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "ThreadActivity/PictureUpdate" {
			for _, h := range handlers {
				if x, ok := h.(ChatUpdateHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleChatUpdate(message.Resource)
					} else {
						go x.HandleChatUpdate(message.Resource)
					}
				}
			}
		} else if message.Resource.MessageType == "Control/Typing" {

		} else if message.Resource.MessageType == "Control/ClearTyping" {

		} else if message.Resource.MessageType == "ThreadActivity/AddMember" {
			for _, h := range handlers {
				if x, ok := h.(ChatUpdateHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandleChatUpdate(message.Resource)
					} else {
						go x.HandleChatUpdate(message.Resource)
					}
				}
			}
		} else {
			fmt.Println()
			fmt.Printf("unknown message type0: %+v", message)
			fmt.Println()
		}
	} else if message.ResourceType == "ThreadUpdate" {
		//resource := Resource{}
		//_ = json.Unmarshal([]byte(message.Resource), &resource)
		//resource, ok := (message.Resource).(Resource)
		//if !ok {
		//	return
		//}
		ConversationLinkArr := strings.Split(message.ResourceLink, "/threads/")
		t, _ := time.Parse(time.RFC3339, message.Time)
		message.Resource.Jid = ConversationLinkArr[1]
		message.Resource.Timestamp = t.Unix()
		fmt.Println()
		fmt.Println("ThreadUpdate")
		fmt.Println()
		//if message.Resource.MessageType == "ThreadActivity/TopicUpdate" {
		//	for _, h := range handlers {
		//		if x, ok := h.(ChatUpdateHandler); ok {
		//			if wac.shouldCallSynchronously(h) {
		//				x.HandleChatUpdate(message.Resource)
		//			} else {
		//				go x.HandleChatUpdate(message.Resource)
		//			}
		//		}
		//	}
		//} else {
		//	fmt.Println()
		//	fmt.Printf("unknown message type1: %+v", message)
		//	fmt.Println()
		//}
	} else if message.ResourceType == "ConversationUpdate" {

	} else if message.ResourceType == "UserPresence" {
		ResourceLinkArr := strings.Split(message.ResourceLink, "/contacts/")
		if len(ResourceLinkArr) < 2 {
			return
		}
 		LinkArr := strings.Split(ResourceLinkArr[1], "/presenceDocs/")
		t, _ := time.Parse(time.RFC3339, message.Time)
		if len(LinkArr) < 1 {
			return
		}
		message.Resource.SendId = LinkArr[0]
		message.Resource.Timestamp = t.Unix()
		if message.Resource.Type == "UserPresenceDoc" {
			for _, h := range handlers {
				if x, ok := h.(UserPresenceHandler); ok {
					if c.shouldCallSynchronously(h) {
						x.HandlePresence(message.Resource)
					} else {
						go x.HandlePresence(message.Resource)
					}
				}
			}
		}
	}else {
		fmt.Println()
		fmt.Printf("unknown message type2: %+v", message)
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

func (c *Conn) handleChats(chats interface{}) {

	//var chatList []Chat
	//c, ok := chats.([]interface{})
	//if !ok {
	//	return
	//}
	//for _, chat := range c {
	//	chatNode, ok := chat.(binary.Node)
	//	if !ok {
	//		continue
	//	}
	//
	//	jid := strings.Replace(chatNode.Attributes["jid"], "@c.us", "@s.whatsapp.net", 1)
	//	chatList = append(chatList, Chat{
	//		jid,
	//		chatNode.Attributes["name"],
	//		chatNode.Attributes["count"],
	//		chatNode.Attributes["t"],
	//		chatNode.Attributes["mute"],
	//		chatNode.Attributes["spam"],
	//	})
	//}
	//for _, h := range wac.handler {
	//	if x, ok := h.(ChatListHandler); ok {
	//		if wac.shouldCallSynchronously(h) {
	//			x.HandleChatList(chatList)
	//		} else {
	//			go x.HandleChatList(chatList)
	//		}
	//	}
	//}
}



