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
	cli.Subscribes()
	cli.Poll()
}
