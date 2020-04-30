package controllers

import (
  "apigo/models"
  "encoding/json"
  "github.com/astaxie/beego/validation"
  "strconv"
  "log"
)

// Operations about object
type TopicController struct {
  BaseController
}

func (this *TopicController) NestPrepare() {
  method := this.Ctx.Request.Method
  // No limit for get method
  if method == "GET" {
    return
  }

  if this.isLogin == false {
    this.Data["json"] = "Auth failed"
    this.Ctx.Output.SetStatus(401)
    this.ServeJSON()
    return
  }
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (o *TopicController) Post() {
  var topic models.Topic
  var err error
  if err = json.Unmarshal(o.Ctx.Input.RequestBody, &topic); err == nil {
    valid := validation.Validation{}
    valid.Valid(&topic)
    if valid.HasErrors() {
      for _, err := range valid.Errors {
        log.Println(err.Key, err.Message)
      }
      o.Data["json"] = valid.Errors;
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }

    topic.User = &o.user
    uid, err := models.AddTopic(&topic)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
    } else {
      o.Data["json"] = map[string]int64{"uid": uid}
    }
  } else {
    o.Ctx.Output.SetStatus(400)
    o.Data["json"] = err.Error()
  }
  o.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *TopicController) Get() {
  objectId := o.Ctx.Input.Param(":objectId")
  if objectId != "" {
    id, _ := strconv.Atoi(objectId)
    topic := models.GetTopicListById(id)
    o.Data["json"] = topic[0]
  }
  o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *TopicController) GetAll() {
  page := o.GetString("page")
  pageSize := o.GetString("page_size")

  var p int
  var size int
  if (page == "" || pageSize == "") {
    p = 1
    size = 10
  } else {
    p ,_ = strconv.Atoi(page)
    size ,_ = strconv.Atoi(pageSize)
  }

  topics, count := models.GetAllTopics(p, size)
  o.Data["json"] = map[string]interface{}{"topics":topics, "count":count}
  o.ServeJSON()
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *TopicController) Put() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if o.user.IsAdmin == false {
    t, errors := models.GetTopicById(id)
    if (errors != nil || t.User.Id != o.user.Id) {
      o.Data["json"] = "You do not have permission to perform this action."
      o.Ctx.Output.SetStatus(401)
      o.ServeJSON()
      return
    }
  }

  var topic models.Topic
  json.Unmarshal(o.Ctx.Input.RequestBody, &topic)
  err := models.UpdateTopic(id, &topic)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Update success!"
  }
  o.ServeJSON()
}

func (o *TopicController) AgreeCount() {
  reqType := o.GetString("type")
  if (reqType != "add" && reqType != "decline") {
    o.Data["json"] = "Bad request"
    o.Ctx.Output.SetStatus(400)
    o.ServeJSON()
    return
  }

  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)
  topic, _ := models.GetTopicById(id)

  if reqType == "add" {
    topic.AgreeCount += 1
  } else {
    if topic.AgreeCount <= 0 {
      o.Data["json"] = "Success"
      o.ServeJSON()
      return
    }
    topic.AgreeCount -= 1
  }

  err := models.UpdateAgreeCount(id, &topic)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Update success!"
  }
  o.ServeJSON()
}

func (o *TopicController) ViewCount() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)
  err := models.AddViewCount(id)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Update success!"
  }
  o.ServeJSON()
}

func (o *TopicController) ReplyCount() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)
  err := models.AddReplyCount(id)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Update success!"
  }
  o.ServeJSON()
}

// @Title Delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (o *TopicController) Delete() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if o.user.IsAdmin == false {
    t, errors := models.GetTopicById(id)
    if (errors != nil || t.User.Id != o.user.Id) {
      o.Data["json"] = "You do not have permission to perform this action."
      o.Ctx.Output.SetStatus(401)
      o.ServeJSON()
      return
    }
  }

  err := models.DeleteTopic(id)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Delete success!"
  }
  o.ServeJSON()
}

