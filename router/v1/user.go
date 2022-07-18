package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type user struct {
	
}

func NewUser() *user {
	return new(user)
}

type userQueryReq struct {
	req.PageReq
	StartDate      string		`query:"start_date" json:"start_date"`
	EndDate        string		`query:"end_date" json:"end_date"`
	Name   string		`query:"name" json:"name"`
	Mobile string		`query:"mobile" json:"mobile"`
	Sex    int		`query:"sex" json:"sex"`
	Old            bool		`query:"old" json:"old"`
	TrafficId      int		`query:"traffic_id" json:"traffic_id"`
	MarketId       int		`query:"market_id" json:"market_id"`
	Id             string		`query:"id" json:"id"`
}

func (a *user)UserQuery(c *fiber.Ctx) error {
	input := new(userQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.User).Page("u.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}


func (a *user)Details(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	info := new(model.User)
	info.One(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c,map[string]interface{}{
		"info":info,
	})
}
//
//type setUserInfoReq struct {
//	TrueName string `json:"true_name"`
//	IdNo string `json:"id_no"`
//	IdPath string `json:"id_path"`
//}
//func (a *user)SetUpload(c *fiber.Ctx) error {
//	uId := c.Locals("userId").(int)
//	userOne := new(model.User)
//	userOne.One(fmt.Sprintf("id = %d ", uId))
//	filePath := ""
//	if userOne.Id > 0{
//		path := fmt.Sprintf("static/upload/%d/", uId)
//		os.MkdirAll(path, os.ModePerm)
//		file, _ := c.FormFile("file")
//		fileArr := strings.Split(file.Filename, ".")
//		ext := fileArr[len(fileArr)-1]
//		fileName := tools.Md5(tools.GetFormatTime() +file.Filename )
//		filePath = fmt.Sprintf("%s%s.%s", path, fileName, ext)
//		c.SaveFile(file, filePath)
//		filePath = "/"+filePath
//	}
//	return  resp.OK(c,filePath)
//}
//func (a *user)SetUserInfo(c *fiber.Ctx) error {
//	input := new(setUserInfoReq)
//	if err := tools.ParseBody(c, input); err != nil {
//		return resp.Err(c, 1, err.Error())
//	}
//
//	uId := c.Locals("userId").(int)
//	userOne := new(model.UserProfile)
//	userOne.One(fmt.Sprintf("user_id = %d ", uId))
//	if userOne.AuthState == 1{
//		return resp.Err(c, 1, "正在审核中")
//	}
//	if userOne.AuthState == 2{
//		return resp.Err(c, 1, "认证已经完成")
//	}
//	if userOne.UserId > 0{
//		//userOne := new(model.UserProfile)
//		userOne.TrueName = input.TrueName
//		userOne.IdNo = input.IdNo
//		userOne.IdPath = input.IdPath
//		userOne.AuthState = 1
//		userOne.Update("true_name,id_no,id_path,auth_state", fmt.Sprintf("user_id = %d", uId))
//	}
//	return resp.OK(c,userOne)
//}
//
//type setUserBindReq struct {
//	HcnId string `json:"hcn_id"`
//}
//func (a *user)SetUserBind(c *fiber.Ctx) error {
//	input := new(setUserBindReq)
//	if err := tools.ParseBody(c, input); err != nil {
//		return resp.Err(c, 1, err.Error())
//	}
//
//	uId := c.Locals("userId").(int)
//	userOne := new(model.UserProfile)
//	userOne.One(fmt.Sprintf("user_id = %d ", uId))
//	if userOne.AuthState != 2{
//		return resp.Err(c, 1, "请先完成实名认证")
//	}
//	if userOne.HcnState == 1{
//		return resp.Err(c, 1, "正在审核中")
//	}
//	if userOne.HcnState == 2{
//		return resp.Err(c, 1, "认证已经完成")
//	}
//
//	if userOne.UserId > 0{
//		//userOne := new(model.UserProfile)
//		userOne.HcnState = 1
//		userOne.HcnId = input.HcnId
//		userOne.Update("true_name,id_no,id_path,auth_state", fmt.Sprintf("user_id = %d", uId))
//	}
//	return resp.OK(c,userOne)
//}
//
//type couponUseIdReq struct {
//	Id int `json:"id"`
//}
//func (a *user)Coupon(c *fiber.Ctx) error {
//
//	uId := c.Locals("userId").(int)
//	lists, _ := new(model.Coupon).Gets(fmt.Sprintf("user_id = %d ", uId))
//	return resp.OK(c, map[string]interface{}{
//		"lists":lists,
//	})
//}
//
//func (a *user)CouponUse(c *fiber.Ctx) error {
//	uId := c.Locals("userId").(int)
//	input := new(couponUseIdReq)
//	if err := tools.ParseBody(c, input); err != nil {
//		return resp.Err(c, 1, err.Error())
//	}
//
//	coupon := new(model.Coupon)
//	coupon.One(fmt.Sprintf("user_id = %d and id = %d", uId, input.Id))
//	if coupon.Id == 0{
//		return resp.Err(c, 1, "没有找到记录")
//	}
//	if coupon.Status != 0{
//		return resp.Err(c, 1, "已经被使用了")
//	}
//	t := tools.GetFormatTime()
//	coupon.Status = 1
//	coupon.UseTime = t
//	coupon.Update(fmt.Sprintf("id = %d", input.Id))
//
//	return resp.OK(c, map[string]interface{}{
//	})
//}
//
//
//type couponAssignIdReq struct {
//	Id int `json:"id"`
//	Name string `json:"name"`
//}
//func (a *user)CouponAssign(c *fiber.Ctx) error {
//	uId := c.Locals("userId").(int)
//	input := new(couponAssignIdReq)
//	if err := tools.ParseBody(c, input); err != nil {
//		return resp.Err(c, 1, err.Error())
//	}
//
//	coupon := new(model.Coupon)
//	coupon.One(fmt.Sprintf("user_id = %d and id = %d", uId, input.Id))
//	if coupon.Id == 0{
//		return resp.Err(c, 1, "没有找到记录")
//	}
//	if coupon.Status != 0{
//		return resp.Err(c, 1, "已经被使用了")
//	}
//	//判断用户是否存在
//	userOne := new(model.User)
//	userOne.One(fmt.Sprintf("username='%s'", input.Name))
//	if userOne.Id == 0{
//		return resp.Err(c, 1, "未找到用户")
//	}
//	t := tools.GetFormatTime()
//	coupon.Status = 2
//	coupon.UseTime = t
//	coupon.Update(fmt.Sprintf("id = %d", input.Id))
//
//	coupon.Id = 0
//	coupon.FromUserId = uId
//	coupon.UserId = userOne.Id
//	coupon.UseTime =""
//	coupon.Status =0
//	return resp.OK(c, map[string]interface{}{
//	})
//}
