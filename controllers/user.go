package controllers

import (
  "apigo/models"
  "encoding/json"
  "github.com/astaxie/beego/validation"
  "strconv"
)

// Operations about object
type UserController struct {
  BaseController
}

func (this *UserController) NestPrepare() {
  method := this.Ctx.Request.Method
  // No limit for post method
  if method == "POST" {
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
func (o *UserController) Post() {
  var user models.User
  var err error
  if err = json.Unmarshal(o.Ctx.Input.RequestBody, &user); err == nil {
    // Check name, email, password format
    valid := validation.Validation{}
    valid.Valid(&user)
    if valid.HasErrors() {
      o.Data["json"] = valid.Errors;
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }

    // Check name exists or not
    _, errName := models.GetUserByName(user.Name)
    if errName == nil {
      o.Data["json"] = "This name already exists";
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }

    // Check email exists or not
    _, errEmail := models.GetUserByEmail(user.Email)
    if errEmail == nil {
      o.Data["json"] = "Email Address is Already Registered";
      o.Ctx.Output.SetStatus(400)
      o.ServeJSON()
      return
    }

    // Add into db
    uid, err := models.AddUser(&user)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
    } else {
      //o.Data["json"] = user
      id := int(uid)
      user, _:= models.GetUserById(id)
      jwtToken := models.GenerateJwtToken(user)
      o.Data["json"] = TokenResponse{ user.Id, user.Name, jwtToken }
    }
  } else {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  }
  o.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *UserController) Get() {
  if o.user.IsAdmin == false {
    o.Data["json"] = "You do not have permission to perform this action."
    o.Ctx.Output.SetStatus(401)
    o.ServeJSON()
    return
  }

  objectId := o.Ctx.Input.Param(":objectId")
  if objectId != "" {
    id, _ := strconv.Atoi(objectId)
    user, err := models.GetUserById(id)
    if err != nil {
      o.Data["json"] = err.Error()
      o.Ctx.Output.SetStatus(400)
    } else {
      o.Data["json"] = user
    }
  }
  o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *UserController) GetAll() {
  // TODO: pagination
  //.OrderBy("-InTime").Limit(size).Offset((p - 1) * size).All(&list)
  if o.user.IsAdmin == false {
    o.Data["json"] = "You do not have permission to perform this action."
    o.Ctx.Output.SetStatus(401)
    o.ServeJSON()
    return
  }

  users := models.GetAllUsers()
  o.Data["json"] = users
  o.ServeJSON()
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *UserController) Put() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if (o.user.IsAdmin == false && o.user.Id != id) {
    o.Data["json"] = "You do not have permission to perform this action."
    o.Ctx.Output.SetStatus(401)
    o.ServeJSON()
    return
  }

  var user models.User
  json.Unmarshal(o.Ctx.Input.RequestBody, &user)
  err := models.UpdateUser(id, &user)
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
func (o *UserController) Delete() {
  objectId := o.Ctx.Input.Param(":objectId")
  id, _ := strconv.Atoi(objectId)

  if (o.user.IsAdmin == false && o.user.Id != id) {
    o.Data["json"] = "You do not have permission to perform this action."
    o.Ctx.Output.SetStatus(401)
    o.ServeJSON()
    return
  }

  err := models.DeleteUser(id)
  if err != nil {
    o.Data["json"] = err.Error()
    o.Ctx.Output.SetStatus(400)
  } else {
    o.Data["json"] = "Delete success!"
  }
  o.ServeJSON()
}
