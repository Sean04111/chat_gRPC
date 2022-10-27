package server_init

import (
	"chat_gRPC/server/models"
	pb "chat_gRPC/server/protocol"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"sync"
)

type Server struct {
	pb.UnimplementedMainChatServer
	MainChatList Chat
}
type Chat struct {
	ChatMain *pb.MainChatList
	Locker   sync.Locker
}

func (this *Server) SendAll(stream pb.MainChat_SendAllServer) error {
	this.Listen(stream)
	return this.WriteReturn(stream)
}

//向数据库里输入数据是否为必要的？
func (this *Server) Listen(stream pb.MainChat_SendAllServer) {
	var n int = 0
	for {
		mess, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("[Server] Recv() error : ", err)
			break
		}
		//取消的并发互斥锁
		//	this.MainChatList.Locker.Lock()
		this.MainChatList.ChatMain.MainChatList[n] = mess
		n++
		messfororm := new(models.MessageInfor)
		messfororm.MessageId = n
		messfororm.Content = mess.Message
		messfororm.SpeakerId = mess.Speaker.Id
		messfororm.SpeakerName = mess.Speaker.Name
		//暂时取消对通话记录的数据库录入（没什么用）
		//	err = this.InsertMessage(messfororm)
		if err != nil {
			fmt.Println("[orm] Insert error : ", err)
		}
		//	this.MainChatList.Locker.Unlock()
	}
}
func (this *Server) WriteReturn(stream pb.MainChat_SendAllServer) error {
	return stream.SendAndClose(this.MainChatList.ChatMain)
}

func (this *Server) InsertMessage(mess *models.MessageInfor) error {
	o := orm.NewOrm()
	_, err := o.Insert(mess)
	return err
}
func (this *Server) GetMessage(speakerid int) models.MessageInfor {
	var resp models.MessageInfor
	o := orm.NewOrm()
	qsu := o.QueryTable("messageinfor")
	_, err := qsu.Filter("speakerid", speakerid).All(&resp)
	if err != nil {
		fmt.Println("[orm] Filter error : ", err)
	}
	return resp
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
