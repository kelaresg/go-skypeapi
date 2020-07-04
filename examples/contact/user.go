package main

import (
	"fmt"
	"github.com/kelaresg/go-skypeapi"
	"github.com/spf13/viper"
)

func main() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	_, err  = cli.Login(username, pwd)
	skypetoken := cli.LoginInfo.SkypeToken
	user := skype.User{}
	// eg1
	user.GetProfile(skypetoken, "live:love.kimi_2")


	// eg2
	//fmt.Printf("\niuput newName and enter to update:")
	//inputReader := bufio.NewReader(os.Stdin)
	//input, err := inputReader.ReadString('\n')
	//if err != nil {
	//	fmt.Printf("err: %s\n", err)
	//	return
	//}
	//inputArr := strings.Split(input, " ")
	//firstName := inputArr[0]
	//inputArr = inputArr[1:]
	//lastName := strings.Join(inputArr, " ")
	//
	//user.UpdateName(skypetoken, firstName, lastName)
}
