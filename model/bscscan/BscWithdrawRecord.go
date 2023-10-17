package bscscan

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type BscWithdrawRecord struct {
	Id               int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UserId           int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	Usdt             decimal.Decimal `json:"usdt" form:"usdt" gorm:"column:usdt;comment:提现金额;type:decimal"`
	Address          string          `json:"address" form:"address" gorm:"column:address;comment:钱包地址;type:varchar(255);"`
	CommissionCharge decimal.Decimal `json:"commissionCharge" form:"commissionCharge" gorm:"column:commission_charge;comment:提现手续费;type:decimal"`
	Status           int             `json:"status" form:"status" gorm:"column:status;comment:状态 0待审核 1审核通过 2已打款 3审核失败;type:int"`
	CreateTime       common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime       common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscWithdrawRecord) TableName() string {
	return "tb_bsc_withdraw_record"
}
