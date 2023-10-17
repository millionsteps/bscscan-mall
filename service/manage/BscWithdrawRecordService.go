package manage

import (
	"main.go/global"
	"main.go/model/bscscan"
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
