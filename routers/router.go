// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"superindo_backend/controllers"

	"fmt"
	"strings"
	"superindo_backend/models"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

var authFilter = func(context *context.Context) {
	header := strings.Split(context.Input.Header("Authorization"), " ")
	if len(header) != 2 {
		context.Abort(403, "Unauthorized")
	}

	if err := models.ValidateToken(header[1]); err != nil {
		context.Abort(403, err.Error())
	}

	fmt.Println(context.Input.URL())
}

func init() {

	nameSpace := beego.NewNamespace("/api",
		beego.NSRouter("/users/login", &controllers.UsersController{}, "post:Login"),
		beego.NSRouter("/users/register", &controllers.UsersController{}, "post:Register"),
		beego.NSRouter("/users/otp", &controllers.UsersController{}, "get:VerifyOtp"),

		beego.NSNamespace("/v1",
			beego.NSBefore(authFilter),
			beego.NSRouter("/products", &controllers.ProductsController{}, "get:Get"),
			beego.NSRouter("/products/add", &controllers.ProductsController{}, "post:Add"),
			beego.NSRouter("/products/user/:userid", &controllers.ProductsController{}, "get:GetByUserId"),
			beego.NSRouter("/products/delete/:productid", &controllers.ProductsController{}, "get:RemoveProductById"),
			beego.NSRouter("/products/update", &controllers.ProductsController{}, "post:UpdateProduct"),

			beego.NSRouter("/users/update", &controllers.UsersController{}, "post:Update"),
			beego.NSRouter("/users/:userid", &controllers.UsersController{}, "get:GetUserById"),
			beego.NSRouter("/users/avatar/:avatarid", &controllers.UsersController{}, "*:GetUserAvatar"),

			beego.NSRouter("/users/cart/add", &controllers.UsersCartController{}, "post:Add"),
			beego.NSRouter("/users/cart/:userid", &controllers.UsersCartController{}, "get:GetCart"),
			beego.NSRouter("/users/cart/delete", &controllers.UsersCartController{}, "get:RemoveCartUsers"),
			beego.NSRouter("/users/cart/deleteall", &controllers.UsersCartController{}, "get:RemoveAllCartUsers"),
		),
	)

	beego.Router("/products/picture/:pictureid", &controllers.ProductsController{}, "*:GetProductsPicture")

	beego.AddNamespace(nameSpace)
}
