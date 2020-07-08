package skype
import (
	"errors"
	"fmt"
	"github.com/gogf/gf/encoding/gurl"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	ErrMediaDownloadFailedWith404 = errors.New("download failed with status code 404")
	ErrMediaDownloadFailedWith410 = errors.New("download failed with status code 410")
	ErrNoURLPresent       = errors.New("no url present")
	ErrFileLengthMismatch = errors.New("file length does not match")
	ErrTooShortFile       = errors.New("file too short")
)

func Download(url string, ce *Conn, fileLength int) ([]byte, error) {
	if url == "" {
		return nil, ErrNoURLPresent
	}
	file, err := ce.downloadMedia(url)

	if err != nil {
		return nil, err
	}
	//if len(file) != fileLength {
	//	return nil, ErrFileLengthMismatch
	//}
	return file, nil
}

func (c *Conn) downloadMedia(url string) (file []byte, err error) {
	fmt.Println("downloadMedia:", url)
	client := &http.Client{
		//Timeout: 20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	headers := map[string]string{
		"skypetoken_asm":    c.LoginInfo.SkypeToken, // "skype_token " + Conn.Session.SkypeToken,
	}
	cookies := map[string]string{
		"skypetoken_asm":    c.LoginInfo.SkypeToken, // "skype_token " + Conn.Session.SkypeToken,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	//add other cookie
	u, err := gurl.ParseURL(url, 2)
	MaxAge := time.Hour * 24 / time.Second
	if len(cookies) > 0 {
		var newCookies []*http.Cookie
		jar, _ := cookiejar.New(nil)
		for cK, cV := range cookies {
			newCookies = append(newCookies, &http.Cookie{
				Name:     cK,
				Value:    cV,
				Path:     "/",
				Domain:   u["host"],
				MaxAge:   int(MaxAge),
				HttpOnly: false,
			})
		}
		jar.SetCookies(req.URL, newCookies)
		client.Jar = jar
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 302 {
			location := resp.Header.Get("Location")
			return c.downloadMedia(location)
		}
		if resp.StatusCode == 404 {
			return nil, ErrMediaDownloadFailedWith404
		}
		if resp.StatusCode == 410 {
			return nil, ErrMediaDownloadFailedWith410
		}
		fmt.Printf("%+v",resp)
		return nil, fmt.Errorf("download failed with status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	if resp.ContentLength <= 10 {
		return nil, ErrTooShortFile
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}