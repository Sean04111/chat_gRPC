package main

import (
	"bufio"
	pb "chat_gRPC/service/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

type User struct {
	ClientChat   pb.ChatClient
	ClientNotify pb.NotifyClient
	Self         *pb.User
}

func (this *User) SaytoAll() {
	fmt.Println("Say something : ")
	for {
		a := bufio.NewReader(os.Stdin)
		input, err := a.ReadString('\n')
		if err == nil {
			_, err2 := this.ClientChat.SendAll(context.Background(), &pb.Message{Id: this.Self.Id,
				Content:     input,
				Speakername: this.Self.Name,
			})
			if err2 != nil {
				return
			}
		} else {
			continue
		}
	}
}
func (this *User) Notify() {
	notifycode, _ := this.ClientNotify.Notify(context.Background(), &pb.UserId{Id: this.Self.Id})
	for {
		code, err := notifycode.Recv()
		if err != nil {
			fmt.Println("Notify error : ", err)
			return
		}
		if code.Code == 200 {
			break
		}
	}
	this.GetNews()
}
func (this *User) GetNews() {
	news, err := this.ClientChat.HaveAll(context.Background(), &pb.UserId{Id: this.Self.Id})
	if err != nil {
		fmt.Println("HaveAll error : ", err)
	}
	fmt.Println(news.Speakername, " : ", news.Content)
}
func main() {
	conn1, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial error : ", err)
	}
	cliC := pb.NewChatClient(conn1)
	conn2, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial error : ", err)
	}
	cliN := pb.NewNotifyClient(conn2)
	user := new(pb.User)
	user.Id = 1
	user.Name = "Sean"
	Sean := new(User)
	Sean.ClientChat = cliC
	Sean.ClientNotify = cliN
	Sean.Self = user
	Sean.SaytoAll()
	go Sean.Notify()
}
