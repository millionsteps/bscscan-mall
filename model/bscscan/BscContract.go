package bscscan

import (
	"main.go/model/common"
)

type BscContract struct {
	Id              int             `json:"id" gorm:"primarykey;AUTO_INCREMENT"`
	ContractAddress string          `json:"contractAddress" form:"contractAddress" gorm:"column:contract_address;comment:合约地址;type:varchar(255)"`
	ContractAbi     string          `json:"contractAbi" form:"contractAbi" gorm:"column:contract_abi;comment:合约abi;type:text"`
	ToAddress       string          `json:"toAddress" form:"toAddress" gorm:"column:to_address;comment:收款地址;type:varchar(255)"`
	CreateTime      common.JSONTime `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"`
	UpdateTime      common.JSONTime `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;type:datetime"`
}

func (BscContract) TableName() string {
	return "tb_bsc_contract"
}
