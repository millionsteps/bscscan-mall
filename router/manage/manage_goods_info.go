package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
)

type ManageGoodsInfoRouter struct {
}

func (m *ManageGoodsInfoRouter) InitManageGoodsInfoRouter(Router *gin.RouterGroup) {
	mallGoodsInfoRouter := Router.Group("v1")
	var mallGoodsInfoApi = v1.ApiGroupApp.ManageApiGroup.ManageGoodsInfoApi
	var bscWithdrawRecordController = v1.ApiGroupApp.ManageApiGroup.BscWithdrawRecordController
	{
		mallGoodsInfoRouter.POST("goods", mallGoodsInfoApi.CreateGoodsInfo)                    // 新建MallGoodsInfo
		mallGoodsInfoRouter.DELETE("deleteMallGoodsInfo", mallGoodsInfoApi.DeleteGoodsInfo)    // 删除MallGoodsInfo
		mallGoodsInfoRouter.PUT("goods/status/:status", mallGoodsInfoApi.ChangeGoodsInfoByIds) // 上下架
		mallGoodsInfoRouter.PUT("goods", mallGoodsInfoApi.UpdateGoodsInfo)                     // 更新MallGoodsInfo
		mallGoodsInfoRouter.GET("goods/:id", mallGoodsInfoApi.FindGoodsInfo)                   // 根据ID获取MallGoodsInfo
		mallGoodsInfoRouter.GET("goods/list", mallGoodsInfoApi.GetGoodsInfoList)               // 获取MallGoodsInfo列表
	}
	{
		mallGoodsInfoRouter.GET("withdraw/list", bscWithdrawRecordController.GetSelectList)           // 获取提现待审核列表
		mallGoodsInfoRouter.POST("withdraw/auditWithdraw", bscWithdrawRecordController.AuditWithdraw) // 审核接口
		mallGoodsInfoRouter.POST("withdraw/remitWithdraw", bscWithdrawRecordController.RemitWithdraw) // 打款接口
	}
}
