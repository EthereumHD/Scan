package stats

import (
	"github.com/EthereumHD/Scan/src/config"
	"github.com/EthereumHD/Scan/src/statistics/model"
	"net"
	"net/http"
	"qoobing.com/utillib.golang/log"
)

func Start() {
	log.Debugf("stats start ...")
	ethstats := &Server{
		Name: model.ID(config.Config().Stats.ServerId),
	}

	addr := config.Config().Stats.StatAddr
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("failed to parse stats.StatAddr, err:%s", err.Error())
		panic("failed to parse address")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", ethstats.WebsocketHandler)
	mux.HandleFunc("/", ethstats.APIHandler)

	log.Debugf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("stats ListenAndServe error:%s", err.Error())
		panic("stats ListenAndServe error")
	}
}
