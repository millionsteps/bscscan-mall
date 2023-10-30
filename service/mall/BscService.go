package mall

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/nanmu42/etherscan-api"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"main.go/global"
	"main.go/model/bscscan"
	"main.go/model/bscscan/dto"
	"main.go/model/bscscan/vo"
	"main.go/model/common"
	"main.go/model/mall"
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

func CheckOrder(txHash string, fromAddress string, toAddress string) (err error, isSuccess bool) {
	client := etherscan.NewCustomized(etherscan.Customization{
		Timeout: 15 * time.Second,
		Key:     "IQTP9I9KRAJ45EH2WUBI4WJ7PK7NUUIVGC",
		BaseURL: "https://api.bscscan.com/api?",
		// BaseURL: "https://api-testnet.bscscan.com/api?",
		Verbose: false,
	})

	contractAddress := "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd"
	startblock := 0
	endblock := 999999999

	//todo 30个一页不够的话就是加到100
	txs, err := client.ERC20Transfers(&contractAddress, &toAddress, &startblock, &endblock, 1, 100, true)
	if err != nil {
		global.GVA_LOG.Error("请求ERC20Transfers接口失败", zap.Error(err))
		return
	}
	hashIsExist := false
	for _, tx := range txs {
		if tx.Hash == txHash {
			if tx.To == strings.ToLower(toAddress) && tx.From == strings.ToLower(fromAddress) {
				hashIsExist = true
			}
		}
	}

	if hashIsExist == false {
		isSuccess = false
		return
	}
	//是否执行成功
	executionStatus, err := client.ExecutionStatus(txHash)
	if err != nil {
		global.GVA_LOG.Error("请求是否执行成功接口失败", zap.Error(err))
		return
	}
	fmt.Println(executionStatus.IsError)
	//是否成功交易
	receiptStatus, err := client.ReceiptStatus(txHash)
	if err != nil {
		global.GVA_LOG.Error("请求是否成功交易接口失败", zap.Error(err))
		return
	}
	fmt.Println(receiptStatus)
	if executionStatus.IsError == 0 && receiptStatus == 1 {
		isSuccess = true
		return
	}
	return
}

func (b *BscService) Withdraw(token string, withdrawDTO dto.WithdrawDTO) (err error) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	var userAccount bscscan.BscMallUserAccount
	err = global.GVA_DB.Where("user_id = ?", userToken.UserId).First(&userAccount).Error
	if err != nil {
		return errors.New("用户账户不存在")
	}
	thisUserUsdt := userAccount.Usdt
	if thisUserUsdt.Cmp(withdrawDTO.Usdt) == -1 {
		return errors.New("余额不足不可提现！")
	}
	var record bscscan.BscWithdrawRecord
	record.UserId = userToken.UserId
	record.Usdt = withdrawDTO.Usdt
	//计算手续费
	usdt := withdrawDTO.Usdt
	commissionCharge := usdt.Mul(decimal.NewFromFloat32(0.05))
	record.CommissionCharge = commissionCharge
	record.Address = withdrawDTO.Address
	record.Status = 0
	record.CreateTime = common.JSONTime{Time: time.Now()}
	//生成提现记录
	tx := global.GVA_DB.Begin()
	if err = tx.Save(&record).Error; err != nil {
		tx.Rollback()
		return errors.New("提现失败！")
	}
	resultUsdt := thisUserUsdt.Sub(usdt)
	//保存账户剩余余额
	userAccount.Usdt = resultUsdt
	userAccount.UpdateTime = common.JSONTime{Time: time.Now()}
	if err = tx.Save(&userAccount).Error; err != nil {
		tx.Rollback()
		return errors.New("提现失败,保存账户失败！")
	}
	//添加明细
	var detail bscscan.BscMallAccountDetail
	detail.UserId = userToken.UserId
	detail.Usdt = withdrawDTO.Usdt
	detail.SourceType = 3
	detail.SourceContent = "提现"
	detail.UpdateTime = common.JSONTime{time.Now()}
	detail.CreateTime = common.JSONTime{time.Now()}
	detail.Type = 1
	if err = tx.Save(&detail).Error; err != nil {
		return errors.New("提现失败,保存账户明细失败")
	}
	tx.Commit()
	return
}

func (b *BscService) GetBonusList(token string, pageSize int, pageNumber int) (err error, bscWithdrawBonusList []bscscan.BscWithdrawBonus, total int64) {
	var userToken mall.MallUserToken
	err = global.GVA_DB.Where("token =?", token).First(&userToken).Error
	if err != nil {
		return errors.New("不存在的用户"), bscWithdrawBonusList, total
	}
	limit := pageSize
	offset := pageSize * (pageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&bscscan.BscWithdrawBonus{})
	db.Where("dao_user_id = ?", userToken.UserId)
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&bscWithdrawBonusList).Error
	return
}
