package manage

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"main.go/global"
	"main.go/model/bscscan"
	"main.go/model/common"
)

type BscWithdrawRecordService struct {
}

func (b *BscWithdrawRecordService) SelectList(pageSize int, pageNumber int) (err error, bscWithdrawRecords []bscscan.BscWithdrawRecord, total int64) {
	limit := pageSize
	offset := pageSize * (pageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&bscscan.BscWithdrawRecord{})
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&bscWithdrawRecords).Error
	return
}

func (b *BscWithdrawRecordService) AuditWithdraw(id int, statusParam int) (err error) {
	var bscWithdrawRecord bscscan.BscWithdrawRecord
	if err = global.GVA_DB.Where("id = ?", id).First(&bscWithdrawRecord).Error; err != nil {
		return errors.New("审核记录不存在！")
	}
	status := bscWithdrawRecord.Status
	if status != 0 {
		return errors.New("请重复审核！")
	}
	bscWithdrawRecord.Status = statusParam
	bscWithdrawRecord.UpdateTime = common.JSONTime{time.Now()}
	err = global.GVA_DB.Save(&bscWithdrawRecord).Error
	if err != nil {
		return errors.New("保存失败！")
	}
	return
}

func (b *BscWithdrawRecordService) RemitWithdraw(id int) (err error) {
	var bscWithdrawRecord bscscan.BscWithdrawRecord
	if err = global.GVA_DB.Where("id = ?", id).First(&bscWithdrawRecord).Error; err != nil {
		return errors.New("审核记录不存在！")
	}
	status := bscWithdrawRecord.Status
	if status != 1 {
		return errors.New("请重复打款！")
	}
	bscWithdrawRecord.Status = 2
	bscWithdrawRecord.UpdateTime = common.JSONTime{time.Now()}
	err = global.GVA_DB.Save(&bscWithdrawRecord).Error
	if err != nil {
		return errors.New("保存失败！")
	}

	//计算节点数量
	var userAccountList []bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("dao_flag = 1").Find(&userAccountList).Error
	if err != nil {
		return errors.New("保存失败！")
	}
	size := len(userAccountList)
	commissionCharge := bscWithdrawRecord.CommissionCharge
	totalBonus := commissionCharge.Mul(decimal.NewFromFloat32(0.5))
	bonus := totalBonus.Div(decimal.NewFromInt(int64(size)))

	//生成分红记录
	for _, account := range userAccountList {
		var bscWithdrawBonus bscscan.BscWithdrawBonus
		bscWithdrawBonus.Usdt = bscWithdrawRecord.Usdt
		bscWithdrawBonus.UserId = bscWithdrawRecord.UserId
		bscWithdrawBonus.DaoUserId = account.UserId
		bscWithdrawBonus.DaoNum = size
		bscWithdrawBonus.Bonus = bonus
		bscWithdrawBonus.BonusFreeze = bonus
		bscWithdrawBonus.CreateTime = common.JSONTime{time.Now()}
		bscWithdrawBonus.UpdateTime = common.JSONTime{time.Now()}
		err = global.GVA_DB.Save(&bscWithdrawBonus).Error
		if err != nil {
			return errors.New("保存失败！")
		}
	}
	return
}
