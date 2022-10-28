package server

import (
	"chat_gRPC/service/model"
	pb "chat_gRPC/service/proto"
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

var MessageNum int = 1

type MessageServer struct {
	pb.UnimplementedChatServer
	Masseges []*pb.Message
}

func (this *MessageServer) SendAll(ctx context.Context, mess *pb.Message) (*pb.Message, error) {
	o := orm.NewOrm()
	InsertMessa := model.Messages{MessageId: int(mess.Id), SpeakerName: mess.Speakername, Content: mess.Content, Time: mess.Time}
	_, err := o.Insert(&InsertMessa)
	if err != nil {
		fmt.Println("[orm]Insert error : ", err)
	}
	MessageNum++
	this.Masseges[mess.Id] = mess
	return mess, nil
}
func (this *MessageServer) GetMessNum(ctx context.Context, id *pb.UserId) (*pb.MessageNum, error) {
	return &pb.MessageNum{Messnum: int64(MessageNum)}, nil
}
func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Println("[orm] Register Driver error : ", err)
	}
	err = orm.RegisterDataBase("default", "mysql", "root:Tjq216828@tcp(127.0.0.1:3306)/chat_gRPC?charset=utf8")
	if err != nil {
		fmt.Println("[orm] Register Data Base error : ", err)
	}
	err = orm.RunSyncdb("default", true, false)
	if err != nil {
		fmt.Println("[orm] Create Table error : ", err)
	}
}
