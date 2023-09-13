package bscscan

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type BscMallUserAccount struct {
	Id            int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UserId        int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	ParentId      int             `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:父节点人id;type:bigint"`
	VipLevel      int             `json:"vipLevel" form:"vipLevel" gorm:"column:vip_level;comment:弱测等级 0没有等级 1一级;type:int"`
	TotalUsdt     decimal.Decimal `json:"totalUsdt" form:"totalUsdt" gorm:"column:total_usdt;comment:当前个人累计业绩;type:decimal"`
	TotalUsdtDown decimal.Decimal `json:"totalUsdtDown" form:"totalUsdtDown" gorm:"column:total_usdt_down;comment:弱侧伞下累计业绩;type:decimal"`
	Dao           decimal.Decimal `json:"dao" form:"dao" gorm:"column:dao;comment:节点数;type:decimal"`
	Usdt          decimal.Decimal `json:"usdt" form:"usdt" gorm:"column:usdt;comment:可提虚拟货币;type:decimal"`
	UsdtFreeze    decimal.Decimal `json:"usdtFreeze" form:"usdtFreeze" gorm:"column:usdt_freeze;comment:冻结虚拟货币;type:decimal"`
	CreateTime    common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime    common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscMallUserAccount) TableName() string {
	return "tb_bsc_mall_user_account"
}
