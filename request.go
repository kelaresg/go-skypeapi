package skype

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/encoding/gurl"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

/**
 *  get and post curl function
 */

type Request struct {
	timeout time.Duration
}

/**
get login info, if not login, return login prompt , else return user login sign or token
*/
func (req *Request) getAuthorization() (sign string, dateTime string, err error) {
	//
	return
}

func (req *Request) requestReturnResponse(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, header map[string]string) (response *http.Response, err error) {
	u, err := gurl.ParseURL(reqUrl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defaultDomain := u["host"]
	//获得每次登录的信息  然后通过token 请求 skype 的官方接口
	if req.timeout == 0 {
		req.timeout = 10
	}

	client := &http.Client{
		Timeout: req.timeout * time.Second, //set timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req1, err := http.NewRequest(method, reqUrl, reqBody) //set body
	if err != nil {
		return
	}
	agent := "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.117 Safari/537.36"
	//add commom header
	req1.Header.Set("Host", defaultDomain)
	req1.Header.Set("User-Agent", agent)
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Content-Type", "application/json")

	for k, v := range header {
		req1.Header.Set(k, v)
	}
	if strings.Index(reqUrl, "ppsecure/post") > -1 {
		// add other cookie
		MaxAge := time.Hour * 24 / time.Second
		if len(cookies) > 0 {
			var newCookies []*http.Cookie
			jar, _ := cookiejar.New(nil)
			for cK, cV := range cookies {
				newCookies = append(newCookies, &http.Cookie{
					Name:     cK,
					Value:    cV,
					Path:     "/",
					Domain:   defaultDomain,
					MaxAge:   int(MaxAge),
					HttpOnly: false,
				})
			}
			jar.SetCookies(req1.URL, newCookies)
			client.Jar = jar
		}
	}

	response, err = client.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	return response, nil
}

/**
 requestWithReturnCookie
*/
func (req *Request) requestWithReturnCookie(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, header map[string]string) (body string, err error, status int, needCookies map[string]string) {
	resp, err := req.requestReturnResponse(method, reqUrl, reqBody, cookies, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	needCookies = map[string]string{}
	headerCookie := resp.Header.Values("Set-Cookie")
	for _, item := range headerCookie {
		itemArr := strings.Split(item, "; ")
		if len(itemArr) > 1 {
			needCookieArr := strings.Split(itemArr[0], "=");
			if len(needCookieArr) > 1 {
				needCookies[needCookieArr[0]] = needCookieArr[1]
			}
		}

	}
	content, err := ioutil.ReadAll(resp.Body) // content, err := ioutil.ReadAll(resp.Header)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
	}
	status = resp.StatusCode
	return
}

/**
 request
 */
func (req *Request) request(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, header map[string]string) (body string, err error, status int) {
	resp, err := req.requestReturnResponse(method, reqUrl, reqBody, cookies, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
	}
	status = resp.StatusCode
	return
}

/**
 * request can return headers
 */
func (req *Request) requestWithCookies(method string, reqUrl string, reqBody io.Reader, cookies map[string]string) (body string, err error, response *http.Response) {
	resp, err := req.requestReturnResponse(method, reqUrl, reqBody, cookies, nil)
	if err != nil {
		return
	}
	//判断是否跳转
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
	}
	response = resp
	return
}

func (req *Request) requestWithCookiesReturnIdValue(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, id string, selector string) (body string, err error, r io.Reader, tValue string) {
	u, err := gurl.ParseURL(reqUrl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defaultDomain := u["host"]
	//获得每次登录的信息  然后通过token 请求 skype 的官方接口
	sign, dateTime, err := req.getAuthorization()
	if err != nil {
		return
	}
	//默认超时
	if req.timeout == 0 {
		req.timeout = 10
	}

	client := &http.Client{
		Timeout: req.timeout * time.Second, //set timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req1, err := http.NewRequest(method, reqUrl, reqBody) //set body
	if err != nil {
		return
	}

	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	//add commom header
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Accept-Charset", "utf-8;")
	req1.Header.Set("Host", defaultDomain)
	req1.Header.Set("X-Date", dateTime)
	req1.Header.Set("Content-Type", "application/html")
	req1.Header.Set("Authorization", sign)
	req1.Header.Set("User-Agent", agent)

	//add other cookie
	MaxAge := time.Hour * 24 / time.Second
	if len(cookies) > 0 {
		var newCookies []*http.Cookie
		jar, _ := cookiejar.New(nil)
		for cK, cV := range cookies {
			newCookies = append(newCookies, &http.Cookie{
				Name:     cK,
				Value:    cV,
				Path:     "/",
				Domain:   defaultDomain,
				MaxAge:   int(MaxAge),
				HttpOnly: false,
			})
		}
		jar.SetCookies(req1.URL, newCookies)
		client.Jar = jar
	}
	resp, err := client.Do(req1)
	if err != nil {
		return
	}
	//判断是否跳转
	defer resp.Body.Close()
	r = resp.Body
	if id != "" && selector != "" {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			idt, _ := s.Attr("id")
			if idt == id {
				tValue, _ = s.Attr("value")
			}

		})
	}
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("resp.StatusCode", resp.StatusCode)
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
	}
	return
}

/**
底层的请求封装
*/
func (req *Request) requestWithLogininfo(method string, reqUrl string, reqBody io.Reader, header map[string]string) (body string, err error, status int, skypetken, expires_in string) {
	resp, err := req.requestReturnResponse(method, reqUrl, reqBody, nil, header)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取登录信息值
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		attrName, _ := s.Attr("name")
		attrVlue, _ := s.Attr("value")
		if attrName == "skypetoken" {
			skypetken = attrVlue
		}
		if attrName == "expires_in" {
			expires_in = attrVlue
		}
	})

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)

	}
	status = resp.StatusCode
	return
}

/**
GET function

*/
func (req *Request) HttpGetJson(path string, params url.Values) (body string, err error, http_status int) {
	//组装
	reqUrl := fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	//请求
	body, err, http_status = req.request("GET", reqUrl, nil, nil, nil)
	return
}

func (req *Request) HttpPostBase(path string, reqBody io.Reader, header map[string]string) (body string, err error, http_code int, skype_token, expires_in string) {
	body, err, http_code, skype_token, expires_in = req.requestWithLogininfo("POST", path, reqBody, header)
	return
}

func (req *Request) HttpGetJsonBackResponse(path string, params url.Values) (body string, err error, res *http.Response) {
	//组装
	reqUrl := fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	//请求
	body, err, res = req.requestWithCookies("GET", reqUrl, nil, nil)
	return
}

/**
POST
@params example
gin.H = type H map[string]interface{}
searchCon := gin.H{
	"pageNum":  page,
	"pageSize": limit,
	"data":     search,
}
params, _ := json.Marshal(searchCon)
*/
func (req *Request) HttpPostJson(path string, params string, cookies map[string]string) (body string, err error, res *http.Response) {
	//请求
	body, err, res = req.requestWithCookies("POST", path, strings.NewReader(params), cookies)
	return
}

/**
add post request with params and data
*/
func (req *Request) HttpPostWithParamAndDataWithIdt(path string, params url.Values, data string, cookies map[string]string, id string) (body string, err error, res io.Reader, tValue string) {
	reqUrl := fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	body, err, res, tValue = req.requestWithCookiesReturnIdValue("POST", reqUrl, strings.NewReader(data), cookies, id, "input")
	return
}

func (req *Request) HttpPostRegistrationToken(path string, data string, header map[string]string) (registrationToken, location string, err error) {
	//获得  resgistration token 信息
	resp, err := req.requestReturnResponse("POST", path, strings.NewReader(data), nil, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	registrationToken = resp.Header.Get("Set-Registrationtoken")
	location = resp.Header.Get("Location")
	return
}

func (req *Request) HttpGetWitHeaderAndCookiesJson(path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, err error) {
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	body, err, _ = req.request("GET", reqUrl, nil, cookies, headers)
	return
}

func (req *Request) HttpPostWitHeaderAndCookiesJson(path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, err error) {
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	body, err, _ = req.request("POST", reqUrl, strings.NewReader(data), cookies, headers)
	return
}

func (req *Request) HttpPutWitHeaderAndCookiesJson (path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, httpCode int, err error){
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	body, err, httpCode = req.request("PUT", reqUrl, strings.NewReader(data), cookies, headers)
	return
}