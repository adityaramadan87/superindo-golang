package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["superindo_backend/controllers:ProductsController"] = append(beego.GlobalControllerRouter["superindo_backend/controllers:ProductsController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/api/lelang",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
