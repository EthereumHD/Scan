package transaction

type InputReq struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"` //范围起点
	PageSize  int `json:"pageSize" form:"pageSize"`   //范围重点
}

type OutputRsp struct {
	ErrNo        int       `json:"err_no"`
	ErrMsg       string    `json:"err_msg"`
	Count        int64     `json:"count"` //个数
	Transactions TransList `json:"transactions"`
}

type TransInfo struct {
	TXHash      string `json:"tx_hash"`
	BlockNumber int64  `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	TxFee       string `json:"txfee"`
	Nonce       int64  `json:"nonce"`
	TxType      int64  `json:"tx_type"`
	TxTypeExt   string `json:"tx_type_ext"`
}

type TransList []TransInfo

type InputAddrReq struct {
	Addr      string `json:"addr" form:"addr"`
	PageIndex int    `json:"pageIndex" form:"pageIndex"` //范围起点
	PageSize  int    `json:"pageSize" form:"pageSize"`   //范围重点
}

type OutputAddrRsp struct {
	ErrNo        int       `json:"err_no"`
	ErrMsg       string    `json:"err_msg"`
	Count        int64     `json:"count"` //个数
	Transactions TransList `json:"transactions"`
}

type InputHeightReq struct {
	Height    int64 `json:"height" form:"height"`
	PageIndex int   `json:"pageIndex" form:"pageIndex"` //范围起点
	PageSize  int   `json:"pageSize" form:"pageSize"`   //范围重点
}

type OutputHeightRsp struct {
	ErrNo        int       `json:"err_no"`
	ErrMsg       string    `json:"err_msg"`
	Count        int64     `json:"count"` //个数
	Transactions TransList `json:"transactions"`
}
