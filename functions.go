package skype

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/**
+ SKYPEWEB_LOCKANDKEY_SECRET
*/
func getMac256Hash(secs string) string {
	clearText := secs + SKYPEWEB_LOCKANDKEY_APPID
	zeroNum := (8 - len(clearText)%8)
	for i := 0; i < zeroNum; i++ {
		clearText += "0"
	}
	//开始  加密 getMac256Hash
	cchClearText := len(clearText) / 4
	pClearText := make([]int, cchClearText)
	for i := 0; i < cchClearText; i++ {
		mib := 0
		for pos := 0; pos < 4; pos++ {
			len1 := 4*i + pos
			b := int([]rune(clearText[len1 : len1+1])[0])
			mi := int(math.Pow(256, float64(pos)))
			mib += mi * b
		}
		pClearText[i] = mib
	}
	sha256Hash := []int{
		0, 0, 0, 0,
	}
	//
	screact_key_str := secs + SKYPEWEB_LOCKANDKEY_SECRET
	h := sha256.New()
	h.Write([]byte(screact_key_str))
	sum := h.Sum(nil)
	hash_str := strings.ToUpper(string(hex.EncodeToString(sum)))
	sha256len := len(sha256Hash)
	for s := 0; s < sha256len; s++ {
		sha256Hash[s] = 0
		for pos := 0; pos < 4; pos++ {
			dpos := 8*s + pos*2
			mi1 := int(math.Pow(256, float64(pos)))
			inthash := hash_str[dpos : dpos+2]
			inthash1, _ := strconv.ParseInt(inthash, 16, 64)
			sha256Hash[s] += int(inthash1) * mi1
		}
	}
	qwMAC, qwSum := cs64(pClearText, sha256Hash)
	macParts := []int{
		qwMAC,
		qwSum,
		qwMAC,
		qwSum,
	}
	scans := []int{0, 0, 0, 0}
	for i, sha := range sha256Hash {
		scans[i] = int64Xor(sha, macParts[i])
	}
	//scan := int64Xor(sha256Hash, macParts)
	hexString := ""
	for _, scan := range scans {
		hexString += int32ToHexString(scan)
	}
	return hexString
}

func int32ToHexString(n int) (hexString string) {
	hexChars := "0123456789abcdef"
	for i := 0; i < 4; i++ {
		num1 := (n >> (i*8 + 4)) & 15
		num2 := (n >> (i * 8)) & 15
		hexString += hexChars[num1 : num1+1]
		hexString += hexChars[num2 : num2+1]
	}
	return
}

func int64Xor(a int, b int) (sc int) {
	sA := fmt.Sprintf("%b", a)
	sB := fmt.Sprintf("%b", b)
	sC := ""
	sD := ""
	diff := math.Abs(float64(len(sA) - len(sB)))
	for d := 0; d < int(diff); d++ {
		sD += "0"
	}
	if len(sA) < len(sB) {
		sD += sA
		sA = sD
	} else if len(sB) < len(sA) {
		sD += sB
		sB = sD
	}
	for a := 0; a < len(sA); a++ {
		if sA[a] == sB[a] {
			sC += "0"
		} else {
			sC += "1"
		}
	}
	to2, _ := strconv.ParseInt(sC, 2, 64)
	xor, _ := strconv.Atoi(fmt.Sprintf("%d", to2))
	return xor
}

func cs64(pdwData, pInHash []int) (qwMAC int, qwSum int) {
	MODULUS := 2147483647
	CS64_a := pInHash[0] & MODULUS
	CS64_b := pInHash[1] & MODULUS
	CS64_c := pInHash[2] & MODULUS
	CS64_d := pInHash[3] & MODULUS
	CS64_e := 242854337
	pos := 0
	qwDatum := 0
	qwMAC = 0
	qwSum = 0
	pdwLen := len(pdwData) / 2
	for i := 0; i < pdwLen; i++ {
		qwDatum = int(pdwData[pos])
		pos += 1
		qwDatum *= CS64_e
		qwDatum = qwDatum % MODULUS
		qwMAC += qwDatum
		qwMAC *= CS64_a
		qwMAC += CS64_b
		qwMAC = qwMAC % MODULUS
		qwSum += qwMAC
		qwMAC += int(pdwData[pos])
		pos += 1
		qwMAC *= CS64_c
		qwMAC += CS64_d
		qwMAC = qwMAC % MODULUS
		qwSum += qwMAC
	}
	qwMAC += CS64_b
	qwMAC = qwMAC % MODULUS
	qwSum += CS64_d
	qwSum = qwSum % MODULUS
	return qwMAC, qwSum
}

func GetConfigYaml() {
	viper.SetConfigName("config")
	viper.AddConfigPath("examples")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
}

func GetConfigYamlForBuildExample() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = getParentDirectory(dir)
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	fmt.Println(viper.ReadInConfig())
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dir string) string {
	return substr(dir, 0, strings.LastIndex(dir, "/"))
}

func UriObject(content, filetype, url, thumb, title, desc string, duration_ms int,  values map[string]string) string  {
	titleTag := "<Title/>"
	descTag := "<Description/>"
	thumbAttr := ""
	valTags := ""
	Durationms := ""
	if len(title) > 0 {
		titleTag = fmt.Sprintf(`<Title>Title: %s</Title>`, title)
	}
	if len(desc) >0  {
		descTag = fmt.Sprintf(`<Description>Description: %s</Description>`, desc)
	}
	if len(thumb) > 0 {
		thumbAttr = fmt.Sprintf(`url_thumbnail="%s"`, thumb)
	}
	if len(values) > 0 {
		for k,v := range values {
			valTags += fmt.Sprintf(`<%s v="%s"/>`, k, v)
		}
	}
	if duration_ms > 0 {
		Durationms = fmt.Sprintf(`duration_ms="%s"`, strconv.Itoa(duration_ms))
	}
	objStr := fmt.Sprintf(`<URIObject type="%s" uri="%s" %s %s>%s%s%s%s</URIObject>`, filetype, url, Durationms, thumbAttr, titleTag, descTag, valTags, content)
	return objStr
}