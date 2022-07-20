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
	productHaddle := v1.NewProduct()
	productDelayConfigHaddle := v1.NewProductDelayConfig()
	smsTemplateHaddle := v1.NewSmsTemplate()
	remindHaddle := v1.NewRemind()
	urgeHaddle := v1.NewUrge()
	app := fiber.New()
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
	//短信模板
	v.Get("/sms_template", smsTemplateHaddle.SmsTemplate)
	v.Post("/sms_template/create_or_update", smsTemplateHaddle.SmsTemplateCreateOrUpdate)
	//产品运营
	v.Get("/product", productHaddle.Product)
	v.Post("/product/create_or_update", productHaddle.ProductCreateOrUpdate)
	//产品配置
	v.Get("/productDelayConfig", productDelayConfigHaddle.ProductDelayConfig)
	v.Post("/productDelayConfig/create_or_update", productDelayConfigHaddle.ProductDelayConfigCreateOrUpdate)
	//预提醒管理
	v.Get("/remind_company", remindHaddle.RemindCompany)
	v.Post("/remind_company/create", remindHaddle.RemindCompanyCreate)
	v.Post("/remind_company/update", remindHaddle.RemindCompanyUpdate)
	v.Get("/remind_group", remindHaddle.RemindGroup)
	v.Post("/remind_group/create", remindHaddle.RemindGroupCreate)
	v.Post("/remind_group/update", remindHaddle.RemindGroupUpdate)
	v.Get("/remind_admin", remindHaddle.RemindAdmin)
	v.Get("/remind_rules", remindHaddle.RemindRules)
	v.Post("/remind_rules/create_or_update", remindHaddle.RemindRulesCreateOrUpdate)
	//催收管理
	v.Get("/urge_company", urgeHaddle.UrgeCompany)
	v.Post("/urge_company/create", urgeHaddle.UrgeCompanyCreate)
	v.Post("/urge_company/update", urgeHaddle.UrgeCompanyUpdate)
	v.Get("/urge_group", urgeHaddle.UrgeGroup)
	v.Post("/urge_group/create", urgeHaddle.UrgeGroupCreate)
	v.Post("/urge_group/update", urgeHaddle.UrgeGroupUpdate)
	v.Get("/urge_admin", urgeHaddle.UrgeAdmin)
	v.Get("/urge_rules", urgeHaddle.UrgeRules)
	v.Post("/urge_rules/create_or_update", urgeHaddle.UrgeRulesCreateOrUpdate)
	app.Listen(fmt.Sprintf(":%d", global.C.Http.Port))
}
