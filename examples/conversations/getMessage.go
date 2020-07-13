package main

import (
	"fmt"
	"github.com/kelaresg/go-skypeapi"
	"github.com/spf13/viper"
)

/**
Retrieve the join URL for a group conversation, if it is currently public.
 */
func main() {
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
	// "19:0be6022fd0d843b4916cf5c0492c3412@thread.skype"
	//fmt.Printf("\niuput conversationId and enter to get url:")
	//inputReader := bufio.NewReader(os.Stdin)
	//input, err := inputReader.ReadString('\n')
	//if err != nil {
	//	fmt.Printf("err: %s\n", err)
	//	return
	//}
	////"19:0be6022fd0d843b4916cf5c0492c3412@thread.skype"
	//inputArr := strings.Split(input, " ")
	//conversationId := inputArr[0]
	//fmt.Println()

	res, err := cli.GetMessages("19:b44b8a9b030e4fe4a2400d517c9f31c8@thread.skype", "", "10")
	next := res.Metadata.SyncState
	res2, err := cli.GetMessages("19:b44b8a9b030e4fe4a2400d517c9f31c8@thread.skype", next, "10")
	fmt.Println("-----------------------------start res2-------------------------------")
	fmt.Println("res2 分页面：", res2)
	fmt.Println("-----------------------------end-------------------------------")
}
