package manage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan/dto"
	"main.go/model/common/request"
	"main.go/model/common/response"
)

type BscWithdrawRecordController struct {
}

// GetSelectList 分页获取提现列表
func (m *BscWithdrawRecordController) GetSelectList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := bscWithdrawRecordService.SelectList(pageInfo.PageSize, pageInfo.PageNumber); err != nil {
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

func (m *BscWithdrawRecordController) AuditWithdraw(c *gin.Context) {
	var auditWithdrawDTO dto.AuditWithdrawDTO
	_ = c.ShouldBindJSON(&auditWithdrawDTO)
	if err := bscWithdrawRecordService.AuditWithdraw(auditWithdrawDTO.Id, auditWithdrawDTO.Status); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败"+err.Error(), c)
	} else {
		response.OkWithMessage("操作成功", c)
	}
}

func (m *BscWithdrawRecordController) RemitWithdraw(c *gin.Context) {
	var auditWithdrawDTO dto.AuditWithdrawDTO
	_ = c.ShouldBindJSON(&auditWithdrawDTO)
	if err := bscWithdrawRecordService.RemitWithdraw(auditWithdrawDTO.Id); err != nil {
		global.GVA_LOG.Error("操作失败!", zap.Error(err))
		response.FailWithMessage("操作失败"+err.Error(), c)
	} else {
		response.OkWithMessage("操作成功", c)
	}
}
