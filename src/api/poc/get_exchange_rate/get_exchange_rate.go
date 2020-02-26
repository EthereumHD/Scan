package get_exchange_rate

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/EhdChain/log"
	"github.com/labstack/echo"
)

type InputReq struct {
}

type OutputRsp struct {
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
	Rates  RateList `json:"rates"`
}

type RateInfo struct {
	Currency    string  `json:"currency"`    //货币类型 ,RMB,USD,KRW
	Rate        float64 `json:"rate"`        //如 0.05 1 poc = 0.05RMB
	Significand int     `json:"significand"` //资产显示小数点尾数
	Symbol      string  `json:"symbol"`      //资产符号
}

type RateList []RateInfo

func Main(cc echo.Context) error {
	log.Debug("get_exchange_rate Main1")
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	log.Debug("get_exchange_rate Main2")
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	//查询区块,数据库查询
	//包装参数
	//get currency
	rsp.Rates = append(rsp.Rates, RateInfo{Currency: "USD", Rate: 0.32, Significand: 2, Symbol: "$"})
	rsp.Rates = append(rsp.Rates, RateInfo{Currency: "RMB", Rate: 2.07, Significand: 2, Symbol: "¥"})
	rsp.Rates = append(rsp.Rates, RateInfo{Currency: "ETH", Rate: 0.00003, Significand: 5, Symbol: "ETH"})
	rsp.Rates = append(rsp.Rates, RateInfo{Currency: "BTC", Rate: 0.0000002, Significand: 8, Symbol: "BTC"})

	//返回结果
	return c.RESULT(rsp)
}
