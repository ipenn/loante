package cron

import (
	"fmt"
	"loante/global"
	"loante/service/model"
	"loante/tools"
	"time"
)

//OrdersPayTimeOut 付款超时的订单 默认30分钟
func OrdersPayTimeOut()  {
	t := tools.GetFormatTime()
	orders := new(model.Orders)
	where := fmt.Sprintf("invalid_time < '%s' and repaid_status = 2", t)
	lists,_ := orders.Page(where, 1, 20)
	for _, item := range lists{
		item.RepaidStatus = 0
		item.Update(fmt.Sprintf("%s and id = %d", where, item.Id))
		global.Log.Info("OrdersPayTimeOut 付款超时的订单状态标记为失败 order_id = %d",item.Id)
	}
}

//BorrowToRemind 扫描订单进入预提醒
func BorrowToRemind(){
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		return
	}
	t := tools.GetFormatTime()[0:7]
	borrow := new(model.Borrow)
	lists,_ := borrow.Page("status = 5",1,1000)
	for _, item := range lists{
		//nowTime, _ := time.Parse("2006-01-02", t)
		borrowTime, _ := time.Parse("2006-01-02", borrow.EndTime)
		//取出商户分配规则数据
		rules, _ := new(model.RemindRules).Page(fmt.Sprintf("mch_id = %d", item.MchId),1,1000)
		//商户是否匹配
		//天数匹配
		for _, ruleItem := range rules{
			tpdMin, _ := time.ParseDuration(fmt.Sprintf("%dd", ruleItem.MinDay))
			minDate := borrowTime.Add(tpdMin).Format("2006-01-02")
			tpdMax, _ := time.ParseDuration(fmt.Sprintf("%dd", ruleItem.MaxDay+1))
			maxDate := borrowTime.Add(tpdMax).Format("2006-01-02")
			if t >=minDate && t< maxDate{
				//分配 borrow_id,mch_id,product_id,remind_id,remind_company_id,remind_group_id,remind_last_time,remind_assign_time,urge_id,urge_company_id,urge_group_id,urge_last_time,urge_assign_time
				visit :=new(model.BorrowVisit)
				visit.One(fmt.Sprintf("borrow_id = %d", item.Id))
				visit.MchId = borrow.MchId
				visit.ProductId = borrow.ProductId
				visit.RemindId = 0
				visit.RemindCompanyId = ruleItem.CompanyId
				if ruleItem.IsAutoApply == 1{ //是否自动派单
					visit.RemindGroupId = ruleItem.GroupId
				}
				visit.RemindAssignTime = tools.GetFormatTime()
				if visit.BorrowId > 0{
					visit.Update(fmt.Sprintf("borrow_id = %d", item.Id))
				}else{
					visit.BorrowId = item.Id
					visit.Insert()
				}
				global.Log.Info("BorrowToRemind 预提醒分配 borrow_id = %d",item.Id)
				break
			}
		}
	}
}

//BorrowToUrge 扫描订单进入催收
func BorrowToUrge()  {
	
}

//BorrowRepaymentTimeout 扫描订单进入预期中并生成一笔滞纳金
func BorrowRepaymentTimeout()  {
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		return
	}
	t := tools.GetFormatTime()[0:10]
	borrow := new(model.Borrow)
	where := fmt.Sprintf("status = 5 and left(end_time,10) <= '%s'", t)
	lists,_ := borrow.Page(where, 1, 1000)
	for _, item := range lists{
		item.Status = 7
		fee := int(float64(item.LoanAmount) * item.LatePaymentFeeRate)
		item.LatePaymentFee += fee
		item.Update(fmt.Sprintf("id = %d", item.Id))
		global.Log.Info("BorrowRepaymentTimeout 订单进入预期中并生成一笔滞纳金 borrow_id = %d fee=%d",item.Id, fee)
	}
}

//BorrowLatePaymentFee 每天计算一次滞纳金
func BorrowLatePaymentFee()  {
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		return
	}
	time.Sleep(5*time.Second)
	t0 := tools.GetFormatTime()[0:10]
	where := fmt.Sprintf("status = 7 and left(end_time,10) < '%s'", t0)
	res, err := global.C.DB.Exec(fmt.Sprintf("update borrow set late_payment_fee = late_payment_fee + late_payment_fee_rate*loan_amount where %s", where))
	if err != nil{
		global.Log.Info("每天计算一次滞纳金 err=%v", err.Error())
		return
	}
	ret,err := res.RowsAffected()
	global.Log.Info("每天计算一次滞纳金 影响=%v err=%v", ret, err.Error())
}