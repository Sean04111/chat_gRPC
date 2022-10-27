package main

import (
	pb "chat_gRPC/service/proto"
	"chat_gRPC/service/server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	MassServer := new(server.MessageServer)
	MassServer.Masseges = make([]*pb.Message, 50)
	NotiManager := new(server.NotifyManager)
	NotiManager.Code = &pb.NotiCode{Code: 400}
	MassServer.Manager = *NotiManager
	conn1, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Listen error :", err)
	}
	m := grpc.NewServer()
	pb.RegisterChatServer(m, MassServer)
	m.Serve(conn1)
	conn2, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Listen error :", err)
	}
	n := grpc.NewServer()
	pb.RegisterNotifyServer(n, NotiManager)
	m.Serve(conn2)
}
