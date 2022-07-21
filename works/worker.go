package main

import (
	"fmt"
	"github.com/robfig/cron"
	_ "loante/global"
	"loante/tools"
	"os"
	"os/signal"
)
func main()  {
	c := cron.New()
	//1 每分钟 查找支付超时的订单
	c.AddFunc("@every 1m", func() {
		fmt.Println(tools.GetFormatTime())
	})

	c.AddFunc("* 1 * * * *", func() {
		fmt.Println(tools.GetFormatTime())
	})

	//2 每天 扫描进入预提醒订单

	//3 每天 扫描进入催收订单

	//4 每天 扫描进入预期中的订单

	//5 每天 扫描 计算滞纳金

	c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}


