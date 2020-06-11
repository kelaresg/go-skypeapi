package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"github.com/kelaresg/go-skypeapi"
	"strings"
)

func main7() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	c := skype.ConversationsClient{}

	userId := cli.UserProfile.Username
	fmt.Printf("\ne.g:conversationName(no spaces) memberUsername1 memberUsername2")
	fmt.Printf("\nenter to create:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	inputArr := strings.Split(input, " ")
	Topic := inputArr[0]
	inputArr = inputArr[1:]
	Members := skype.Members{}

	// The user who created the group must be in the Members and have "Admin" rights
	member2 := skype.Member{
		Id: "8:"+userId,
		Role: "Admin",
	}

	Members.Members = append(Members.Members, member2)
	Members.Properties = skype.Properties {
		HistoryDisclosed: "true",
		Topic: Topic,
	}
	c.CreateConversationGroup(cli.LoginInfo.LocationHost, cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr, Members)

	Members = skype.Members{}
	for _, memberId := range inputArr {
		Members.Members = append(Members.Members, skype.Member{
			Id: "8:"+memberId,
			Role: "Admin",
		})
	}
	c.AddMember(cli.LoginInfo.LocationHost, cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr, Members, "")
	fmt.Println("-----------------------------end-------------------------------")
}
