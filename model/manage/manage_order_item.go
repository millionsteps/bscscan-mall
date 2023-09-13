package manage

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type MallOrderItem struct {
	OrderItemId   int             `json:"orderItemId" gorm:"primarykey;AUTO_INCREMENT"`
	OrderId       int             `json:"orderId" form:"orderId" gorm:"column:order_id;;type:bigint"`
	GoodsId       int             `json:"goodsId" form:"goodsId" gorm:"column:goods_id;;type:bigint"`
	UserId        int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	GoodsName     string          `json:"goodsName" form:"goodsName" gorm:"column:goods_name;comment:商品名;type:varchar(200);"`
	GoodsCoverImg string          `json:"goodsCoverImg" form:"goodsCoverImg" gorm:"column:goods_cover_img;comment:商品主图;type:varchar(200);"`
	SellingPrice  int             `json:"sellingPrice" form:"sellingPrice" gorm:"column:selling_price;comment:商品实际售价;type:int"`
	GoodsCount    int             `json:"goodsCount" form:"goodsCount" gorm:"column:goods_count;;type:bigint"`
	CreateTime    common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	ReleaseFlag   int             `json:"releaseFlag" form:"releaseFlag" gorm:"column:release_flag;comment:释放状态 0订单未完成 1释放中 2全部释放完成;type:tinyint"`
	UsdtFreeze    decimal.Decimal `json:"usdtFreeze" form:"usdtFreeze" gorm:"column:usdt_freeze;comment:冻结虚拟货币;type:decimal"`
	UsdtAble      decimal.Decimal `json:"usdtAble" form:"usdtAble" gorm:"column:usdt_able;comment:可释放;type:decimal"`
	Usdt          decimal.Decimal `json:"usdt" form:"usdt" gorm:"column:usdt;comment:已释放;type:decimal"`
	ReleaseRate   decimal.Decimal `json:"releaseRate" form:"releaseRate" gorm:"column:release_rate;comment:释放比例;type:decimal"`
}

func (MallOrderItem) TableName() string {
	return "tb_newbee_mall_order_item"
}
