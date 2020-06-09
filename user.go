package skype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type User struct {
	Payload *Payload
}

type Payload struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

/**
 * Retrieve the join URL for a group conversation, if it is currently public.
 */
func (c *User)GetProfile(skypeToken string, UserId string) () {
	//join url e.g https://join.skype.com/rIYuOwyiqUIu
	path := fmt.Sprintf("%s/users/%s/profile", API_USER, UserId)
	req := Request{timeout: 30}
	headers := map[string]string{
		"X-Skypetoken":      skypeToken,
		//"Authentication":    "skypetoken=" + skypeToken,
		//"RegistrationToken": regToken,
		//"BehaviorOverride":  "redirectAs404",
		//"X-Skypetoken":     skypeToken,
	}
	//data := map[string]string{
	//	"shortId": "rIYuOwyiqUIu",
	//	"type":   "wl",
	//}
	//params, _ := json.Marshal(data)
	//body, err, _ := req.request("post", path, strings.NewReader(string(params)), nil, headers)
	body, err, _ := req.request("GET", path, nil, nil, headers)
	if err != nil {
		fmt.Println("getProfile err: ", err)
	}
	fmt.Println("getProfile resp: ", body)
	conInfo := JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	fmt.Println(conInfo.Resource)
	return
}

func (c *User)UpdateName(skypeToken string, firstName string, lastName string) () {
	//join url e.g https://join.skype.com/rIYuOwyiqUIu
	//path := fmt.Sprintf("%s/users/%s/profile", API_USER, UserId)
	path := "https://edge.skype.com/profile/v1/users/self/profile/partial"
	req := Request{timeout: 30}
	headers := map[string]string{
		"x-skypetoken":      skypeToken,
	}
	data := map[string]interface{}{
		"payload": Payload{
			firstName,
			lastName,
		},
	}
	params, _ := json.Marshal(data)
	body, err, _ := req.request("POST", path, strings.NewReader(string(params)), nil, headers)
	if err != nil {
		fmt.Println("getProfile err: ", err)
	}
	fmt.Println("getProfile resp: ", body)
	conInfo := JoinToConInfo{}
	json.Unmarshal([]byte(body), &conInfo)
	fmt.Println(conInfo.Resource)
	return
}