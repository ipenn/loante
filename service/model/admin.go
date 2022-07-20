package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type AdminLittle struct {
	bun.BaseModel `bun:"table:admin,alias:a"`
	Id	int	`json:"id" bun:",pk"`
	AdminName	string	`json:"admin_name"`
}

type Admin struct {
	bun.BaseModel `bun:"table:admin,alias:a"`
	Id	int	`json:"id" bun:",pk"`
	AdminName	string	`json:"admin_name"`
	Password	string	`json:"-"`
	Salt	string	`json:"-"`
	CreateTime	string	`json:"create_time"`
	Status	int	`json:"status"`
	RoleId	int	`json:"role_id"`
	LoginTime	string	`json:"login_time"`
	LoginIp	string	`json:"login_ip"`
	Email	string	`json:"email"`
	Mobile	string	`json:"mobile"`
	MchId	int	`json:"-" bun:"mch_id"`
	RemindId	int	`json:"-"`
	UrgeId        int `json:"-"`
	RemindGroupId int             `json:"-"`
	UrgeGroupId   int             `json:"-"`
	Merchant      *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	RemindGroup   	*RemindGroup   `json:"remind_group" bun:"rel:belongs-to,join:remind_group_id=id"`
	RemindCompany 	*RemindCompany `json:"remind_company" bun:"rel:belongs-to,join:remind_id=id"`
	UrgeGroup  	 	*UrgeGroup   `json:"urge_group" bun:"rel:belongs-to,join:urge_group_id=id"`
	Company 		*UrgeCompany `json:"urge_company" bun:"rel:belongs-to,join:urge_id=id"`
}

func (a *Admin)Insert()  {
	a.CreateTime = tools.GetFormatTime()
	a.Salt = tools.InviteCode(5)
	a.Password = tools.Md5(fmt.Sprintf("%s%s", a.Password, a.Salt))
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Admin)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Admin)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Admin)Page(where string, page, limit int) ([]Admin, int) {
	var datas []Admin
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("RemindGroup").Relation("RemindCompany").Relation("UrgeGroup").Relation("UrgeCompany").Where(where).Order(fmt.Sprintf("a.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}