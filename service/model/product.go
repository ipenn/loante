package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ProductLittle struct {
	bun.BaseModel `bun:"table:product,alias:p"`
	Id            int    `json:"id" bun:",pk"`
	AppNo         string `json:"app_no"`
	ProductName   string `json:"product_name"`
}
type Product struct {
	bun.BaseModel `bun:"table:product,alias:p"`
	ProductLittle
	IconPath            string          `json:"icon_path"`
	MchId               int             `json:"mch_id"`
	DayMaxApply         int             `json:"day_max_apply"`
	MaxApply            int             `json:"max_apply"`
	DayApplyPass        int             `json:"day_apply_pass"`
	StartAmount         int             `json:"start_amount"`
	TodayApplyCount     int             `json:"today_apply_count"`
	ApplyStartTime      string          `json:"apply_start_time"`
	ApplyEndTime        string          `json:"apply_end_time"`
	TotalMaxApplyCount  int             `json:"total_max_apply_count"`
	UpTime              string          `json:"up_time"`
	DownTime            string          `json:"down_time"`
	IsAutoLending       int             `json:"is_auto_lending"`
	IsRejectNew         int             `json:"is_reject_new"`
	IsRejectOld         int             `json:"is_reject_old"`
	IsStopLending       int             `json:"is_stop_lending"`
	Status              int             `json:"status"`
	CreateTime          string          `json:"create_time"`
	Description         string          `json:"description"`
	RateNormalInterest  float64         `json:"rate_normal_interest"`
	RateOverdueInterest float64         `json:"rate_overdue_interest"`
	RateService         float64         `json:"rate_service"`
	RateTax             float64         `json:"rate_tax"`
	Merchant            *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
}

func (a *Product) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Product) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Product) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Product) Page(where string, page, limit int) ([]Product, int) {
	var d []Product
	count, _ := global.C.DB.NewSelect().Model(&d).Relation("Merchant").Where(where).Order(fmt.Sprintf("p.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return d, count
}
