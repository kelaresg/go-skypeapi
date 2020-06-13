package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"github.com/kelaresg/go-skypeapi"
	"strings"
)

/**
Retrieve the join URL for a group conversation, if it is currently public.
 */
func main6() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	c := skype.Conn{}

	testUserId := "8:"+cli.UserProfile.Username
	// "19:0be6022fd0d843b4916cf5c0492c3412@thread.skype"
	fmt.Printf("\niuput conversationId and enter to get url:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	//"19:0be6022fd0d843b4916cf5c0492c3412@thread.skype"
	inputArr := strings.Split(input, " ")
	conversationId := inputArr[0]
	c.GetConJoinUrl(cli.LoginInfo.LocationHost, cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr, conversationId, testUserId)

	fmt.Println("-----------------------------end-------------------------------")
}
