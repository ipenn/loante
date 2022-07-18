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
	authHaddle := v1.NewAuth()
	systemHaddle := v1.NewSystem()
	utmHaddle := v1.NewUtmSource()
	merchantHaddle := v1.NewMerchant()
	paymentHaddle := v1.NewPayment()
	userHaddle := v1.NewUser()
	borrowHandle := v1.NewBorrow()
	payHandle := v1.NewPayFlow()
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
	v.Post("/auth/login", authHaddle.Login)

	v.Use(middleware.Auth)
	v.Use(middleware.Logger)
	v.Get("/side_menu", systemHaddle.SideMenu)
	v.Get("/admins", systemHaddle.AdminsList)
	v.Post("/admin/create", systemHaddle.AdminCreate)
	v.Get("/rights", systemHaddle.RightsList)
	v.Get("/roles", systemHaddle.RolesList)
	v.Post("/roles/create", systemHaddle.RoleCreate)
	v.Post("/roles/delete", systemHaddle.RoleDelete)
	v.Get("/system/fields", systemHaddle.SystemFields)

	//渠道
	v.Get("/utm/lists", utmHaddle.Lists)
	v.Post("/utm/create", utmHaddle.Create)
	v.Post("/utm/modify", utmHaddle.Modify)
	//渠道风险配置
	v.Get("/utm/risk_config", utmHaddle.RiskConfig)
	v.Post("/utm/risk_create", utmHaddle.RiskCreate)
	//商户
	v.Get("/merchant/list", merchantHaddle.Lists)
	v.Post("/merchant/create", merchantHaddle.Create)
	v.Post("/merchant/modify", merchantHaddle.Modify)
	v.Post("/merchant/fund/create", merchantHaddle.FundCreate)
	v.Get("/merchant/service_rule", merchantHaddle.ServiceRule) //进件计价规则
	v.Post("/merchant/service_rule/create", merchantHaddle.ServiceRuleCreate) //进件计价规则创建
	v.Post("/merchant/service_rule/del", merchantHaddle.ServiceRuleDel) //进件计价规则删除
	v.Get("/merchant/service_price", merchantHaddle.ServicePrice) //服务定价
	v.Post("/merchant/service_price/create", merchantHaddle.ServicePriceCreate) //服务定价创建
	v.Post("/merchant/service_price/del", merchantHaddle.ServicePriceDel) //服务定价删除

	//支付平臺
	v.Get("/payment/list", paymentHaddle.Lists)
	v.Post("/payment/modify", paymentHaddle.Modify)
	v.Post("/payment/set", paymentHaddle.Set)
	//商户的产品 支付的參數配置
	v.Get("/payment/config", paymentHaddle.ConfigLists)
	v.Post("/payment/config/create", paymentHaddle.ConfigCreate)
	v.Post("/payment/config/del", paymentHaddle.ConfigDel)
	v.Post("/payment/config/set", paymentHaddle.ConfigSet)
	//产品默认收 放款 支付通道 配置
	v.Get("/payment/default", paymentHaddle.DefaultLists)
	v.Post("/payment/default/create", paymentHaddle.DefaultCreate)
	v.Post("/payment/default/del", paymentHaddle.DefaultDel)
	//客户管理
	v.Get("/user/list", userHaddle.UserQuery)
	v.Get("/user/detail", userHaddle.Details)

	//借贷
	v.Get("/borrow/list", borrowHandle.Query)

	//还款
	v.Get("/pay_flow/repayments", payHandle.Repayments) //还款记录
	v.Get("/pay_flow/repayments/export", payHandle.RepaymentsExport) //导出还款记录
	v.Get("/pay_flow/reconciliation", payHandle.Reconciliation)//平账
	v.Get("/pay_flow/deposit", payHandle.Deposits)//入账
	v.Get("/pay_flow/loan", payHandle.Loans) //放款
	v.Get("/pay_flow/batch_loan", payHandle.BatchLoans) //批量重放款


	app.Listen(fmt.Sprintf(":%d", global.C.Http.Port))
}