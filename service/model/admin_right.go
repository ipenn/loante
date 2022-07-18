package model

import (
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type AdminRight struct {
	bun.BaseModel `bun:"table:admin_right,alias:ar"`
	Id	int	`json:"id"`
	RoleName	string	`json:"role_name"`
	Rights	string	`json:"-"`
	IsDelete	string	`json:"is_delete"`
	CreateTime	string	`json:"create_time"`
	UpdateTime	string	`json:"update_time"`
	IsOpenContact	string	`json:"is_open_contact"`
	IsOpenSms	string	`json:"is_open_sms"`
	IsOpenApp	string	`json:"is_open_app"`
}

func (a *AdminRight)Insert()  {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *AdminRight)Gets(where string) ([]AdminRight, int) {
	var datas []AdminRight
	count, err := global.C.DB.NewSelect().Model(&datas).Where(where).Order("id asc").ScanAndCount(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
	return datas, count
}

func (a *AdminRight)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}


func (a *AdminRight)Del(where string)  {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *AdminRight)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
