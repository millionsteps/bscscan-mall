package mall

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
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
