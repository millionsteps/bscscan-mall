package mall

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan/dto"
	"main.go/model/common/request"
	"main.go/model/common/response"
)

type BscApi struct {
}

func (b *BscApi) GetContract(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if err, contract := bscService.GetBscConfigById(id); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		response.OkWithData(contract, c)
	}
}

func (b *BscApi) Withdraw(c *gin.Context) {
	var withdrawDTO dto.WithdrawDTO
	_ = c.ShouldBindJSON(&withdrawDTO)
	token := c.GetHeader("token")
	if err := bscService.Withdraw(token, withdrawDTO); err != nil {
		global.GVA_LOG.Error("提现申请失败", zap.Error(err))
		response.FailWithMessage("提现申请失败，"+err.Error(), c)
	} else {
		response.OkWithMessage("申请成功！", c)
	}
}

// GetBonusList 分页获取分红列表
func (m *BscApi) GetBonusList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	token := c.GetHeader("token")
	if err, list, total := bscService.GetBonusList(token, pageInfo.PageSize, pageInfo.PageNumber); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else if len(list) < 1 {
		// 前端项目这里有一个取数逻辑，如果数组为空，数组需要为[] 不能是Null
		response.OkWithDetailed(response.PageResult{
			List:       make([]interface{}, 0),
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "SUCCESS", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}
