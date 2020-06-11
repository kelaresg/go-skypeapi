package main

import (
	"fmt"
	"github.com/kelaresg/go-skypeapi"
	"github.com/spf13/viper"
)

func main() {
	cli, err := skype.NewClient()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	fmt.Printf("\nSkypeToken : %+v\n", cli.LoginInfo.SkypeToken)
	fmt.Printf("SkypeExpires : %+v\n", cli.LoginInfo.SkypeExpires)
	fmt.Printf("RegistrationToken : %+v\n", cli.LoginInfo.RegistrationToken)
	fmt.Printf("EndpointId : %+v\n", cli.LoginInfo.EndpointId)
	fmt.Printf("RegistrationExpires : %+v\n", cli.LoginInfo.RegistrationExpires)
	fmt.Printf("LocationHost : %+v\n", cli.LoginInfo.LocationHost)
	fmt.Printf("UserId : %+v\n", cli.UserProfile.Username)
	fmt.Printf("RegistrationtokensStr : %+v\n", cli.LoginInfo.RegistrationtokensStr)
}
