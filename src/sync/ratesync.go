package sync

import (
	"github.com/EthereumHD/Scan/src/config"
	"github.com/EthereumHD/Scan/src/model"
	"errors"
	"qoobing.com/utillib.golang/log"
	"time"
)

type RateRecord struct {
	Exchange string  `json:"exchange"`
	Rate     float64 `json:"rate"`
	Weight   float64 `json:"weight"`
}

type RateManager struct {
	EthRateRecord map[string]RateRecord
	BtcRateRecord map[string]RateRecord
	USDRateRecord map[string]RateRecord
	KWRRateRecord map[string]RateRecord
}

var manager RateManager

func (m *RateManager) reset() {
	m.EthRateRecord = make(map[string]RateRecord)
	m.BtcRateRecord = make(map[string]RateRecord)
	m.USDRateRecord = make(map[string]RateRecord)
	m.KWRRateRecord = make(map[string]RateRecord)
}

func (m *RateManager) Add(currency string, r RateRecord) error {

	log.Debugf("query token rate,%s:%+v", currency, r)
	switch currency {
	case "eth":
		m.EthRateRecord[r.Exchange] = r
	case "btc":
		m.BtcRateRecord[r.Exchange] = r
	case "usd":
		m.USDRateRecord[r.Exchange] = r
	case "kwr":
		m.KWRRateRecord[r.Exchange] = r
	default:
		log.Fatalf("not find currency:%s", currency)
		return errors.New("currency is not exist")
	}

	return nil
}

func (m *RateManager) countrate() (rate model.UbbeyRate) {

	var weight_eth float64
	for _, v := range m.EthRateRecord {
		weight_eth = weight_eth + v.Weight
	}
	for _, v := range m.EthRateRecord {
		rate.Eth = rate.Eth + v.Rate*(v.Weight/weight_eth)
	}

	return rate
}

func StartSyncRate() {

	for {
		//manager.reset()
		//getRateFromIndex()
		//
		//time.Sleep(time.Second * 4)
		//rate := manager.countrate()
		//
		//err := model.SetRate(c.Redis(), rate)
		//if err != nil {
		//	log.Fatalf("SetRate error:%s", err.Error())
		//	continue
		//}
		//
		//r := model.Rate{
		//	F_eth: rate.Eth,
		//	F_btc: rate.Btc,
		//	F_usd: rate.USD,
		//	F_kwr: rate.KWR,
		//}
		//err = r.CreateRate(c.Mysql())
		//if err != nil {
		//	log.Fatalf("CreateRate error:%s", err.Error())
		//	continue
		//}

		time.Sleep(time.Second * time.Duration(config.Config().RateSyncInterval))
	}

}

func getRateFromIndex() {

}
