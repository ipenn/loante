package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"loante/global"
	v1 "loante/router/v1"
	"loante/service/middleware"
)

func Init() {
	authHandle := v1.NewAuth()
	systemHandle := v1.NewSystem()
	utmHandle := v1.NewUtmSource()
	merchantHandle := v1.NewMerchant()
	paymentHandle := v1.NewPayment()
	userHandle := v1.NewUser()
	borrowHandle := v1.NewBorrow()
	payHandle := v1.NewPayFlow()
	visitHandle := v1.NewVisit()
	appHandle := v1.NewApp()

	productHandle := v1.NewProduct()
	productDelayConfigHandle := v1.NewProductDelayConfig()
	smsTemplateHandle := v1.NewSmsTemplate()
	remindHandle := v1.NewRemind()
	urgeHandle := v1.NewUrge()
	payBackHandle := v1.NewPayBack()
	BlackHandle := v1.NewBlack()
	whitePhoneHandle := v1.NewWhitePhone()
	uploadHandle := v1.NewUpload()
	StatTrafficHandle := v1.NewStatTraffic()

	app := fiber.New(fiber.Config{
		//Prefork:       true,
		CaseSensitive: true,
		//StrictRouting: true,
		//ServerHeader:  "Fiber",
	})
	app.Use(cors.New())
	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v := app.Group("/v1")
	//登录
	v.Post("/auth/login", authHandle.Login)
	//给APP后端调用
	v.Post("/pay_notify", payBackHandle.PayNotify) //收款回调
	v.Post("/out_notify", payBackHandle.OutNotify) //放款回调
	v.Post("/sms/send", appHandle.SmsSend)         //发送短信
	v.Post("/pay/in", payHandle.PayPartial)        //代收
	v.Post("/pay/out", payHandle.PayOut)           //代付
	//上传
	v.Post("/upload", uploadHandle.Upload)              //普通上传
	v.Post("/upload/base64", uploadHandle.UploadBase64) //上传Base64

	//中间件
	v.Use(middleware.Auth)
	v.Use(middleware.Logger)

	//权限和系统设置
	v.Get("/side_menu", systemHandle.SideMenu)
	v.Get("/admins", systemHandle.AdminsList)
	v.Post("/admin/create", systemHandle.AdminCreate)
	v.Post("/admin/pwd_reset", systemHandle.PwdMchReset)
	v.Get("/rights", systemHandle.RightsList)
	v.Get("/roles", systemHandle.RolesList)
	v.Post("/roles/create", systemHandle.RoleCreate)
	v.Post("/roles/delete", systemHandle.RoleDelete)
	v.Get("/system/fields", systemHandle.SystemFields)
	v.Get("/system_setting", systemHandle.SystemSettingList)
	v.Get("/admin_log", systemHandle.AdminLogList)
	v.Post("/system_setting/update_value", systemHandle.SystemSettingUpdateValue)
	v.Get("/increase_rule", systemHandle.IncreaseRuleList)
	v.Post("/increase_rule/create_or_update", systemHandle.IncreaseRuleCreateOrUpdate)
	v.Post("/increase_rule/del", systemHandle.IncreaseRuleDel)
	v.Get("/packages", systemHandle.Packages)
	v.Get("/package/little", systemHandle.PackageLittle)

	//统计报表
	v.Get("/home_page_data", StatTrafficHandle.HomepageData)
	v.Get("/stat_traffic", StatTrafficHandle.StatTrafficList)
	v.Get("/after_report_unified", StatTrafficHandle.AfterReportUnifiedList)
	v.Get("/report_pay_in", StatTrafficHandle.ReportPayIn)
	v.Get("/report_pay_out", StatTrafficHandle.ReportPayOut)

	//财务管理
	v.Get("/expend_statistics", StatTrafficHandle.ExpendStatisticsList)
	v.Get("/merchant_statistics", StatTrafficHandle.MerchantStatisticsList)

	//渠道
	v.Get("/utm/lists", utmHandle.Lists)
	v.Post("/utm/create", utmHandle.Create)
	v.Post("/utm/modify", utmHandle.Modify)
	//渠道风险配置
	v.Get("/utm/risk_config", utmHandle.RiskConfig)
	v.Post("/utm/risk_create", utmHandle.RiskCreate)
	//商户
	v.Get("/merchant/list", merchantHandle.Lists)
	v.Post("/merchant/create", merchantHandle.Create)
	v.Post("/merchant/modify", merchantHandle.Modify)
	v.Get("/merchant/funds", merchantHandle.Funds)                              //商户资金列表 财务 -> 流水
	v.Post("/merchant/fund/create", merchantHandle.FundCreate)                  //商户充值和退款
	v.Get("/merchant/service_rule", merchantHandle.ServiceRule)                 //进件计价规则
	v.Post("/merchant/service_rule/create", merchantHandle.ServiceRuleCreate)   //进件计价规则创建
	v.Post("/merchant/service_rule/del", merchantHandle.ServiceRuleDel)         //进件计价规则删除
	v.Get("/merchant/service_price", merchantHandle.ServicePrice)               //服务定价
	v.Post("/merchant/service_price/create", merchantHandle.ServicePriceCreate) //服务定价创建
	v.Post("/merchant/service_price/del", merchantHandle.ServicePriceDel)       //服务定价删除

	//支付平臺
	v.Get("/payment/list", paymentHandle.Lists)
	v.Post("/payment/modify", paymentHandle.Modify)
	v.Post("/payment/set", paymentHandle.Set)
	v.Get("/payment/info", paymentHandle.Info)
	//商户的产品 支付的參數配置
	v.Get("/payment/config", paymentHandle.ConfigLists)
	v.Post("/payment/config/create", paymentHandle.ConfigCreate)
	v.Post("/payment/config/del", paymentHandle.ConfigDel)
	v.Post("/payment/config/set", paymentHandle.ConfigSet)
	//产品默认收 放款 支付通道 配置
	v.Get("/payment/default", paymentHandle.DefaultLists)
	v.Post("/payment/default/create", paymentHandle.DefaultCreate)
	v.Post("/payment/default/del", paymentHandle.DefaultDel)
	//客户管理
	v.Get("/user/list", userHandle.UserQuery)
	v.Get("/user/detail", userHandle.Details)
	v.Post("/user/black/set", userHandle.SetBlack)
	v.Get("/user/contact", userHandle.Contact) //获取用户的通讯录
	v.Get("/user/sms", userHandle.Sms)         //获取用户的短信
	v.Get("/user/app", userHandle.App)         //获取用户的APP
	v.Get("/user/call", userHandle.Call)       //获取用户的通话记录

	v.Get("/customer_feedback", userHandle.CustomerFeedBack)                            //客户反馈
	v.Post("/customer_feedback/update_status", userHandle.CustomerFeedBackUpdateStatus) //客户反馈

	v.Get("/visit/reminds", visitHandle.RemindBorrowAll)     //预提醒订单列表
	v.Get("/visit/reminding", visitHandle.RemindBorrowing)   //预提醒中订单
	v.Get("/visit/reminded", visitHandle.RemindBorrowed)     //预提醒完成的订单
	v.Get("/visit/remind_detail", visitHandle.RemindDetail)  //预提醒记录 一笔借贷可能会有多条记录
	v.Post("/visit/remind/action", visitHandle.RemindAction) //新增预提醒
	v.Post("/visit/remind/assign", visitHandle.RemindAssign) //预提醒分配

	v.Get("/visit/urges", visitHandle.UrgeBorrowAll)         //催收订单列表
	v.Get("/visit/urging", visitHandle.UrgeBorrowing)        //催收中订单
	v.Get("/visit/urged", visitHandle.UrgeBorrowed)          //催收完成的订单
	v.Get("/visit/urge_detail", visitHandle.UrgeDetail)      //催收记录 一笔借贷可能会有多条记录
	v.Get("/visit/urge_report", visitHandle.UrgeReport)      //催收业绩
	v.Post("/visit/urge/action", visitHandle.UrgeAction)     //新增催记
	v.Post("/visit/remind/action", visitHandle.RemindAction) //新增预提醒
	v.Post("/visit/utr/create", visitHandle.UtrCreate)       //utr新增记录
	v.Post("/visit/utr/examine", visitHandle.UtrExamine)     //utr对账单审核

	//借贷
	v.Get("/borrow/list", borrowHandle.Query)                     //获取借贷信息列表
	v.Get("/borrow/export", borrowHandle.QueryExport)             //获取借贷信息导出的功能
	v.Post("/borrow/reconciliation", borrowHandle.Reconciliation) //平账操作
	v.Post("/borrow/deposit", borrowHandle.Deposit)               //入账操作
	v.Post("/borrow/funds", borrowHandle.Funds)                   //费用变更
	v.Post("/borrow/set_loan/fail", borrowHandle.SetLoanFail)     //设置放款失败/进入重放款

	//还款
	v.Get("/pay_flow/repayments", payHandle.Repayments)              //还款记录
	v.Get("/pay_flow/repayments/export", payHandle.RepaymentsExport) //导出还款记录

	v.Get("/pay_flow/reconciliation", payHandle.Reconciliation) //平账
	v.Get("/pay_flow/deposit", payHandle.Deposits)              //入账
	v.Get("/pay_flow/loan", payHandle.Loans)                    //放款
	v.Get("/pay_flow/batch_loan", payHandle.BatchLoans)         //批量重放款
	v.Get("/pay_flow/utr", payHandle.Utrs)                      //UTR对账单
	v.Get("/pay_flow/utr_dismissed", payHandle.UtrsDismissed)   //UTR驳回列表 UTR对账单验证失败的
	v.Post("/pay_flow/utr_verify", payHandle.UtrsVerify)        //UTR对账单验证
	v.Post("/pay_flow/pay_partial", payHandle.PayPartial)       //生成 部分还款链接

	//短信模板
	v.Get("/sms_template", smsTemplateHandle.SmsTemplate)
	v.Post("/sms_template/create_or_update", smsTemplateHandle.SmsTemplateCreateOrUpdate)
	//产品运营
	v.Get("/product", productHandle.Product)
	v.Post("/product/create_or_update", productHandle.ProductCreateOrUpdate)
	//提额规则
	v.Get("/product/precept", productHandle.ProductPrecept)
	v.Post("/product/precept_create", productHandle.ProductPreceptCreate)
	v.Post("/product/precept_del", productHandle.ProductPreceptDel)
	v.Post("/product/create_or_update", productHandle.ProductCreateOrUpdate) //修改产品
	v.Post("/product/update_For_mch", productHandle.ProductUpdateForMch)     //修改产品(商户)
	//产品配置
	v.Get("/productDelayConfig", productDelayConfigHandle.ProductDelayConfig)
	v.Post("/productDelayConfig/create_or_update", productDelayConfigHandle.ProductDelayConfigCreateOrUpdate)

	//预提醒管理
	v.Get("/remind_company", remindHandle.RemindCompany)
	v.Post("/remind_company/create", remindHandle.RemindCompanyCreate)
	v.Post("/remind_company/update", remindHandle.RemindCompanyUpdate)
	v.Get("/remind_group", remindHandle.RemindGroup)
	v.Post("/remind_group/create", remindHandle.RemindGroupCreate)
	v.Post("/remind_group/update", remindHandle.RemindGroupUpdate)
	v.Get("/remind_admin", remindHandle.RemindAdmin)
	v.Get("/remind_rules", remindHandle.RemindRules)
	v.Post("/remind_rules/create_or_update", remindHandle.RemindRulesCreateOrUpdate)
	//催收管理
	v.Get("/urge_company", urgeHandle.UrgeCompany)
	v.Post("/urge_company/create", urgeHandle.UrgeCompanyCreate)
	v.Post("/urge_company/update", urgeHandle.UrgeCompanyUpdate)
	v.Get("/urge_group", urgeHandle.UrgeGroup)
	v.Post("/urge_group/create", urgeHandle.UrgeGroupCreate)
	v.Post("/urge_group/update", urgeHandle.UrgeGroupUpdate)
	v.Get("/urge_admin", urgeHandle.UrgeAdmin)
	v.Get("/urge_rules", urgeHandle.UrgeRules)
	v.Post("/urge_rules/create_or_update", urgeHandle.UrgeRulesCreateOrUpdate)
	//黑名单管理
	v.Post("/user_black/create", BlackHandle.UserBlackCreate)
	v.Post("/user_black/del", BlackHandle.UserBlackDel)
	v.Get("/user_black", BlackHandle.UserBlackList)
	v.Get("/regional_black", BlackHandle.RegionalBlack) //区域黑名单
	//白名单管理
	v.Post("/white_phone/create", whitePhoneHandle.WhitePhoneCreate)
	v.Post("/white_phone/del", whitePhoneHandle.WhitePhoneDel)
	v.Get("/white_phone", whitePhoneHandle.WhitePhoneList)
	app.Listen(fmt.Sprintf(":%d", global.C.Http.Port))
}
