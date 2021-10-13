package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func TimeNow() (t string) {
	return time.Now().Format("02/01/2006 15:04:05")
}
