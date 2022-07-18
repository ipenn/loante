package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
	"strings"
)

type system struct {}

func NewSystem() *system {
	return new(system)
}

type adminMenu struct {
	Id int	`json:"id"`
	Name string	`json:"name"`
	Path string	`json:"path"`
	Icon string	`json:"icon"`
	ParentId int	`json:"-"`
	Children []adminMenu	`json:"routes"`
}
//SideMenu 返回后台的导航
func (a *system)SideMenu(c *fiber.Ctx) error {
	var menus []adminMenu
	var parentMenu []*adminMenu
	roleId := c.Locals("roleId").(string)
	row := new(model.AdminRight)
	row.One(fmt.Sprintf("id = '%s'", roleId))
	//获取菜单
	adminMenus := new(model.AdminMenu).Gets(fmt.Sprintf("id > 0"))
	for _, item := range adminMenus{
		if item.ParentId == 0{
			parentMenu = append(parentMenu, &adminMenu{
				Id: item.Id,
				Name: item.Name,
				Path: item.Path,
				Icon: item.Icon,
				ParentId: 0,
				Children: []adminMenu{},
			})
		}
	}
	for _, item := range adminMenus{
		for key, parent := range parentMenu{
			if item.ParentId != parent.Id{
				continue
			}
			if row.Rights != "*"{
				if strings.Index(row.Rights, item.Rights) == -1{
					continue
				}
			}
			parentMenu[key].Children = append(parentMenu[key].Children, adminMenu{
				Id: item.Id,
				Name:item.Name,
				Path:item.Path,
				Icon: item.Icon,
				ParentId: item.ParentId,
			})
		}
	}
	for _, item := range parentMenu{
		if len(item.Children) > 0{
			menus = append(menus, *item)
		}
	}
	return resp.OK(c, menus)
}


type adminReq struct {
	req.PageReq
	MchId	int	`json:"mch_id" query:"mch_id"`
	AdminName	string	`json:"admin_name" query:"admin_name"`
	RoleId	int	`json:"role_id" query:"role_id"`
	Valid	int	`json:"valid" query:"valid"`
	UserType	int	`json:"user_type" query:"user_type"`
	StartTime	string	`json:"start_time" query:"start_time"`
	EndTime	string	`json:"end_time" query:"end_time"`
}
//AdminsList 获取管理人员
func (a *system)AdminsList(c *fiber.Ctx) error {
	input := new(adminReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "a.id > 0"
	if input.MchId > 0{
		where = fmt.Sprintf("%s and a.mch_id = %d", where, input.MchId)
	}
	if len(input.AdminName) > 0{
		where = fmt.Sprintf("%s and a.admin_name = '%s'", where, input.AdminName)
	}
	lists, count := new(model.Admin).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

func (a *system)RolesList(c *fiber.Ctx) error {
	input := new(req.PageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.AdminRight).Gets("id > 0")
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type rightsResp struct {
	Id int `json:"id"`
	Name string `json:"name" query:"name"`
	Right bool `json:"right" query:"right"`
	ParentId int `json:"parent_id" query:"parent_id"`
}
func (a *system)RightsList(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//获取拥有的权限资源
	right := model.AdminRight{}
	right.One(fmt.Sprintf("id = %d", input.Id))
	if input.Id > 0 && right.Id == 0{
		return resp.Err(c, 1, "没有找到角色")
	}
	if right.Id == 1{
		return resp.Err(c, 1, "超级管理员不能修改权限")
	}
	var data []rightsResp
	menus := new(model.AdminMenu).Gets("id > 0")
	for _, item := range menus{
		rr := rightsResp{
			Id:item.Id,
			Name:item.Name,
			Right:false,
			ParentId:item.ParentId,
		}
		if len(item.Path) > 0 && strings.Index(right.Rights, item.Path) > -1{
			rr.Right = true
		}
		data = append(data, rr)
	}
	return resp.OK(c, data)
}


type adminCreateReq struct {
	AdminName string `json:"admin_name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=3,max=32"`
	MchId	int	`json:"mch_id"`
	Mobile	string	`json:"mobile"`
	Email	string	`json:"email"`
	RoleId	int	`json:"role_id"`
	Id	int	`json:"roleId"`
}
func (a *system)AdminCreate(c *fiber.Ctx) error {
	input := new(adminCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//检查唯一性
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("admin_name = '%s'", input.AdminName))
	if admin.Id > 0 && admin.Id != input.Id{
		return resp.Err(c, 1, "管理员已经存在")
	}
	admin.AdminName = input.AdminName
	admin.Password = input.Password
	admin.RoleId = input.RoleId
	admin.MchId = input.MchId
	admin.Mobile = input.Mobile
	admin.Email = input.Email
	admin.Insert()
	return resp.OK(c, "")
}

type roleCreateReq struct {
	req.IdReq
	RoleName string `json:"role_name" validate:"required,min=3,max=32"`
	Right    []int  `json:"right" validate:"required"`
}
//RoleCreate 角色创建
func (a *system)RoleCreate(c *fiber.Ctx) error {
	input := new(roleCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	right := model.AdminRight{}
	if input.Id > 0{
		right.One(fmt.Sprintf("id = %d", input.Id))
		if right.Id == 0{
			return resp.Err(c, 1, "没有找到角色")
		}
	}else{
		right.RoleName = input.RoleName
		right.One(fmt.Sprintf("role_name = '%s'", input.RoleName))
		if right.Id > 0{
			return resp.Err(c, 1, "角色名称已经存在")
		}
	}
	//处理权限码
	menus := new(model.AdminMenu).GetIds(input.Right)
	rights := ""
	for _, item := range menus{
		rights = fmt.Sprintf("%s@%s@%s", rights, item.Path, item.Rights)
	}
	//更新到权限码
	right.Rights = rights
	right.UpdateTime = tools.GetFormatTime()
	if input.Id > 0{
		right.Update(fmt.Sprintf("id = %d", input.Id))
	}else{
		right.Insert()
	}
	return resp.OK(c, "")
}

//RoleDelete 删除角色
func (a *system)RoleDelete(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	right := model.AdminRight{}
	right.One(fmt.Sprintf("id = %d", input.Id))
	if right.Id == 0{
		return resp.Err(c, 1, "没有找到角色")
	}
	if right.Id < 100{
		return resp.Err(c, 1, "系统角色不能删除")
	}
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("role_id = %d", input.Id))
	if admin.Id > 0{
		return resp.Err(c, 1, "还存在管理员不能删除")
	}
	right.Del(fmt.Sprintf("id = %d", right.Id))
	return resp.OK(c, "")
}

//SystemFields 系统特定字段
func (a *system)SystemFields(c *fiber.Ctx) error {
	return resp.OK(c, global.C.Maps)
}