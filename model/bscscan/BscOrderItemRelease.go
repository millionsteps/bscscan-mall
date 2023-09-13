package bscscan

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type BscOrderItemRelease struct {
	Id           int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	OrderId      int             `json:"orderId" form:"orderId" gorm:"column:order_id;;type:bigint"`
	OrderItemId  int             `json:"orderItemId" form:"orderItemId" gorm:"column:order_item_id;;type:bigint"`
	UserId       int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	ReleaseState int             `json:"releaseState" form:"releaseState" gorm:"column:release_state;comment:释放状态 0未释放 1已释放;type:tinyint"`
	UsdtFreeze   decimal.Decimal `json:"usdtFreeze" form:"usdtFreeze" gorm:"column:usdt_freeze;comment:冻结虚拟货币;type:decimal"`
	UsdtBegin    decimal.Decimal `json:"usdtBegin" form:"usdtBegin" gorm:"column:usdt_begin;comment:期初;type:decimal"`
	UsdtEnd      decimal.Decimal `json:"usdtEnd" form:"usdtEnd" gorm:"column:usdt_end;comment:期初;type:decimal"`
	ThisUsdt     decimal.Decimal `json:"thisUsdt" form:"thisUsdt" gorm:"column:this_usdt;comment:本次释放;type:decimal"`
	ReleaseRate  decimal.Decimal `json:"releaseRate" form:"releaseRate" gorm:"column:release_rate;comment:释放比例;type:decimal"`
	RelesaeDate  string          `json:"relesaeDate" form:"relesaeDate" gorm:"column:relesae_date;comment:释放日期 yyyy-MM-dd;type:char(10);"`
	CreateTime   common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime   common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscOrderItemRelease) TableName() string {
	return "tb_bsc_order_item_release"
}
