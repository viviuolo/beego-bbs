package models

import (
  "github.com/astaxie/beego/orm"
  "errors"
  "time"
)

type Topic struct {
  Id             int       `orm:"pk;auto"`
  Title          string    `orm:"unique" valid:"Required; MaxSize(500)"`
  Content        string    `orm:"type(text);null" valid:"MaxSize(1000)"`
  User           *User     `orm:"rel(fk)"`
  AgreeCount     int       `orm:"default(0)"`
  ViewCount      int       `orm:"default(0)"`
  ReplyCount     int       `orm:"default(0)"`
  LastReplyUser  *User     `orm:"rel(fk);null"`
  LastReplyTime  time.Time `orm:"auto_now_add;type(datetime)"`
  CreatedTime    time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
  orm.RegisterModel(new(Topic))
}

func GetAllTopics(page int, pageSize int)( []orm.Params, int64) {
  o := orm.NewOrm()
  var topics []orm.Params
  o.QueryTable("topic").Limit(pageSize).Offset((page - 1) * pageSize).OrderBy("-CreatedTime").RelatedSel().Values(&topics, "id", "title", "content", "user", "user__name", "agree_count", "view_count", "reply_count", "created_time")

  cnt, err := o.QueryTable("topic").Count()
  if err != nil {
    cnt = 0
  }

  return topics, cnt
}

func GetTopicListById (id int) (u []orm.Params) {
  o := orm.NewOrm()
  var topic []orm.Params
  o.QueryTable("topic").Filter("id", id).RelatedSel().Values(&topic, "id", "title", "content", "user", "user__name", "agree_count", "view_count", "reply_count", "created_time")
  return topic
}

func GetTopicById (id int) (u Topic, err error) {
  o := orm.NewOrm()
  var topic Topic
  err = o.QueryTable("topic").Filter("id", id).RelatedSel().One(&topic)
  if err == nil {
    return topic, nil
	}
  return topic, errors.New("Topic not exists")
}

func AddTopic(item *Topic) (id int64, err error) {
  o := orm.NewOrm()
  id, err = o.Insert(item)
  if err == nil {
      return id, nil
  }
  return 0, err
}

func DeleteTopic(id int) (err error) {
  o := orm.NewOrm()
  topic, errors := GetTopicById(id)
  if errors != nil {
    return errors
  }

  if _, errors := o.Delete(&topic); errors == nil {
      return nil
  }
  return errors
}

func UpdateTopic(id int, item *Topic) (err error) {
  o := orm.NewOrm()
  item.Id = id
  if _, err := o.Update(item, "Title", "Content"); err == nil {
    return nil
  }
  return err
}

func UpdateAgreeCount(id int, item *Topic) (err error) {
  o := orm.NewOrm()
  item.Id = id
  if _, err := o.Update(item, "AgreeCount"); err == nil {
    return nil
  }
  return err
}

func AddViewCount(id int) (err error) {
  o := orm.NewOrm()
  _, err = o.QueryTable("topic").Filter("id", id).Update(orm.Params{
    "ViewCount": orm.ColValue(orm.ColAdd, 1),
  })
  return err
}

func AddReplyCount(id int) (err error) {
  o := orm.NewOrm()
  _, err = o.QueryTable("topic").Filter("id", id).Update(orm.Params{
    "ReplyCount": orm.ColValue(orm.ColAdd, 1),
  })
  return err
}