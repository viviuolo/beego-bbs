package controllers

import (
  "github.com/astaxie/beego"
  "apigo/models"
  "github.com/dgrijalva/jwt-go"
)

type NestPreparer interface {
  NestPrepare()
}

type BaseController struct {
  beego.Controller
  user     models.User
  isLogin  bool
}

func (this *BaseController) Prepare() {
  tokenString := this.Ctx.Input.Header("Authorization")
  token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return models.SigningKey, nil
  })

  if err == nil {
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
      id := claims.Id
      name := claims.Name

      user, errors := models.GetUserById(id)
      if (errors != nil ||  user.Name != name) {
        this.isLogin = false
      } else {
        this.user = user
        this.isLogin = true
      }
    } else {
      this.isLogin = false
    }
  } else {
    this.isLogin = false
  }

  if app, ok := this.AppController.(NestPreparer); ok {
    app.NestPrepare()
  }
}
