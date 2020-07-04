package skype

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/encoding/gurl"
	"net/url"
	"strings"
)

type Presence string

const (
	PresenceOnline  Presence = "Online"
	PresenceOffline Presence = "Offline"
	PresenceIdle    Presence = "Idle"
	PresenceAway    Presence = "Away"
	PresenceHidden  Presence = "Hidden"
)

//type Location struct {
//	Country string `json:"country"` // almost certainly an enum...
//	City    string `json:"city"`
//}
//
//type Phone struct {
//	Number string `json:"number"` // pattern: /^+\getAuthorizationState+$/  (with country code)
//	Type   int64  `json:"type"`   // enum, seen: 2
//}

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
	Id          string              `json:"id"`        // username
	PersonId    string              `json:"person_id"` // [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	Type        string              `json:"type"`      // "skype" | "agent" | string; // enum ?
	DisplayName string              `json:"display_name"`
	Authorized  bool                `json:"authorized"` //? accepted contact request ?
	Suggested   bool                `json:"suggested"`  //?
	Mood        bool                `json:"mood"`       //?
	Blocked     bool                `json:"blocked"`
	AvatarUrl   string              `json:"avatar_url"` // Canonical form: https://api.skype.com/users/{id}/profile/avatar
	Locations   []Location          `json:"locations"`
	Phones      []Phone             `json:"phones"`
	Name        UserInfoProfileName `json:"name"`
	agent       interface{}
}

type Location struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Type    string `json:"type"` // home|
}

type Phone struct {
	Number string `json:"number"`
	Type   int64  `json:"type"`
}

type Agent struct {
	Trust string `json:"trust"`
	Type   int64  `json:"type"`
	Info   struct{
		Trusted bool `json:"trusted"`
		Type   int64  `json:"type"`
	}  `json:"info"`
}

type ContactGroup struct {
	id         string
	name       string
	isFavorite bool
}
type ContactsList struct {
	Contacts []Contact `json:"contacts"`
	Count    int        `json:"count"`
	Scope    string     `json:"scope"`
}

type Contact struct {
	PersonId            string                      `json:"person_id"` // 8:xxxxxx
	Mri                 string                      `json:"mri"`
	DisplayName         string                      `json:"display_name"`
	DisplayNameSource   string                      `json:"display_name_source"`
	Profile             UserInfoProfile             `json:"profile"`
	Authorized          bool                        `json:"authorized"`
	Blocked             bool                        `json:"blocked"`
	Explicit            bool                        `json:"explicit"`
	CreationTime        string                      `json:"creation_time"`
	Agent               Agent                       `json:"agent"`
	RelationshipHistory UserInfoRelationshipHistory `json:"relationship_history"`
}

type UserInfoProfile struct {
	AvatarUrl   string              `json:"avatar_url"`
	Birthday    string              `json:"birthday"`
	Gender      string              `json:"gender"`
	Locations   []Location          `json:"locations"`
	Phones      []Phone             `json:"phones"`
	Name        UserInfoProfileName `json:"name"`
	Language    string              `json:"language"` // en|
	SkypeHandle string              `json:"skype_handle"`
	Mood        string              `json:"mood"`
	About       string              `json:"about"`
}

//type UserInfoProfileLocations struct {
//	Type    string `json:"type"`
//	Country string `json:"country"`
//}

type UserInfoProfileName struct {
	First    string `json:"first"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

type UserInfoRelationshipHistory struct {
	Sources []UserInfoRelationshipHistorySources `json:"sources"`
}

type UserInfoRelationshipHistorySources struct {
	Type    string `json:"type"`
	Time    string `json:"time"`
	Subtype string `json:"subtype"`
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

func (c *Conn) ContactList(id string) (err error) {
	url := fmt.Sprintf("%s/users/%s/contacts", API_CONTACTS, id)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": c.LoginInfo.SkypeToken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	if err != nil {
		return err
	}
	//fmt.Println()
	//fmt.Println("contacts list", body)
	list := ContactsList{}
	json.Unmarshal([]byte(body), &list)
	c.updateContacts(list.Contacts)
	return
}

func (c *Conn) ContactGroupList(id string) (err error) {
	url := fmt.Sprintf("%s/users/%s/groups", API_CONTACTS, id)
	//fmt.Println(url)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": c.LoginInfo.SkypeToken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(url, nil, "", nil, headers)
	//fmt.Println(body)
	if err != nil {
		return err
	}
	list := GroupsList{}
	json.Unmarshal([]byte(body), &list)
	c.ContactClient.Groups = &list
	return
}

func (c *Conn) GetAllContactInfo(id, skypetoken string) (err error) {
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
	c.ContactClient.Groups = &list
	return
}

func (c *Conn) BlockList(id, skypetoken string) (err error) {
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
	c.ContactClient.Blocks = &list
	return
}


type Block struct {
	UiVersion     string `json:"ui_version"`
	ReportAbuse   bool `json:"report_abuse"`
	DeleteContact bool `json:"delete_contact"`
	ReportContext string `json:"report_context"`
}
/**
 * BlockContact
 * id: live:xxxxxxxxx
 * otherId: 8:live:xxxxxxxx
 */
func (c *Conn)BlockContact(skypeToken string, id string, otherId string, report bool, deleteContact bool) (err error, conInfo JoinToConInfo) {
	idEncode := gurl.Encode(id)
	otherIdEncode := gurl.Encode(otherId)
	path := fmt.Sprintf("%s/users/%s/contacts/blocklist/%s", API_CONTACTS, idEncode, otherIdEncode)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypeToken,
	}
	data := Block{
		"skype.com",
		report,
		deleteContact,
		"profile",
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("PUT", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("BlockContact err: ", err)
	}
	fmt.Println("BlockContact resp: ", body)
	return
}

/**
 * UnBlockContact
 * id: live:xxxxxxxxx
 * otherId: 8:live:xxxxxxxx
 */
func (c *Conn)UnBlockContact(skypeToken string, id string, otherId string) (err error, conInfo JoinToConInfo) {
	idEncode := gurl.Encode(id)
	otherIdEncode := gurl.Encode(otherId)
	path := fmt.Sprintf("%s/users/%s/contacts/blocklist/%s", API_CONTACTS, idEncode, otherIdEncode)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypeToken,
	}
	body, err, _ := req.request("DELETE", path, nil, nil, headers)
	if err != nil {
		fmt.Println("UnBlockContact err: ", err)
	}
	fmt.Println("UnBlockContact resp: ", body)
	return
}


/**
POST https://contacts.skype.com/contacts/v2/users/live%3A1163765691/contacts
request payload
{greeting: ""
mri: "8:live:.cid.xxxxxxxxxx"
send_invite: false}

id: live:xxxxxxxxx
otherId: 8:live:xxxxxxxxxxxxxx
 */
func (c *Conn)AddContact(skypeToken string, id string, otherId string)  {
	idEncode := gurl.Encode(id)
	path := fmt.Sprintf("%s/users/%s/contacts", API_CONTACTS, idEncode)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		//"Authentication":    "skypetoken=" + skypeToken,
		//"RegistrationToken": regToken,
		//"BehaviorOverride":  "redirectAs404",
		"X-Skypetoken":     skypeToken,
	}
	data := map[string]interface{}{
		"greeting": "",
		"mri":   otherId,
		"send_invite": false,
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("AddContact err: ", err)
	}
	fmt.Println("AddContact resp: ", body)
	return
}
/**
PUT https://azwus1-client-s.gateway.messenger.live.com/v1/users/ME/contacts/8:live:.cid.xxxxxxxxxxx
formdata:v1/users/ME/contacts/8:live:.cid.xxxxxxxxxxxxx
 * Add a user to the current user’s contact list. This has no effect on auth status, which must be approved by accepting an invite.
 * userId: 8:live:.cid.xxxxxxxxxxxx – user thread identifier of not-yet-contact
 */
func (c *Conn)AddContact2(apiHost string ,skypeToken string, regToken string, userId string)  {
	path := fmt.Sprintf("%s/v1/users/ME/contacts/%s", apiHost, userId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
	}
	body, err, _ := req.request("PUT", path, nil, nil, headers)
	if err != nil {
		fmt.Println("AddContact2 err: ", err)
	}
	fmt.Println("AddContact2 resp: ", body)
	return
}

/**
 * Add a user to the current user’s contact list. This has no effect on auth status, which must be approved by accepting an invite.
 */
func (c *Conn)RemoveUser(apiHost string, skypeToken string, regToken string, conversationId string, userId string)  {
	path := fmt.Sprintf("%s", API_JOIN_URL)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		//"Authentication":    "skypetoken=" + skypeToken,
		//"RegistrationToken": regToken,
		//"BehaviorOverride":  "redirectAs404",
		"X-Skypetoken":     skypeToken,
	}
	data := map[string]string{
		"baseDomain": "https://join.skype.com/launch/",
		"threadId":   conversationId,
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("get join url err: ", err)
	}
	fmt.Println("get join url resp: ", body)
	return
}

/**
 * Delete contact
 */
func (c *Conn)DeleteContact(skypeToken string, id string, otherId string)  {
	//DELETE https://contacts.skype.com/contacts/v2/users/live%3Axxxxxx/contacts/8%3Alive%3A.cid.xxxxxxxxxxxxx
	idEncode := gurl.Encode(id)
	otherIdEncode := gurl.Encode(otherId)
	path := fmt.Sprintf("%s/users/%s/contacts/%s", API_CONTACTS, idEncode, otherIdEncode)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":     skypeToken,
	}
	body, err, _ := req.request("DELETE", path, nil, nil, headers)
	if err != nil {
		fmt.Println("DeleteContact err: ", err)
	}
	fmt.Println("DeleteContact resp: ", body)
	return
}

/**
 * Delete contact
 */
func (c *Conn)ddDeleteUser(skypeToken string, conversationId string, id string, otherId string)  {
	//https://contacts.skype.com/contacts/v2/users/live%3A1163765691/contacts/8%3Alive%3A.cid.d3feb90dceeb51cc
	idEncode := gurl.Encode(id)
	otherIdEncode := gurl.Encode(otherId)
	path := fmt.Sprintf("%s/users/%s/contacts/%s", API_CONTACTS, idEncode, otherIdEncode)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":     skypeToken,
	}
	data := map[string]string{
		"baseDomain": "https://join.skype.com/launch/",
		"threadId":   conversationId,
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("get join url err: ", err)
	}
	fmt.Println("get join url resp: ", body)
	return
}

type ContactsStatusRes struct {
	Responses ContactStatus
}

type ContactStatus struct {
	Status  int8
	Contact string
	Payload struct {
		Id                       string   `json:"id"`
		Type                     string   `json:"type"`
		SelfLink                 string   `json:"selfLink"`
		Availability             string   `json:"availability"`
		Status                   string   `json:"status"`
		Capabilities             string   `json:"capabilities"`
		EndpointPresenceDocLinks []string `json:"endpointPresenceDocLinks"`
	}
}
/**
Get current status of group contacts
Request URL: https://azcus1-client-s.gateway.messenger.live.com/v1/users/ME/contacts/ALL/presenceDocs/messagingService?cMri=8:live:xxxxxx&cMri=8:xxxxxxx&cMri=8:live:xxxxxx
Request Method: GET
RegistrationToken:
Authentication:
Origin: https://web.skype.com
Referer: https://web.skype.com/
@params
ids []"8:*****"
 */
func (c *Conn)GetContactsCurrentStatus(ids []string) (content *ContactsStatusRes, err error) {
	if len(ids) < 1 {
		err = errors.New("ids is empty")
		return nil, err
	}
	path := fmt.Sprintf("%s/v1/users/ME/contacts/ALL/presenceDocs/messagingService", c.LoginInfo.LocationHost)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"BehaviorOverride":  "redirectAs404",
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
	}
	params := url.Values{}
	for _, id := range ids {
		params.Set("cMri", id)
	}

	body, err := req.HttpGetWitHeaderAndCookiesJson(path, params, "", nil, headers)
	if err != nil {
		fmt.Println("GetContactsCurrentStatus err: ", err)
	}
	content = &ContactsStatusRes{}
	err = json.Unmarshal([]byte(body), content)
	if err != nil {
		fmt.Println("GetContactsCurrentStatus Unmarshal err: ", err)
	}
	fmt.Println("GetContactsCurrentStatus resp: ", body)
	return content, nil
}