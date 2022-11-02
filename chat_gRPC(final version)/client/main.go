package main

import (
	"bufio"
	pb "chat_gRPC/service/proto"
	"context"
	"database/sql"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
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

var (
	DB_database string = "sql_news"
	DB_username string = "sql_news"
	DB_password string = "666"
	DB_host     string = "121.36.131.50"
	DB_port     string = "3306"
)
var Port string = ":8080"

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
			this.GetAllMass()
		} else {
			continue
		}
	}
}
func (this *User) GetMassFromId(id int) *pb.Message {
	find := new(pb.Message)
	o := orm.NewOrm()
	qsu := o.QueryTable("messages")
	_, err := qsu.Filter("message_id", id).All(find)
	if err != nil {
		fmt.Println("[orm] Filter error : ", err)
		return nil
	}
	return find
}
func (this *User) GetMassFromName(speakername string) *pb.Message {
	find := new(pb.Message)
	o := orm.NewOrm()
	qsu := o.QueryTable("messages")
	_, err := qsu.Filter("speakername", speakername).All(find)
	if err != nil {
		fmt.Println("[orm] Filter error : ", err)
		return nil
	}
	return find
}
func (this *User) GetAllMass() {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&timeout=5000ms",
			DB_username, DB_password, DB_host, DB_port, DB_database))
	if err != nil {
		fmt.Println("herr", err)
	}
	defer db.Close()
	fmt.Println("-------------聊天框--------------")
	for i := 1; i < 10; i++ {
		var kv pb.Message
		err = db.QueryRow("SELECT message_id, speakername, content, time FROM messages where message_id = ? ", i).Scan(&kv.Id, &kv.Speakername, &kv.Content, &kv.Time)
		if err != nil {
			continue
		}
		fmt.Println("["+kv.Speakername+"]", " : ", kv.Content, " @", kv.Time)
	}
	fmt.Println("---------------------------------")
}
func main() {
	conn1, err := grpc.Dial("127.0.0.1"+Port, grpc.WithInsecure())
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
	cliC := pb.NewChatClient(conn1)
	fmt.Println(conn1.GetState().String())
	user := new(pb.User)
	user.Id = 1
	user.Name = "Sean"
	Sean := new(User)
	Sean.ClientChat = cliC
	Sean.Self = user
	Sean.SaytoAll()
}

func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Println("[orm] Register Driver error : ", err)
	}
	err = orm.RegisterDataBase("default", "mysql", "sql_news:666@tcp(121.36.131.50:3306)/sql_news?charset=utf8")
	if err != nil {
		fmt.Println("[orm] Register Data Base error : ", err)
	}
	/*
		err = orm.RunSyncdb("default", true, false)
		if err != nil {
			fmt.Println("[orm] Create Table error : ", err)
		}
	*/
}
