package stats

import (
	"testing"
	"time"

	"github.com/EthereumHD/Scan/src/statistics/model"
	"github.com/EthereumHD/Scan/src/statistics/stats"
)

func TestCollector(t *testing.T) {
	col := collector{}
	if err := col.Collect(model.PingReport{"foo", time.Now()}); err != ErrNodeNotAuthorized {
		t.Errorf("collected unauthorized report: err=%q", err)
	}

	if err := col.Collect(model.AuthReport{ID: "foo"}); err != nil {
		t.Errorf("failed to collect auth: err=%q", err)
	}

	if err := col.Collect(model.PingReport{"foo", time.Now()}); err != nil {
		t.Errorf("failed to collect ping after auth: err=%q", err)
	}

}
