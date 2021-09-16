package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
)

type Session struct {
	Id        int
	Token     string    `json:"token"`
	IsValid   bool      `json:"is_valid"`
	ExpiresAt time.Time `json:"expires_at"`
	UserId    int       `json:"user_id"`
}

type SessionClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

func init() {
	orm.RegisterModel(new(Session))
}

func CreateToken(userId int) (*Session, error) {
	o := orm.NewOrm()

	var expAT = 24 * 30 * time.Hour
	key := os.Getenv("KEY_SUPERINDO")

	claims := SessionClaims{
		userId,
		jwt.StandardClaims{
			Issuer:    "api",
			ExpiresAt: time.Now().Add(expAT).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	data := Session{
		Token:     tokenString,
		IsValid:   true,
		ExpiresAt: time.Now().Add(expAT),
		UserId:    userId,
	}

	if _, err := o.Insert(&data); err != nil {
		return nil, err
	} else {
		return &data, nil
	}

}

func ValidateToken(tokenString string) (err error) {

	key := os.Getenv("KEY_LELANG")

	token, err := jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*SessionClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
		return nil
	}

	return errors.New("Claims or token null")
}
