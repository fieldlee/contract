package module

// erc init
type InitParam struct {
	Issuer      string `json:"issuer"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint   `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
}

// 交易信息
type TransferParam struct {
	Name  string `json:"name"`
	From  string `json:"from"`
	To    string `json:"to"` // 操作内容
	Value uint64 `json:"value"`
}

// 查询信息
type QueryParam struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
