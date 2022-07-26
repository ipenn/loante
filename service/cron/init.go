package cron

import (
	"fmt"
	"loante/global"
	"loante/service/model"
	"loante/tools"
	"time"
)

const TimeLayout = "2006-01-02"
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
	t := tools.GetFormatTime()[0:10]
	borrow := new(model.Borrow)
	lists,_ := borrow.Page("b.status = 5",1,1000)
	for _, item := range lists{
		borrowTime, _ := time.Parse(TimeLayout, item.EndTime[0:10])
		//取出商户分配规则数据
		rules, _ := new(model.RemindRules).Page(fmt.Sprintf("rr.mch_id = %d", item.MchId),1,1000)
		//商户是否匹配
		//天数匹配
		for _, ruleItem := range rules{
			minDate := borrowTime.Add(time.Duration(ruleItem.MinDay) * time.Hour * 24).Format(TimeLayout)
			maxDate := borrowTime.Add(time.Duration(ruleItem.MaxDay+1) * time.Hour * 24).Format(TimeLayout)
			//tpdMin, _ := time.ParseDuration(fmt.Sprintf("%dd", ruleItem.MinDay))
			//minDate := borrowTime.Add(tpdMin).Format("2006-01-02")
			//tpdMax, _ := time.ParseDuration(fmt.Sprintf("%dd", ruleItem.MaxDay+1))
			//maxDate := borrowTime.Add(tpdMax).Format("2006-01-02")
			if t >=minDate && t< maxDate{
				//分配 borrow_id,mch_id,product_id,remind_id,remind_company_id,remind_group_id,remind_last_time,remind_assign_time,urge_id,urge_company_id,urge_group_id,urge_last_time,urge_assign_time
				visit :=new(model.BorrowVisit)
				visit.One(fmt.Sprintf("borrow_id = %d", item.Id))
				visit.MchId = item.MchId
				visit.ProductId = item.ProductId
				visit.UserId = item.Uid
				visit.RemindId = 0
				visit.RemindCompanyId = ruleItem.CompanyId
				if ruleItem.IsAutoApply == 1{ //是否自动派单
				}
				visit.RemindGroupId = ruleItem.GroupId
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
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		return
	}
	t := tools.GetFormatTime()[0:10]
	borrow := new(model.Borrow)
	lists,_ := borrow.Page("b.status >= 5 and b.status <=7 ",1,1000)
	for _, item := range lists{
		borrowTime, _ := time.Parse(TimeLayout, item.EndTime[0:10])
		//取出商户分配规则数据
		rules, _ := new(model.UrgeRules).Page(fmt.Sprintf("ur.mch_id = %d", item.MchId),1,1000)
		//商户是否匹配
		//天数匹配
		for _, ruleItem := range rules{
			minDate := borrowTime.Add(time.Duration(ruleItem.MinDay) * time.Hour * 24).Format(TimeLayout)
			maxDate := borrowTime.Add(time.Duration(ruleItem.MaxDay+1) * time.Hour * 24).Format(TimeLayout)
			if t >=minDate && t< maxDate{
				visit :=new(model.BorrowVisit)
				visit.One(fmt.Sprintf("borrow_id = %d", item.Id))
				visit.MchId = item.MchId
				visit.ProductId = item.ProductId
				visit.UserId = item.Uid
				visit.UrgeId = 0
				visit.UrgeCompanyId = ruleItem.CompanyId
				if ruleItem.IsAutoApply == 1{ //是否自动派单
				}
				visit.UrgeGroupId = ruleItem.GroupId
				visit.UrgeAssignTime = tools.GetFormatTime()
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

//BorrowRepaymentTimeout 扫描订单进入预期中并生成一笔滞纳金
func BorrowRepaymentTimeout()  {
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		return
	}
	t := tools.GetFormatTime()[0:10]
	borrow := new(model.Borrow)
	where := fmt.Sprintf("b.status = 5 and left(b.end_time,10) <= '%s'", t)
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
	ret,_ := res.RowsAffected()
	global.Log.Info("每天计算一次滞纳金 影响=%v err=%v", ret)
}

//BorrowExpireDay 更新剩余天数或者逾期天数
func BorrowExpireDay()  {
	if !(time.Now().Hour() == 0 && time.Now().Minute() == 0){
		//return
	}
	t0 := tools.GetFormatTime()[0:10]
	where := fmt.Sprintf("status > 4 and status < 8")
	res, err := global.C.DB.Exec(fmt.Sprintf("update borrow set expire_day = TIMESTAMPDIFF(DAY,end_time, '%s') where %s", t0, where))
	if err != nil{
		global.Log.Info("更新逾期天数 剩余天数 err=%v", err.Error())
		return
	}
	ret,_ := res.RowsAffected()
	global.Log.Info("更新逾期天数 剩余天数 影响=%v err=%v", ret)
	//
	//res, err = global.C.DB.Exec(fmt.Sprintf("update borrow set expire_day = TIMESTAMPDIFF(DAY, '%s',end_time) where %s and left(end_time,10) >= '%s'", t0, where, t0))
	//if err != nil{
	//	global.Log.Info("更新剩余天数 err=%v", err.Error())
	//	return
	//}
	//ret,_ = res.RowsAffected()
	//global.Log.Info("更新剩余天数 影响=%v err=%v", ret)
}

//RepaymentSMSNotify 还款提醒
func RepaymentSMSNotify()  {
	if time.Now().Hour() != 12 {
		return
	}
	lists,_ := new(model.Borrow).Page("b.expire_day > -2 and b.expire_day < 3", 1, 1000)
	global.Log.Info("发送短信开始")
	for _, item := range lists{
		t2 := 5
		var ps []string
		if item.ExpireDay == -1{
			t2 = 3
		}else if item.ExpireDay == 0{
			t2 = 4
		}
		//查询号码
		user := new(model.User)
		user.One(fmt.Sprintf("id = %d", item.Uid))
		if new(model.SmsTemplate).Send(t2, user.Phone, ps) {
			//判断商户是否存在
			merchantData := new(model.Merchant)
			merchantData.One(fmt.Sprintf("id = %d", item.MchId))
			if merchantData.Id > 0 {
				merchantData.AddService(1, 1) //扣费
			}
		}
	}
}

//StatUrgeTask 催收业绩统计
func StatUrgeTask() {
	//if time.Now().Hour() != 1 {
	//	return
	//}
	global.Log.Info("催收业绩统计开始")
	t0 := tools.ToAddDay(-1)[0:10]
	t := tools.GetFormatTime()[0:10]
	//更新历史的在库案件数
	new(model.StatUrge).Set("borrow_count", 0, "id > 0")
	//获取所有的催收员
	admins := []model.Admin{}
	global.C.DB.NewSelect().Model((*model.Admin)(nil)).
		Where("role_id = 3 or role_id = 7").Scan(global.C.Ctx, &admins)
	var stats []model.StatUrge
	for _, admin := range admins{
		gType := 0
		gUrgeId := admin.RemindId
		gUrgeGroupId := admin.RemindGroupId
		if admin.RoleId == 7{
			gType = 1
			gUrgeId = admin.UrgeId
			gUrgeGroupId = admin.UrgeGroupId
		}
		stats = append(stats, model.StatUrge{
			Type: gType,
			MchId: admin.MchId,
			UrgeId: admin.Id,
			UrgeCompanyId: gUrgeId,
			UrgeGroupId: gUrgeGroupId,
			StatTime: t,
		})
	}
	var ts []model.StatUrge
	//获取在库的催收业绩
	//全部催收
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as borrow_count").
		ColumnExpr("bv.mch_id").
		ColumnExpr("bv.urge_company_id").
		ColumnExpr("bv.urge_group_id").
		ColumnExpr("bv.urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where("urge_id > 0 and status = 7").GroupExpr("urge_id").Scan(global.C.Ctx, &ts)
	for _, item := range ts{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].BorrowCount = item.BorrowCount
			}
		}
	}

	var ts2 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).ColumnExpr("count(*) as borrow_count").
		ColumnExpr("bv.mch_id").
		ColumnExpr("bv.remind_company_id as urge_company_id").
		ColumnExpr("bv.remind_group_id as urge_group_id").
		ColumnExpr("bv.remind_id as urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where("bv.urge_id = 0 and borrow.status = 5").GroupExpr("remind_id").Scan(global.C.Ctx, &ts2)

	for _, item := range ts2{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].BorrowCount = item.BorrowCount
			}
		}
	}

	//新增案件数
	var ts3 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as borrow_new").
		ColumnExpr("urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where(fmt.Sprintf("urge_id > 0 and borrow.status = 7 and urge_assign_time = '%s'", t0)).GroupExpr("urge_id").Scan(global.C.Ctx, &ts3)
	for _, item := range ts3{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].BorrowNew = item.BorrowNew
			}
		}
	}
	var ts4 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as borrow_new").
		ColumnExpr("remind_id as urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where(fmt.Sprintf("urge_id = 0 and borrow.status = 5 and remind_assign_time = '%s'", t0)).GroupExpr("remind_id").Scan(global.C.Ctx, &ts4)
	for _, item := range ts4{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].BorrowNew = item.BorrowNew
			}
		}
	}
	//完成案件数
	var ts5 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as urge_closed_count").ColumnExpr("urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where(fmt.Sprintf("urge_id > 0 and borrow.status = 9 and left(complete_time,10) = '%s'", t0)).GroupExpr("urge_id").Scan(global.C.Ctx, &ts5)
	for _, item := range ts5{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].UrgeClosedCount = item.UrgeClosedCount
			}
		}
	}
	var ts6 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as urge_closed_count").ColumnExpr("remind_id as urge_id").
		Join("LEFT JOIN borrow").JoinOn("bv.borrow_id = borrow.id").
		Where(fmt.Sprintf("urge_id = 0 and borrow.status = 8 and left(complete_time, 10) = '%s'", t0)).GroupExpr("remind_id").Scan(global.C.Ctx, &ts6)
	for _, item := range ts6{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].UrgeClosedCount = item.UrgeClosedCount
			}
		}
	}
	//urge_closed_amount 回收金额
	var ts7 []model.Orders
	global.C.DB.NewSelect().Model((*model.Orders)(nil)).
		ColumnExpr("sum(actual_amount) as actual_amount").
		ColumnExpr("urge_id").
		Where(fmt.Sprintf("repaid_status = 1 and left(create_time, 10) = '%s'", t0)).GroupExpr("urge_id").Scan(global.C.Ctx, &ts7)
	for _, item := range ts7{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].UrgeClosedAmount = float64(item.ActualAmount)
			}
		}
	}

	//visit_count 跟进案件量
	var ts9 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as visit_count").ColumnExpr("urge_id").
		Where(fmt.Sprintf("urge_id > 0 and left(urge_last_time,10) = '%s'", t0)).GroupExpr("urge_id").Scan(global.C.Ctx, &ts9)
	for _, item := range ts9{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].VisitCount = item.VisitCount
			}
		}
	}

	var ts10 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisit)(nil)).
		ColumnExpr("count(*) as visit_count").ColumnExpr("remind_id as urge_id").
		Where(fmt.Sprintf("urge_id = 0 and left(remind_last_time,10) = '%s'", t0)).GroupExpr("remind_id").Scan(global.C.Ctx, &ts10)
	for _, item := range ts10{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].VisitCount = item.VisitCount
			}
		}
	}

	//visit_detail_count 催记记录次数
	var ts11 []model.StatUrge
	global.C.DB.NewSelect().Model((*model.BorrowVisitDetail)(nil)).
		ColumnExpr("count(*) as visit_detail_count").ColumnExpr("urge_id").
		Where(fmt.Sprintf("left(create_time,10) = '%s'", t0)).GroupExpr("urge_id").Scan(global.C.Ctx, &ts11)
	for _, item := range ts11{
		for k, it := range stats{
			if item.UrgeId == it.UrgeId{
				stats[k].VisitDetailCount = item.VisitDetailCount
			}
		}
	}
	for _, item := range ts{
		item.Insert()
	}
	global.Log.Info("催收业绩统计结束")
}