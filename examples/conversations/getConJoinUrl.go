package main

import (
	"bufio"
	"fmt"
	"github.com/kelaresg/go-skypeapi"
	"github.com/spf13/viper"
	"os"
	"strings"
)

/**
Retrieve the join URL for a group conversation, if it is currently public.
 */
func main444() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	_, err  = cli.Login(username, pwd)
	//c := skype.Conn{}

	//testUserId := "8:"+cli.UserProfile.Username
	//fmt.Printf("\niuput conversationId and enter to get url:")
	//inputReader := bufio.NewReader(os.Stdin)
	//input, err := inputReader.ReadString('\n')
	//if err != nil {
	//	fmt.Printf("err: %s\n", err)
	//	return
	//}
	//inputArr := strings.Split(input, " ")
	//conversationId := inputArr[0]
	//fmt.Println()
	fmt.Printf("\niuput conversationId and enter to get url:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	inputArr := strings.Split(input, " ")
	cli.GetConJoinUrl(inputArr[0])

	fmt.Println("-----------------------------end-------------------------------")
}
