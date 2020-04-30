// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
  "apigo/controllers"
  "github.com/astaxie/beego"
)

func init() {
  ns := beego.NewNamespace("/api",
    beego.NSNamespace("/user",
      beego.NSInclude(
        &controllers.UserController{},
      ),
    ),
    beego.NSNamespace("/topic",
      beego.NSInclude(
        &controllers.TopicController{},
      ),
    ),
    beego.NSNamespace("/reply",
      beego.NSInclude(
        &controllers.ReplyController{},
      ),
    ),
  )
  beego.AddNamespace(ns)

  beego.Router("/api/login", &controllers.LoginController{}, "post:Login")
  beego.Router("/api/authtoken", &controllers.AuthTokenController{}, "get:AuthToken")
  beego.Router("/api/topicstar/:objectId([0-9]+)", &controllers.TopicController{}, "put:AgreeCount")
  beego.Router("/api/topicview/:objectId([0-9]+)", &controllers.TopicController{}, "put:ViewCount")
}
