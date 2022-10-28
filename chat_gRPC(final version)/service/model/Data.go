package model

import "github.com/beego/beego/v2/client/orm"

type Messages struct {
	MessageId   int    `orm:"auto" orm:"column(messageid)"`
	SpeakerName string `orm:"column(speakername)"`
	Content     string `orm:"column(content)"`
	Time        string `orm:"column(time)"`
}
type Users struct {
	Id   int64  `orm:"cloumn(id)"`
	Name string `orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(Messages), new(Users))
}
