package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"superindo_backend/helper"
	"superindo_backend/models"

	_ "github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type UsersCartController struct {
	beego.Controller
}

func (a *UsersCartController) Add() {
	var usersCart models.UsersCart
	json.Unmarshal(a.Ctx.Input.RequestBody, &usersCart)

	a.Ctx.Input.Bind(&usersCart.ProductId, "product_id")
	a.Ctx.Input.Bind(&usersCart.UserId, "user_id")

	rc, msg := models.AddCart(usersCart)
	helper.Response(rc, msg, nil, a.Controller)
}

func (a *UsersCartController) GetCart() {
	userId, _ := strconv.Atoi(a.GetString(":userid"))

	rc, msg, data := models.GetListCart(userId)
	helper.Response(rc, msg, data, a.Controller)
}

func (a *UsersCartController) RemoveCartUsers() {
	userId, _ := strconv.Atoi(a.GetString("userid"))
	productID, _ := strconv.Atoi(a.GetString("productid"))

	log.Print("USER ID ", userId)
	log.Print("PRODUCT ID", productID)

	rc, msg := models.RemoveCart(userId, productID)
	helper.Response(rc, msg, nil, a.Controller)
}

func (a *UsersCartController) RemoveAllCartUsers() {
	userId, _ := strconv.Atoi(a.GetString("userid"))

	log.Print("USER ID ", userId)

	rc, msg := models.RemoveCartUser(userId)
	helper.Response(rc, msg, nil, a.Controller)
}
