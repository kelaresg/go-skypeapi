package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"github.com/kelaresg/go-skypeapi"
	"strings"
)

func main3() {
	cli, err := skype.NewClient()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYamlForBuildExample()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	skypetoken := cli.LoginInfo.SkypeToken

	fmt.Printf("\niuput newName and enter to update:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	inputArr := strings.Split(input, " ")
	firstName := inputArr[0]
	inputArr = inputArr[1:]
	lastName := strings.Join(inputArr, " ")
	user := skype.User{}
	//user.GetProfile(skypetoken, "self")
	user.UpdateName(skypetoken, firstName, lastName)
}
