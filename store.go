package skype

type Store struct {
	Contacts map[string]Contact
	Chats    map[string]ContactGroup
}

func newStore() *Store {
	return &Store{
		make(map[string]Contact),
		make(map[string]ContactGroup),
	}
}

func (c *Conn) updateContacts(contacts []Contact) {
	for _, contact := range contacts {
		PersonId := contact.PersonId + "@s.skype.net"
		c.Store.Contacts[PersonId] = contact
	}
}

func (c *Conn) updateChats(chats interface{}) {
	ch, ok := chats.([]interface{})
	if !ok {
		return
	}

	for _, chat := range ch {
		chatNode, ok := chat.(ContactGroup)
		if !ok {
			continue
		}

		c.Store.Chats[chatNode.id] = ContactGroup{
			chatNode.id,
			chatNode.name,
			chatNode.isFavorite,
		}
	}
}