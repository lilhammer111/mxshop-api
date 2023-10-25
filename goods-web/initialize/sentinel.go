package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func TrafficLimitRules() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v", err)
	}

	// todo 从nacos中读取限流规则配置，there is an official support by sentinel
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //匀速通过
			Threshold:              3,           //100ms只能就已经来了1W的并发， 1s就是10W的并发
			StatIntervalInMs:       6000,
		},
	})

	if err != nil {
		zap.S().Fatalf("加载规则失败: %v", err)
	}
}
