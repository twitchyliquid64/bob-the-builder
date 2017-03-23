package auth

import (
  "errors"
  "time"
  "github.com/hoisie/web"
  "net/http"
  "math/rand"
)

var ErrCantCheckCookieUserPassword = errors.New("Cannot check the password of a user authed by cookie")

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}


type session struct {
  Username string
  Created time.Time
}

type userCookie struct{
  Username string
  Created time.Time
}

func (u *userCookie) Name() string{
  return u.Username
}

func (u *userCookie) CheckPassword(pass string)(bool, error) {
  return false, ErrCantCheckCookieUserPassword
}


func CookieAuth(src UserSource)Auther{
  return &CookieAuther{userSource: src, sessions: map[string]session{}}
}

type CookieAuther struct {
  userSource UserSource
  sessions map[string]session
}

func (d *CookieAuther) AuthInfo(ctx *web.Context)(*AuthInfo, error) {
  cookie, err := ctx.Request.Cookie("sass")
  if err != nil {
    if err == http.ErrNoCookie {
      return nil, ErrNotAuthenticated
    }
    return nil, err
  }

  session, ok := d.sessions[cookie.Value]
	if !ok {
		return nil, ErrNotAuthenticated
	}

  return &AuthInfo{User: &userCookie{Username: session.Username, Created: session.Created}}, nil
}



func (d *CookieAuther) DoLogin(ctx *web.Context)(*AuthInfo, error){
  username := ctx.Request.FormValue("username")
  pwd := ctx.Request.FormValue("password")


  usr, lookupErr := d.userSource.GetUserByUsername(username)
  if lookupErr == ErrUserDoesntExist{
    return nil, ErrNotAuthenticated
  } else if lookupErr != nil {
    return nil, lookupErr
  }

  passwordOk, passwdErr := usr.CheckPassword(pwd)
  if passwdErr != nil{
    return nil, passwdErr
  }
  if !passwordOk{
    return nil, ErrNotAuthenticated
  }

  //make session
  SID := RandStringRunes(22)
  d.sessions[SID] = session{Username: username, Created: time.Now()}
  ctx.SetCookie(&http.Cookie{Name: "sass", Value: SID})

  return &AuthInfo{User: usr}, nil
}
