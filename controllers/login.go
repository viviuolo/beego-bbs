package controllers

import (
  "github.com/astaxie/beego"
  "encoding/json"
  "apigo/models"
)

type LoginController struct {
  beego.Controller
}

type AuthTokenController struct {
  BaseController
}

type TokenResponse struct {
  Id     int
  Name   string
  Token  string
}

func (this *LoginController) Login() {
  var ob models.User
  var err error
  if err = json.Unmarshal(this.Ctx.Input.RequestBody, &ob); err != nil {
    this.Data["json"] = err.Error()
    this.Ctx.Output.SetStatus(400)
    this.ServeJSON()
    return
  }

  username := ob.Name
  password := ob.Password
  if (username == "" || password == "" ) {
    this.Data["json"] = "Invalid parameter!"
    this.Ctx.Output.SetStatus(401)
    this.ServeJSON()
    return
  }

  user, isLogin := models.CheckPassword(username, password)
  if isLogin != true {
    this.Data["json"] = "Auth failed!"
    this.Ctx.Output.SetStatus(401)
  } else {
    // Return JWT token
    jwtToken := models.GenerateJwtToken(user)
    this.Data["json"] = TokenResponse{ user.Id, user.Name, jwtToken }
  }
  this.ServeJSON()
}

func (this *AuthTokenController) AuthToken() {
  if this.isLogin == false {
    this.Data["json"] = "Auth failed"
    this.Ctx.Output.SetStatus(401)
  } else {
    jwtToken := models.GenerateJwtToken(this.user)
    this.Data["json"] = TokenResponse{ this.user.Id, this.user.Name, jwtToken }
  }
  this.ServeJSON()
}
