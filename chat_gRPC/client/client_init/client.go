package client_init

import (
	"bufio"
	pb "chat_gRPC/client/protocol"
	"context"
	"fmt"
	"os"
)

type Client struct {
	Clienty  pb.MainChatClient
	MainChat *pb.MainChatList
	Self     *pb.User
}

func (this *Client) Show() {
	fmt.Println("世界消息：")
	for _, b := range this.MainChat.MainChatList {
		if len(b.Message) != 0 {
			fmt.Println("[", b.Speaker.Name, "]", " : ", b.Message)
		}
	}
}
func (this *Client) SendAndGet() {
	stream, err := this.Clienty.SendAll(context.Background())
	if err != nil {
		fmt.Println("[Client] SendAll error : ", err)
		return
	}
	var n int = 0
	fmt.Println("说点什么：(你可以一次性发送三条消息！)")
	for {
		a := bufio.NewReader(os.Stdin)
		input, err := a.ReadString('\n')
		if n == 3 {
			break
		}
		if err == nil {
			n = n + 1
			mass := new(pb.Message)
			mass.Message = input
			mass.Speaker = this.Self
			err1 := stream.Send(mass)
			if err1 != nil {
				fmt.Println("[Client] Stream Send error : ", err1)
				break
			} else {
				continue
			}
		}

	}
	MainChat, err2 := stream.CloseAndRecv()
	if err2 != nil {
		fmt.Println("[Client] stream Recv error : ", err2)
	}
	this.MainChat.MainChatList = MainChat.MainChatList
}
func NewClient(cling pb.MainChatClient, user *pb.User) Client {
	var cli Client
	cli.Clienty = cling
	cli.MainChat = new(pb.MainChatList)
	cli.MainChat.MainChatList = make([]*pb.Message, 50)
	cli.Self = user
	return cli
}
