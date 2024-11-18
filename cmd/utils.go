package cmd

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

func Execute() {

	// 实例化调度器对象
	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	// 定时解析并更新线上服务器列表
	_, _ = s.NewJob(
		gocron.CronJob(
			fmt.Sprintf("%d 3-20/2 * * *", _GLO_RAND.Intn(60)),
			false, // without Seconds.
		),
		gocron.NewTask(func() { parseServerList() }),
	)

	// 服务维护状态下的端口健康监控检查
	_, _ = s.NewJob(
		gocron.CronJob(
			"*/5 * 10-21 * * *", true, // with Seconds.
		),
		gocron.NewTask(func() { healthCheckWithinPortClosed() }),
	)

	// 服务运行状态下的端口健康监控检查
	_, _ = s.NewJob(
		gocron.CronJob(
			// "*/5 * 10-21 * * *", true, // with Seconds. - for dev testing
			"45-59/1 7-9/1 * * *", false, // without Seconds.
		),
		gocron.NewTask(func() { healthCheckWithinPortOpened() }),
	)

	s.Start() // 开始调度任务
	select {} // 阻塞主 goroutine，直到任务全部执行完成
}
