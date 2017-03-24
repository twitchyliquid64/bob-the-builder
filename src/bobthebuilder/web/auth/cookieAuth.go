package auth

import (
  "errors"
  "time"
  "github.com/hoisie/web"
  "net/http"
  "crypto/rand"
  "sync"
)

const (
  SESSION_ID_CHAR_LENGTH = 48
  EXPIRY_TIME = time.Hour * 6
)

var ErrCantCheckCookieUserPassword = errors.New("Cannot check the password of a user authed by cookie")
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]byte, n)
    o := make([]rune, n)
    rand.Read(b)
    for i := range b {
        o[i] = letterRunes[int(b[i]) % len(letterRunes)]
    }
    return string(o)
}


type session struct {
  Username string
  Created time.Time
  Expiry time.Time
}

func (u *session) Name() string{
  return u.Username
}

func (u *session) CheckPassword(pass string)(bool, error) {
  return false, ErrCantCheckCookieUserPassword
}


func CookieAuth(src UserSource)Auther{
  return &CookieAuther{userSource: src, sessions: map[string]session{}}
}

type CookieAuther struct {
  userSource UserSource
  sessions map[string]session
  lock sync.Mutex
}

func (d *CookieAuther) AuthInfo(ctx *web.Context)(*AuthInfo, error) {
  d.lock.Lock()
  defer d.lock.Unlock()

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

  return &AuthInfo{User: &session}, nil
}


func (d *CookieAuther) cleanupExpired(){
  //assumes lock already held
  toDeleteSessions := []string{}
  for SID, session := range d.sessions {
    if time.Now().After(session.Expiry) {
      toDeleteSessions = append(toDeleteSessions, SID)
    }
  }

  for _, SID := range toDeleteSessions {
    delete(d.sessions, SID)
  }
}


func (d *CookieAuther) DoLogin(ctx *web.Context)(*AuthInfo, error){
  d.lock.Lock()
  defer d.lock.Unlock()
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
  expiry := time.Now().Add(EXPIRY_TIME)
  SID := RandStringRunes(SESSION_ID_CHAR_LENGTH)
  d.sessions[SID] = session{Username: username, Created: time.Now(), Expiry: expiry}
  ctx.SetCookie(&http.Cookie{Name: "sass", Value: SID, Expires: expiry, HttpOnly: true})

  d.cleanupExpired()
  return &AuthInfo{User: usr}, nil
}
