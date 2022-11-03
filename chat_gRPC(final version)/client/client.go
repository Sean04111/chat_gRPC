package client

import (
	"bufio"
	"chat_gRPC/service/model"
	pb "chat_gRPC/service/proto"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

type User struct {
	ClientChat pb.ChatClient
	Self       *pb.User
}

//you can set you database by these vars
var (
	DB_database string = "sql_news"
	DB_username string = "sql_news"
	DB_password string = "666"
	DB_host     string = "121.36.131.50"
	DB_port     string = "3306"
)

func (this *User) SaytoAll() {
	fmt.Println("Say something : ")
	for {
		a := bufio.NewReader(os.Stdin)
		input, err := a.ReadString('\n')
		if err == nil {
			Now := time.Now().Month().String() + "/" + strconv.Itoa(time.Now().Day()) + " " + strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second())
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			MassageNum, _ := this.ClientChat.GetMessNum(ctx, &pb.UserId{Id: this.Self.Id})
			_, err2 := this.ClientChat.SendAll(ctx, &pb.Message{Id: MassageNum.Messnum,
				Content:     input,
				Speakername: this.Self.Name,
				Time:        Now,
			})
			if err2 != nil {
				return
			}
		} else {
			continue
		}
	}
}
func (this *User) GetMassFromId(id int) *pb.Message {
	var Re *pb.Message
	return Re
}
func (this *User) GetMassFromName(speakername string) *pb.Message {
	find := new(pb.Message)

	return find
}
func (this *User) GetAllMass(maxnum int) []model.Messages {
	var Re []model.Messages
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&timeout=5000ms",
			DB_username, DB_password, DB_host, DB_port, DB_database))
	if err != nil {
		fmt.Println("herr", err)
	}
	defer db.Close()
	for i := 1; i <= maxnum; i++ {
		var kv model.Messages
		err = db.QueryRow("SELECT message_id, speakername, content, time FROM messages where message_id = ? ", i).Scan(&kv.MessageId, &kv.SpeakerName, &kv.Content, &kv.Time)
		if err != nil {
			continue
		}
		Re = append(Re, kv)
	}
	return Re
}

func NewUser(id int, name string, grpcclientconn *grpc.ClientConn) *User {
	Sean := new(User)
	self := new(pb.User)
	self.Id = int64(id)
	self.Name = name
	Sean.Self = self
	cliC := pb.NewChatClient(grpcclientconn)
	Sean.ClientChat = cliC
	return Sean
}

func main() {
	conn1, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	fmt.Println(conn1.Target())
	if err != nil {
		fmt.Println("Dial error : ", err)
	}
	if conn1 == nil {
		fmt.Println("wrong")
		return
	}
	defer func(conn1 *grpc.ClientConn) {
		err := conn1.Close()
		if err != nil {
			fmt.Println("Failed to close the conn")
		}
	}(conn1)
	s := NewUser(12, "Atoi", conn1)
	fmt.Println(s.GetAllMass(2))
}
