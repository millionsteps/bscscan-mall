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
	UsdtFreeze       decimal.Decimal `json:"usdtFreeze"`
	BonusFlag        int             `json:"bonusFlag"`
	TotalUsdtDownA   decimal.Decimal `json:"totalUsdtDownA"`
	TotalUsdtDownB   decimal.Decimal `json:"totalUsdtDownB"`
}

type MallUserBonusDetailResponse struct {
	TotalUsdt    decimal.Decimal `json:"totalUsdt"`
	WithdrawUsdt decimal.Decimal `json:"withdrawUsdt"`
	DaoNum       int             `json:"daoNum"`
	Usdt         decimal.Decimal `json:"usdt"`
}
