package common

var ERR = map[string]string{
	"NONE": "000",
}

var STATUS = map[string]string{
	"Init":     "Inited",
	"Transfer": "Transfering",
	"Close":    "Closed",
}

const (
	//下划线
	ULINE = "_"
	//合同信息
	CONTRACT_INFO = "CONTRACT_INFO"
	// 合同交易信息
	CONTRACT_TRANSFER = "CONTRACT_TRANSFER"
)
