package skype

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/encoding/gurl"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	loggedIn    bool //has logged in or not
	session     *Session
	Store       *Store
	Handlers    []Handler
	LoginInfo   *LoginInfo
	UserProfile *UserProfile
}

type LoginInfo struct {
	SkypeToken            string
	SkypeExpires          string
	RegistrationToken     string
	RegistrationExpires   string
	LocationHost          string
	EndpointId            string
	RegistrationtokensStr string
}

/**
{"about":null,"avatarUrl":null,"birthday":null,"city":null,"country":null,"emails":["zhaosl@shinetechchina.com"],
"firstname":"lyle","gender":"0","homepage":null,"jobtitle":null,"language":null,"lastname":"zhao","mood":null,
"phoneHome":null,"phoneMobile":null,"phoneOffice":null,"province":null,"richMood":null,"username":"live:zhaosl_4"}
*/
type UserProfile struct {
	About       string   `json:"about"`
	AvatarUrl   string   `json:"avatarUrl"`
	Birthday    string   `json:"birthday"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Emails      []string `json:"emails"`
	Firstname   string   `json:"firstname"`
	Gender      string   `json:"gender"`
	Homepage    string   `json:"homepage"`
	Jobtitle    string   `json:"jobtitle"`
	Language    string   `json:"language"`
	Lastname    string   `json:"lastname"`
	Mood        string   `json:"mood"`
	PhoneHome   string   `json:"phone_home"`
	PhoneOffice string   `json:"phone_office"`
	Province    string   `json:"province"`
	RichMood    string   `json:"rich_mood"`
	Username    string   `json:"username"`
}

func NewClient() (cli *Client, err error) {
	c := &Client{
		loggedIn: false,
		session:  nil,
	}
	return c, nil
}

/**
login Skype by web auth
*/
func (c *Client) Login(usernanme, password string) (err error) {
	MSPRequ, MSPOK, PPFT, err := c.getParams()

	if err != nil {
		return errors.New("params get error")
	}
	//发送密码登录等
	_, err, t_value := c.sendCreds(usernanme, password, MSPRequ, MSPOK, PPFT)
	if err != nil {
		return errors.New("sendCreds get error")
	}
	//最后以部获得   token RegistrationExpires
	err = c.getToken(t_value)
	if err != nil {
		return errors.New("token get error")
	}
	//获得用户SkypeRegistrationTokenProvider
	c.LoginInfo.LocationHost = API_MSGSHOST
	err = c.SkypeRegistrationTokenProvider(c.LoginInfo.SkypeToken)
	if err != nil {
		return errors.New("SkypeRegistrationTokenProvider get error")
	}
	//请求获得用户的id （类型  string）
	err = c.GetUserId(c.LoginInfo.SkypeToken)
	if err != nil {
		return errors.New("GetUserId get error")
	}
	return
}

/**
获得用户的id
*/
func (c *Client) GetUserId(skypetoken string) (err error) {
	//params := url.Values{}
	//params.Set("auth", skypetoken)
	req := Request{
		timeout: 30,
	}
	headers := map[string]string{
		"x-skypetoken": skypetoken,
	}
	body, err := req.HttpGetWitHeaderAndCookiesJson(fmt.Sprintf("%s/users/self/profile", API_USER), nil, "", nil, headers)
	//解析参数
	if err != nil {
		return errors.New("get userId err")
	}

	userProfile := UserProfile{}
	json.Unmarshal([]byte(body), &userProfile)
	c.UserProfile = &userProfile
	return
}

/**
    Request a new registration token using a current Skype token.
	Args:
		skypeToken (str): existing Skype token
	Returns:
		(str, datetime.datetime, str, SkypeEndpoint) tuple: registration token, associated expiry if known,
															resulting endpoint hostname, endpoint if provided
	Raises:
		.SkypeAuthException: if the login request is rejected
		.SkypeApiExce`ption: if the login form can't be processed
 * Value used for the `ClientInfo` header of the request for the registration token.
*/
func (c *Client) SkypeRegistrationTokenProvider(skypetoken string) (err error) {
	secs := strconv.Itoa(int(time.Now().Unix()))
	lockAndKeyResponse := getMac256Hash(secs)
	LockAndKey := "appId=" + SKYPEWEB_LOCKANDKEY_APPID + "; time=" + secs + "; lockAndKeyResponse=" + lockAndKeyResponse
	req := Request{
		timeout: 30,
	}
	header := map[string]string{
		"Authentication":   "skypetoken=" + skypetoken,
		"LockAndKey":       LockAndKey,
		"BehaviorOverride": "redirectAs404",
	}
	data := map[string]interface{}{
		"endpointFeatures": "Agent",
	}
	params, _ := json.Marshal(data)
	//fmt.Println("https://client-s.gateway.messenger.live.com/v1/users/" + DEFAULT_USER + "/endpoints")
	fmt.Println(c.LoginInfo.LocationHost + "/v1/users/" + DEFAULT_USER + "/endpoints")
	registrationTokenStr, location, err := req.HttpPostRegistrationToken(c.LoginInfo.LocationHost+"/v1/users/"+DEFAULT_USER+"/endpoints", string(params), header)
	println("registrationTokenStr: ", registrationTokenStr)
	println("location: ", location)
	locationArr := strings.Split(location, "/v1")
	c.storeInfo(registrationTokenStr, c.LoginInfo.LocationHost)
	if locationArr[0] != c.LoginInfo.LocationHost {
		newRegistrationToken, _, err := req.HttpPostRegistrationToken(location, string(params), header)
		if err == nil {
			c.storeInfo(newRegistrationToken, locationArr[0])
		}
	}
	return
}

func (c *Client) storeInfo(registrationTokenStr string, locationHost string) {
	regArr := strings.Split(registrationTokenStr, ";")
	registrationToken := ""
	registrationExpires := ""
	if len(regArr) > 0 {
		for _, v := range regArr {
			v = strings.Replace(v, " ", "", -1)
			if len(v) > 0 {
				if strings.Index(v, "registrationToken=") > -1 {
					vv := strings.Split(v, "registrationToken=")
					registrationToken = vv[1]
				} else {
					vv := strings.Split(v, "=")
					if vv[0] == "expires" {
						registrationExpires = vv[1]
					}
					if vv[0] == "endpointId" {
						if vv[1] != "" {
							c.LoginInfo.EndpointId = vv[1]
						}
					}
				}

			}
		}
	}
	c.LoginInfo.LocationHost = locationHost
	c.LoginInfo.RegistrationToken = registrationToken
	c.LoginInfo.RegistrationExpires = registrationExpires
	//println("new registrationToken2: ", registrationTokenStr)
	if strings.Index(registrationTokenStr, "endpointId=") == -1 {
		registrationTokenStr = registrationTokenStr + "; endpointId=" + c.LoginInfo.EndpointId
		//println("new registrationToken3: ", registrationTokenStr)
	} else {
		c.LoginInfo.RegistrationtokensStr = registrationTokenStr
	}
	return
}

func (c *Client) Subscribes() {
	req := Request{
		timeout: 60,
	}

	subscribePath := c.SubscribePath()
	data := map[string]interface{}{
		"interestedResources": []string{
			"/v1/threads/ALL",
			"/v1/users/ME/contacts/ALL",
			"/v1/users/ME/conversations/ALL/messages",
			"/v1/users/ME/conversations/ALL/properties",
		},
		"template":    "raw",
		"channelType": "httpLongPoll",
	}
	// fmt.Println("c.LoginInfo.RegistrationtokensStr1: ", c.LoginInfo.RegistrationtokensStr)
	header := map[string]string{
		"Authentication": "skypetoken=" + c.LoginInfo.SkypeToken,
		//"RegistrationToken":  "registrationToken=" + c.LoginInfo.RegistrationToken+"; expires=" + c.LoginInfo.RegistrationExpires + "; endpointId=" + c.LoginInfo.EndpointId,
		"RegistrationToken": c.LoginInfo.RegistrationtokensStr,
		"BehaviorOverride":  "redirectAs404",
	}
	params, _ := json.Marshal(data)
	_, err, _ := req.request("post", subscribePath, strings.NewReader(string(params)), nil, header)
	if err != nil {
		fmt.Println("Subscribes request err: ", err)
	}
}

func (c *Client) Poll() {
	req := Request{
		timeout: 60,
	}
	pollPath := c.PollPath()
	// fmt.Println("c.LoginInfo.RegistrationtokensStr2: ", c.LoginInfo.RegistrationtokensStr)
	//return
	header := map[string]string{
		"Authentication": "skypetoken=" + c.LoginInfo.SkypeToken,
		//"RegistrationToken":  "registrationToken=" + c.LoginInfo.RegistrationToken+"; expires=" + c.LoginInfo.RegistrationExpires + "; endpointId=" + c.LoginInfo.EndpointId,
		"RegistrationToken": c.LoginInfo.RegistrationtokensStr,
		"BehaviorOverride":  "redirectAs404",
	}
	data := map[string]interface{}{
		"endpointFeatures": "Agent",
	}
	fmt.Println()
	fmt.Println("The message listener is ready")
	fmt.Println()
	params, _ := json.Marshal(data)
Loop:
	for i := 0; i <= 1000; i++ {
		if i > 1000 {
			goto Loop
		}
		body, err, _ := req.request("post", pollPath, strings.NewReader(string(params)), nil, header)
		if err != nil {
			fmt.Println("poller err: ", err)
		}
		fmt.Println("poller body: ", body)
		if body != "" {
			var bodyContent struct {
				EventMessages []Conversation `json:"eventMessages"`
			}
			err = json.Unmarshal([]byte(body), &bodyContent)
			if err != nil {
				// fmt.Println("poller body: ", body)
				fmt.Println("json.Unmarshal poller body err: ", err)
			}
			fmt.Printf("%v", bodyContent)
			if len(bodyContent.EventMessages) > 0 {
				for _, message := range bodyContent.EventMessages {
					if message.Type == "EventMessage" {

					}
				}
			}
		}
	}
}

func (c *Client) PollPath() (path string) {
	path = c.LoginInfo.LocationHost + "/v1/users/ME/endpoints/" + c.LoginInfo.EndpointId + "/subscriptions/0/poll"
	return
}

func (c *Client) SubscribePath() (path string) {
	path = c.LoginInfo.LocationHost + "/v1/users/ME/endpoints/" + c.LoginInfo.EndpointId + "/subscriptions"
	return
}

func (c *Client) getToken(t string) (err error) {

	// # Now pass the login credentials over.
	params_map := url.Values{}
	params_map.Set("client_id", "578134")
	params_map.Set("redirect_uri", "https://web.skype.com")

	req := Request{
		timeout: 30,
	}
	data := map[string]interface{}{
		"t":            t,
		"client_id":    "578134",
		"oauthPartner": "999",
		"site_name":    "lw.skype.com",
		"redirect_uri": "https://web.skype.com",
	}
	query, _ := json.Marshal(data)
	_, err, _, token, exprise := req.HttpPostBase(fmt.Sprintf("%s/microsoft?%s", API_LOGIN, gurl.BuildQuery(params_map)), string(query))
	c.LoginInfo = &LoginInfo{
		SkypeToken:   token,
		SkypeExpires: exprise,
	}
	return
}

func (c *Client) sendCreds(username, pwd, MSPRequ, MSPOK, PPFT string) (body string, err error, t_value string) {
	// # Now pass the login credentials over.
	params_map := url.Values{}
	params_map.Set("wa", "wsignin1.0")
	params_map.Set("wp", "MBI_SSL")
	params_map.Set("wreply", "https://lw.skype.com/login/oauth/proxy?client_id=578134&site_name=lw.skype.com&redirect_uri=https%3A%2F%2Fweb.skype.com%2F")
	req := Request{
		timeout: 30,
	}
	cookies := map[string]string{
		"MSPRequ": MSPRequ,
		"MSPOK":   MSPOK,
		"CkTst":   strconv.Itoa(time.Now().Second() * 1000),
	}
	params_map.Add("login", username)
	params_map.Add("passwd", pwd)
	params_map.Add("PPFT", PPFT)
	query, _ := json.Marshal(params_map)
	body, err, _, t_value = req.HttpPostWithParamAndDataWithIdt(fmt.Sprintf("%s/ppsecure/post.srf", API_MSACC), params_map, string(query), cookies, "t")
	return
}

func (c *Client) getParams() (MSPRequ, MSPOK, PPFT string, err error) {
	params := url.Values{}
	params.Set("client_id", "578134")
	params.Set("redirect_uri", "https://web.skype.com")
	if err != nil {
		return "", "", "", errors.New("parameters is not right！")
	}
	req := Request{
		timeout: 30,
	}
	//第一步, 302重定向跳转
	//fmt.Println(fmt.Sprintf("%s/oauth/microsoft", API_LOGIN))
	redirect_url, err, _ := req.HttpGetJson(fmt.Sprintf("%s/oauth/microsoft", API_LOGIN), params)
	//请求跳转的链接
	if err != nil {
		return "", "", "", errors.New("error redirect url at first step")
	}
	lgoin_srf_param := url.Values{}
	login_srf_body, err, login_srf_response := req.HttpGetJsonBackResponse(redirect_url, lgoin_srf_param)
	//从 内容中匹配出来  PPFT
	buf := `<input.*?name="PPFT".*?value="(.*?)` + `\"`
	reg := regexp.MustCompile(buf)
	ppfts := reg.FindAllString(login_srf_body, -1)
	var ppft_byte []byte
	var ppft_str string
	if len(ppfts) > 0 {
		for k, v := range ppfts {
			if k == 0 {
				ppftbbf := `value=".*?"`
				ppftreg := regexp.MustCompile(ppftbbf)
				ppftsppft := ppftreg.FindAllString(v, -1)
				ppft_byte = []byte(ppftsppft[0])[7:]
				ppft_str = string(ppft_byte[0 : len(ppft_byte)-1])
			}
		}
	}
	for _, v := range login_srf_response.Cookies() {
		if v.Name == "MSPRequ" {
			MSPRequ = v.Value
		}
		if v.Name == "MSPOK" {
			MSPOK = v.Value
		}
	}
	//发送账号密码  判定是否存在次账号
	return MSPRequ, MSPOK, ppft_str, nil
}
