package main

import (
	"fmt"
	"github.com/kelaresg/go-skypeapi"
	"github.com/spf13/viper"
)

func main4() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYamlForBuildExample()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	// eg1
	//user.GetProfile(skypetoken, "live:xxxxxx")
	//user.GetContactsProfile(skypetoken)
	// username
	// like live:xxxxx will be more precise
	cli.NameSearch( "keyword")


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
