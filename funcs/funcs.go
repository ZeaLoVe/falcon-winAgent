package funcs

import (
	"github.com/ZeaLoVe/falcon-winAgent/g"
	"github.com/open-falcon/common/model"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				MemMetrics,
				DiskIOMetrics,
				DeviceMetrics,
				ProcMetrics,
			},
			Interval: interval,
		},
	}
}
