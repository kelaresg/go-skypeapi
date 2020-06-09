package main

import (
	"fmt"
	"github.com/spf13/viper"
	"skype"
)

func main() {
	cli, err := skype.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	cli.Subscribes()
	cli.Poll()
}
