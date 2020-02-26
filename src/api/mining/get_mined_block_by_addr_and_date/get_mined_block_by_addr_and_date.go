package get_mined_block_by_addr_and_date

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"github.com/labstack/echo"
	"time"
)

type Input struct {
	Addr       string `form:"addr" validate:"required"`
	StartDate  string `form:"start_date" validate:"required"`
	EndDate    string `form:"end_date" validate:"required"`
	OffsetTime int64  `json:"offset_time" form:"offset_time"`
}

type Output struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Blocks []BlockInfo `json:"blocks"`
}

type BlockInfo struct {
	BlockNumber int64  `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`
	Txn         int64  `json:"txn"`
	BlockFees   string `json:"block_fees"`
	BlockReward string `json:"block_reward"`
	GasUsed     string `json:"gas_used"`
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

	//check date ,不能超过100 天，必须是凌晨整点
	withSecond := "2006-01-02 15:04:05"
	start, err := time.ParseInLocation(withSecond, input.StartDate, time.UTC)
	if err != nil {
		return c.RESULT_PARAMETER_ERROR(input.StartDate + " time error")
	}

	end, err := time.ParseInLocation(withSecond, input.EndDate, time.UTC)
	if err != nil {
		return c.RESULT_PARAMETER_ERROR(input.EndDate + " time error")
	}

	if end.Unix()-start.Unix() > 8640000 { //100天
		return c.RESULT_PARAMETER_ERROR("day range should less than 100 day")
	}

	if start != time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC) ||
		end != time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC) {
		return c.RESULT_PARAMETER_ERROR("date should be like 2006-01-02 00:00:00")
	}

	//查询每天
	for day := start; day.Unix() <= end.Unix(); day = day.Add(time.Hour * 24) {
		blocks, err := (&Block{}).FindMinedBlockByAddrAndTime(c.Mysql(), input.Addr, day.Unix()+input.OffsetTime, day.Add(time.Hour*24).Unix()-1+input.OffsetTime)
		if err != nil && err.Error() != DATA_NOT_EXIST {
			return c.RESULT_ERROR(ERR_DATABASE_ERROR, err.Error())
		}

		for _, block := range blocks {
			block_info := BlockInfo{
				BlockNumber: block.F_block,
				Timestamp:   block.F_timestamp,
				Txn:         block.F_txn,
				BlockFees:   block.F_fees,
				BlockReward: block.F_reward,
				GasUsed:     block.F_gas_used,
			}

			output.Blocks = append(output.Blocks, block_info)

		}
	}

	// return
	output.ErrNo = 0
	output.ErrMsg = "success"

	return c.RESULT(output)
}
