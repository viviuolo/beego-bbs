package main

import (
  "github.com/astaxie/beego"
  "github.com/astaxie/beego/orm"
  "log"
  "fmt"
  _ "github.com/lib/pq"
  _ "apigo/routers"
)

func init() {
  host := beego.AppConfig.String("sqlhost")
  port, _ := beego.AppConfig.Int("sqlport")
  user := beego.AppConfig.String("sqluser")
  password := beego.AppConfig.String("sqlpass")
  dbname := beego.AppConfig.String("sqldbname")
  dbConfig := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

  // Using postgresql, todo: you can replace it
  orm.RegisterDriver("postgres", orm.DRPostgres)
  orm.RegisterDataBase("default", "postgres", dbConfig)

  // Create tables
  //orm.RunSyncdb("default", false, true)
  // Database alias.
  name := "default"

  // Drop table and re-create.
  force := false

  // Print log.
  verbose := true

  // Error.
  err := orm.RunSyncdb(name, force, verbose)
  if err != nil {
    log.Println(err)
  }

  orm.Debug = true
}

func main() {
  if beego.BConfig.RunMode == "dev" {
    beego.BConfig.WebConfig.DirectoryIndex = true
    beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
  }

  beego.Run()
}
