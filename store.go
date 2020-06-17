package skype

type Store struct {
	Contacts map[string]ContactInfo
	Chats    map[string]ContactGroup
}

func newStore() *Store {
	return &Store{
		make(map[string]ContactInfo),
		make(map[string]ContactGroup),
	}
}

func (c *Conn) updateContacts(contacts interface{}) {
	ch, ok := contacts.([]interface{})
	if !ok {
		return
	}

	for _, contact := range ch {
		contactNode, ok := contact.(ContactInfo)
		if !ok {
			continue
		}

		c.Store.Contacts[contactNode.Id] = ContactInfo{
			contactNode.Id,
			contactNode.PersonId,
			contactNode.Type,
			contactNode.DisplayName,
			contactNode.Authorized,
			contactNode.Suggested,
			contactNode.Mood,
		}
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