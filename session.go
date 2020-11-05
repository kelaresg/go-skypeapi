package skype

type Session struct {
	Username             string
	Password             string
	SkypeToken           string
	SkypeExpires         string
	RegistrationToken    string
	RegistrationTokenStr string
	RegistrationExpires  string
	LocationHost         string
	EndpointId           string
}

//type Session struct {
//	SkypeToken           string
//	SkypeExpires         string
//	RegistrationToken    string
//	RegistrationExpires  string
//	LocationHost         string
//	EndpointId           string
//	RegistrationTokenStr string
//}