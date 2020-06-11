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
func main() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	c := skype.ContactClient{}

	fmt.Printf("\ne.g:live:******")
	fmt.Printf("\niuput skypeUsername and enter to remove:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	inputArr := strings.Fields(input)
	otherid := "8:"+inputArr[0]
	c.DeleteContact(cli.LoginInfo.SkypeToken, cli.UserProfile.Username, otherid)
	//c.AddContact(cli.LoginInfo.SkypeToken, cli.UserProfile.Username, otherid)
	fmt.Println("-----------------------------end-------------------------------")
}
