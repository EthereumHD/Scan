package mining

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"github.com/labstack/echo"
	"math/big"
	"time"
)

type Input struct {
	Addr       string `json:"addr" form:"addr" validate:"required"`
	StartDate  string `json:"start_date" form:"start_date" validate:"required"`
	EndDate    string `json:"end_date" form:"end_date" validate:"required"`
	OffsetTime int64  `json:"offset_time" form:"offset_time"`
}

type dayinfo struct {
	Date   string `json:"date"`
	Fees   string `json:"fees"`
	Reward string `json:"reward"`
	Count  uint64 `json:"count"`
}

type Output struct {
	ErrNo        int       `json:"err_no"`
	ErrMsg       string    `json:"err_msg"`
	TotalReward  string    `json:"total_reward"`
	LastXhReward string    `json:"last_x_reward"`
	LastXhFees   string    `json:"last_x_fees"`
	DayInfo      []dayinfo `json:"day_info"`
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
	withDate := "2006-01-02"
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

		rewards := big.NewInt(0)
		fees := big.NewInt(0)
		count := uint64(0)
		for _, block := range blocks {

			reward, success := big.NewInt(0).SetString(block.F_reward, 10)
			if !success {
				return c.RESULT_ERROR(ERR_INNER_ERROR, "Failed to convert %s to BigInt"+block.F_reward)
			}

			fee, success := big.NewInt(0).SetString(block.F_fees, 10)
			if !success {
				return c.RESULT_ERROR(ERR_INNER_ERROR, "Failed to convert %s to BigInt"+block.F_reward)
			}

			rewards = big.NewInt(0).Add(rewards, reward)
			fees = big.NewInt(0).Add(fees, fee)
			count++
		}

		if rewards.Cmp(big.NewInt(0)) == 0 && fees.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		output.DayInfo = append(output.DayInfo,
			dayinfo{
				Date:   day.Format(withDate),
				Fees:   fees.String(),
				Reward: rewards.String(),
				Count:  count,
			},
		)
	}

	miner_reward, err := (&MinerReward{}).FindRewardByMiner(c.Mysql(), input.Addr)
	if err != nil {
		if err.Error() != DATA_NOT_EXIST {
			return c.RESULT_ERROR(ERR_DATABASE_ERROR, err.Error())
		} else {
			miner_reward.F_total_reward = "0"
			miner_reward.F_total_fees = "0"
		}
	}

	//查询过去24小时
	var (
		tnow           = time.Now()
		laststart      = tnow.Add(-24 * time.Hour).Unix()
		lastend        = tnow.Unix()
		last_x_rewards = big.NewInt(0)
		last_x_fees    = big.NewInt(0)
	)
	blocks, err := (&Block{}).FindMinedBlockByAddrAndTime(c.Mysql(), input.Addr, laststart, lastend)
	if err != nil && err.Error() != DATA_NOT_EXIST {
		return c.RESULT_ERROR(ERR_DATABASE_ERROR, err.Error())
	}

	for _, block := range blocks {

		reward, success := big.NewInt(0).SetString(block.F_reward, 10)
		if !success {
			return c.RESULT_ERROR(ERR_INNER_ERROR, "Failed to convert %s to BigInt"+block.F_reward)
		}

		fees, success := big.NewInt(0).SetString(block.F_fees, 10)
		if !success {
			return c.RESULT_ERROR(ERR_INNER_ERROR, "Failed to convert %s to BigInt"+block.F_reward)
		}

		last_x_rewards = big.NewInt(0).Add(last_x_rewards, reward)
		last_x_fees = big.NewInt(0).Add(last_x_fees, fees)
	}

	// return
	output.ErrNo = 0
	output.ErrMsg = "success"
	output.TotalReward = miner_reward.F_total_reward
	output.LastXhFees = last_x_fees.String()
	output.LastXhReward = last_x_rewards.String()

	return c.RESULT(output)
}
