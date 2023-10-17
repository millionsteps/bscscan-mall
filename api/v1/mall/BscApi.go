package mall

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan/dto"
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
		response.FailWithMessage("提现申请失败", c)
	} else {
		response.OkWithMessage("申请成功！", c)
	}
}
