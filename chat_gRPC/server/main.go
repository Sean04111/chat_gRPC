package main

import (
	pb "chat_gRPC/server/protocol"
	"chat_gRPC/server/server_init"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	Address := "127.0.0.1:8080"
	conn, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("[net] Listen error : ", err)
	}
	var ChatList = new(server_init.Chat)
	ChatList.ChatMain = new(pb.MainChatList)
	ChatList.ChatMain.MainChatList = make([]*pb.Message, 50)
	var s = grpc.NewServer()
	pb.RegisterMainChatServer(s, &server_init.Server{MainChatList: *ChatList})
	//先取消并发模式看看
	err = s.Serve(conn)
	if err != nil {
		fmt.Println("[net] Serve error : ", err)
	}
}
