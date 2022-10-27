package models

import "github.com/beego/beego/v2/client/orm"

type MessageInfor struct {
	MessageId   int    `orm:"auto" orm:"column(messageid)"`
	SpeakerId   int64  `orm:"column(speakerid)"`
	SpeakerName string `orm:"column(speakername)"`
	Content     string `orm:"column(content)"`
}
type Users struct {
	Id   int64  `orm:"cloumn(id)"`
	Name string `orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(MessageInfor), new(Users))
}
