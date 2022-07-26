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
	//c2.OrdersPayTimeOut()
	//c2.BorrowToRemind()
	//c2.BorrowToUrge()
	//c2.BorrowRepaymentTimeout()
	//c2.BorrowLatePaymentFee()
	//c2.BorrowExpireDay()
	c2.StatUrgeTask()

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

	//6 每天零点更新剩余天数
	c.AddFunc("@every 1m", c2.BorrowExpireDay)

	//7 发送还款的通知短信
	c.AddFunc("@every 1h", c2.RepaymentSMSNotify)

	//8 催收业绩统计
	c.AddFunc("@every 1h", c2.StatUrgeTask)



	c.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}


