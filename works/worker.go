package main

import (
	"github.com/robfig/cron"
	_ "loante/global"
	c2 "loante/service/cron"
	"os"
	"os/signal"
)
func main()  {
	c := cron.New()
	//c.AddFunc("* 1 * * * *", func() {
	//	fmt.Println(tools.GetFormatTime())
	//})

	//1 每分钟 查找支付超时的订单
	c.AddFunc("@every 1m", c2.OrdersPayTimeOut)

	//2 每天 扫描进入预提醒订单
	c.AddFunc("@every 1m", c2.BorrowToRemind)

	//3 每天 扫描进入催收订单
	c.AddFunc("@every 1m", c2.BorrowToUrge)

	//4 每天 扫描进入预期中的订单
	c.AddFunc("@every 1m", c2.BorrowRepaymentTimeout)

	//5 每天 扫描 计算滞纳金
	c.AddFunc("@every 1m", c2.BorrowLatePaymentFee)

	c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}


