package v1

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type auth struct {
}

func NewAuth() *auth {
	return new(auth)
}

type loginForm struct {
	AdminName string `json:"admin_name"`
	Password string `json:"password"`
}

func (a *auth)Login(c *fiber.Ctx) error {
	input := new(loginForm)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	admin := new(model.Admin)
	admin.One(fmt.Sprintf("admin_name = '%s' and password = md5(CONCAT('%s', salt))", input.AdminName, input.Password))
	if admin.Id == 0{
		return resp.Err(c, 1, "用户名密码不匹配")
	}
	if admin.Status == 0{
		return resp.Err(c, 1, "账户被冻结")
	}

	//获取role_type
	right := new(model.AdminRight)
	right.One(fmt.Sprintf("id = %d", admin.RoleId))
	if right.Id == 0{
		return resp.Err(c, 1, "没有找到对应权限")
	}
	//生成 TOKEN
	token := jwt.New(jwt.SigningMethodHS256)
	exp := tools.GetUnixTime() + int64(global.C.Safety.Expired)
	claims := token.Claims.(jwt.MapClaims)
	claims["AdminName"] = admin.AdminName
	claims["Id"] = admin.Id
	claims["RoleId"] = fmt.Sprintf("%d", admin.RoleId)
	claims["AdminType"] = right.RoleName
	claims["MchId"] = admin.MchId
	claims["RemindId"] = admin.RemindId
	claims["UrgeId"] = admin.UrgeId
	claims["RemindGroupId"] = admin.RemindGroupId
	claims["UrgeGroupId"] = admin.UrgeGroupId
	claims["Exp"] = exp
	t, err := token.SignedString([]byte(global.C.Safety.Secret))
	if err != nil {
		return resp.Err(c, 1, "签名失败！")
	}
	return resp.OK(c,map[string]interface{}{
		"admin_name":admin.AdminName,
		"token":t,
		"exp":exp,
		"admin_type":right.RoleName,
		"role_id": admin.RoleId,
		"mch_id": admin.MchId,
	})
}
