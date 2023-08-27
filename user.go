package skype

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type User struct {
	Payload *Payload
}

type Payload struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Mood      string `json:"mood"`
	RichMood  string `json:"richMood"`
}

/**
 * get user self profile.
 */
func (c *Conn)GetProfile(UserId string) () {
	path := fmt.Sprintf("%s/users/%s/profile", API_USER, UserId)
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":      c.LoginInfo.SkypeToken,
	}

	body, err, _ := req.request("GET", path, nil, nil, headers)
	if err != nil {
		fmt.Println("getProfile err: ", err)
	}
	fmt.Println("getProfile resp: ", body)
	conInfo := JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	fmt.Println(conInfo.Resource)
	return
}

/**
 * get user self profile.
 */
func (c *Conn)GetContactsProfile(ids []string) error {
	path := fmt.Sprintf("%s/users/self/contacts/profiles", API_USER)
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":      c.LoginInfo.SkypeToken,
	}

	data := map[string] []string{
		"contacts": {},
	}
	for _, id := range ids {
		data["contacts"] = append(data["contacts"], id)
	}
	params, _ := json.Marshal(data)
	fmt.Println()
	fmt.Println("GetContactsProfile: ", data)
	fmt.Println()
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("GetContactsProfile err: ", err)
		return err
	}
	fmt.Println("GetContactsProfile resp: ", body)
	conInfo := JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	fmt.Println(conInfo.Resource)
	return nil
}

type NameSearchRsp struct {
	RequestId string `json:"requestId"`
	Results [] struct{
		NodeProfileData struct{
			SkypeId string `json:"skypeId"`
			SkypeHandle string `json:"skypeHandle"`
			Name string `json:"name"`
			AvatarUrl string `json:"avatarUrl"`
			CountryCode string `json:"countryCode"`
			Gender string `json:"gender"`
			ContactType string `json:"contactType"`
		} `json:"nodeProfileData"`
	} `json:"results"`
}

// keyWord
// like "live:xxxxx" will be more precise
func (c *Conn)NameSearch(keyWord string) (*NameSearchRsp, error) {
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":      c.LoginInfo.SkypeToken,
	}
	params := url.Values{}
	params.Set("searchstring", keyWord)
	params.Set("requestId", "8:"+ keyWord)

	body, err := req.HttpGetWitHeaderAndCookiesJson(API_DIRECTORY, params, "", nil, headers)

	if err != nil {
		fmt.Println("NameSearch err: ", err)
		return nil, err
	}
	fmt.Println("NameSearch resp: ", body)
	nameSearchRsp := &NameSearchRsp{}
	json.Unmarshal([]byte(body), nameSearchRsp)
	fmt.Println(nameSearchRsp)
	for _, node := range nameSearchRsp.Results {
		personId := "8:" + node.NodeProfileData.SkypeId + "@s.skype.net"
		if _, ok := c.Store.Contacts[personId]; !ok {
			contact := Contact{}
			contact.DisplayName = node.NodeProfileData.Name
			contact.Profile.AvatarUrl = fmt.Sprintf("https://avatar.skype.com/v1/avatars/%s/public?returnDefaultImage=false", node.NodeProfileData.SkypeId)
			c.Store.Contacts[personId] = contact
			fmt.Println("NameSearch Result for2: ", c.Store.Contacts[personId])
		}
	}
	return nameSearchRsp, nil
}

/**
 Update username
 */
func (c *Conn)UpdateName(skypeToken string, firstName string, lastName string) () {
	path := "https://edge.skype.com/profile/v1/users/self/profile/partial"
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken":      skypeToken,
	}
	data := map[string]interface{}{
		"payload": Payload{
			FirstName: firstName,
			LastName:  lastName,
		},
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("getProfile err: ", err)
	}
	fmt.Println("getProfile resp: ", body)
	conInfo := JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	fmt.Println(conInfo.Resource)
	return
}

// Set or update profile mood.
//
// on successful update, it returns response string, otherwise error is returned
//
// Example:
//
// emoticon: (wfh)
//
// mood: Working from home
//
// Official emoticon list can be found here: https://support.skype.com/en/faq/FA12330/what-is-the-full-list-of-emoticons
func (c *Conn) SetMood(skypeToken, emoticon, mood string) (string, error) {
	emot := emoticon
	if len(emoticon) > 2 {
		emot = emoticon[1 : len(emoticon)-1]
	}

	path := "https://edge.skype.com/profile/v1/users/self/profile/partial"
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken": skypeToken,
	}
	data := map[string]interface{}{
		"payload": Payload{
			Mood:     fmt.Sprintf("%s %s", emoticon, mood),
			RichMood: fmt.Sprintf("<ss type=\"%s\">%s</ss> %s", emot, emoticon, mood),
		},
	}

	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		return "", err
	}

	return body, nil
}