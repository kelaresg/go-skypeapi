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
* block user
 */
func main2() {
	cli, err := skype.NewClient()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	c := skype.ContactClient{}

	fmt.Printf("\ne.g:live:****** block/unblock")
	fmt.Printf("\niuput skypeUsername and enter to action:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	inputArr := strings.Fields(input)
	otherid := "8:"+inputArr[0]
	if inputArr[1] == "unblock" {
		c.UnBlockContact(cli.LoginInfo.SkypeToken, cli.UserProfile.Username, otherid)
	} else {
		c.BlockContact(cli.LoginInfo.SkypeToken, cli.UserProfile.Username, otherid, false, false)
	}
	fmt.Println("-----------------------------end-------------------------------")
}
