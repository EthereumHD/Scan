package transaction

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"fmt"
	"github.com/labstack/echo"
	"qoobing.com/utillib.golang/log"
)

func Get_by_addr(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputAddrRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	argc := new(InputAddrReq)

	if err := c.BindInput(argc); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}
	log.Debugf("receive Get_Blocks: %+v", argc)
	//检查参数
	if argc.PageIndex < 1 || argc.PageSize <= 0 {
		log.Debugf("param error")
		return c.RESULT_ERROR(ERR_PARAMETER_INVALID, "param error")
	}
	//查询数据库
	//sql := "(F_from = '" + argc.Addr + "' or F_to = '" + argc.Addr + "') "
	count, err := GetTransactionsCountByAddr(c.Mysql(), argc.Addr)
	if err != nil {
		log.Debugf("GetTransactionsCountByAddr error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(TRANSACTION_COUNT_ERROR, fmt.Sprintf("GetTransactionsCountByAddr error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	rsp.Count = count

	offset := (argc.PageIndex - 1) * argc.PageSize
	size := argc.PageSize
	transList, err := GetTransactionsByAddr(c.Mysql(), argc.Addr, offset, size)
	if err != nil {
		log.Debugf("GetTransactionsByAddr error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(GET_TRANSACTIONS_ERROR, fmt.Sprintf("GetTransactionsByAddr error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	//包装参数
	for _, trans := range transList {
		var transInfo TransInfo
		transInfo.TXHash = trans.F_tx_hash
		transInfo.BlockNumber = trans.F_block
		transInfo.From = trans.F_from
		transInfo.To = trans.F_to
		transInfo.Value = trans.F_value
		transInfo.TxFee = trans.F_tx_fee
		transInfo.Timestamp = trans.F_timestamp
		rsp.Transactions = append(rsp.Transactions, transInfo)
	}

	//返回结果
	return c.RESULT(rsp)
}
