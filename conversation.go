package skype

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	//"github.com/pkg/errors"
	"net/url"
	"strconv"
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

type Conversation struct {
	TargetLink                string                 `json:"targetLink"`
	ResourceLink              string                 `json:"resourceLink"`
	ResourceType              string                 `json:"resourceType"`
	ThreadProperties          ThreadProperties       `json:"threadProperties"`
	Id                        interface{}                 `json:"id"`      //string | int?
	Type                      string                 `json:"type"`    // "Conversation" | string;
	Version                   int64                  `json:"version"` // a timestamp ? example: 1464030261015
	Properties                ConversationProperties `json:"properties"`
	LastMessage               LastMessage            `json:"lastMessage"`
	Messages                  string                 `json:"message"`
	LastUpdatedMessageId      int64                  `json:"lastUpdatedMessageId"`
	LastUpdatedMessageVersion int64                  `json:"lastUpdatedMessageVersion"`
	Resource                  Resource           `json:"resource"`
	Time                      string                 `json:"time"`
}

type ThreadProperties struct {
	Topic       string `json:"topic"`
	Lastjoinat  string `json:"lastjoinat"` // ? a timestamp ? example: "1421342788493"
	Version     string `json:"version"`    //? a timestamp ? example: "1464029299838"
	Members     string `json:"members"`
	Membercount string `json:"membercount"`
}

type ConversationProperties struct {
	ConversationStatusProperties string `json:"conversationstatusproperties"` // ?
	ConsumptionHorizonPublished  string `json:"consumptionhorizonpublished"`  // ?
	OneToOneThreadId             string `json:"onetoonethreadid"`             // ?
	LastImReceivedTime           string `json:"lastimreceivedtime"`           // ?
	ConsumptionHorizon           string `json:"consumptionhorizon"`           // ?
	ConversationStatus           string `json:"conversationstatus"`           // ?
	IsEmptyConversation          string `json:"isemptyconversation"`          // ?
	IsFollowed                   string `json:"isfollowed"`                   // ?
}

type LastMessage struct {
	Id                  string `json:"id"`                  // ?
	OriginContextId     string `json:"origincontextid"`     // ?
	OriginalArrivalTime string `json:"originalarrivaltime"` // ?
	MessageType         string `json:"messagetype"`         // ?
	Version             string `json:"version"`             // ?
	ComposeTime         string `json:"composetime"`         // ?
	ClientMessageId     string `json:"clientmessageid"`     // ?
	ConversationLink    string `json:"conversationLink"`    // ?
	Content             string `json:"content"`             // ?
	Type                string `json:"type"`                // ?
	ConversationId      string `json:"conversationid"`      // ?
	From                string `json:"from"`                // ?
}

type ChatTopicContent struct {
	XMLName  xml.Name `xml:"topicupdate"` //
	EventTime string `xml:"eventtime"` //
	Initiator string `xml:"initiator"`
	Value string `xml:"value"`
}

type ChatPictureContent struct {
	XMLName  xml.Name `xml:"pictureupdate"` //
	EventTime string `xml:"eventtime"` //
	Initiator string `xml:"initiator"`
	Value string `xml:"value"`
}

type Resource struct {
	ConversationLink      string      `json:"conversationLink"`
	Type                  string      `json:"type"`
	EventId               string      `json:"eventId"`
	From                  string      `json:"from"`
	ClientMessageId       string      `json:"clientmessageid"`
	Version               interface{} `json:"version"` // string|number
	MessageType           string      `json:"messagetype"`
	CounterPartyMessageId string      `json:"counterpartymessageid"`
	ImDisplayName         string      `json:"imdisplayname"`
	Content               string      `json:"content"`
	ComposeTime           string      `json:"composetime"`
	OriginContextId       string      `json:"origincontextid"`
	OriginalArrivalTime   string      `json:"originalarrivaltime"`
	AckRequired           string      `json:"ackrequired"`
	ContentType           string      `json:"contenttype"`
	IsVideoCall           string      `json:"isVideoCall"` // "FALSE|TRUE"
	IsActive              bool        `json:"isactive"`
	ThreadTopic           string      `json:"threadtopic"`
	ContentFormat         string      `json:"contentformat"`
	Id                    string      `json:"id"`
	Jid                   string      `json:"jid"`       // conversation id(custom filed)
	SendId                string      `json:"sendid"`    // send id id(custom filed)
	Timestamp             int64       `json:"timestamp"` // custom filed
	UserPresence
	EndpointPresence
	Amsreferences []string `json:"amsreferences"`
}

type UserPresence struct {
	Id                       string   `json:"id"`
	Type                     string   `json:"type"`
	SelfLink                 string   `json:"selfLink"`
	Availability             string `json:"availability"`
	Status                   Presence `json:"status"`
	Capabilities             string   `json:"capabilities"`
	LastSeenAt               string   `json:"lastSeenAt"`
	EndpointPresenceDocLinks []string `json:"endpointPresenceDocLinks"`
}

type EndpointPresence struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	SelfLink   string `json:"selfLink"`
	PublicInfo struct {
		Capabilities     string `json:"capabilities"`
		NodeInfo         string `json:"nodeInfo"`
		SkypeNameVersion string `json:"skypeNameVersion"`
		Typ              string `json:"typ"`
		Version          string `json:"version"`
	} `json:"publicInfo"`
	PrivateInfo struct {
		EpName string `json:"epname"`
	} `json:"privateInfo"`
}

type MediaMessageContent struct {
	XMLName      xml.Name `xml:"URIObject"` //
	Uri          string   `xml:"uri,attr"`
	DurationMs   string   `xml:"duration_ms,attr"`
	UrlThumbnail string   `xml:"url_thumbnail,attr"`
	Type         string   `xml:"type,attr"`
	DocId        string   `xml:"doc_id,attr"`
	Width        string   `xml:"width,attr"`
	Height       string   `xml:"height,attr"`
	A            struct {
		Href string `xml:"href,attr"`
	} `xml:"a"`
	OriginalName struct {
		V string `xml:"v,attr"`
	} `xml:"OriginalName"`
	FileSize struct {
		V string `xml:"v,attr"`
	} `xml:"FileSize"`
	Meta struct {
		Type         string `xml:"type,attr"`
		OriginalName string `xml:"originalName,attr"`
	} `xml:"meta"`
}

func (Re *Resource) Download(ce *Conn, mediaType string) (data []byte, mediaMessage *MediaMessageContent, err error) {
	mediaMessage = &MediaMessageContent{}
	err = xml.Unmarshal([]byte(Re.Content), mediaMessage)
	if err != nil {
		return nil, nil, err
	}

	var mediaUrl string
	switch mediaType {
	case "RichText/UriObject":
		mediaUrl = mediaMessage.Uri + "/views/imgpsh_mobile_save_anim"
	case "RichText/Media_GenericFile":
		mediaUrl = mediaMessage.Uri + "/views/original"
	case "RichText/Media_Video":
		mediaUrl = mediaMessage.Uri + "/views/video"
	case "RichText/Media_AudioMsg":
		mediaUrl = mediaMessage.Uri + "/views/audio"
	}

	fmt.Println("content.A.Href", mediaMessage.Uri)
	fmt.Println("content.FileSize", mediaMessage.FileSize.V)
	fileLength, err := strconv.Atoi(mediaMessage.FileSize.V)
	if err != nil {
		return nil, nil, err
	}
	data, err = Download(mediaUrl, ce, fileLength)
	return
}

type ContactMessageContent struct {
	XMLName xml.Name `xml:"contacts"` //
	C       struct {
		T string `xml:"t,attr"`
		S string `xml:"s,attr"` // live:xxxxx
		F string `xml:"f,attr"` // username
	} `xml:"c"`
}

func (Re *Resource) ParseContact() (contactMessage *ContactMessageContent, err error) {
	contactMessage = &ContactMessageContent{}
	err = xml.Unmarshal([]byte(Re.Content), contactMessage)
	if err != nil {
		return nil, err
	}
	return contactMessage, nil
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
func (c *Conn) GetConversations() (err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/users/ME/conversations", c.LoginInfo.LocationHost)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
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
	c.updateChats(data.Conversations)
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
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken, // "skypetoken=" + skypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
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


type ConsumptionHorizonsRsp struct {
	Id                  string               `json:"id"`
	Version             string               `json:"version"`
	ConsumptionHorizons []ConsumptionHorizon `json:"consumptionhorizons"`
}
type ConsumptionHorizon struct {
	ConsumptionHorizon string `json:"consumptionhorizon"`
	Id                 string `json:"id"`
}
/**
Fetch all members in conversation
@params
conId : conversation id
*/
func (c *Conn) GetConsumptionHorizons(conId string) (content *ConsumptionHorizonsRsp, err error) {
	path := fmt.Sprintf("%s/v1/threads/%s/consumptionhorizons", c.LoginInfo.LocationHost, conId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
		"BehaviorOverride":  "redirectAs404",
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "cross-site",
	}
	for k, v := range headers {
		fmt.Println(k, ":", v)
	}

	body, err := req.HttpGetWitHeaderAndCookiesJson(path, nil, "", nil, headers)
	content = &ConsumptionHorizonsRsp{}
	err = json.Unmarshal([]byte(body), content)
	fmt.Println("GetConsumptionHorizons detail: ", body)
	return
}

/**
.Create a new group conversation.
*/
func (c *Conn) CreateConversationGroup(members Members) (err error) {
	//API_MSGSHOST
	path := fmt.Sprintf("%s/v1/threads", c.LoginInfo.LocationHost)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
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
func (c *Conn) AddMember(members Members, conversationId string) (err error) {
	//API_MSGSHOST
	//https://client-s.gateway.messenger.live.com/v1/threads/4323A0b5463022fd0d43b4916cf5c6492c3412%40thread.skype/members
	path := fmt.Sprintf("%s/v1/threads/%s/members", c.LoginInfo.LocationHost, conversationId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
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
func (c *Conn)RemoveMember(conversationId string, userId string) (err error) {
	//DELETE Request URL: https://client-s.gateway.messenger.live.com/v1/threads/1434A0b436022fd0d84342916c3435c0432c3412%40thread.skype/members/8:live:.cid.db9****2b51cc
	path := fmt.Sprintf("%s/v1/threads/%s/members/%s", c.LoginInfo.LocationHost, conversationId, userId)
	fmt.Println(path)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
		"BehaviorOverride":  "redirectAs404",
	}
	body, err, _ := req.request("delete", path, nil, nil, headers)
	if err != nil {
		fmt.Println("RemoveMember err: ", err)
		return err
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


