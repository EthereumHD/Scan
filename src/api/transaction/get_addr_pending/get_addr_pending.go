package get_addr_pending

import (
	"github.com/labstack/echo"
	//"time"
	"github.com/EthereumHD/Scan/src/api/transaction"
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"go-web3"
	"go-web3/providers"
	"math/big"
	"qoobing.com/utillib.golang/log"
	"strings"
)

type Input struct {
	Addr string `json:"addr" form:"addr" validate:"required"`
}

type Output struct {
	ErrNo        int                     `json:"err_no"`
	ErrMsg       string                  `json:"err_msg"`
	PenddingList []transaction.TransInfo `json:"penddings"`
}

func Main(cc echo.Context) error {

	//Step 1. init x
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Redis()
	c.Mysql()

	//Step 2. parameters initial
	var (
		input  Input
		output Output
	)
	if err := c.BindInput(&input); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}

	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	content, err := webthree.Txpool.Content()
	if err != nil {
		log.Fatalf("Content error:%s", err.Error())
		return err
	}

	for _, txmap := range content.Pending {
		for _, tx := range txmap {
			if tx.From == strings.ToLower(input.Addr) || tx.To == strings.ToLower(input.Addr) {

				tx_fee := big.NewInt(1).Mul(tx.GasPrice, tx.Gas)

				trans := transaction.TransInfo{
					TXHash:      tx.Hash,
					BlockNumber: 0,
					From:        tx.From,
					To:          tx.To,
					Value:       tx.Value.String(),
					TxFee:       tx_fee.String(),
					Nonce:       tx.Nonce.Int64(),
				}

				output.PenddingList = append(output.PenddingList, trans)
			}
		}
	}

	////find pending
	//pending_list, err := (&model.Pending{}).FindPendingByAddr(c.Mysql(), input.Addr)
	//if err != nil && err.Error() != DATA_NOT_EXIST {
	//	return c.RESULT_ERROR(ERR_DATABASE_ERROR, err.Error())
	//}
	//
	////return
	//for _, pending := range pending_list {
	//
	//	withNanos := "2006-01-02 15:04:05"
	//	t, err := time.ParseInLocation(withNanos, pending.F_create_time, time.Local)
	//	if err != nil {
	//		return c.RESULT_ERROR(ERR_DATABASE_ERROR, err.Error())
	//	}
	//
	//	trans := transaction.TransInfo{
	//		pending.F_tx_hash,
	//		0,
	//		t.Unix(),
	//		pending.F_from,
	//		pending.F_to,
	//		pending.F_value,
	//		pending.F_tx_fee,
	//		0,
	//	}
	//
	//	output.PenddingList = append(output.PenddingList, trans)
	//}
	output.ErrNo = 0
	output.ErrMsg = "success"

	return c.RESULT(output)
}
