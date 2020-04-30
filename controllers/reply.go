package controllers

import (
  "apigo/models"
  "encoding/json"
  "github.com/astaxie/beego/validation"
  "strconv"
)

// Operations about object
type ReplyController struct {
  BaseController
}

type RequestReply struct {
  Topic   int
  Content string
}


func (this *ReplyController) NestPrepare() {
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
func (o *ReplyController) Post() {
  var reqReply RequestReply
  var reply models.Reply
  var err error
  if err = json.Unmarshal(o.Ctx.Input.RequestBody, &reqReply); err == nil {
    valid := validation.Validation{}
    valid.Valid(&reqReply)
    if valid.HasErrors() {
      o.Data["json"] = valid.Errors;
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }
  
    // Get topic
    topicId := reqReply.Topic
    topic, err := models.GetTopicById(topicId)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }

    reply.User = &o.user
    reply.Content = reqReply.Content
    reply.Topic = &topic
    uid, err := models.AddReply(&reply)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
    } else {
      o.Data["json"] = map[string]int64{"uid": uid}
      models.AddReplyCount(topicId)
    }
  } else {
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
func (o *ReplyController) Get() {
  objectId := o.Ctx.Input.Param(":objectId")
  if objectId != "" {
    id, _ := strconv.Atoi(objectId)
    reply, err := models.GetReplyById(id)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
    } else {
      o.Data["json"] = reply
    }
  }
  o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *ReplyController) GetAll() {
  objectId := o.GetString("topic_id")
  if (objectId == "") {
    o.Data["json"] = "Invalid parameter."
    o.Ctx.Output.SetStatus(400)
    o.ServeJSON()
    return
  }
  topicId, _ := strconv.Atoi(objectId)

  page := o.GetString("page")
  pageSize := o.GetString("page_size")

  var p int
  var size int
  if (page == "" || pageSize == "") {
    p = 1
    size = 10
  } else {
    p, _ = strconv.Atoi(page)
    size, _ = strconv.Atoi(pageSize)
  }

  replies, count := models.GetAllReplies(topicId, p, size)
  o.Data["json"] =  map[string]interface{}{"replies":replies, "count":count}
  o.ServeJSON()
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *ReplyController) Put() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if o.user.IsAdmin == false {
    reply, errors := models.GetReplyById(id)
    if (errors != nil || reply.User.Id != o.user.Id) {
      o.Data["json"] = "You do not have permission to perform this action."
      o.Ctx.Output.SetStatus(401)
      o.ServeJSON()
      return
    }
  }

  var reply models.Reply
  json.Unmarshal(o.Ctx.Input.RequestBody, &reply)

  err := models.UpdateReply(id, &reply)
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
func (o *ReplyController) Delete() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if o.user.IsAdmin == false {
    reply, errors := models.GetReplyById(id)
    if (errors != nil || reply.User.Id != o.user.Id) {
      o.Data["json"] = "You do not have permission to perform this action."
      o.Ctx.Output.SetStatus(401)
      o.ServeJSON()
      return
    }
  }

  err := models.DeleteReply(id)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Delete success!"
  }
  o.ServeJSON()
}

