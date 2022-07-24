package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type PackageLittle struct {
	bun.BaseModel  `bun:"table:package,alias:b"`
	Id                     int    `json:"id" bun:",pk"`
	Name                   string	`json:"name"`
	AppNo                  string	`json:"app_no"`
	AppType                string	`json:"app_type"`
	Status                 string	`json:"status"`
}

type Package struct {
	bun.BaseModel          `bun:"table:package,alias:p"`
	PackageLittle
	PayReturnUrl          string	`json:"pay_return_url"`
	UpdateH5Url           string	`json:"update_h5_url"`
	RepaymentInfoUrl      string	`json:"repayment_info_url"`
	VersionCode          string	`json:"version_code"`
	Version              string	`json:"version"`
	IsMandatoryUpdate    string	`json:"is_mandatory_update"`
	IsNeedUpdate         string	`json:"is_need_update"`
	Whatsapp             string	`json:"whatsapp"`
	AppGpUrl             string	`json:"app_gp_url"`
	CurrentUrl           string	`json:"current_url"`
	UpdateInfo           string	`json:"update_info"`
	Firebase             string	`json:"firebase"`
	RegisterAgreementUrl string	`json:"register_agreement_url"`
	PrivacyAgreementUrl  string	`json:"privacy_agreement_url"`
	FacebookId           string	`json:"facebook_id"`
	FacebookKey          string	`json:"facebook_key"`
	UpdateTime           string	`json:"update_time"`
	Remark               string	`json:"remark"`
	CreateTime           string	`json:"create_time"`
}


func (a *Package)Insert()  {
	global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
}

func (a *Package)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Package)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Package)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Package)Page(where string, page, limit int) ([]Package, int) {
	var datas []Package
	count, _ := global.C.DB.NewSelect().Model(&datas).
		Where(where).Order(fmt.Sprintf("id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *Package)PageLittles(where string) ([]PackageLittle, int) {
	var datas []PackageLittle
	count, _ := global.C.DB.NewSelect().Model(&datas).
		Where(where).Order(fmt.Sprintf("id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}