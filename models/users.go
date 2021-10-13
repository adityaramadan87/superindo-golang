package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"superindo_backend/helper"

	"github.com/beego/beego/v2/client/orm"
	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/bcrypt"
)

type S_user struct {
	Id           int    `form:"id" json:"id" pg:"id"`
	Name         string `form:"name" json:"name" pg:"name"`
	Email        string `form:"email" json:"email"`
	Phone        string `form:"phone" json:"phone"`
	Password     string `form:"password" json:"password"`
	Verified     bool   `form:"verified" json:"verified"`
	Role         string `form:"role" json:"role"`
	BirthDate    string `form:"birth_date" json:"birth_date"`
	CreateDate   string `form:"create_date" json:"create_date"`
	VerifiedDate string `form:"verified_date" json:"verified_date"`
	Avatar       string `form:"avatar" json:"avatar"`
}

type UserAvatar struct {
	Id     int    `json:"id"`
	Avatar string `json:"avatar"`
	UserId int    `json:"user_id"`
}

type UserOtp struct {
	Id        int    `json:"id"`
	OtpNumber string `json:"otp_number"`
}

type MessageData struct {
	To      string            `json:"to"`
	From    string            `json:"from"`
	Type    string            `json:"type"`
	Content map[string]string `json:"content"`
}

func init() {
	orm.RegisterModel(new(S_user), new(UserAvatar), new(UserOtp))
}

func Register(u S_user) (us *S_user, msg string, rc int) {
	o := orm.NewOrm()

	var userAvailable S_user
	o.Raw("SELECT * FROM s_user WHERE phone = ?", u.Phone).QueryRow(&userAvailable)
	if userAvailable != (S_user{}) {
		return nil, "number phone already exists", 1
	}

	o.Raw("SELECT * FROM s_user WHERE email = ?", u.Email).QueryRow(&userAvailable)
	if userAvailable != (S_user{}) {
		return nil, "email already exists", 1
	}

	id, error := o.Insert(&u)
	if error != nil {
		return nil, error.Error(), 1
	}

	var user S_user
	errQuery := o.Raw("SELECT * FROM s_user WHERE id = ?", id).QueryRow(&user)
	if errQuery != nil {
		return nil, "Error while query " + errQuery.Error(), 1
	}

	return &user, "Success", 0
}

func LoginUser(email string, password string) (rc int, msg string, ss *Session) {
	o := orm.NewOrm()

	log.Print("he " + email + "  " + password)

	var users S_user
	if err := o.Raw("SELECT * FROM s_user WHERE email = ?", email).QueryRow(&users); err != nil {
		log.Print(errors.New("Error while query " + err.Error()))
		return 1, "Error while query " + err.Error(), nil
	}

	if users == (S_user{}) {
		return 1, "user not found please make sure your email is correct", nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))
	if err != nil {
		log.Print(err)
		return 1, "Password not match", nil
	} else {

		if ss, err := CreateToken(users.Id); err != nil {
			return 1, err.Error(), nil
		} else {
			return 0, "Success", ss
		}
	}

}

func UserPhoneVerification(phoneNumber string) (rc int, msg string) {
	o := orm.NewOrm()

	phoneNmbr, _ := strconv.Atoi(phoneNumber)

	otp := helper.EncodeToString(6)

	var userOtp UserOtp
	o.Raw("SELECT * FROM user_otp WHERE otp_number = ?", otp).QueryRow(&userOtp)
	if userOtp != (UserOtp{}) {
		fmt.Println("Same Otp : ", otp)
		UserPhoneVerification(phoneNumber)
	}

	data := MessageData{
		To:   "+62" + strconv.Itoa(phoneNmbr),
		From: "36ac2e2092b84ae2a081d86acd83ffc6",
		Type: "text",
		Content: map[string]string{
			"text": "Hallo berhasil hit api otp : " + otp,
		},
	}

	client := resty.New()

	if _, err := o.Insert(&UserOtp{
		OtpNumber: otp,
	}); err != nil {
		fmt.Println(err.Error())

		return 1, "Error while insert otp"
	}

	response, err := client.R().
		SetHeaders(map[string]string{
			"Authorization": "AccessKey YOUR_API_KEY",
			"Content-Type":  "application/json",
		}).
		SetBody(&data).
		Post("https://conversations.messagebird.com/v1/send")

	if err != nil {
		log.Print(errors.New("Error while send otp " + err.Error()))
		return 1, "error while send Otp"
	} else {
		fmt.Println(response, err)
		if response.StatusCode() == 202 {
			var marshal map[string]interface{}

			json.Unmarshal(response.Body(), &marshal)
			fmt.Println("ID : ", marshal["id"])

			return 0, "Success"
		}

		return 1, "Failed Status Code : " + strconv.Itoa(response.StatusCode())
	}
}

func UserById(userId int) (rc int, msg string, su *S_user) {
	o := orm.NewOrm()

	var users S_user
	if err := o.Raw("SELECT * FROM s_user WHERE id = ?", userId).QueryRow(&users); err != nil {
		log.Print(errors.New("Error while query " + err.Error()))
		return 1, "Error while query " + err.Error(), nil
	}

	return 0, "Success", &users

}

func UpdateUsers(users *S_user) (data *S_user, msg string, rc int) {
	if strconv.Itoa(users.Id) != "" {
		o := orm.NewOrm()

		var us S_user
		if er := o.Raw("SELECT * FROM s_user WHERE id = ?", users.Id).QueryRow(&us); er != nil {
			log.Print("Error while query " + er.Error())
		}

		if us == (S_user{}) {
			return nil, "Data with id " + strconv.Itoa(users.Id) + " not found", 1
		}

		if users.Avatar != "" {
			var userAvatar UserAvatar
			if er := o.Raw("SELECT * FROM user_avatar WHERE user_id = ?", us.Id).QueryRow(&userAvatar); er == nil {
				userAvatar.Avatar = users.Avatar
				if id, err := o.Update(&userAvatar); err == nil {
					us.Avatar = "http://localhost:8080/users/avatar/" + strconv.Itoa(int(id))
				} else {
					return nil, "Error while update user avatar " + err.Error(), 1
				}
			} else {
				userAvatar.UserId = users.Id
				userAvatar.Avatar = users.Avatar

				if id, err := o.Insert(&userAvatar); err == nil {
					us.Avatar = "http://localhost:8080/users/avatar/" + strconv.Itoa(int(id))
				} else {
					return nil, "Error while insert user avatar " + err.Error(), 1
				}
			}
		}
		if users.Name != "" {
			us.Name = users.Name
		}
		if users.Email != "" {
			us.Email = users.Email
		}

		_, erro := o.Update(&us)
		if erro != nil {
			return nil, "Error when Update data " + erro.Error(), 1
		}

		var userDone S_user
		if err := o.Raw("SELECT * FROM s_user WHERE id = ?", users.Id).QueryRow(&userDone); err != nil {
			return nil, "Error when select user " + err.Error(), 1
		}

		return &userDone, "Success", 0

	}

	return nil, "Id not null", 1

}
