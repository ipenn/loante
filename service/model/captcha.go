package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type Captcha struct {
	bun.BaseModel `bun:"table:captchas,alias:c"`
	Id	int	`json:"id"`
	Address		string		`json:"address"`
	Code		string		`json:"code"`
	CreateTime		string		`json:"create_time"`
	CreateAt		int64		`json:"create_at"`
	Used		int		`json:"used"`

}

func (a *Captcha)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Captcha)Update(col string,where string)  {
	global.C.DB.NewUpdate().Model(a).Column(col).Where(where).Exec(global.C.Ctx)
}

func (a *Captcha)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Captcha)Gets(where string) ([]Captcha, int) {
	var datas []Captcha
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Captcha)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Captcha)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}


//Check 判断验证码是否正确
func (a *Captcha)Check(phone, code string) bool {
	where := fmt.Sprintf("address = '%s'", phone)
	global.C.DB.NewSelect().Model(a).Where(where).Order("id desc").Scan(global.C.Ctx)
	if a.Id == 0 || a.Used == 1 || a.Code != code || tools.GetUnixTime() - 10*60 > a.CreateAt{
		return false
	}
	return true
}