package models

import (
  "github.com/astaxie/beego/orm"
  "errors"
  "time"
  "golang.org/x/crypto/bcrypt"
  "github.com/dgrijalva/jwt-go"
  "log"
  //"strconv"
)

type User struct {
  Id          int        
  Name        string     `orm:"unique" valid:"Required; MaxSize(100)"` 
  Password    string     `valid:"Required; MaxSize(100)"` 
  Email       string     `orm:"unique" valid:"Email; MaxSize(100)"`
  Avatar      string
  IsActive    bool       `orm:"default(true)"`
  IsAdmin     bool       `orm:"default(false)"`
  CreatedTime time.Time  `orm:"auto_now_add;type(datetime)"`
}

type CustomClaims struct {
  Id     int
  Name   string
  jwt.StandardClaims
}
var SigningKey = []byte("AllYourBase")

/*
type UserInfo struct {
  Id          int
  Name        string
  Email       string
  Avatar      string
  IsAdmin     bool
}
*/

func init() {
  orm.RegisterModel(new(User))
}

func GetAllUsers() []orm.Params {
  o := orm.NewOrm()
  var users []orm.Params
  o.QueryTable("user").Values(&users)
  return users
}

func GetUserById (id int) (u User, err error) {
  o := orm.NewOrm()
  var user User
  err = o.QueryTable("user").Filter("id", id).One(&user)
  if err == nil {
    return user, nil
	}
  return user, errors.New("User not exists")
}

func GetUserByName (name string) (u User, err error) {
  o := orm.NewOrm()
  var user User
  err = o.QueryTable("user").Filter("name", name).One(&user)
  if err == nil {
    return user, nil
	}
  return user, errors.New("User not exists")
}

func GetUserByEmail (email string) (u User, err error) {
  o := orm.NewOrm()
  var user User
  err = o.QueryTable("user").Filter("email", email).One(&user)
  if err == nil {
    return user, nil
	}
  return user, errors.New("User not exists")
}

func AddUser(item *User) (id int64, err error) {
  o := orm.NewOrm()

  password := item.Password
  saltedBytes := []byte(password)
  hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
  if err != nil {
		return -1, err
  }
  item.Password = string(hashedBytes[:])

  id, err = o.Insert(item)
  if err == nil {
      return id, nil
  }
  return -1, err
}

func DeleteUser(id int) (err error) {
  o := orm.NewOrm()
  user, errors := GetUserById(id)
  if errors != nil {
    return errors
  }

  if _, errors := o.Delete(&user); errors == nil {
      return nil
  }
  return errors
}

func UpdateUser(id int, item *User) (err error) {
  o := orm.NewOrm()
  user, errors := GetUserById(id)
  if errors != nil {
    return errors
  }

  item.Id = user.Id
  if _, errors := o.Update(item); errors == nil {
      return nil
  }
  return errors
}

func CheckPassword(username string, password string) (User, bool) {
  user, errors := GetUserByName(username)
  if errors != nil {
    return user, false
  }

  byteHash := []byte(user.Password)
  plainPwd := []byte(password)
  err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
  if err != nil {
      log.Println(err)
      return user, false
  }

  return user, true
}

func GenerateJwtToken(user User) string {
  claims := CustomClaims{
    user.Id,
    user.Name,
    jwt.StandardClaims{
      ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), //2 hours
      Issuer: "API",
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  ss, err := token.SignedString(SigningKey)
  log.Printf("%v %v %v", ss, err, time.Now().Add(time.Hour * 2).Unix())
  return ss
}
