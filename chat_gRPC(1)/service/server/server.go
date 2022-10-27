package server

import (
	pb "chat_gRPC/service/proto"
	"context"
	"fmt"
)

type MessageServer struct {
	pb.UnimplementedChatServer
	Masseges []*pb.Message
	Manager  NotifyManager
}
type NotifyManager struct {
	pb.UnimplementedNotifyServer
	Code *pb.NotiCode
}

func (this *MessageServer) SendAll(ctx context.Context, mess *pb.Message) (*pb.Message, error) {
	this.Masseges[mess.Id] = mess
	this.Manager.Code.Code = 200
	return mess, nil
}
func (this *MessageServer) HaveAll(ctx context.Context, id *pb.UserId) (*pb.Message, error) {
	for n, m := range this.Masseges {
		if len(this.Masseges[n+1].Content) == 0 {
			return m, nil
		} else {
			continue
		}
	}
	return nil, nil
}
func (this *NotifyManager) Notify(id *pb.UserId, stream pb.Notify_NotifyServer) error {
	for {
		err := stream.Send(this.Code)
		if err != nil {
			fmt.Println("[Notify Manager]:Send error : ", err)
			break
		}
	}
	return nil
}
