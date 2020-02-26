package transaction

import (
	"github.com/labstack/echo"
	."github.com/EthereumHD/Scan/src/apicontext"
	"qoobing.com/utillib.golang/log"
	."github.com/EthereumHD/Scan/src/const"
	."github.com/EthereumHD/Scan/src/model"
	"fmt"
)

func Get_transactions(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	argc := new(InputReq)

	if err := c.BindInput(argc); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}
	log.Debugf("receive Get_Blocks: %+v", argc)
	//检查参数
	if argc.PageIndex < 1 || argc.PageSize <= 0{
		log.Debugf("param error")
		return c.RESULT_ERROR(ERR_PARAMETER_INVALID, "param error")
	}
	//查询数据库
	count,err := GetTransactionsCount(c.Mysql())
	if  err != nil{
		log.Debugf("GetTransactionsCount error:%s",err.Error())
		return c.RESULT_ERROR(TRANSACTION_COUNT_ERROR,fmt.Sprintf("GetTransactionsCount error:%s",err.Error()))//return c.RESULT(rsp)
	}
	rsp.Count = count

	offset := (argc.PageIndex - 1) * argc.PageSize
	size := argc.PageSize
	transList,err := GetTransactions(c.Mysql(),offset,size);
	if  err != nil{
		log.Debugf("GetRecentBlocks error:%s",err.Error())
		return c.RESULT_ERROR(GET_TRANSACTIONS_ERROR,fmt.Sprintf("GetRecentBlocks error:%s",err.Error()))//return c.RESULT(rsp)
	}
	//包装参数
	for _,trans := range transList{
		var transInfo  TransInfo
		transInfo.TXHash = trans.F_tx_hash
		transInfo.BlockNumber = trans.F_block
		transInfo.From = trans.F_from
		transInfo.To = trans.F_to
		transInfo.Value = trans.F_value
		transInfo.TxFee = trans.F_tx_fee
		transInfo.Timestamp = trans.F_timestamp
		rsp.Transactions = append(rsp.Transactions,transInfo)
	}

	//返回结果
	return c.RESULT(rsp)
}
