package bscscan

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type BscMallAccountDetail struct {
	Id            int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UserId        int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	Usdt          decimal.Decimal `json:"usdt" form:"usdt" gorm:"column:usdt;comment:变动余额;type:decimal"`
	Type          int             `json:"type" form:"type" gorm:"column:type;comment:类型 0收入 1支出;type:tinyint"`
	SourceType    int             `json:"sourceType" form:"sourceType" gorm:"column:source_type;comment:来源类型 0购买节点产品 1购买普通产品 2下级购买产品分红;type:int"`
	SourceContent string          `json:"sourceContent" form:"sourceContent" gorm:"column:source_content;comment:来源中文;type:varchar(128);"`
	BeginUsdt     decimal.Decimal `json:"beginUsdt" form:"beginUsdt" gorm:"column:begin_usdt;comment:期初资金;type:decimal"`
	EndUsdt       decimal.Decimal `json:"endUsdt" form:"endUsdt" gorm:"column:end_usdt;comment:期末资金;type:decimal"`
	SubUserId     int             `json:"subUserId" form:"subUserId" gorm:"column:sub_user_id;comment:下级id;type:bigint"`
	CreateTime    common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime    common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscMallAccountDetail) TableName() string {
	return "tb_bsc_mall_account_detail"
}
