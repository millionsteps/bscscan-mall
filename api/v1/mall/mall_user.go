package mall

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/common/request"
	"main.go/model/common/response"
	mallReq "main.go/model/mall/request"
	"main.go/utils"
)

type MallUserApi struct {
}

func (m *MallUserApi) UserRegister(c *gin.Context) {
	var req mallReq.RegisterUserParam
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.MallUserRegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := mallUserService.RegisterUser(req); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
	}
	response.OkWithMessage("创建成功", c)
}

func (m *MallUserApi) UserInfoUpdate(c *gin.Context) {
	var req mallReq.UpdateUserInfoParam
	_ = c.ShouldBindJSON(&req)
	token := c.GetHeader("token")
	if err := mallUserService.UpdateUserInfo(token, req); err != nil {
		global.GVA_LOG.Error("更新用户信息失败", zap.Error(err))
		response.FailWithMessage("更新用户信息失败"+err.Error(), c)
	}
	response.OkWithMessage("更新成功", c)
}

func (m *MallUserApi) GetUserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	if err, userDetail := mallUserService.GetUserDetail(token); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		response.OkWithData(userDetail, c)
	}
}

func (m *MallUserApi) GetUserBonusInfo(c *gin.Context) {
	token := c.GetHeader("token")
	if err, userDetail := mallUserService.GetUserBonusInfo(token); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		response.OkWithData(userDetail, c)
	}
}

func (m *MallUserApi) UserLogin(c *gin.Context) {
	var req mallReq.UserLoginParam
	_ = c.ShouldBindJSON(&req)
	if err, _, adminToken := mallUserService.UserLogin(req); err != nil {
		response.FailWithMessage("登陆失败", c)
	} else {
		response.OkWithData(adminToken.Token, c)
	}
}

func (m *MallUserApi) UserAddressLogin(c *gin.Context) {
	var req mallReq.UserAddressLoginParam
	_ = c.ShouldBindJSON(&req)
	if err, _, adminToken := mallUserService.UserAddressLogin(req); err != nil {
		response.FailWithMessage("登陆失败", c)
	} else {
		response.OkWithData(adminToken.Token, c)
	}
}

func (m *MallUserApi) UserLogout(c *gin.Context) {
	token := c.GetHeader("token")
	if err := mallUserTokenService.DeleteMallUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
}

func (m *MallUserApi) GetUserTeamList(c *gin.Context) {
	token := c.GetHeader("token")
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := mallUserService.GetUserTeamList(pageInfo.PageNumber, token); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (m *MallUserApi) GetAccountDetailList(c *gin.Context) {
	token := c.GetHeader("token")
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := mallUserService.GetAccountDetailList(pageInfo.PageNumber, token); err != nil {
		global.GVA_LOG.Error("未查询到记录", zap.Error(err))
		response.FailWithMessage("未查询到记录", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}
