package main

import (
	pb "chat_gRPC/service/proto"
	"chat_gRPC/service/server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var Port string = ":8080"
var MessageNum int = 100

func main() {
	fmt.Println("starting the server.....")
	MassServer := new(server.MessageServer)
	MassServer.Masseges = make([]*pb.Message, MessageNum)
	conn1, err := net.Listen("tcp", Port)

	if err != nil {
		fmt.Println("Listen error :", err)
	}
	m := grpc.NewServer()
	pb.RegisterChatServer(m, MassServer)
	err_ := m.Serve(conn1).Error()
	fmt.Println(err_)
}
