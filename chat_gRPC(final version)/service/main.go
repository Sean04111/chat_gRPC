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
	conn1, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Listen error :", err)
	}
	m := grpc.NewServer()
	pb.RegisterChatServer(m, MassServer)
	err = m.Serve(conn1)
	if err != nil {
		fmt.Println("[server] Serve error : ", err)
	}

}
