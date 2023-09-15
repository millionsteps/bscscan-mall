package mall

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/copier"
	"main.go/global"
	"main.go/model/bscscan"
	"main.go/model/bscscan/vo"
)

type BscService struct {
}

//根据id查询合约信息
func (b *BscService) GetBscConfigById(id int) (err error, contractVO vo.ContractVO) {
	var contract bscscan.BscContract
	err = global.GVA_DB.Where("id =?", id).First(&contract).Error
	if err != nil {
		return errors.New("不存在的合约"), contractVO
	}
	copier.Copy(&contractVO, &contract)
	var abiData interface{}
	jsonData := []byte(contract.ContractAbi)
	json.Unmarshal(jsonData, &abiData)
	contractVO.ContractAbiJson = abiData
	return
}
