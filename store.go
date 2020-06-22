package skype

type Store struct {
	Contacts map[string]Contact
	Chats    map[string]Conversation
}

func newStore() *Store {
	return &Store{
		make(map[string]Contact),
		make(map[string]Conversation),
	}
}

func (c *Conn) updateContacts(contacts []Contact) {
	for _, contact := range contacts {
		PersonId := contact.PersonId + "@s.skype.net"
		c.Store.Contacts[PersonId] = contact
	}
}

func (c *Conn) updateChats(chats []Conversation) {
	//ch, ok := chats.([]interface{})
	//if !ok {
	//	return
	//}

	for _, chat := range chats {
		//chatNode, ok := chat.(Conversation)
		//if !ok {
		//	continue
		//}
		c.Store.Chats[chat.Id] = chat
	}
}