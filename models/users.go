package models

import (
	"errors"
	"log"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
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

func init() {
	orm.RegisterModel(new(S_user), new(UserAvatar))
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
