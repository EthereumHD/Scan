package transaction

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"github.com/EthereumHD/Scan/src/sync"
	"fmt"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	"math/big"
	"qoobing.com/utillib.golang/log"
	"strings"
	"time"
)

type InputAddrTypeReq struct {
	InputAddrReq
	TxType int64 `json:"txType" form:"txType"` //请求的交易类型， 0 所有交易；1 转账交易；2 收款交易；3 抵押交易；4 赎回交易
}

func Get_by_addr_and_type(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputAddrRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	argc := new(InputAddrTypeReq)

	if err := c.BindInput(argc); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}
	log.Debugf("receive Get_Blocks: %+v", argc)
	//检查参数
	if argc.PageIndex < 1 || argc.PageSize <= 0 {
		log.Debugf("param error")
		return c.RESULT_ERROR(ERR_PARAMETER_INVALID, "param error")
	}
	//查链获取pending数据
	allList := []Transaction{}
	allCount := 0
	txtype := argc.TxType
	offset := (argc.PageIndex - 1) * argc.PageSize
	size := argc.PageSize
	pendingList, err := get_pending(argc.Addr, argc.TxType)
	if err != nil {
		log.Debugf("get_pending error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(GET_TRANSACTIONS_ERROR, fmt.Sprintf("get_pending error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	pendingLen := len(pendingList)
	allCount += pendingLen
	if pendingLen <= offset {
		log.Debugf("pendingLen <= offset")
		offset = offset - pendingLen
		size = size
		allList = []Transaction{}
	} else if pendingLen > offset && pendingLen <= offset+size {
		log.Debugf("pendingLen > offset && pendingLen <= offset+size")
		offset = 0
		size = size - pendingLen
		allList = pendingList[offset:]
	} else if pendingLen > offset+size {
		log.Debugf("pendingLen > offset+size")
		offset = 0
		size = offset + size - pendingLen
		allList = pendingList[offset : offset+size]
	}

	//查询数据库
	dbTransList, count, err := GetTransactionsByAddrAndType(c.Mysql(), argc.Addr, txtype, offset, size)
	if err != nil {
		log.Debugf("GetTransactionsByAddr error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(GET_TRANSACTIONS_ERROR, fmt.Sprintf("GetTransactionsByAddr error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	allList = concat(allList, dbTransList)

	//包装参数
	rsp.Count = count
	for _, trans := range allList {
		var transInfo TransInfo
		transInfo.TXHash = trans.F_tx_hash
		transInfo.BlockNumber = trans.F_block
		transInfo.From = trans.F_from
		transInfo.To = trans.F_to
		transInfo.Value = trans.F_value
		transInfo.TxFee = trans.F_tx_fee
		transInfo.Timestamp = trans.F_timestamp
		transInfo.TxTypeExt = trans.F_tx_type_ext
		if trans.F_tx_type == 0 && strings.ToUpper(trans.F_from) == strings.ToUpper(argc.Addr) {
			transInfo.TxType = TX_TYPE_FROM_ME
			transInfo.TxTypeExt = trans.F_value
		} else if trans.F_tx_type == 0 && strings.ToUpper(trans.F_to) == strings.ToUpper(argc.Addr) {
			transInfo.TxType = TX_TYPE_TO_ME
			transInfo.TxTypeExt = trans.F_value
		} else {
			transInfo.TxType = trans.F_tx_type
		}
		rsp.Transactions = append(rsp.Transactions, transInfo)
	}

	//返回结果
	return c.RESULT(rsp)
}

func get_pending(addr string, txtype int64) (pendingList []Transaction, err error) {
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	content, err := webthree.Txpool.Content()
	if err != nil {
		log.Fatalf("Content error:%s", err.Error())
		return pendingList, err
	}
	log.Debugf("get_pending, webthree.Txpool:%+v", content)

	for _, txmap := range content.Pending {
		for _, tx := range txmap {
			if strings.ToLower(tx.From) == strings.ToLower(addr) || strings.ToLower(tx.To) == strings.ToLower(addr) {
				tx_fee := big.NewInt(1).Mul(tx.GasPrice, tx.Gas)
				newTx := Transaction{}
				newTx.F_tx_type, newTx.F_tx_type_ext = sync.CalcTransactionType(tx)
				newTx.F_from = tx.From
				newTx.F_to = tx.To
				newTx.F_block = -2
				newTx.F_modify_time = time.Now().Format("2006-01-02 15:04:05.000")
				newTx.F_create_time = newTx.F_modify_time
				newTx.F_id = 0
				newTx.F_timestamp = 0
				newTx.F_tx_fee = tx_fee.String()
				newTx.F_tx_hash = tx.Hash
				newTx.F_value = tx.Value.String()

				if txtype == TX_TYPE_ALL ||
					(txtype == TX_TYPE_ME_MORTGAGE && newTx.F_tx_type == TX_TYPE_ME_MORTGAGE) ||
					(txtype == TX_TYPE_ME_REDEEM && newTx.F_tx_type == TX_TYPE_ME_REDEEM) ||
					(txtype == TX_TYPE_QUERY_3OR4 && (newTx.F_tx_type == TX_TYPE_ME_REDEEM || newTx.F_tx_type == TX_TYPE_ME_MORTGAGE)) ||
					(txtype == TX_TYPE_FROM_ME && strings.ToLower(tx.From) == strings.ToLower(addr)) ||
					(txtype == TX_TYPE_TO_ME && strings.ToLower(tx.To) == strings.ToLower(addr)) {
					//add to result
					pendingList = append(pendingList, newTx)
				}
			}
		}
	}
	return pendingList, nil
}

func concat(arr1 []Transaction, arr2 []Transaction) []Transaction {
	for _, o := range arr2 {
		arr1 = append(arr1, o)
	}
	return arr1
}
