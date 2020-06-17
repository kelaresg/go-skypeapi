package skype

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Member struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}
type Properties struct {
	HistoryDisclosed  string `json:"historydisclosed"` // true|false
	Topic  string `json:"topic"`
}
type Members struct {
	Members []Member `json:"members"`
	Properties  Properties `json:"properties"`
}

type ThreadProperties struct {
	Topic       string `json:"topic"`
	Lastjoinat  string `json:"lastjoinat"` // ? a timestamp ? example: "1421342788493"
	Version     string `json:"version"`    //? a timestamp ? example: "1464029299838"
	Members     string `json:"members"`
	Membercount string `json:"membercount"`
}

type Join struct {
	Blob     string
	Id       string
	JoinUrl  string
	ThreadId string
}

type JoinToConInfo struct {
	Action     string
	ChatBlob       string
	FlowId  string
	Id string
	Resource string
}

type LastMessage struct {
	Id                  string `json:"id"`                  // ?
	OriginalArrivalTime string `json:"originalarrivaltime"` // ?
	MessageType         string `json:"messagetype"`         // ?
	Version             string `json:"version"`             // ?
	ComposeTime         string `json:"composetime"`         // ?
	ClientMessageiId    string `json:"clientmessageid"`     // ?
	ConversationLink    string `json:"conversationLink"`    // ?
	Content             string `json:"content"`             // ?
	Type                string `json:"type"`                // ?
	ConversationId      string `json:"conversationid"`      // ?
	From                string `json:"from"`                // ?
}

type Conversation struct {
	// https://{host}/v1/threads/{19:threadId} or // https://{host}/v1/users/ME/contacts/{8:contactId}
	TargetLink       string            `json:"targetLink"`
	ResourceLink     string            `json:"resourceLink"`
	ResourceType     string            `json:"resourceType"`
	ThreadProperties ThreadProperties `json:"threadProperties"`
	Id               interface{}       `json:"id"`      //string 或者 int
	Type             string            `json:"type"`    // "Conversation" | string;
	Version          int64             `json:"version"` // a timestamp ? example: 1464030261015
	Properties       struct {
		ConversationStatusProperties string `json:"conversationstatusproperties"` // ?
		OneToOneThreadId             string `json:"onetoonethreadid"`             // ?
		LastImReceivedTime           string `json:"lastimreceivedtime"`           // ?
		ConsumptionHorizon           string `json:"consumptionhorizon"`           // ?
		ConversationStatus           string `json:"conversationstatus"`           // ?
		IsEmptyConversation          string `json:"isemptyconversation"`          // ?
	} `json:"properties"`
	LastMessage               LastMessage `json:"lastMessage"`
	Messages                  string       `json:"message"`
	LastUpdatedMessageId      int64        `json:"lastUpdatedMessageId"`
	LastUpdatedMessageVersion int64        `json:"lastUpdatedMessageVersion"`
	Resource                  Resource    `json:"resource"`
}

type Resource struct {
	ConversationLink      string `json:"conversationLink"`
	Type                  string `json:"type"`
	EventId               string `json:"eventId"`
	From                  string `json:"from"`
	ClientMessageId       string `json:"clientmessageid"`
	Version               interface{} `json:"version"` // string|number
	MessageType           string `json:"messagetype"`
	CounterPartyMessageId string `json:"counterpartymessageid"`
	ImDisplayName         string `json:"imdisplayname"`
	Content               string `json:"content"`
	ComposeTime           string `json:"composetime"`
	OriginContextId       string `json:"origincontextid"`
	OriginalArrivalTime   string `json:"originalarrivaltime"`
	AckRequired           string `json:"ackrequired"`
	IsVideoCall           string `json:"isVideoCall"` // "FALSE|TRUE"
	IsActive              bool   `json:"isactive"`
	Id                    string   `json:"id"`
	Jid                   string   `json:"jid"` // conversation id
	SendId                string   `json:"sendid"`
	Timestamp             int64   `json:"timestamp"`
}

func (Re *Resource)GetFromMe (ce *Conn) bool{
	if Re.ConversationLink != "" {
		ConversationLinkArr := strings.Split(Re.ConversationLink, "/conversations/")
		Re.Jid = ConversationLinkArr[1]
	}
	if Re.From != "" {
		FromArr := strings.Split(Re.From, "/contacts/")
		Re.SendId = FromArr[1]
	}
	if ce.UserProfile != nil && ce.UserProfile.Username != "" && ce.UserProfile.Username == Re.SendId {
		fmt.Println()
		fmt.Println("GetFromMe true: ", ce.UserProfile.Username, Re.SendId)
		fmt.Println()
		return true
	} else {
		fmt.Println()
		fmt.Println("GetFromMe false: ", ce.UserProfile, Re.SendId)
		fmt.Println()
	}
	return false
}

type ConversationsList struct {
	Conversations []Conversation `json:"conversations"`
	Metadata      Metadata       `json:"_metadata"`
}

type Metadata struct {
	TotalCount   int    `json:"totalCount"`
	ForwardLink  string `json:"forwardLink"`
	BackwardLink string `json:"backwardLink"`
	SyncState    string `json:"syncState"`
}

type ConversationsClient struct {
	ConversationsList *ConversationsList
}

/**
This returns an array of conversations that the current user has most recently interacted with
*/
func (c *Conn) GetConversations(apiHost string, skypeToken string, regToken string) (err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/users/ME/conversations", apiHost)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "cross-site",
	}
	for k, v := range headers {
		fmt.Println(k, ":", v)
	}

	params := url.Values{}
	params.Set("startTime", "0")
	params.Set("view", "msnp24Equivalent")
	params.Set("targetType", "Passport|Skype|Lync|Thread")
	//params.Set("pageSize", "30")
	body, err := req.HttpGetWitHeaderAndCookiesJson(path, params, "", nil, headers)
	data := &ConversationsList{}
	json.Unmarshal([]byte(body), data)
	c.ConversationsList = data
	return
}

/**
Retrieve details about a conversation.
*/
func (c *Conn) GetConversation(id string) (conversation *Conversation, err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/users/ME/conversations/%s", c.LoginInfo.LocationHost, id)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,// "skypetoken=" + skypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationtokensStr,
		"BehaviorOverride":  "redirectAs404",
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "cross-site",
	}
	for k, v := range headers {
		fmt.Println(k, ":", v)
	}

	params := url.Values{}
	params.Set("startTime", "0")
	params.Set("view", "msnp24Equivalent")
	params.Set("targetType", "Passport|Skype|Lync|Thread")
	//params.Set("pageSize", "30")
	fmt.Println(params)
	body, err := req.HttpGetWitHeaderAndCookiesJson(path, params, "", nil, headers)
	fmt.Println("conversation detail: ", body)
	data := &Conversation{}
	json.Unmarshal([]byte(body), data)
	return data, nil
}

/**
Fetch additional group-specific information, including the members and admins of the chat, topic, and join permissions.
*/
func (c *Conn) GetConversationThreads(apiHost string, skypeToken string, regToken string, id string) (err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/threads/%s", apiHost, id)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "cross-site",
	}
	for k, v := range headers {
		fmt.Println(k, ":", v)
	}

	params := url.Values{}
	params.Set("startTime", "0")
	params.Set("view", "msnp24Equivalent")
	params.Set("targetType", "Passport|Skype|Lync|Thread")
	fmt.Println(params)
	body, err := req.HttpGetWitHeaderAndCookiesJson(path, params, "", nil, headers)
	fmt.Println("conversation detail: ", body)
	return
}

/**
.Create a new group conversation.
*/
func (c *Conn) CreateConversationGroup(apiHost string, skypeToken string, regToken string, members Members) (err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/threads", apiHost)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
		"Location": "/",
	}

	data := members
	params, _ := json.Marshal(data)
	fmt.Println("params: ")
	fmt.Println(members)
	//return
	body, err, _ := req.request("post", path, strings.NewReader(string(params)), nil, headers)
	fmt.Println("CreateConversationGroup resp: ", body)
	return
}

/**
add a member to a group conversation.
*/
func (c *Conn) AddMember(apiHost string, skypeToken string, regToken string, members Members, conversationId string) (err error) {
	//API_MSGSHOST
	//https://client-s.gateway.messenger.live.com/v1/threads/4323A0b5463022fd0d43b4916cf5c6492c3412%40thread.skype/members
	path := fmt.Sprintf("%s/v1/threads/%s/members", apiHost, conversationId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
	}
	data := members
	params, _ := json.Marshal(data)
	fmt.Println("params: ")
	fmt.Println(members)
	//return
	body, err, _ := req.request("post", path, strings.NewReader(string(params)), nil, headers)
	fmt.Println("CreateConversationGroup resp: ", body)
	return
}

/**
 * Remove Member From Conversation
 */
func (c *Conn)RemoveMember(apiHost string, skypeToken string, regToken string, conversationId string, userId string)  {
	//DELETE Request URL: https://client-s.gateway.messenger.live.com/v1/threads/1434A0b436022fd0d84342916c3435c0432c3412%40thread.skype/members/8:live:.cid.db9****2b51cc
	path := fmt.Sprintf("%s/v1/threads/%s/members/%s", apiHost, conversationId, userId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
	}
	body, err, _ := req.request("delete", path, nil, nil, headers)
	if err != nil {
		fmt.Println("RemoveMember err: ", err)
	}
	fmt.Println("RemoveMember resp: ", body)
	return
}

/**
 * Retrieve the join URL for a group conversation, if it is currently public.
 */
func (c *Conn)GetConJoinUrl(apiHost string, skypeToken string, regToken string, conversationId string, userId string)  {
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
 * Retrieve the join URL for a group conversation, if it is currently public.
 */
func (c *Conn)JoinConByCode(skypeToken string, regToken string, joinUrl string) (err error, conInfo JoinToConInfo) {
	joinUrlArr := strings.Split(joinUrl, ".com/")
	//join url e.g https://join.skype.com/IYu****iqUIu
	path := fmt.Sprintf("%s/api/v2/conversation/", API_JOIN)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": regToken,
		"BehaviorOverride":  "redirectAs404",
	}
	data := map[string]string{
		"shortId": joinUrlArr[1],
		"type":   "wl",
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("post", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("join by code err: ", err)
	}
	fmt.Println("join by code resp: ", body)
	conInfo = JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	return
}
