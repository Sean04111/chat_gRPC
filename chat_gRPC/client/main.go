package main

import (
	"chat_gRPC/client/client_init"
	pb "chat_gRPC/client/protocol"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("[net] Client Dial error : ", err)
	}
	defer conn.Close()
	client := pb.NewMainChatClient(conn)
	Bill := new(pb.User)
	Bill.Name = "Bill"
	Bill.Id = 123
	c := client_init.NewClient(client, Bill)
	for {
		c.SendAndGet()
		c.Show()
	}
}
