package models

import (
	"log"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type ProductsItem struct {
	Id             int       `json:"id"`
	Picture        string    `json:"picture"`
	Stock          int       `json:"Stock"`
	StockCondition string    `json:"stock_condition"`
	Description    string    `json:"description"`
	IdSeller       int       `json:"id_seller"`
	Price          int       `json:"price"`
	CreateAt       time.Time `json:"create_at"`
	UpdateAt       time.Time `json:"update_at"`
	ItemName       string    `json:"item_name"`
}

type ProductsPicture struct {
	Id    int    `json:"id"`
	Image string `json:"image"`
}

func init() {
	orm.RegisterModel(new(ProductsItem), new(ProductsPicture))
}

func AddProduct(a ProductsItem) (rc int, msg string) {
	o := orm.NewOrm()

	if a != (ProductsItem{}) {

		a.CreateAt = time.Now()

		var user S_user
		log.Print(a.IdSeller)
		o.Raw("SELECT * FROM s_user WHERE id = ?", a.IdSeller).QueryRow(&user)
		if user == (S_user{}) {
			return 1, "Seller ID not found"
		}

		if a.Picture != "" {
			var prodPicture ProductsPicture
			prodPicture.Image = a.Picture

			if id, err := o.Insert(&prodPicture); err == nil {
				a.Picture = "/products/picture/" + strconv.Itoa(int(id))
			} else {
				return 1, "Error while insert Product pict " + err.Error()
			}
		} else {
			return 1, "Picture is required"
		}

		_, err := o.Insert(&a)
		if err != nil {
			return 1, "Error when insert data " + err.Error()
		}

		return 0, "Insert data success"

	}

	return 1, "Auction Item null"
}

func UpdateProducts(products *ProductsItem) (data *ProductsItem, msg string, rc int) {
	if strconv.Itoa(products.Id) != "" {
		o := orm.NewOrm()

		var us ProductsItem
		if er := o.Raw("SELECT * FROM products_item WHERE id = ?", products.Id).QueryRow(&us); er != nil {
			log.Print("Error while query " + er.Error())
		}

		if us == (ProductsItem{}) {
			return nil, "Data with id " + strconv.Itoa(products.Id) + " not found", 1
		}

		if products.Picture != "" {
			var productsPicture ProductsPicture

			productsPicture.Image = products.Picture

			if id, err := o.Insert(&productsPicture); err == nil {
				us.Picture = "/products/picture/" + strconv.Itoa(int(id))
			} else {
				return nil, "Error while insert Product picture " + err.Error(), 1
			}
		}

		us.ItemName = products.ItemName
		us.Stock = products.Stock
		us.Price = products.Price
		us.StockCondition = products.StockCondition
		us.Description = products.Description

		us.UpdateAt = time.Now()

		_, erro := o.Update(&us)
		if erro != nil {
			return nil, "Error when Update data " + erro.Error(), 1
		}

		return nil, "Success", 0

	}

	return nil, "Id not null", 1

}

func GetAllProduct() (rc int, msg string, data []ProductsItem) {
	o := orm.NewOrm()

	_, err := o.Raw(
		"SELECT * FROM products_item WHERE stock > 2 ORDER BY create_at DESC").QueryRows(&data)

	if err != nil {
		log.Print(err.Error())
		return 1, err.Error(), nil
	}

	return 0, "Success", data

}

func GetProductByUserId(userID int) (rc int, msg string, data []ProductsItem) {
	o := orm.NewOrm()

	_, err := o.Raw(
		"SELECT * FROM products_item WHERE id_seller = ? ORDER BY create_at DESC", userID).QueryRows(&data)

	if err != nil {
		log.Print(err.Error())
		return 1, err.Error(), nil
	}

	return 0, "Success", data
}

func RemoveProduct(productId int) (rc int, msg string) {
	o := orm.NewOrm()

	var products ProductsItem
	o.Raw("SELECT * FROM products_item WHERE id = ?", productId).QueryRow(&products)

	if products == (ProductsItem{}) {
		return 1, "Products not found"
	}

	if _, err := o.Delete(&products); err != nil {
		return 1, err.Error()
	}

	return 0, "Success delete product"
}
