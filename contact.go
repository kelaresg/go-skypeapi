package skype

import (
	"encoding/json"
	"fmt"
)

type Location struct {
	Country string `json:"country"` // almost certainly an enum...
	City    string `json:"city"`
}

type Phone struct {
	Number string `json:"number"` // pattern: /^+\getAuthorizationState+$/  (with country code)
	Type   int64  `json:"type"`   // enum, seen: 2
}

type SearchContact struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Country     string `json:"country"`
	City        string `json:"city"`
	AvatarUrl   string `json:"avatarUrl"`
	Aisplayname string `json:"displayname"`
	Username    string `json:"username"`
	Mood        string `json:"mood"`
	Emails      []string
	Gender      string `json:"gender"` // its numeric it seems
}

type ContactInfo struct {
	Id          string `json:"id"`        // username
	PersonId    string `json:"person_id"` // [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	Type        string `json:"type"`      // "skype" | "agent" | string; // enum ?
	DisplayName string `json:"display_name"`
	Authorized  bool   `json:"authorized"` //? accepted contact request ?
	Suggested   bool   `json:"suggested"`  //?
	Mood        bool   `json:"mood"`       //?
	Blocked     bool   `json:"blocked"`
	AvatarUrl   string `json:"avatar_url"` // Canonical form: https://api.skype.com/users/{id}/profile/avatar
	Locations   []struct {
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
	} `json:"locations"`
	Phones []struct {
		Number string `json:"number"`
		Type   int64  `json:"type"`
	} `json:"phones"`
	Name struct {
		First    string `json:"first"`
		Surname  string `json:"surname"`  //? also last-name ?
		Nickname string `json:"nickname"` // username, it is NOT the local nickname that you can modify
	} `json:"name"`
	agent interface{}
}

type ContactGroup struct {
	id         string
	name       string
	isFavorite bool
}
type ContactsList struct {
	Contacts []UserInfo `json:"contacts"`
	Count    int        `json:"count"`
	Scope    string     `json:"scope"`
}

type UserInfo struct {
	PersonId            string                      `json:"person_id"`
	Mri                 string                      `json:"mri"`
	DisplayName         string                      `json:"display_name"`
	DisplayNameSource   string                      `json:"display_name_source"`
	Profile             UserInfoProfile             `json:"profile"`
	Authorized          string                      `json:"authorized"`
	Blocked             string                      `json:"blocked"`
	Explicit            string                      `json:"explicit"`
	CreationTime        string                      `json:"creation_time"`
	RelationshipHistory UserInfoRelationshipHistory `json:"relationship_history"`
}
type UserInfoProfile struct {
	Gender      string                     `json:"gender"`
	Locations   []UserInfoProfileLocations `json:"locations"`
	Name        UserInfoProfileName        `json:"name"`
	SkypeHandle string                     `json:"skype_handle"`
}
type UserInfoProfileLocations struct {
	Type    string `json:"type"`
	Country string `json:"country"`
}
type UserInfoProfileName struct {
	First   string `json:"first"`
	Surname string `json:"surname"`
}
type UserInfoRelationshipHistory struct {
	Sources []UserInfoRelationshipHistorySources `json:"sources"`
}
type UserInfoRelationshipHistorySources struct {
	Type string `json:"type"`
	Time string `json:"time"`
}
type GroupsList struct {
	Count  int         `json:"count"`
	Groups []GroupInfo `json:"groups"`
	Scope  string
}
type GroupInfo struct {
	Contacts   []string `json:"contacts"`
	Id         string   `json:"id"`
	IsFavorite bool     `json:"is_favorite"`
	Name       string   `json:"name"`
}

type Blocks struct {
	Blocklist []struct {
		Mri string `json:"mri"`
	} `json:"blocklist"`
	Scope string `json:"scope"`
	Count int    `json:"count"`
}

type ContactClient struct {
	Users  *ContactsList
	Groups *GroupsList
	Blocks *Blocks
}

func (c *ContactClient) ContactList(id, skypetoken string) (err error) {
	//fmt.Println("string:id", id)
	//fmt.Println(API_CONTACTS)
	url := fmt.Sprintf("%s/users/%s/contacts", API_CONTACTS, id)
	//fmt.Println(url)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypetoken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	if err != nil {
		return err
	}
	//fmt.Println("contacts list", body)
	list := ContactsList{}
	json.Unmarshal([]byte(body), &list)
	c.Users = &list
	return
}

func (c *ContactClient) ContactGroupList(id, skypetoken string) (err error) {
	url := fmt.Sprintf("%s/users/%s/groups", API_CONTACTS, id)
	//fmt.Println(url)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypetoken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	//fmt.Println(body)
	if err != nil {
		return err
	}
	list := GroupsList{}
	json.Unmarshal([]byte(body), &list)
	c.Groups = &list
	return
}

func (c *ContactClient) GetAllContactInfo(id, skypetoken string) (err error) {
	url := fmt.Sprintf("%s/users/%s", API_CONTACTS, id)
	//fmt.Println(url)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypetoken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	//fmt.Println(body)
	if err != nil {
		return err
	}
	list := GroupsList{}
	json.Unmarshal([]byte(body), &list)
	c.Groups = &list
	return
}

func (c *ContactClient) BlockList(id, skypetoken string) (err error) {
	url := fmt.Sprintf("%s/users/%s/blocklist", API_CONTACTS, id)
	fmt.Println(url)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypetoken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	//fmt.Println(body)
	if err != nil {
		return err
	}
	list := Blocks{}
	json.Unmarshal([]byte(body), &list)
	c.Blocks = &list
	return
}
