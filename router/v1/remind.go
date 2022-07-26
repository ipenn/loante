package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type remind struct{}

func NewRemind() *remind {
	return new(remind)
}

type remindCompanyCreate struct {
	CompanyName string `json:"companyName"`
	MchId       string `json:"mchId"`
	UserName    string `json:"userName"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

// RemindCompanyCreate 创建预提醒公司
func (a *remind) RemindCompanyCreate(c *fiber.Ctx) error {
	input := new(remindCompanyCreate)
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
	admin.RoleId = 5
	admin.MchId = tools.ToInt(input.MchId)
	admin.Mobile = input.Mobile
	admin.Email = input.Email
	admin.Insert()
	if admin.Id == 0 {
		return resp.Err(c, 500, "主管创建失败")
	}
	remindCompany := new(model.RemindCompany)
	remindCompany.AdminId = admin.Id
	remindCompany.UserName = admin.AdminName
	remindCompany.MchId = tools.ToInt(input.MchId)
	remindCompany.CreateTime = tools.GetFormatTime()
	remindCompany.CompanyName = input.CompanyName
	remindCompany.Description = input.Description
	remindCompany.Insert()
	return resp.OK(c, "")
}

type remindCompanyList struct {
	MchId string `query:"mchId" json:"mch_id"`
	Id    string `query:"id" json:"id"`
	Page  int    `query:"page" json:"page"`
	Size  int    `query:"size" json:"size"`
}

// RemindCompany 预提醒公司列表
func (a *remind) RemindCompany(c *fiber.Ctx) error {
	input := new(remindCompanyList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "rc.id>0"
	if input.MchId != "" {
		where += " and rc.mch_id=" + input.MchId
	}
	if input.Id != "" {
		where += " and rc.id =" + input.Id
	}
	lists, count := new(model.RemindCompany).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type RemindCompanyUpdate struct {
	Id          string `json:"id"`
	AdminId     string `json:"adminId"`
	MchId       string `json:"mchId"`
	CompanyName string `json:"companyName"`
	Description string `json:"description"`
}

// RemindCompanyUpdate 预提醒公司修改
func (a *remind) RemindCompanyUpdate(c *fiber.Ctx) error {
	input := new(RemindCompanyUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}

	remindCompany := new(model.RemindCompany)
	remindCompany.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	if remindCompany.Id == 0 {
		return resp.Err(c, 1, "没有找到预提醒公司")
	}
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
	if remindCompany.AdminId != tools.ToInt(input.AdminId) {
		if admin.Id == 0 {
			return resp.Err(c, 1, "没有找到角色")
		}
		//新主管添加权限
		admin.MchId = tools.ToInt(input.MchId)
		admin.RoleId = 5
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))

		//旧主管取消权限
		oldAdmin := new(model.Admin)
		oldAdmin.One(fmt.Sprintf("id=%d", remindCompany.AdminId))
		oldAdmin.RoleId = 7
		admin.Update(fmt.Sprintf("id=%d", oldAdmin.Id))
	}

	remindCompany.Id = tools.ToInt(input.Id)
	remindCompany.AdminId = tools.ToInt(input.AdminId)
	remindCompany.MchId = tools.ToInt(input.MchId)
	remindCompany.UserName = admin.AdminName
	remindCompany.CompanyName = input.CompanyName
	remindCompany.Description = input.Description
	remindCompany.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}

type remindAdmin struct {
	Page          int    `query:"page" json:"page"`
	Size          int    `query:"size" json:"size"`
	RemindId      string `query:"remind_id" json:"remind_id"`
	RemindGroupId string `query:"remind_group_id" json:"remind_group_id"`
}

// RemindAdmin 预提醒专员列表
func (a *remind) RemindAdmin(c *fiber.Ctx) error {
	input := new(remindAdmin)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := " a.id>0 and a.role_id in (5,6,7) "
	if input.RemindGroupId != "" {
		where += fmt.Sprintf(" and a.remind_group_id=%d", tools.ToInt(input.RemindGroupId))
	}
	if input.RemindId != "" {
		where += fmt.Sprintf(" and a.remind_id=%d", tools.ToInt(input.RemindId))
	}
	admin := new(model.Admin)
	lists, count := admin.Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type remindGroupCreate struct {
	CompanyId string `json:"company_id"`
	MchId     string `json:"mch_id"`
	AdminId   string `json:"admin_id"`
	GroupName string `json:"group_name"`
	Status    string `json:"status"`
}

// RemindGroupCreate 创建预提醒分组
func (a *remind) RemindGroupCreate(c *fiber.Ctx) error {
	input := new(remindGroupCreate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	group := new(model.RemindGroup)
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
		admin.RoleId = 6
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
	}
	return resp.OK(c, "")
}

type remindGroupUpdate struct {
	Id        string `json:"id"`
	AdminId   string `json:"admin_id"`
	GroupName string `json:"group_name"`
	Status    string `json:"status"`
}

// RemindGroupUpdate 修改预提醒分组
func (a *remind) RemindGroupUpdate(c *fiber.Ctx) error {
	input := new(remindGroupUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	group := new(model.RemindGroup)
	group.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	if tools.ToInt(input.AdminId) != 0 {
		admin := new(model.Admin)
		admin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		admin.MchId = group.MchId
		admin.RoleId = 6
		admin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		if group.AdminId != 0 {
			oldAdmin := new(model.Admin)
			oldAdmin.One(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
			oldAdmin.MchId = group.MchId
			oldAdmin.RoleId = 7
			oldAdmin.Update(fmt.Sprintf("id=%d", tools.ToInt(input.AdminId)))
		}
	}
	group.AdminId = tools.ToInt(input.AdminId)
	group.GroupName = input.GroupName
	group.Status = tools.ToInt(input.Status)
	group.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}

type remindGroup struct {
	CompanyId string `query:"company_id" json:"company_id"`
	Id        string `query:"id" json:"id"`
	Status    string `query:"status" json:"status"`
	Page      int    `query:"page" json:"page"`
	Size      int    `query:"size" json:"size"`
}

// RemindGroup 预提醒分组列表
func (a *remind) RemindGroup(c *fiber.Ctx) error {
	input := new(remindGroup)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "rg.id>0"
	if input.CompanyId != "" {
		where += fmt.Sprintf(" rg.company_id=%d", tools.ToInt(input.CompanyId))
	} else if input.Id != "" {
		where += fmt.Sprintf(" rg.id=%d", tools.ToInt(input.Id))
	}
	if input.Status != "" {
		where += fmt.Sprintf(" rg.status=%d", tools.ToInt(input.Status))
	}
	lists, count := new(model.RemindGroup).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type remindRulesCreateOrUpdate struct {
	Id          string `json:"id"`
	CompanyId   string `json:"company_id"`
	GroupId     string `json:"group_id"`
	MaxDay      string `json:"max_day"`
	MinDay      string `json:"min_day"`
	Remark      string `json:"remark"`
	IsAutoApply string `json:"is_auto_apply"`
}

// RemindRulesCreateOrUpdate 预提醒规则新建或修改
func (a *remind) RemindRulesCreateOrUpdate(c *fiber.Ctx) error {
	input := new(remindRulesCreateOrUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	if tools.ToInt(input.MinDay) > 0 || tools.ToInt(input.MaxDay) > 0 {
		return resp.Err(c, 1, "预提醒规则天数不能大于0")
	}
	if tools.ToInt(input.MinDay)>tools.ToInt(input.MaxDay){
		return resp.Err(c, 1, "最小天数不能大于最大天数")
	}
	rules := new(model.RemindRules)
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

type remindRules struct {
	CompanyId   string `json:"company_id"`
	GroupId     string `json:"group_id"`
	IsAutoApply string `json:"is_auto_apply"`
	Page        int    `json:"page"`
	Size        int    `json:"size"`
}

// RemindRules 预提醒规则列表
func (a *remind) RemindRules(c *fiber.Ctx) error {
	input := new(remindRules)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "rr.id>0"
	if input.CompanyId != "" {
		where += " and rr.company_id =" + input.CompanyId
	}
	if input.GroupId != "" {
		where += " and rr.group_id =" + input.GroupId
	}
	if input.IsAutoApply != "" {
		where += " and rr.is_auto_apply =" + input.IsAutoApply
	}
	lists, count := new(model.RemindRules).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//RemindCreate 添加预提醒专员
func  (a *remind) RemindCreate(c *fiber.Ctx) error {
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
	com := new(model.RemindCompany)
	com.One(fmt.Sprintf("id = %d", input.CompanyId))
	if com.Id == 0{
		return resp.Err(c, 1, "预提醒公司不存在")
	}
	gr := new(model.RemindGroup)
	gr.One(fmt.Sprintf("id = %d", input.GroupId))
	if gr.Id == 0{
		return resp.Err(c, 1, "预提醒组不存在")
	}
	if gr.CompanyId != input.CompanyId{
		return resp.Err(c, 1, "预提醒公司和预提醒组不匹配")
	}
	//准备数据
	admin.AdminName = input.AdminName
	admin.RoleId = 3
	admin.MchId = com.MchId
	admin.RemindId =  input.CompanyId
	admin.RemindGroupId =  input.GroupId
	admin.Mobile = input.Phone
	admin.Email = input.Email
	if input.Id > 0{
		admin.Update(fmt.Sprintf("id = %d", input.Id))
	}else{
		admin.Insert()
	}
	return resp.OK(c, "")
}