package mall

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
	"main.go/middleware"
)

type MallUserRouter struct {
}

func (m *MallUserRouter) InitMallUserRouter(Router *gin.RouterGroup) {
	mallUserRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	userRouter := Router.Group("v1")
	var mallUserApi = v1.ApiGroupApp.MallApiGroup.MallUserApi
	{
		mallUserRouter.PUT("/user/info", mallUserApi.UserInfoUpdate)                      //修改用户信息
		mallUserRouter.GET("/user/info", mallUserApi.GetUserInfo)                         //获取用户信息
		mallUserRouter.GET("/user/team/list", mallUserApi.GetUserTeamList)                //获取团队信息
		mallUserRouter.GET("/user/account/detail/list", mallUserApi.GetAccountDetailList) //获取用户账户明细
		mallUserRouter.POST("/user/logout", mallUserApi.UserLogout)                       //登出
	}
	{
		userRouter.POST("/user/register", mallUserApi.UserRegister)          //用户注册
		userRouter.POST("/user/login", mallUserApi.UserLogin)                //登陆
		userRouter.POST("/user/address/login", mallUserApi.UserAddressLogin) //钱包地址登陆 包含邀请注册登录
	}

}
