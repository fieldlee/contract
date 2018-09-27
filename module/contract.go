package module

// 交易结构
type Transfer struct {
	TxHash      string `json:"txHash"`
	OperateTime int64  `json:"operateTime"`
	Operation   string `json:"operation"` // 操作内容
	From        string `json:"from"`
	To          string `json:"to"`
	Value       uint64 `json:"value"`
	Operator    string `json:"operator"` // 确权信息操作人
}

/** 操作日志 **/
type TransferLog map[string][]Transfer

// 交易结构
type Allowance map[string]uint64

// ERC结构
type ERC struct {
	Name        string               `json:"name"`
	Symbol      string               `json:"symbol"`
	Decimals    uint                 `json:"decimals"`
	TotalSupply uint64               `json:"totalSupply"`
	BalanceOf   map[string]uint64    `json:"balanceOf"`
	Allowances  map[string]Allowance `json:"allowances"`
}
