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

type SendMessageBack struct {
	OriginalArrivalTime int `json:"OriginalArrivalTime"`
}

type MessageClient struct {
	SendMessageBack *SendMessageBack
	SendFileBack    *SendMessageBack
}

func (c *Conn) SendMsg(chatThreadId, content, clientMessageId string, output chan<- error) (err error) {
	//API_MSGSHOST chat thread identifier
	surl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", c.LoginInfo.LocationHost, chatThreadId)
	req := Request{timeout: 30}
	headers := map[string]string{
		"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
		"RegistrationToken": c.LoginInfo.RegistrationtokensStr,
	}
	//currentTimeNanoStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	//currentTimeNanoStr = currentTimeNanoStr[:len(currentTimeNanoStr)-3]
	//clientMessageId := currentTimeNanoStr + fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	data := map[string]interface{}{
		"contenttype":     "text",
		"clientmessageid": clientMessageId, // A large integer (~20 digits)
		//"composetime":     time.Now().Format(time.RFC3339),
		"messagetype":     "Text",
		"content":         content,
	}
	params, _ := json.Marshal(data)
	body, err := req.HttpPostWitHeaderAndCookiesJson(surl, nil, string(params), nil, headers)
	//back := &SendMessageBack{}
	//json.Unmarshal([]byte(body), back)
	//m.SendMessageBack = back
	fmt.Println("SendMsg rsp body :", body)
	if err != nil {
		output <- fmt.Errorf("message sending responded with %d", err)
	} else {
		output <- nil
	}
	return
}

/**
 send file
`{permissions: {8:live:116xxxx691: ["read"]}, type: "pish/image", filename: "gh_e12cb68793e0_258.jpg"}
`permissions: {8:live:116xxxx691: ["read"]}
`type: "pish/image"
`filename: "gh_e12cb68793e0_258.jpg"
*/
func (c *Conn) SendFile(chatThreadId, filename, fileType string, duration_ms int) (err error) {
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
		"Authorization":    "skype_token " + c.LoginInfo.SkypeToken,
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
			"Authentication":    "skypetoken=" + c.LoginInfo.SkypeToken,
			"RegistrationToken": c.LoginInfo.RegistrationtokensStr,
		}
		requestBody1, _ := json.Marshal(requestBody)
		surl := fmt.Sprintf("%s/v1/users/ME/conversations/%s/messages", c.LoginInfo.LocationHost, chatThreadId)
		body, err := req.HttpPostWitHeaderAndCookiesJson(surl, nil, string(requestBody1), nil, headers1)
		if err != nil {
			return err
		}
		back := &SendMessageBack{}
		json.Unmarshal([]byte(body), back)
		c.SendFileBack = back
	}
	return
}
