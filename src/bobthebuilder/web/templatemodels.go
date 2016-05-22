package web

import (
  "bobthebuilder/builder"
  "bobthebuilder/config"
)

type modelBasic struct{
  Config *config.Config
  Builder *builder.Builder

  Opt1 string
}
