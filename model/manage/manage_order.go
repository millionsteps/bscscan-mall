package manage

import (
	"github.com/shopspring/decimal"
	"main.go/model/common"
)

type MallOrder struct {
	OrderId     int             `json:"orderId" form:"orderId" gorm:"primarykey;AUTO_INCREMENT"`
	OrderNo     string          `json:"orderNo" form:"orderNo" gorm:"column:order_no;comment:订单号;type:varchar(20);"`
	UserId      int             `json:"userId" form:"userId" gorm:"column:user_id;comment:用户主键id;type:bigint"`
	TotalPrice  decimal.Decimal `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单总价;type:int"`
	PayStatus   int             `json:"payStatus" form:"payStatus" gorm:"column:pay_status;comment:支付状态:0.未支付,1.支付成功,-1:支付失败;type:tinyint"`
	PayType     int             `json:"payType" form:"payType" gorm:"column:pay_type;comment:0.无 1.支付宝支付 2.微信支付;type:tinyint"`
	PayTime     common.JSONTime `json:"payTime" form:"payTime" gorm:"column:pay_time;comment:支付时间;type:datetime"`
	OrderStatus int             `json:"orderStatus" form:"orderStatus" gorm:"column:order_status;comment:订单状态:0.待支付 1.已支付 2.配货完成 3:出库成功 4.交易成功 -1.手动关闭 -2.超时关闭 -3.商家关闭;type:tinyint"`
	ExtraInfo   string          `json:"extraInfo" form:"extraInfo" gorm:"column:extra_info;comment:订单body;type:varchar(100);"`
	TxHash      string          `json:"txHash" form:"txHash" gorm:"column:tx_hash;comment:交易成功hash;type:varchar(255);"`
	FromAddress string          `json:"fromAddress" form:"fromAddress" gorm:"column:from_address;comment:付款地址;type:varchar(255);"`
	ToAddress   string          `json:"toAddress" form:"toAddress" gorm:"column:to_address;comment:收款地址;type:varchar(255);"`
	IsDeleted   int             `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识字段(0-未删除 1-已删除);type:tinyint"`
	CreateTime  common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime  common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:最新修改时间;type:datetime"`
}

// TableName MallOrder 表名
func (MallOrder) TableName() string {
	return "tb_newbee_mall_order"
}
