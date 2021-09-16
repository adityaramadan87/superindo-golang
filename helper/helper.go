package helper

import (
	"encoding/base64"
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type Helper struct {
}

type Res struct {
	Rc   int         `json:"rc"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func Response(rc int, msg string, data interface{}, u beego.Controller) {
	u.Data["json"] = Res{rc, msg, data}
	u.ServeJSON()
}

func ImageResponse(data []byte, c beego.Controller) {
	c.Data["json"] = &data
	c.ServeJSON()
}

func EncodeBase64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}

func DecodeBase64(message string) (retour []byte) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	i, _ := base64.StdEncoding.DecodeString(message)
	fmt.Printf("base64: %s\n", i)
	return base64Text
}

func TimeNow() (t string) {
	return time.Now().Format("02/01/2006 15:04:05")
}
