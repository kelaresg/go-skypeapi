package skype

import (
	"encoding/json"
	"errors"
	"fmt"
)

const DEFAULT_USER string = "ME"
const DEFAULT_ENDPOINT string = "SELF"

type message interface {
	Get(url string)
}

type requestOptions struct {
	uri     string
	cookies interface{}
	headers struct {
		RegistrationToken string // registrationToken from https://client-s.gateway.messenger.live.com/v1/users/ME/endpoints
		Authentication    string // skypetoken from https://web.skype.com/login/microsoft rsp
		EndpointId        string // endpointId from https://client-s.gateway.messenger.live.com/v1/users/ME/endpoints
	}
}

type SendMessageBack struct {
	OriginalArrivalTime int `json:"OriginalArrivalTime"`
}

type MessageClient struct {
	SendMessageBack *SendMessageBack
	SendFileBack    *SendMessageBack
}


type SendMessage struct {
	Jid string //receiver id(conversation id)
	ClientMessageId string
	SkypeEditedId string
	Timestamp int64
	Type string
	*SendTextMessage
	*SendMediaMessage
}

type SendTextMessage struct {
	Content string
}

type SendMediaMessage struct {
	FileName string
	FileType string
	FileSize string
	RawData []byte
	Duration int
}

func (c *Conn) SendText(chatThreadId string, content *SendMessage) (err error) {
	surl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", c.LoginInfo.LocationHost, chatThreadId)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
	}
	data := map[string]interface{}{
		"contenttype":     "text",
		"messagetype":     "RichText",
		"content":         content.Content,
	}
	if len(content.SkypeEditedId) > 0 {
		data["skypeeditedid"] = content.SkypeEditedId
	} else {
		data["clientmessageid"] = content.ClientMessageId // A large integer (~20 digits)
	}
	params, _ := json.Marshal(data)
	body, err := req.HttpPostWitHeaderAndCookiesJson(surl, nil, string(params), nil, headers)
	//back := &SendMessageBack{}
	//json.Unmarshal([]byte(body), back)
	//m.SendMessageBack = back
	fmt.Println("SendMsg rsp body :", body)
	return
}

func (c *Conn) UploadFile(chatThreadId string, content *SendMessage) (fullUrl, fileId string, code int, err error) {
	meta := map[string]interface{}{
		"permissions": map[string]interface{}{
			chatThreadId: []string{"read"},
		},
		"filename": content.FileName,
	}
	objType := "imgpsh"
	if content.Type == "m.image" || content.Type == "" {
		meta["type"] = "pish/image"
	} else if content.Type == "m.audio" {
		meta["type"] = "sharing/audio"
		objType = "audio"
	} else if content.Type == "m.video" {
		meta["type"] = "sharing/audio"
		objType = "video"
	} else if content.Type == "avatar/group" {
		meta["type"] = "avatar/group"
		objType = "avatar"
	} else {
		meta["type"] = "sharing/file"
		objType = "original"
	}
	headers := map[string]string{
		"X-Client-Version": "0/0.0.0.0",
		"Authorization":    "skype_token " + c.LoginInfo.SkypeToken,
	}
	data, _ := json.Marshal(meta)
	req := Request{timeout: 30}
	fileIdRsp, _ := req.HttpPostWitHeaderAndCookiesJson("https://api.asm.skype.com/v1/objects", nil, string(data), nil, headers)
	fileIdType := struct {
		ID string `json:"id"`
	}{}
	json.Unmarshal([]byte(fileIdRsp), &fileIdType)
	fmt.Println(fileIdType.ID)
	fileId = fileIdType.ID
	if fileId == "" {
		return "", "", 0, errors.New("get upload file err at first step")
	}
	fullUrl = fmt.Sprintf("https://api.asm.skype.com/v1/objects/%s/content/%s", fileId, objType)
	fmt.Println("fullUrl:", fullUrl)
	headers["Content-Type"] = "application"
	headers["Content-Length"] = content.FileSize

	fmt.Println("message SendimageMsg headers: ", headers)
	fmt.Println()
	fmt.Println("message fileSize:", content.FileSize)
	fmt.Println()
	_, code, err = req.HttpPutWitHeaderAndCookiesJson(fullUrl, nil, string(content.SendMediaMessage.RawData), nil, headers)
	return
}
/**
 send file
`{permissions: {8:live:116xxxx691: ["read"]}, type: "pish/image", filename: "gh_e12cb68793e0_258.jpg"}
`permissions: {8:live:116xxxx691: ["read"]}
`type: "pish/image"
`filename: "gh_e12cb68793e0_258.jpg"
*/
func (c *Conn) SendFile(chatThreadId string, content *SendMessage) (err error) {
	req := Request{timeout: 30}
	_, fileId, code, err := c.UploadFile(chatThreadId, content)
	messageType := "RichText/UriObject"
	if code == 201 || code == 200 {
		fmt.Println("message SendimageMsg1: ")
		imageContent := MediaContentFormat(content.Type, content.FileName, content.FileSize, content.SendMediaMessage.Duration, fileId)
		if content.Type == "m.audio" {
			messageType = "RichText/Media_AudioMsg"
		}
		if content.Type == "m.file" {
			messageType = "RichText/Media_GenericFile"
		}
		if content.Type == "m.video" {
			messageType = "RichText/Media_Video"
		}

		requestBody := map[string]interface{}{
			"clientmessageid":    content.ClientMessageId,
			"content":       imageContent,
			"messagetype":   messageType,
			//"composetime":   time.Now().Format(time.RFC3339),
			"contenttype":   "text",
			//"imdisplayname": "Oliver1 Zhao2",
			//"receiverdisplayname": "test zhao",
			"amsreferences": []string{fileId},
		}
		headers1 := map[string]string{
			"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
			"RegistrationToken": c.LoginInfo.RegistrationTokenStr,
		}
		fmt.Println()
		fmt.Printf("requestBody: %+v", requestBody)
		fmt.Println()
		requestBody1, _ := json.Marshal(requestBody)
		messageUrl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", c.LoginInfo.LocationHost, chatThreadId)
		body, err := req.HttpPostWitHeaderAndCookiesJson(messageUrl, nil, string(requestBody1), nil, headers1)
		if err != nil {
			return err
		}
		back := &SendMessageBack{}
		json.Unmarshal([]byte(body), back)
		c.SendFileBack = back
	}
	return
}

func MediaContentFormat(fileType string, filename string, fileSize string, durationMs int, bodyFileOneId string) string {
	var imageContent string
	fullUrl := "https://api.asm.skype.com/v1/objects/" + bodyFileOneId
	if fileType == "m.image" || fileType == "" {
		viewLink_1 := fmt.Sprintf("https://api.asm.skype.com/s/i?%s", bodyFileOneId)
		viewLink := fmt.Sprintf(`<a href="%s">%s</a>`, viewLink_1, viewLink_1)
		values := map[string]string{
			"OriginalName": filename,
			"FileSize":     fileSize,
		}
		imageContent = UriObject(
			fmt.Sprintf(`%s<meta type="photo" originalName="%s"></meta>`, viewLink, filename),
			"Picture.1",
			bodyFileOneId,
			fullUrl,
			fmt.Sprintf("%s/views/imgt1", fullUrl),
			"",
			"",
			durationMs,
			values,
		)
	} else {
		viewLink_1 := fmt.Sprintf("https://login.skype.com/login/sso?go=webclient.xmm&amp;docid=%s", bodyFileOneId)
		viewLink := fmt.Sprintf(`<a href="%s">%s</a>`, viewLink_1, viewLink_1)
		ffileTypeStr := "File.1"
		values := map[string]string{
			"OriginalName": filename,
			"FileSize":     fileSize,
		}
		thumbnail := ""
		if fileType == "m.audio" {
			//ffileTypeStr = "Audio.1/Message.1" // if need send audio message like skype
			ffileTypeStr = "Audio.1"
			thumbnail = fullUrl + "/views/audio"
		} else if fileType == "m.file" {
			thumbnail = fullUrl + "/views/original"
		}  else if fileType == "m.video" {
			thumbnail = fullUrl + "/views/thumbnail"
		}

		imageContent = UriObject(
			viewLink,
			ffileTypeStr,
			bodyFileOneId,
			fullUrl,
			thumbnail,
			filename,
			filename,
			durationMs,
			values,
		)
	}
	return imageContent
}

