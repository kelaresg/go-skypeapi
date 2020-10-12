package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/kelaresg/go-skypeapi"
)

func main090() {
	cli, err := skype.NewConn()
	if err != nil {
		fmt.Println(err)
	}
	skype.GetConfigYaml()
	username := viper.GetString("user.username")
	pwd := viper.GetString("user.password")
	err  = cli.Login(username, pwd)
	fmt.Println(cli.UserProfile.Username)
	//c := skype.Conn{}

	// eg1
	//cli.GetConversations(cli.session.LocationHost, cli.session.SkypeToken, cli.session.RegistrationTokenStr)
	//fmt.Println("conversations start")
	//fmt.Println("conversations count :",c.ConversationsList.Metadata.TotalCount)
	//fmt.Println("conversations content :",c.ConversationsList.Conversations)
	//
	//for  _,v := range c.ConversationsList.Conversations {
	//	fmt.Printf("\nconversation id: %s", v.Id)
	//	fmt.Printf("\nconversation name: %s\n", v.ThreadProperties.Topic)
	//}

	//eg2
	// cli.GetConversation("19:0be6022fd0d843b4916cf5c0492c3412@thread.skype")

	//eg3
	// cli.GetConversationThreads(cli.session.LocationHost, cli.session.SkypeToken, cli.session.RegistrationTokenStr, "19:0be6022fd0d843b4916cf5c0492c3412@thread.skype")

	//eg4
	//member1 := skype.Member{
	//	Id: "8:live:zhaosl_4",
	//	Role: "User",
	//}
	//// The user who created the group must be in the Members and have "Admin" rights
	//member2 := skype.Member{
	//	Id: "8:live:116xxxx691",
	//	Role: "Admin",
	//}
	//Members := skype.Members{}
	//Members.Members = append(Members.Members, member1)
	//Members.Members = append(Members.Members, member2)
	//cli.CreateConversationGroup(cli.session.LocationHost, cli.session.SkypeToken, cli.session.RegistrationTokenStr, Members)

	//eg5
	//member1 := skype.Member{
	//	Id: "8:live:.cid.d3feb90dceeb51cc",
	//	//Id: "8:live:live:liyu13526435030",
	//	Role: "Admin",
	//}
	//Members := skype.Members{}
	//Members.Members = append(Members.Members, member1)
	//cli.AddMemberToConversation(cli.session.LocationHost, cli.session.SkypeToken, cli.session.RegistrationTokenStr, Members, "19:0be6022fd0d843b4916cf5c0492c3412@thread.skype")

	//eg6
	//testUserId := "8:live:.cid.d3feb90dceeb51cc"
	//cli.RemoveMemberFromConversation(cli.session.LocationHost, cli.session.SkypeToken, cli.session.RegistrationTokenStr, "19:0be6022fd0d843b4916cf5c0492c3412@thread.skype", testUserId)

	//EG7
	//cli.SetConversationThreads("19:3131232131@thread.skype", map[string]string{
	//	"topic": "3132131",
	//})
	fmt.Println("-----------------------------end-------------------------------")
}
