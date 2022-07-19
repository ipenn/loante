package model

import (
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type AdminLog struct {
	bun.BaseModel `bun:"table:admin_log,alias:al"`
	Id	int	`json:"id"`
	AdminId	int	`json:"admin_id"`
	AdminName	string	`json:"admin_name"`
	Method	string	`json:"method"`
	Path    string `json:"path"`
	ReqBody string `json:"req_body"`
	CreateTime	string	`json:"create_time"`
}

func (a *AdminLog)Insert()  {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *AdminLog)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *AdminLog)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *AdminLog)Gets(where string) ([]AdminLog, int) {
	var datas []AdminLog
	count, err := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
	return datas, count
}

func (a *AdminLog)Count(where string) int {
	count, err := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
	return count
}

func (a *AdminLog)Del(where string)  {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
