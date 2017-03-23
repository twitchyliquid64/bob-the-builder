package auth

import "github.com/hoisie/web"

func MultiAuth(auths ...Auther)Auther{
  return &MultiAuther{auths}
}

type MultiAuther struct {
  methods []Auther
}

func (d *MultiAuther) AuthInfo(ctx *web.Context)(*AuthInfo, error) {
  for _, method := range d.methods {
    info, err := method.AuthInfo(ctx)
    if err == nil && info != nil {
      return info, err
    }
    if err != ErrNotAuthenticated {
      return nil, err
    }
  }

  return nil, ErrNotAuthenticated
}

func (d *MultiAuther) DoLogin(ctx *web.Context)(*AuthInfo, error){
  for _, method := range d.methods {
    info, err := method.AuthInfo(ctx)
    if err != nil && info != nil {
      return info, err
    }
  }
  return nil, ErrNotAuthenticated
}
