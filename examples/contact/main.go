package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/kelaresg/go-skypeapi"
)

func main() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	if err != nil {
		fmt.Printf("contact name err: %s", err.Error())
		return
	}

	userId := cli.UserProfile.Username
	fmt.Println(userId)

	cli.ContactList(userId)
	for _,v := range cli.Store.Contacts  {
		fmt.Printf("\ncontact name: %s, contact id: %s", v.DisplayName, v.PersonId)
	}

	//get contact group list
	fmt.Println("\ngroup: Retrieve a list of contact groups defined by the user")
	cli.ContactGroupList(userId)
	for _, gv := range cli.ContactClient.Groups.Groups {
		fmt.Println("contact group :", gv)
	}

	//get blocklist
	fmt.Println("\ngroup: Retrieve a list of blocked users")
	cli.BlockList(userId)
	for _, bv := range cli.ContactClient.Blocks.Blocklist  {
		fmt.Println("block:", bv)
	}

	fmt.Println("\n-----------------------------end-------------------------------")
}
