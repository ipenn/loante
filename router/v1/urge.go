package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type urge struct{}

func NewUrge() *urge {
	return new(urge)
}

type urgeCompanyCreate struct {
	CompanyName string `json:"companyName"`
	MchId       string `json:"mchId"`
	UserName    string `json:"userName"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

// UrgeCompanyCreate 新增催收公司
func (a *urge) UrgeCompanyCreate(c *fiber.Ctx) error {
	input := new(urgeCompanyCreate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("admin_name = '%s'", input.UserName))
	if admin.Id > 0 {
		return resp.Err(c, 1, "管理员已经存在")
	}
	admin.AdminName = input.UserName
	admin.Salt = tools.InviteCode(5)
	admin.Password = tools.Md5(fmt.Sprintf("%s%s", input.Password, admin.Salt))
	admin.RoleId = 2
	admin.MchId = tools.ToInt(input.MchId)
	admin.Mobile = input.Mobile
	admin.Email = input.Email
	admin.Insert()
	if admin.Id == 0 {
		return resp.Err(c, 500, "主管创建失败")
	}
	urgeCompany := new(model.UrgeCompany)
	urgeCompany.AdminId = admin.Id
	urgeCompany.UserName = admin.AdminName
	urgeCompany.MchId = tools.ToInt(input.MchId)
	urgeCompany.CreateTime = tools.GetFormatTime()
	urgeCompany.CompanyName = input.CompanyName
	urgeCompany.Description = input.Description
	urgeCompany.Insert()
	return resp.OK(c, "")
}

type urgeCompany struct {
	MchId string `query:"mchId" json:"mch_id"`
	Id    string `query:"id" json:"id"`
	Page  int    `query:"page" json:"page"`
	Size  int    `query:"size" json:"size"`
}

// UrgeCompany 催收公司列表
func (a *urge) UrgeCompany(c *fiber.Ctx) error {
	input := new(urgeCompany)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "uc.id>0"
	if input.MchId != "" {
		where += " and uc.mch_id=" + input.MchId
	}
	if input.Id != "" {
		where += " and uc.id =" + input.Id
	}
	lists, count := new(model.UrgeCompany).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type urgeCompanyUpdate struct {
	Id          string `json:"id"`
	AdminId     string `json:"adminId"`
	MchId       string `json:"mchId"`
	CompanyName string `json:"companyName"`
	Description string `json:"description"`
}

// UrgeCompanyUpdate 修改催收公司
func (a *urge) UrgeCompanyUpdate(c *fiber.Ctx) error {
	input := new(urgeCompanyUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}

	company := new(model.UrgeCompany)
	company.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	if company.Id == 0 {
		return resp.Err(c, 1, "没有找到催收公司")
	}
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
	if company.AdminId != tools.ToInt(input.AdminId) {
		if admin.Id == 0 {
			return resp.Err(c, 1, "没有找到角色")
		}
		//新主管添加权限
		admin.MchId = tools.ToInt(input.MchId)
		admin.RoleId = 2
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))

		//旧主管取消权限
		oldAdmin := new(model.Admin)
		oldAdmin.One(fmt.Sprintf("id=%d", company.AdminId))
		oldAdmin.RoleId = 3
		admin.Update(fmt.Sprintf("id=%d", oldAdmin.Id))
	}

	//company.Id = tools.ToInt(input.Id)
	company.AdminId = tools.ToInt(input.AdminId)
	company.MchId = tools.ToInt(input.MchId)
	company.UserName = admin.AdminName
	company.CompanyName = input.CompanyName
	company.Description = input.Description
	company.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}

type urgeAdmin struct {
	Page        int    `query:"page" json:"page"`
	Size        int    `query:"size" json:"size"`
	UrgeId      string `query:"Urge_id" json:"Urge_id"`
	UrgeGroupId string `query:"Urge_group_id" json:"Urge_group_id"`
}

// UrgeAdmin 催收员列表
func (a *urge) UrgeAdmin(c *fiber.Ctx) error {
	input := new(urgeAdmin)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := " a.id>0 and a.role_id in (2,3,4) "
	if input.UrgeGroupId != "" {
		where += fmt.Sprintf(" a.urge_group_id=%d", tools.ToInt(input.UrgeGroupId))
	}
	if input.UrgeId != "" {
		where += fmt.Sprintf(" a.urge_id=%d", tools.ToInt(input.UrgeId))
	}
	admin := new(model.Admin)
	lists, count := admin.Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type urgeGroupCreate struct {
	CompanyId string `json:"company_id"`
	MchId     string `json:"mch_id"`
	AdminId   string `json:"admin_id"`
	GroupName string `json:"group_name"`
	Status    string `json:"status"`
}

// UrgeGroupCreate 创建催收员分组
func (a *urge) UrgeGroupCreate(c *fiber.Ctx) error {
	input := new(urgeGroupCreate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	group := new(model.UrgeGroup)
	group.CompanyId = tools.ToInt(input.CompanyId)
	group.MchId = tools.ToInt(input.MchId)
	group.AdminId = tools.ToInt(input.AdminId)
	group.GroupName = input.GroupName
	group.Status = tools.ToInt(input.Status)
	group.Insert()
	if tools.ToInt(input.AdminId) != 0 {
		admin := new(model.Admin)
		admin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		admin.MchId = tools.ToInt(input.MchId)
		admin.RoleId = 4
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
	}
	return resp.OK(c, "")
}

type urgeGroupUpdate struct {
	Id        string `json:"id"`
	AdminId   string `json:"admin_id"`
	GroupName string `json:"group_name"`
	Status    string `json:"status"`
}

// UrgeGroupUpdate 修改催收员分组
func (a *urge) UrgeGroupUpdate(c *fiber.Ctx) error {
	input := new(urgeGroupUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	group := new(model.UrgeGroup)
	group.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	if tools.ToInt(input.AdminId) != 0 {
		admin := new(model.Admin)
		admin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		admin.MchId = group.MchId
		admin.RoleId = 3
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		if group.AdminId != 0 {
			oldAdmin := new(model.Admin)
			oldAdmin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
			oldAdmin.MchId = group.MchId
			oldAdmin.RoleId = 4
			oldAdmin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		}
	}
	group.AdminId = tools.ToInt(input.AdminId)
	group.GroupName = input.GroupName
	group.Status = tools.ToInt(input.Status)
	group.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}

type urgeGroup struct {
	CompanyId string `query:"company_id" json:"company_id"`
	Id        string `query:"id" json:"id"`
	Status    string `query:"status" json:"status"`
	Page      int    `query:"page" json:"page"`
	Size      int    `query:"size" json:"size"`
}

// UrgeGroup 催收员分组列表
func (a *urge) UrgeGroup(c *fiber.Ctx) error {
	input := new(urgeGroup)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "ug.id>0"
	if input.CompanyId != "" {
		where += fmt.Sprintf(" ug.company_id=%d", tools.ToInt(input.CompanyId))
	} else if input.Id != "" {
		where += fmt.Sprintf(" ug.id=%d", tools.ToInt(input.Id))
	}
	if input.Status != "" {
		where += fmt.Sprintf(" ug.status=%d", tools.ToInt(input.Status))
	}
	lists, count := new(model.UrgeGroup).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type urgeRulesCreateOrUpdate struct {
	Id          string `json:"id"`
	CompanyId   string `json:"company_id"`
	GroupId     string `json:"group_id"`
	MaxDay      string `json:"max_day"`
	MinDay      string `json:"min_day"`
	Remark      string `json:"remark"`
	IsAutoApply string `json:"is_auto_apply"`
}

// UrgeRulesCreateOrUpdate 新增或修改催收规则
func (a *urge) UrgeRulesCreateOrUpdate(c *fiber.Ctx) error {
	input := new(urgeRulesCreateOrUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	if tools.ToInt(input.MinDay) < 0 || tools.ToInt(input.MaxDay) < 0 {
		return resp.Err(c, 1, "催收规则天数不能小于0")
	}
	if tools.ToInt(input.MinDay)>tools.ToInt(input.MaxDay){
		return resp.Err(c, 1, "最小天数不能大于最大天数")
	}
	rules := new(model.UrgeRules)
	rules.CompanyId = tools.ToInt(input.CompanyId)
	rules.GroupId = tools.ToInt(input.GroupId)
	rules.MaxDay = tools.ToInt(input.MaxDay)
	rules.MinDay = tools.ToInt(input.MinDay)
	rules.Remark = input.Remark
	rules.IsAutoApply = tools.ToInt(input.IsAutoApply)
	if tools.ToInt(input.Id) == 0 {
		rules.Insert()
	} else {
		rules.Id = tools.ToInt(input.Id)
		rules.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	}
	return resp.OK(c, "")
}

type urgeRules struct {
	CompanyId   string `json:"company_id" query:"company_id"`
	GroupId     string `json:"group_id" query:"group_id"`
	IsAutoApply string `json:"is_auto_apply" query:"is_auto_apply"`
	Page        int    `json:"page" query:"page"`
	Size        int    `json:"size" query:"size"`
}

// UrgeRules 催收规则列表
func (a *urge) UrgeRules(c *fiber.Ctx) error {
	input := new(urgeRules)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "ur.id>0"
	if input.CompanyId != "" {
		where += " and ur.company_id =" + input.CompanyId
	}
	if input.GroupId != "" {
		where += " and ur.group_id =" + input.GroupId
	}
	if input.IsAutoApply != "" {
		where += " and ur.is_auto_apply =" + input.IsAutoApply
	}
	lists, count := new(model.UrgeRules).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type urgeCreateReq struct {
	Id   int `json:"id"`
	CompanyId   int `json:"company_id"`
	GroupId     int `json:"group_id"`
	AdminName string `json:"admin_name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Password string `json:"password"`
	Remark string `json:"remark"`
}
//UrgeCreate 添加催收员
func  (a *urge) UrgeCreate(c *fiber.Ctx) error {
	input := new(urgeCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//检测用户名是否重复
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("admin_name = '%s'", input.AdminName))
	if admin.Id > 0 && admin.Id != input.Id {
		return resp.Err(c, 1, "登录名已经存在")
	}
	//查询催收公司
	com := new(model.UrgeCompany)
	com.One(fmt.Sprintf("id = %d", input.CompanyId))
	if com.Id == 0{
		return resp.Err(c, 1, "催收公司不存在")
	}
	gr := new(model.UrgeGroup)
	gr.One(fmt.Sprintf("id = %d", input.GroupId))
	if gr.Id == 0{
		return resp.Err(c, 1, "催收组不存在")
	}
	if gr.CompanyId != input.CompanyId{
		return resp.Err(c, 1, "催收公司和催收组不匹配")
	}
	//准备数据
	admin.AdminName = input.AdminName
	admin.RoleId = 3
	admin.MchId = com.MchId
	admin.UrgeId =  input.CompanyId
	admin.UrgeGroupId =  input.GroupId
	admin.Mobile = input.Phone
	admin.Email = input.Email
	if input.Id > 0{
		admin.Update(fmt.Sprintf("id = %d", input.Id))
	}else{
		admin.Insert()
	}
	return resp.OK(c, "")
}