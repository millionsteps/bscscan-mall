package mall

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
	"main.go/middleware"
)

type BscConfigController struct {
}

func (m *MallUserRouter) InitBscConfigRouter(Router *gin.RouterGroup) {
	bscConfigRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	var bscApi = v1.ApiGroupApp.MallApiGroup.BscApi
	{
		bscConfigRouter.GET("/contract/info", bscApi.GetContract) //获取合约信息

		bscConfigRouter.POST("/withdraw", bscApi.Withdraw) //提现
	}
}
