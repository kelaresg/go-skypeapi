package skype

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const DEFAULT_USER string = "ME"
const DEFAULT_ENDPOINT string = "SELF"
const POLL_URL = "/v1/users/ME/endpoints/{92c6524e-7a60-454d-a555-06dbe51a419c}/subscriptions/0/poll?ackId=1039"
const MSGS_HOST = "https://azwcus1-client-s.gateway.messenger.live.com"

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

//POST /v1/users/ME/endpoints/%7B5643ecd4-ffff-ffff-6657-715da9e8b5b4%7D/subscriptions/0/poll?ackId=1056 HTTP/1.1
//Host: azwcus1-client-s.gateway.messenger.live.com
//Connection: keep-alive
//Content-Length: 0
//Accept: application/json
//Accept-Encoding: gzip, deflate, br
//Accept-Language: zh-CN
//Authentication: skypetoken=eyJhbGciOiJSUzI1NiIsImtpZCI6IjEwMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA2NDg1NDYsImV4cCI6MTU5MDczNDk0NSwic2t5cGVpZCI6ImxpdmU6MTE2Mzc2NTY5MSIsInNjcCI6OTU2LCJjc2kiOiIxNTkwNjQ4NTQ1IiwiY2lkIjoiNGYyNmJjMmIyODAwOWZkNiIsImFhdCI6MTU5MDEzMjQ5Mn0.c9X1tVdIIv0cow777RH4ClcS-wEXktjn5SMfHwRfkBVk2SMTOf3_BMq9Cz-EQRhrpMKRfLgZw8oHAvs3fdc3PK31AVCR3h-ZymR8k9m4-2tQ_1YPgRrVLywYxV4lnBsnZl2D9YlhXzTPpIxwjMy8iEr-yYHH7L3MIby7H7PTeOUa0ln0AI8a8no7VmqV8Z8c3dObRLN94lDSMd5isxGUNTNVpQZaip2_LJMnH_FksZtPEXpW8aJ3WNZO5Fm5pDmCNj-yxLepLATp5_uIFl4R3VSlVDj0Dx_0pF2daKJH2yo_USOwfLOdaAIqOV2dCYGg2Qwm7w_LZ0oX3f3tzyqsXA
//BehaviorOverride: redirectAs404
//ClientInfo: os=OSX; osVer=10.15; proc=x86; lcid=zh-CN; deviceType=1; country=CN; clientName=skype4life; clientVer=1432/8.60.0.76//skype4life; timezone=Asia/Shanghai
//EndpointId: {5643ecd4-ffff-ffff-6657-715da9e8b5b4}
//RegistrationToken: registrationToken=U2lnbmF0dXJlOjI6Mjg6QVFRQUFBQU0vTGhtcnpmTGVyTDgzTC9qOGRNRDtWZXJzaW9uOjY6MToxO0lzc3VlVGltZTo0OjE5OjUyNDg5NDg3MDg1Mjk1ODY5NjI7RXAuSWRUeXBlOjc6MTo4O0VwLklkOjI6MTU6bGl2ZToxMTYzNzY1NjkxO0VwLkVwaWQ6NTozNjo1NjQzZWNkNC1mZmZmLWZmZmYtNjY1Ny03MTVkYTllOGI1YjQ7RXAuTG9naW5UaW1lOjc6MTowO0VwLkF1dGhUaW1lOjQ6MTk6NTI0ODk0ODcwODUyOTU4Njk2MjtFcC5BdXRoVHlwZTo3OjI6MTU7RXAuRXhwVGltZTo0OjE5OjUyNDg5NDkzMzU4NzczODc5MDQ7VXNyLk5ldE1hc2s6MTE6MToyO1Vzci5YZnJDbnQ6NjoxOjA7VXNyLlJkcmN0RmxnOjI6MDo7VXNyLkV4cElkOjk6MTowO1Vzci5FeHBJZExhc3RMb2c6NDoxOjA7VXNlci5BdGhDdHh0OjI6NDAwOkNsTnJlWEJsVkc5clpXNFBiR2wyWlRveE1UWXpOelkxTmpreEFRTlZhV01VTVM4eEx6QXdNREVnTVRJNk1EQTZNREFnUVUwTVRtOTBVM0JsWTJsbWFXVmsxcDhBS0N1OEprOEFBQUFBQUFCQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUQyeHBkbVU2TVRFMk16YzJOVFk1TVFBQUFBQUFBQUFBQUFkT2IxTmpiM0psQUFBQUFBUUFBQUFBQUFBQUFBQUFBTmFmQUNncnZDWlBBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUJEMnhwZG1VNk1URTJNemMyTlRZNU1RQUFBQUFBNFY3UFhnY0FBQUFJU1dSbGJuUnBkSGtPU1dSbGJuUnBkSGxWY0dSaGRHVUlRMjl1ZEdGamRITU9RMjl1ZEdGamRITlZjR1JoZEdVSVEyOXRiV1Z5WTJVTlEyOXRiWFZ1YVdOaGRHbHZiaFZEYjIxdGRXNXBZMkYwYVc5dVVtVmhaRTl1YkhrQUFBPT07; expires=1590734945; endpointId={5643ecd4-ffff-ffff-6657-715da9e8b5b4}
//Sec-Fetch-Mode: cors
//Sec-Fetch-Site: cross-site
//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Skype/8.60.0.76 Chrome/78.0.3904.130 Electron/7.2.3 Safari/537.36
//X-ECS-Etag: "3ZwBSEJ4CybjtKLUlq2Fd8qg5GgFJfzNkc9nXF0T2qU="


//{
//"eventMessages": [{
//"id": 1057,
//"type": "EventMessage",
//"resourceType": "NewMessage",
//"time": "2020-05-28T14:16:40Z",
//"resourceLink": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/8:live:.cid.d3feb90dceeb51cc/messages/1590675400417",
//"resource": {
//"type": "Message",
//"messagetype": "Control/ClearTyping",
//"originalarrivaltime": "2020-05-28T14:16:40.409Z",
//"version": "1590675400417",
//"eventId": "0",
//"contenttype": "Application/Message",
//"origincontextid": "0",
//"isVideoCall": "False",
//"isactive": true,
//"from": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/contacts/8:live:.cid.d3feb90dceeb51cc",
//"id": "1590675400417",
//"conversationLink": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/8:live:.cid.d3feb90dceeb51cc",
//"counterpartymessageid": "1590675400417",
//"imdisplayname": "live:.cid.d3feb90dceeb51cc",
//"ackrequired": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/ALL/messages/1590675400417/ack",
//"composetime": "2020-05-28T14:16:40.409Z"
//}
//}]
//}

//
//{
//"eventMessages": [{
//"id": 1096,
//"type": "EventMessage",
//"resourceType": "NewMessage",
//"time": "2020-05-28T14:25:04Z",
//"resourceLink": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/8:live:.cid.d3feb90dceeb51cc/messages/1590675903909",
//"resource": {
//"conversationLink": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/8:live:.cid.d3feb90dceeb51cc",
//"type": "Message",
//"eventId": "0",
//"from": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/contacts/8:live:.cid.d3feb90dceeb51cc",
//"clientmessageid": "1478527645777549129",
//"version": "1590675903909",
//"messagetype": "RichText",
//"counterpartymessageid": "1590675903909",
//"imdisplayname": "live:.cid.d3feb90dceeb51cc",
//"content": "5444",
//"composetime": "2020-05-28T14:25:02.938Z",
//"origincontextid": "0",
//"originalarrivaltime": "2020-05-28T14:25:02.938Z",
//"ackrequired": "https://azwcus1-client-s.gateway.messenger.live.com/v1/users/ME/conversations/ALL/messages/1590675903909/ack",
//"contenttype": "text/plain; charset=UTF-8",
//"isVideoCall": "False",
//"isactive": true,
//"id": "1590675903909"
//}
//}]
//}
type SendMessageBack struct {
	OriginalArrivalTime int `json:"OriginalArrivalTime"`
}

type MessageClient struct {
	SendMessageBack *SendMessageBack
	SendFileBack    *SendMessageBack
}

func (m *MessageClient) SendMsg(locationHost, chatThreadId, content, skypeToken, resToken string) (err error) {
	//API_MSGSHOST chat thread identifier
	surl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", locationHost, chatThreadId)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + skypeToken,
		"RegistrationToken": resToken,
	}
	clientmessageid := time.Now().Unix() * 1000
	data := map[string]interface{}{
		"contenttype":     "text",
		"clientmessageid": strconv.Itoa(int(clientmessageid)),
		"messagetype":     "Text",
		"content":         content,
	}
	params, _ := json.Marshal(data)
	body, err := req.HttpPostWitHeaderAndCookiesJson(surl, nil, string(params), nil, headers)
	back := &SendMessageBack{}
	json.Unmarshal([]byte(body), back)
	m.SendMessageBack = back
	return
}

/**
 send file
`{permissions: {8:live:116xxxx691: ["read"]}, type: "pish/image", filename: "gh_e12cb68793e0_258.jpg"}
`permissions: {8:live:116xxxx691: ["read"]}
`type: "pish/image"
`filename: "gh_e12cb68793e0_258.jpg"
*/
func (m *MessageClient) SendFile(locationHost, chatThreadId, filename, skypeToken, resToken, fileType string, duration_ms int) (err error) {
	meta := map[string]interface{}{
		"permissions": map[string]interface{}{
			chatThreadId: []string{"read"},
		},
		"filename": filename,
	}
	objType := "imgpsh"
	if fileType == "image" || fileType == "" {
		meta["type"] = "pish/image"
	} else if fileType == "audio" {
		meta["type"] = "sharing/audio"
		objType = "audio"
	} else {
		meta["type"] = "sharing/file"
		objType = "original"
	}
	headers := map[string]string{
		"X-Client-Version": "0/0.0.0.0",
		"Authorization":    "skype_token " + skypeToken,
	}
	data, _ := json.Marshal(meta)
	req := Request{timeout: 30}
	bodyfile_one, _ := req.HttpPostWitHeaderAndCookiesJson("https://api.asm.skype.com/v1/objects", nil, string(data), nil, headers)
	bodyfile_one_d := struct {
		ID string `json:"id"`
	}{}
	json.Unmarshal([]byte(bodyfile_one), &bodyfile_one_d)
	fmt.Println(bodyfile_one_d.ID)
	if bodyfile_one_d.ID == "" {
		return errors.New("get upload file err at first step!")
	}
	fullUrl := fmt.Sprintf("https://api.asm.skype.com/v1/objects/%s/content/%s", bodyfile_one_d.ID, objType)
	fmt.Println("fullUrl:", fullUrl)
	//Processing file information , it will be deleted when it comes online
	filePath := filename
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fstat, _ := os.Stat(filePath)
	filesize := fstat.Size()

	headers["Content-Type"] = "application"
	headers["Content-Length"] = strconv.Itoa(len(string(f)))
	_, code, err := req.HttpPutWitHeaderAndCookiesJson(fullUrl, nil, string(f), nil, headers)
	messagetype := "RichText/UriObject"
	if code == 201 || code == 200 {
		imageContent := ""
		if fileType == "image" || fileType == "" {
			viewLink_1 := fmt.Sprintf("https://api.asm.skype.com/s/i?%s", bodyfile_one_d.ID)
			viewLink := fmt.Sprintf(`<a href="%s">%s</a>`, viewLink_1, viewLink_1)
			values := map[string]string{
				"OriginalName": filename,
				"FileSize":     strconv.Itoa(int(filesize)),
			}
			imageContent = UriObject(
				fmt.Sprintf(`%s<meta type="photo" originalName="%s"/>`, viewLink, filename),
				"Picture.1",
				fullUrl,
				fmt.Sprintf("%s/views/imgt1", fullUrl),
				"",
				"",
				duration_ms,
				values,
			)
		} else {
			viewLink_1 := fmt.Sprintf("https://login.skype.com/login/sso?go=webclient.xmm&docid=%s", bodyfile_one_d.ID)
			viewLink := fmt.Sprintf(`<a href="%s">%s</a>`, viewLink_1, viewLink_1)
			ffileTypeStr := "File.1"
			values := map[string]string{
				"OriginalName": filename,
				"FileSize":     strconv.Itoa(int(filesize)),
			}
			if fileType == "audio" {
				ffileTypeStr = "Audio.1/Message.1"
				messagetype = "RichText/Media_AudioMsg"
			}

			imageContent = UriObject(
				viewLink,
				ffileTypeStr,
				fullUrl,
				fmt.Sprintf("%s/views/thumbnail", fullUrl),
				filename,
				filename,
				duration_ms,
				values,
			)
		}
		clientmessageid := time.Now().Unix() * 1000
		requestBody := map[string]interface{}{
			"tmessageid":    strconv.Itoa(int(clientmessageid)),
			"content":       imageContent,
			"messagetype":   messagetype,
			"contenttype":   "text",
			"amsreferences": []string{bodyfile_one_d.ID},
		}
		headers1 := map[string]string{
			"Authentication":    "skypetoken=" + skypeToken,
			"RegistrationToken": resToken,
		}
		requestBody1, _ := json.Marshal(requestBody)
		surl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", locationHost, chatThreadId)
		body, err := req.HttpPostWitHeaderAndCookiesJson(surl, nil, string(requestBody1), nil, headers1)
		if err != nil {
			return err
		}
		back := &SendMessageBack{}
		json.Unmarshal([]byte(body), back)
		m.SendFileBack = back
	}
	return
}
