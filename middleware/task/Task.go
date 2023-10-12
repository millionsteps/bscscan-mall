package task

import (
	"fmt"

	cron "github.com/robfig/cron/v3"
	"main.go/service/mall"
)

func InitTask() {
	c := cron.New(cron.WithSeconds()) //精确到秒

	//定时任务
	spec := "0 */1 * * * ?" //cron表达式，每分钟一次
	c.AddFunc(spec, func() {
		fmt.Println("11111---spec")
		mall.ReleaseUsdt()
	})

	specObject := "0 */3 * * * ?" //cron表达式，每三分钟一次
	c.AddFunc(specObject, func() {
		// fmt.Println("11111---specObject")
		// object.CheckObjectSaleExpire()
	})
	c.Start()
	select {} //阻塞主线程停止
}
