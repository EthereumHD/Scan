package get_summary

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	"github.com/EthereumHD/Scan/src/model"
	"github.com/EthereumHD/Scan/src/sync"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"math/big"
	"qoobing.com/utillib.golang/log"
	"time"
)

type InputReq struct {
}

type OutputRsp struct {
	ErrNo         int      `json:"err_no"`
	ErrMsg        string   `json:"err_msg"`
	Difficulty    *big.Int `json:"difficulty"`          //挖矿难度
	Capability    *big.Int `json:"capability"`          //算力容量(byte)
	OnlineMiner   int      `json:"online_miner"`        //在线矿工数
	PbDayReward   *big.Int `json:"pb_day_reward"`       //每PB日均收益(wei)
	BlockCount24H int      `json:"block_count_last24h"` //24小时爆块数
	BlockNumber   *big.Int `json:"block_number"`        //出块总数，几快高
	TotalRewarded *big.Int `json:"total_rewarded"`      //总已挖数量(wei)
	TotalMortgage *big.Int `json:"total_mortgage"`      //总抵押数量(wei)
}

func Main(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()
	c.Web3()

	//Step 2. parameters initial

	rsp := OutputRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	//Step 3. Get LastBlock
	lastbalck := sync.GLastBlock
	rsp.BlockNumber = lastbalck.Number

	//Step 4. Calculate Capability
	b := calcCapability(lastbalck.Difficulty)
	rsp.Difficulty = lastbalck.Difficulty
	rsp.Capability = b

	//Step 5. Get onlineMiner
	rsp.OnlineMiner = 1031
	if n, err := model.GetRecentOneDayBlockNumber(c.Mysql()); err != nil {
		log.Fatalf("GetTotalRewarded error:%s", err.Error())
	} else {
		rsp.BlockCount24H = int(n)
	}

	//Step 6. Get Reward/pb/day
	pb := big.NewInt(1024)
	if pb.Int64() == 0 {
		pb = pb.SetInt64(1)
	}
	reward := calcDayReward(c.Mysql())
	reward.Div(reward, pb)
	rsp.PbDayReward = reward

	//Step 7. Get total
	web3 := c.Web3()
	if totalRewarded, err := web3.Eth.GetTotalRewarded("latest"); err != nil {
		log.Fatalf("GetTotalRewarded error:%s", err.Error())
	} else {
		rsp.TotalRewarded = totalRewarded
	}
	if totalMortgage, err := web3.Eth.GetTotalMortgage("latest"); err != nil {
		log.Fatalf("GetTotalMortgage error:%s", err.Error())
	} else {
		rsp.TotalMortgage = totalMortgage
	}

	//返回结果
	return c.RESULT(rsp)
}

func calcCapability(difficulty *big.Int) (b *big.Int) {
	return new(big.Int).Mul(difficulty, big.NewInt(1456))
}

var cache = struct {
	reward     *big.Int
	updatetime time.Time
}{
	big.NewInt(0), time.Time{},
}

func calcDayReward(db *gorm.DB) *big.Int {
	now := time.Now()
	if !now.After(cache.updatetime.Add(120)) {
		return cache.reward
	}

	blocks, err := model.GetRecentBlocks(db, 0, ONEDAYBLOCK)
	if err != nil {
		log.Fatalf("GetRecentBlocks error:%s", err.Error())
		return big.NewInt(0)
	}

	reward := big.NewInt(0)
	for _, b := range blocks {
		blockreward, _ := big.NewInt(0).SetString(b.F_reward, 10)
		reward.Add(reward, blockreward)
	}
	cache.reward = reward
	cache.updatetime = now
	return reward
}
