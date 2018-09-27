package module

/** 查询信息 **/
type ChanInfo struct {
	ContractName string `json:"contractName"`
	Success      bool   `json:"success"`
	Info         string `json:"info"`
}

/** 查询信息 **/
type ReturnInfo struct {
	Success bool   `json:"success"`
	Info    string `json:"info"`
}

/** 查询信息 **/
type QueryInfo struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Value   uint64 `json:"value"`
	Info    string `json:"info"`
}

/** 查询信息 **/
type QueryLog struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Info    string `json:"info"`
	Actions string `json:"actions"`
}
