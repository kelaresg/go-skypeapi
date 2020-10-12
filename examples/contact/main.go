package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/kelaresg/go-skypeapi"
)

func main4() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	//_, _, _ = cli.LoginApiAuth("zhaosl@shinetechchina.com", "zsl630235")
	err  = cli.Login(username, pwd)
	userId := cli.UserProfile.Username

	//contant := skype.ContactClient{}
	cli.ContactList(userId)
	//users := *contant.Users
	//for _,v := range users.Contacts  {
	//	fmt.Println()
	//	fmt.Printf("contact name :%s, contact id :%s", v.DisplayName, v.PersonId)
	//}

	//get contact group list
	//fmt.Println()
	//fmt.Println()
	//fmt.Println("group: Retrieve a list of contact groups defined by the user")
	//contant.ContactGroupList(userId, skypetoken)
	//for _, gv := range contant.Groups.Groups {
	//	fmt.Println("contact group :", gv)
	//}

	//get blocklist
	//fmt.Println()
	//fmt.Println("group: Retrieve a list of blocked users")
	//contant.BlockList(userId, skypetoken)
	//for _, bv := range contant.Blocks.Blocklist  {
	//	fmt.Println("block:", bv)
	//}
	fmt.Println("-----------------------------end-------------------------------")
	//getAllContactInfo
	//contant.GetAllContactInfo(userId, skypetoken)
	//for _, bv := range contant.Blocks.Blocklist  {
	//	fmt.Println("block:", bv)
	//}
}
