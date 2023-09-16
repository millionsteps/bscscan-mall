package mall

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"
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

func CheckOrder(txHash string, fromAddress string, toAddress string) (err error, isSuccess bool) {
	client := etherscan.NewCustomized(etherscan.Customization{
		Timeout: 15 * time.Second,
		Key:     "IQTP9I9KRAJ45EH2WUBI4WJ7PK7NUUIVGC",
		// BaseURL: "https://api.bscscan.com/api?",
		BaseURL: "https://api-testnet.bscscan.com/api?",
		Verbose: false,
	})

	contractAddress := "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd"
	startblock := 0
	endblock := 999999999

	//todo 30个一页不够的话就是加到100
	txs, err := client.ERC20Transfers(&contractAddress, &toAddress, &startblock, &endblock, 1, 30, true)
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
