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

func (req *Request) requestReturnResopnse(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, header map[string]string) (response *http.Response, err error) {
	u, err := gurl.ParseURL(reqUrl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defaultDomain := u["host"]
	//获得每次登录的信息  然后通过token 请求 skype 的官方接口
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
	agent := "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.117 Safari/537.36"
	//add commom header
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Accept-Charset", "utf-8;")
	req1.Header.Set("Host", defaultDomain)
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("User-Agent", agent)
	for k, v := range header {
		req1.Header.Set(k, v)
	}
	if len(cookies) > 0 {
		for c_k, c_v := range cookies {
			req1.Header.Set(c_k, c_v)
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
底层的请求封装
*/
func (req *Request) request(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, header map[string]string) (body string, err error, status int) {
	resp, err := req.requestReturnResopnse(method, reqUrl, reqBody, cookies, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	fmt.Println(" request body", resp.Request)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
		//fmt.Printf("%s", resp.Header)
		fmt.Println(resp.Header)
	}
	//
	status = resp.StatusCode
	return
}

/**
 * request can return headers
 */
func (req *Request) requestWithCookies(method string, reqUrl string, reqBody io.Reader, cookies map[string]string) (body string, err error, response *http.Response) {
	resp, err := req.requestReturnResopnse(method, reqUrl, reqBody, cookies, nil)
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

func (req *Request) requestWithCookiesReturnIdValue(method string, reqUrl string, reqBody io.Reader, cookies map[string]string, id string, selector string) (body string, err error, r io.Reader, id_value string) {
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

	agent := "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.117 Safari/537.36"
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
		cookiess := []*http.Cookie{}
		jar, _ := cookiejar.New(nil)
		for c_k, c_v := range cookies {
			cookiess = append(cookiess, &http.Cookie{
				Name:     c_k,
				Value:    c_v,
				Path:     "/",
				Domain:   defaultDomain,
				MaxAge:   int(MaxAge),
				HttpOnly: false,
			})
		}
		jar.SetCookies(req1.URL, cookiess)
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
				id_value, _ = s.Attr("value")
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
func (req *Request) requestWithLogininfo(method string, reqUrl string, reqBody io.Reader) (body string, err error, status int, skypetken, expires_in string) {
	resp, err := req.requestReturnResopnse(method, reqUrl, reqBody, nil, nil)
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

func (req *Request) HttpPostBase(path string, params string) (body string, err error, http_code int, skype_token, expires_in string) {
	body, err, http_code, skype_token, expires_in = req.requestWithLogininfo("POST", path, strings.NewReader(params))
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
func (req *Request) HttpPostWithParamAndDataWithIdt(path string, params url.Values, data string, cookies map[string]string, id string) (body string, err error, res io.Reader, t_value string) {
	reqUrl := fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	body, err, res, t_value = req.requestWithCookiesReturnIdValue("POST", reqUrl, strings.NewReader(data), cookies, id, "input")
	return
}

func (req *Request) HttpPostRegistrationToken(path string, data string, header map[string]string) (registrationtoken, location string, err error) {
	//获得  resgistration token 信息
	resp, err := req.requestReturnResopnse("POST", path, strings.NewReader(data), nil, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	registrationtoken = resp.Header.Get("Set-Registrationtoken")
	location = resp.Header.Get("Location")
	return
}

func (req *Request) HttpGetWitHeaderAndCookiesJson(path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, err error) {
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	println("HttpGetWitHeaderAndCookiesJson:", reqUrl)
	body, err, _ = req.request("GET", reqUrl, nil, cookies, headers)
	return
}

func (req *Request) HttpPostWitHeaderAndCookiesJson(path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, err error) {
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	println("HttpPostWitHeaderAndCookiesJson:", reqUrl)
	body, err, _ = req.request("POST", reqUrl, strings.NewReader(data), cookies, headers)
	return
}

func (req *Request) HttpPutWitHeaderAndCookiesJson (path string, params url.Values, data string, cookies map[string]string, headers map[string]string) (body string, httpCode int, err error){
	reqUrl := path
	if len(params) >0 {
		reqUrl = fmt.Sprintf("%s?%s", path, gurl.BuildQuery(params))
	}
	resp, err := req.requestReturnResopnse("PUT", reqUrl, strings.NewReader(data), cookies, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	fmt.Println(" request body", resp.Request)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		body = location
	} else {
		body = string(content)
		//fmt.Printf("%s", resp.Header)
		fmt.Println(resp.Header)
	}
	//
	httpCode = resp.StatusCode
	return
}