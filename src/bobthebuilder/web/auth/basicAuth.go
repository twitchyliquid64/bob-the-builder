package auth

import (
  "errors"
  "github.com/hoisie/web"
)

var ErrNoLoginForBasicAuth = errors.New("Logging in is not a valid operation for Basic Authentication")

func BasicAuth(src UserSource)Auther{
  return &BasicAuther{src}
}

type BasicAuther struct {
  userSource UserSource
}

func (d *BasicAuther) AuthInfo(ctx *web.Context)(*AuthInfo, error) {
  username, pwd, ok := ctx.Request.BasicAuth()
	if !ok {
		return nil, ErrNotAuthenticated
	}

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

  return &AuthInfo{User: usr}, nil
}



func (d *BasicAuther) DoLogin(ctx *web.Context)(*AuthInfo, error){
  return nil, ErrNoLoginForBasicAuth
}
