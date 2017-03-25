package web

import (
  "bobthebuilder/builder"
  "bobthebuilder/config"
  "bobthebuilder/web/auth"
)

type modelBasic struct{
  Config *config.Config
  Builder *builder.Builder

  Opt1 string
  Auth *auth.AuthInfo
}
