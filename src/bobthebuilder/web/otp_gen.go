package web

import (
  "github.com/pquerna/otp/totp"
  "github.com/pquerna/otp"
  "bobthebuilder/logging"

	"bytes"
	"image/png"
  "encoding/base64"
  "github.com/hoisie/web"
)

type OTPPageModel struct {
  QR_DATA string
  Key *otp.Key
}

func OTPUtilHandler(ctx *web.Context) {
  key, err := totp.Generate(totp.GenerateOpts{
  		Issuer: ctx.Params["issuer"],
  		AccountName: ctx.Params["account"],
  })

  if err != nil {
    ctx.ResponseWriter.Write([]byte("Err: " + err.Error()))
    return
  }

  // Convert TOTP key into a QR code encoded as a PNG image.
  var buf bytes.Buffer
  img, err := key.Image(200, 200)
  if err != nil {
    ctx.ResponseWriter.Write([]byte("Err: " + err.Error()))
    return
  }

  png.Encode(&buf, img)
  QR_string := base64.StdEncoding.EncodeToString(buf.Bytes())

  t := templates.Lookup("otpgen")
  if t == nil {
    logging.Error("web", "No template found.")
    return
  }

  err = t.Execute(ctx.ResponseWriter, OTPPageModel{QR_DATA: QR_string, Key: key})
  if err != nil{
    logging.Error("views-otp", err)
  }
}
