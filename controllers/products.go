package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/jpeg"
	"image/png"
	"log"
	"strconv"
	"strings"
	"superindo_backend/helper"
	"superindo_backend/models"

	"github.com/chai2010/webp"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "golang.org/x/crypto/bcrypt"
)

type ProductsController struct {
	beego.Controller
}

func (a *ProductsController) Add() {
	var productsItem models.ProductsItem
	json.Unmarshal(a.Ctx.Input.RequestBody, &productsItem)

	a.Ctx.Input.Bind(&productsItem.Picture, "picture")
	a.Ctx.Input.Bind(&productsItem.Price, "price")
	a.Ctx.Input.Bind(&productsItem.StockCondition, "stock_condition")
	a.Ctx.Input.Bind(&productsItem.Description, "description")
	a.Ctx.Input.Bind(&productsItem.Stock, "stock")
	a.Ctx.Input.Bind(&productsItem.IdSeller, "id_seller")
	a.Ctx.Input.Bind(&productsItem.ItemName, "item_name")

	rc, msg := models.AddProduct(productsItem)
	helper.Response(rc, msg, nil, a.Controller)

	a.ServeJSON()
}

// @Title getAllAuction
// @Summary getAllAuction
// @Description get all auction
// @Success 200 {object} models.AuctionItem
// @Failure 403 body is empty
// @router /api/lelang [get]
func (a *ProductsController) Get() {
	rc, msg, data := models.GetAllProduct()

	helper.Response(rc, msg, data, a.Controller)
}

func (a *ProductsController) GetByUserId() {
	userID, _ := strconv.Atoi(a.GetString(":userid"))

	rc, msg, data := models.GetProductByUserId(userID)
	helper.Response(rc, msg, data, a.Controller)
}

func (a *ProductsController) RemoveProductById() {
	productID, _ := strconv.Atoi(a.GetString(":productid"))

	rc, msg := models.RemoveProduct(productID)
	helper.Response(rc, msg, nil, a.Controller)
}

func (a *ProductsController) UpdateProduct() {
	var products models.ProductsItem
	json.Unmarshal(a.Ctx.Input.RequestBody, &products)

	a.Ctx.Input.Bind(&products.Id, "id")
	a.Ctx.Input.Bind(&products.Picture, "picture")
	a.Ctx.Input.Bind(&products.ItemName, "item_name")
	a.Ctx.Input.Bind(&products.Price, "price")
	a.Ctx.Input.Bind(&products.Stock, "stock")
	a.Ctx.Input.Bind(&products.Description, "description")
	a.Ctx.Input.Bind(&products.StockCondition, "stock_condition")

	_, msg, rc := models.UpdateProducts(&products)
	helper.Response(rc, msg, nil, a.Controller)

}

func (a *ProductsController) GetProductsPicture() {
	imgID, _ := strconv.Atoi(a.GetString(":pictureid"))

	o := orm.NewOrm()

	var img string
	if err := o.Raw("SELECT image FROM products_picture WHERE id = ?", imgID).QueryRow(&img); err != nil {
		helper.Response(1, "Error while get Image "+err.Error(), nil, a.Controller)
		return
	}

	coI := strings.Index(img, ",")
	rawImage := img[coI+1:]

	// Encoded Image DataUrl //
	unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))
	res := bytes.NewReader(unbased)

	switch strings.TrimSuffix(img[5:coI], ";base64") {
	case "image/png":
		pngImg, _ := png.Decode(res)

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, pngImg, nil); err != nil {
			log.Println("unable to encode image.")
		}

		a.Ctx.Request.Header.Set("Content-Type", "image/png")
		a.Ctx.Request.Header.Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := a.Ctx.ResponseWriter.Write(buffer.Bytes()); err != nil {
			log.Print("errrrrorr", err.Error())
		}
	case "image/jpeg":
		jpgImg, _ := jpeg.Decode(res)

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, jpgImg, nil); err != nil {
			log.Println("unable to encode image.")
		}

		a.Ctx.Request.Header.Set("Content-Type", "image/png")
		a.Ctx.Request.Header.Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := a.Ctx.ResponseWriter.Write(buffer.Bytes()); err != nil {
			log.Print("errrrrorr", err.Error())
		}
	case "image/webp":
		jpgImg, _ := webp.Decode(res)

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, jpgImg, nil); err != nil {
			log.Println("unable to encode image.")
		}

		a.Ctx.Request.Header.Set("Content-Type", "image/webp")
		a.Ctx.Request.Header.Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := a.Ctx.ResponseWriter.Write(buffer.Bytes()); err != nil {
			log.Print("errrrrorr", err.Error())
		}
	}
}
