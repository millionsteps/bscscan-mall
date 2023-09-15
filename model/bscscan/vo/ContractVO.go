package vo

type ContractVO struct {
	Id              int         `json:"id"`
	ContractAddress string      `json:"contractAddress"`
	ContractAbiJson interface{} `json:"contractAbiJson"`
	ToAddress       string      `json:"toAddress"`
}
