package cron

import (
	"github.com/ZeaLoVe/falcon-winAgent/funcs"
	"github.com/ZeaLoVe/falcon-winAgent/g"
	"github.com/open-falcon/common/model"
	"time"
)

func InitDataHistory() {
	for {
		//如果有数据需要两次调用才能采集的，需要在这里提前初始化一次
		funcs.UpdateIfStat()
		time.Sleep(time.Duration(g.Config().Transfer.Interval) * time.Second)
	}
}

func Collect() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if g.Config().Transfer.Addr == "" {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {

	for {
	REST:
		time.Sleep(time.Duration(sec) * time.Second)

		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		mvs := []*model.MetricValue{}
		ignoreMetrics := g.Config().IgnoreMetrics

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				if b, ok := ignoreMetrics[mv.Metric]; ok && b {
					continue
				} else {
					mvs = append(mvs, mv)
				}
			}
		}

		now := time.Now().Unix()
		for j := 0; j < len(mvs); j++ {
			mvs[j].Step = sec
			mvs[j].Endpoint = hostname
			mvs[j].Timestamp = now
		}

		g.SendToTransfer(mvs)

	}
}
