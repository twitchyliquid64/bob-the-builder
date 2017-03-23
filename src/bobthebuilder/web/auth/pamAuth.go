package auth

import (
  "os/exec"
  "strings"
  "errors"
  "github.com/hoisie/web"
)

var ErrNoLoginForPAMAuth = errors.New("Logging in is not a valid operation for Basic Authentication")

type userPAM struct{
  Username string
}

func (u *userPAM) Name() string{
  return u.Username
}

func (u *userPAM) CheckPassword(pass string)(bool, error) {
  out, err := exec.Command("python", "auth/pam-auth.py", u.Username, pass).Output()
  if err != nil {
    return false, err
  }
  if strings.HasPrefix(string(out), "OK") {
    return true, nil
  }
  return false, nil
}

type PAMAuther struct {
}

func (d *PAMAuther) AuthInfo(ctx *web.Context)(*AuthInfo, error) {
  username, pwd, ok := ctx.Request.BasicAuth()
	if !ok {
		return nil, ErrNotAuthenticated
	}
  if ok, err := (&userPAM{username}).CheckPassword(pwd); ok && err == nil {
    return &AuthInfo{User: &userPAM{username}}, nil
  } else if err != nil {
    return nil, err
  }
  return nil, ErrNotAuthenticated
}

func (d *PAMAuther) DoLogin(ctx *web.Context)(*AuthInfo, error){
  return nil, ErrNoLoginForPAMAuth
}
