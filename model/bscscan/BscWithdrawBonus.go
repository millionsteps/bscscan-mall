package bscscan

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type BscWithdrawBonus struct {
	Id              int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UserId          int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	WithdrawAddress string          `json:"withdrawAddress" form:"withdrawAddress" gorm:"column:withdraw_address;comment:提现用户钱包地址;type:varchar(255);"`
	DaoUserId       int             `json:"daoUserId" form:"daoUserId" gorm:"column:dao_user_id;comment:节点用户主键id;type:bigint"`
	Usdt            decimal.Decimal `json:"usdt" form:"usdt" gorm:"column:usdt;comment:提现金额;type:decimal"`
	DaoNum          int             `json:"daoNum" form:"daoNum" gorm:"column:dao_num;comment:节点用户数;type:int"`
	BonusFreeze     decimal.Decimal `json:"bonusFreeze" form:"bonusFreeze" gorm:"column:bonus_freeze;comment:分红金额;type:decimal"`
	Bonus           decimal.Decimal `json:"bonus" form:"bonus" gorm:"column:bonus;comment:待稀释的分红金额;type:decimal"`
	CreateTime      common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime      common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscWithdrawBonus) TableName() string {
	return "tb_bsc_withdraw_bonus"
}
