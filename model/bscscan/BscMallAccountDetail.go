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
	SourceType    int             `json:"sourceType" form:"sourceType" gorm:"column:source_type;comment:来源类型 0购买节点产品 1购买普通产品 2下级购买产品分红 3直推奖励10%;type:int"`
	SourceContent string          `json:"sourceContent" form:"sourceContent" gorm:"column:source_content;comment:来源中文;type:varchar(128);"`
	BeginUsdt     decimal.Decimal `json:"beginUsdt" form:"beginUsdt" gorm:"column:begin_usdt;comment:期初资金;type:decimal"`
	EndUsdt       decimal.Decimal `json:"endUsdt" form:"endUsdt" gorm:"column:end_usdt;comment:期末资金;type:decimal"`
	SubUserId     int             `json:"subUserId" form:"subUserId" gorm:"column:sub_user_id;comment:下级id;type:bigint"`
	CreateTime    common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime    common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
	ReleaseFlag   int             `json:"releaseFlag" form:"releaseFlag" gorm:"column:release_flag;comment:释放状态 1释放中 2全部释放完成;type:tinyint"`
	UsdtFreeze    decimal.Decimal `json:"usdtFreeze" form:"usdtFreeze" gorm:"column:usdt_freeze;comment:冻结虚拟货币;type:decimal"`
	UsdtAble      decimal.Decimal `json:"usdtAble" form:"usdtAble" gorm:"column:usdt_able;comment:可释放;type:decimal"`
	UsdtRelease   decimal.Decimal `json:"usdtRelease" form:"usdtRelease" gorm:"column:usdt_release;comment:已释放;type:decimal"`
	ReleaseRate   decimal.Decimal `json:"releaseRate" form:"releaseRate" gorm:"column:release_rate;comment:释放比例;type:decimal"`
}

func (BscMallAccountDetail) TableName() string {
	return "tb_bsc_mall_account_detail"
}
