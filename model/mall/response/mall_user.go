package response

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type MallUserDetailResponse struct {
	NickName         string          `json:"nickName"`
	LoginName        string          `json:"loginName"`
	IntroduceSign    string          `json:"introduceSign"`
	UserId           int             `json:"userId"`
	BscAddress       string          `json:"bscAddress"`
	LoginType        int             `json:"loginType"`
	VipLevel         int             `json:"vipLevel"`
	Usdt             decimal.Decimal `json:"usdt"`
	CardNum          int             `json:"cardNum"`
	CardUsdt         decimal.Decimal `json:"cardUsdt"`
	CreateTime       common.JSONTime `json:"createTime"`
	EmailAddress     string          `json:"emailAddress"`
	ParentBscAddress string          `json:"parentBscAddress"`
}
