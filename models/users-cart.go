package models

import (
	"log"

	"github.com/beego/beego/v2/client/orm"
)

type UsersCart struct {
	Id        int `form:"id" json:"id"`
	UserId    int `form:"user_id" json:"user_id"`
	ProductId int `form:"product_id" json:"product_id"`
}

type ListCart struct {
	Qty       int    `json:"qty"`
	UserId    int    `json:"user_id"`
	ProductId int    `json:"product_id"`
	ItemName  string `json:"item_name"`
	Price     int    `json:"price"`
	Picture   string `json:"picture"`
}

func init() {
	orm.RegisterModel(new(UsersCart))
}

func AddCart(usersCart UsersCart) (rc int, msg string) {
	o := orm.NewOrm()

	if _, err := o.Insert(&usersCart); err != nil {
		log.Print(err.Error())
		return 1, "Error while insert cart"
	}

	return 0, "Success"
}

func GetListCart(userId int) (rc int, msg string, data []ListCart) {
	o := orm.NewOrm()

	_, err := o.Raw(`SELECT COUNT(*) as qty, a.user_id, a.product_id, c.item_name, c.price, c.picture 
		FROM users_cart a 
		LEFT JOIN products_item c ON c.id = a.product_id 
		WHERE a.user_id = ?
		GROUP BY user_id, product_id, c.item_name, c.price, c.picture`, userId).QueryRows(&data)
	if err != nil {
		log.Print(err.Error())
		return 1, "Error while get data", nil
	}

	return 0, "Success", data
}

func RemoveCart(userId int, productId int) (rc int, msg string) {
	o := orm.NewOrm()

	_, err := o.Raw("DELETE FROM users_cart WHERE user_id = ? AND product_id = ?", userId, productId).Exec()
	if err != nil {
		return 1, err.Error()
	}

	return 0, "Success delete product"
}

func RemoveCartUser(userId int) (rc int, msg string) {
	o := orm.NewOrm()

	_, err := o.Raw("DELETE FROM users_cart WHERE user_id = ?", userId).Exec()
	if err != nil {
		return 1, err.Error()
	}

	return 0, "Success delete product"
}
