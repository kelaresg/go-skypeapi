package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/kelaresg/go-skypeapi"
)

func main() {
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}

	err = cli.Login(username, pwd)

	c := skype.ConversationsClient{}
	c.GetConversations(cli.LoginInfo.LocationHost, cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr)
	fmt.Println("conversations:", c.ConversationsList)
	for _, v := range c.ConversationsList.Conversations {
		fmt.Println("conversation id :", v.Id)
		//fmt.Println("conversation LastMessage :", v.LastMessage)
	}
	fmt.Println()
	fmt.Println("The message sender is ready")

	/**\
	eg 1 : get send message params
	*/
	m := skype.MessageClient{}
	fmt.Println("-------------------------------------------")
	//for {
	//	fmt.Printf("\neg: <ConversationId> <message content>")
	//	fmt.Printf("\nEnter to send:")
	//	inputReader := bufio.NewReader(os.Stdin)
	//	input, err := inputReader.ReadString('\n')
	//	if err != nil {
	//		fmt.Printf("err: %s\n", err)
	//		return
	//	}
	//	inputArr := strings.Split(input, " ")
	//	ChatThreadId := inputArr[0]
	//	contentArr := inputArr[1:]
	//	m.SendMsg(cli.LoginInfo.LocationHost, ChatThreadId, strings.Join(contentArr, " "), cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr)
	//}
	/**
	eg 2 send file  ,
	example :
	*/
	for {
		var ChatThreadId, filename, filetype string
		var duration_ms int
		fmt.Printf("\n ChatThreadId filename filetype \n (filetype:image audio other)（example: 8:live:116xxxx691 aaa.png image; 8:live:116xxxx691 aaa.txt other;  8:live:116xxxx691 aaa.mp3 audio 4006(ms)）: ")
		fmt.Scanln(&ChatThreadId, &filename, &filetype, &duration_ms)
		fmt.Println(ChatThreadId, filename, filetype)
		fmt.Println("mp3时长：",   duration_ms)
		m.SendFile(cli.LoginInfo.LocationHost, ChatThreadId, filename, cli.LoginInfo.SkypeToken, cli.LoginInfo.RegistrationtokensStr, filetype, duration_ms)
		fmt.Println("send success")
		fmt.Println("go in next send logic")
		fmt.Println("-------------------------------")
	}

	/**
	eg 3 : send a
	 */
}