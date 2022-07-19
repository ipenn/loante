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

	productHandle := v1.NewProduct()
	productDelayConfigHandle := v1.NewProductDelayConfig()
	smsTemplateHandle := v1.NewSmsTemplate()
	remindHandle := v1.NewRemind()
	urgeHandle := v1.NewUrge()
	UserBlackHandle := v1.NewUserBlack()
	whitePhoneHandle := v1.NewWhitePhone()

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

	//中间件
	v.Use(middleware.Auth)
	v.Use(middleware.Logger)

	//权限和系统设置
	v.Get("/side_menu", systemHandle.SideMenu)
	v.Get("/admins", systemHandle.AdminsList)
	v.Post("/admin/create", systemHandle.AdminCreate)
	v.Get("/rights", systemHandle.RightsList)
	v.Get("/roles", systemHandle.RolesList)
	v.Post("/roles/create", systemHandle.RoleCreate)
	v.Post("/roles/delete", systemHandle.RoleDelete)
	v.Get("/system/fields", systemHandle.SystemFields)

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
	v.Post("/merchant/fund/create", merchantHandle.FundCreate)
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

	//借贷
	v.Get("/borrow/list", borrowHandle.Query)

	//还款
	v.Get("/pay_flow/repayments", payHandle.Repayments)              //还款记录
	v.Get("/pay_flow/repayments/export", payHandle.RepaymentsExport) //导出还款记录
	v.Get("/pay_flow/reconciliation", payHandle.Reconciliation)      //平账
	v.Get("/pay_flow/deposit", payHandle.Deposits)                   //入账
	v.Get("/pay_flow/loan", payHandle.Loans)                         //放款
	v.Get("/pay_flow/batch_loan", payHandle.BatchLoans)              //批量重放款

	//短信模板
	v.Get("/sms_template", smsTemplateHandle.SmsTemplate)
	v.Post("/sms_template/create_or_update", smsTemplateHandle.SmsTemplateCreateOrUpdate)
	//产品运营
	v.Get("/product", productHandle.Product)
	v.Post("/product/create_or_update", productHandle.ProductCreateOrUpdate)
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
	v.Post("/user_black/create", UserBlackHandle.UserBlackCreate)
	v.Post("/user_black/del", UserBlackHandle.UserBlackDel)
	v.Get("/user_black", UserBlackHandle.UserBlackList)
	//白名单管理
	v.Post("/white_phone/create", whitePhoneHandle.WhitePhoneCreate)
	v.Post("/white_phone/del", whitePhoneHandle.WhitePhoneDel)
	v.Get("/white_phone", whitePhoneHandle.WhitePhoneList)
	app.Listen(fmt.Sprintf(":%d", global.C.Http.Port))
}
