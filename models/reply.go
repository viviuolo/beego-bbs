package models

import (
	"github.com/astaxie/beego/orm"
  "time"
  "errors"
)

type Reply struct {
  Id      int       `orm:"pk;auto"`
  Topic   *Topic    `orm:"rel(fk)" valid:"Required;MaxSize(500)"`
  Content string    `orm:"type(text)" valid:"Require;MaxSize(1000)"`
  User    *User     `orm:"rel(fk)"`
  CreatedTime    time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
  orm.RegisterModel(new(Reply))
}

func GetAllReplies(topicId int, page int, pageSize int) ([]orm.Params, int64) {
  o := orm.NewOrm()
  var replies []orm.Params
  //o.QueryTable("reply").Values(&replies)
  o.QueryTable("reply").Filter("Topic", topicId).Limit(pageSize).Offset((page - 1) * pageSize).OrderBy("-CreatedTime").RelatedSel().Values(&replies, "id", "topic", "content", "user", "user__name", "created_time")

  cnt, err := o.QueryTable("reply").Filter("Topic", topicId).Count()
  if err != nil {
    cnt = 0
  }

  return replies, cnt
}

func GetReplyById (id int) (u Reply, err error) {
  o := orm.NewOrm()
  var reply Reply
  err = o.QueryTable("reply").Filter("id", id).One(&reply)
  if err == nil {
    return reply, nil
	}
  return reply, errors.New("Reply not exists")
}

func GetReplyByTopicId (id int) (u Reply, err error) {
  o := orm.NewOrm()
  var reply Reply
  err = o.QueryTable("reply").Filter("id", id).One(&reply)
  if err == nil {
    return reply, nil
	}
  return reply, errors.New("Reply not exists")
}

func AddReply(item *Reply) (id int64, err error) {
  o := orm.NewOrm()
  id, err = o.Insert(item)
  if err == nil {
      return id, nil
  }
  return 0, err
}

func DeleteReply(id int) (err error) {
  o := orm.NewOrm()
  reply, errors := GetReplyById(id)
  if errors != nil {
    return errors
  }

  if _, errors := o.Delete(&reply); errors == nil {
      return nil
  }
  return errors
}

func UpdateReply(id int, item *Reply) (err error) {
  o := orm.NewOrm()
  reply, errors := GetReplyById(id)
  if errors != nil {
    return errors
  }

  item.Id = reply.Id
  if _, errors := o.Update(item); errors == nil {
      return nil
  }
  return errors
}
